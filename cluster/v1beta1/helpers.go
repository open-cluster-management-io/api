package v1beta1

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	v1 "open-cluster-management.io/api/cluster/v1"
)

type ManagedClustersGetter interface {
	List(selector labels.Selector) (ret []*v1.ManagedCluster, err error)
}

type ManagedClusterSetsGetter interface {
	List(selector labels.Selector) (ret []*ManagedClusterSet, err error)
}

type ManagedClusterSetBindingsGetter interface {
	List(namespace string, selector labels.Selector) (ret []*ManagedClusterSetBinding, err error)
}

// GetClustersFromClusterSet return the ManagedClusterSet's managedClusters
func GetClustersFromClusterSet(clusterSet *ManagedClusterSet,
	clustersGetter ManagedClustersGetter) ([]*v1.ManagedCluster, error) {
	var clusters []*v1.ManagedCluster

	if clusterSet == nil {
		return nil, nil
	}

	clusterSelector, err := BuildClusterSelector(clusterSet)
	if err != nil {
		return nil, err
	}
	if clusterSelector == nil {
		return nil, fmt.Errorf("failed to build ClusterSelector with clusterSet: %v", clusterSet)
	}
	clusters, err = clustersGetter.List(clusterSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to list ManagedClusters: %w", err)
	}
	return clusters, nil
}

// GetClusterSetsOfClusterByCluster return the managedClusterSets of a managedCluster
func GetClusterSetsOfCluster(cluster *v1.ManagedCluster,
	clusterSetsGetter ManagedClusterSetsGetter) ([]*ManagedClusterSet, error) {
	var returnClusterSets []*ManagedClusterSet

	if cluster == nil {
		return nil, nil
	}

	allClusterSets, err := clusterSetsGetter.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, clusterSet := range allClusterSets {
		clusterSelector, err := BuildClusterSelector(clusterSet)
		if err != nil {
			return nil, err
		}
		if clusterSelector == nil {
			return nil, fmt.Errorf("failed to build ClusterSelector with clusterSet: %v", clusterSet)
		}
		if clusterSelector.Matches(labels.Set(cluster.Labels)) {
			returnClusterSets = append(returnClusterSets, clusterSet)
		}
	}
	return returnClusterSets, nil
}

func BuildClusterSelector(clusterSet *ManagedClusterSet) (labels.Selector, error) {
	if clusterSet == nil {
		return nil, nil
	}
	selectorType := clusterSet.Spec.ClusterSelector.SelectorType

	switch selectorType {
	case "", LegacyClusterSetLabel:
		return labels.SelectorFromSet(labels.Set{
			ClusterSetLabel: clusterSet.Name,
		}), nil
	case LabelSelector:
		return metav1.LabelSelectorAsSelector(clusterSet.Spec.ClusterSelector.LabelSelector)
	default:
		return nil, fmt.Errorf("selectorType is not right: %s", clusterSet.Spec.ClusterSelector.SelectorType)
	}
}

// GetBoundManagedClusterSetBindings returns all bindings that are bounded to clustersets in the given namespace.
func GetBoundManagedClusterSetBindings(namespace string,
	clusterSetBindingsGetter ManagedClusterSetBindingsGetter) ([]*ManagedClusterSetBinding, error) {
	// get all clusterset bindings under the namespace
	bindings, err := clusterSetBindingsGetter.List(namespace, labels.Everything())
	if err != nil {
		return nil, err
	}

	boundBindings := []*ManagedClusterSetBinding{}
	for _, binding := range bindings {
		if meta.IsStatusConditionTrue(binding.Status.Conditions, ClusterSetBindingBoundType) {
			boundBindings = append(boundBindings, binding)
		}
	}

	return boundBindings, nil
}

type PlacementDecisionGetter interface {
	List(selector labels.Selector, namespace string) (ret []*PlacementDecision, err error)
}

// +k8s:deepcopy-gen=false
type PlacementDecisionClustersTracker struct {
	placement                      *Placement
	placementDecisionGetter        PlacementDecisionGetter
	existingScheduledClusters      sets.Set[string]
	existingScheduledClusterGroups map[GroupKey]sets.Set[string]
	clustesGroupsIndexToName       map[int32]string
	clustesGroupsNameToIndex       map[string][]int32
	lock                           sync.RWMutex
}

// +k8s:deepcopy-gen=false
type GroupKey struct {
	GroupName  string `json:"groupName,omitempty"`
	GroupIndex int32  `json:"groupIndex,omitempty"`
}

func NewPlacementDecisionClustersTracker(placement *Placement, pdl PlacementDecisionGetter, existingScheduledClusters sets.Set[string], existingScheduledClusterGroups map[GroupKey]sets.Set[string]) *PlacementDecisionClustersTracker {
	pdct := &PlacementDecisionClustersTracker{
		placement:                      placement,
		placementDecisionGetter:        pdl,
		existingScheduledClusters:      existingScheduledClusters,
		existingScheduledClusterGroups: existingScheduledClusterGroups,
	}
	pdct.generateGroupsNameIndex()
	return pdct
}

// Get updates the tracker's decisionClusters and returns added and deleted cluster names.
func (pdct *PlacementDecisionClustersTracker) Get() (sets.Set[string], sets.Set[string], error) {
	pdct.lock.Lock()
	defer pdct.lock.Unlock()

	if pdct.placement == nil || pdct.placementDecisionGetter == nil {
		return nil, nil, nil
	}

	// Get the generated PlacementDecisions
	decisionSelector := labels.SelectorFromSet(labels.Set{
		PlacementLabel: pdct.placement.Name,
	})
	decisions, err := pdct.placementDecisionGetter.List(decisionSelector, pdct.placement.Namespace)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list PlacementDecisions: %w", err)
	}

	// Get the decision cluster names and groups
	newScheduledClusters := sets.New[string]()
	newScheduledClusterGroups := map[GroupKey]sets.Set[string]{}
	for _, d := range decisions {
		groupKey, err := parseGroupKeyFromDecision(d)
		if err != nil {
			return nil, nil, err
		}

		if _, exist := newScheduledClusterGroups[groupKey]; !exist {
			newScheduledClusterGroups[groupKey] = sets.New[string]()
		}

		for _, sd := range d.Status.Decisions {
			newScheduledClusters.Insert(sd.ClusterName)
			newScheduledClusterGroups[groupKey].Insert(sd.ClusterName)
		}
	}

	// Compare the difference
	added := newScheduledClusters.Difference(pdct.existingScheduledClusters)
	deleted := pdct.existingScheduledClusters.Difference(newScheduledClusters)

	// Update the existing decision cluster names and groups
	pdct.existingScheduledClusters = newScheduledClusters
	pdct.existingScheduledClusterGroups = newScheduledClusterGroups
	pdct.generateGroupsNameIndex()

	return added, deleted, nil
}

func (pdct *PlacementDecisionClustersTracker) generateGroupsNameIndex() {
	pdct.clustesGroupsIndexToName = map[int32]string{}
	pdct.clustesGroupsNameToIndex = map[string][]int32{}

	for groupkey := range pdct.existingScheduledClusterGroups {
		// index to name
		pdct.clustesGroupsIndexToName[groupkey.GroupIndex] = groupkey.GroupName
		// name to index
		if index, exist := pdct.clustesGroupsNameToIndex[groupkey.GroupName]; exist {
			pdct.clustesGroupsNameToIndex[groupkey.GroupName] = append(index, groupkey.GroupIndex)
		} else {
			pdct.clustesGroupsNameToIndex[groupkey.GroupName] = []int32{groupkey.GroupIndex}
		}
	}

	// sort index order
	for _, index := range pdct.clustesGroupsNameToIndex {
		sort.Slice(index, func(i, j int) bool {
			return index[i] < index[j]
		})
	}
}

// Existing() returns the tracker's existing decision cluster names of groups listed in groupKeys.
// Return empty set when groupKeys is empty.
func (pdct *PlacementDecisionClustersTracker) Existing(groupKeys []GroupKey) (sets.Set[string], []GroupKey, map[GroupKey]sets.Set[string]) {
	existingGroupKeys, existingClusterGroups := pdct.ExistingClusterGroups(groupKeys)

	existingClusters := sets.New[string]()
	for _, clusterSets := range existingClusterGroups {
		existingClusters.Insert(clusterSets.UnsortedList()...)
	}

	return existingClusters, existingGroupKeys, existingClusterGroups
}

// ExistingBesides returns the tracker's existing decision cluster names except cluster groups listed in groupKeys.
// Return all the clusters when groupKeys is empty.
func (pdct *PlacementDecisionClustersTracker) ExistingBesides(groupKeys []GroupKey) (sets.Set[string], []GroupKey, map[GroupKey]sets.Set[string]) {
	existingGroupKeys, existingClusterGroups := pdct.ExistingClusterGroupsBesides(groupKeys)

	existingClusters := sets.New[string]()
	for _, clusterSets := range existingClusterGroups {
		existingClusters.Insert(clusterSets.UnsortedList()...)
	}

	return existingClusters, existingGroupKeys, existingClusterGroups
}

// ExistingClusterGroups returns the tracker's existing decision cluster names for groups listed in groupKeys.
// Return empty set when groupKeys is empty.
func (pdct *PlacementDecisionClustersTracker) ExistingClusterGroups(groupKeys []GroupKey) ([]GroupKey, map[GroupKey]sets.Set[string]) {
	pdct.lock.RLock()
	defer pdct.lock.RUnlock()

	resultClusterGroups := make(map[GroupKey]sets.Set[string])
	resultGroupKeys := []GroupKey{}

	includeGroupKeys := pdct.fulfillGroupKeys(groupKeys)
	for _, groupKey := range includeGroupKeys {
		if clusters, found := pdct.existingScheduledClusterGroups[groupKey]; found {
			resultClusterGroups[groupKey] = clusters
			resultGroupKeys = append(resultGroupKeys, groupKey)
		}
	}

	return resultGroupKeys, resultClusterGroups
}

// ExistingClusterGroupsBesides returns the tracker's existing decision cluster names except cluster groups listed in groupKeys.
// Return all the clusters when groupKeys is empty.
func (pdct *PlacementDecisionClustersTracker) ExistingClusterGroupsBesides(groupKeys []GroupKey) ([]GroupKey, map[GroupKey]sets.Set[string]) {
	pdct.lock.RLock()
	defer pdct.lock.RUnlock()

	resultClusterGroups := make(map[GroupKey]sets.Set[string])
	resultGroupKeys := []GroupKey{}

	excludeGroupKeys := pdct.fulfillGroupKeys(groupKeys)
	includeGroupKeys := pdct.getOrderedGroupKeysBesides(excludeGroupKeys)
	for _, groupKey := range includeGroupKeys {
		if clusters, found := pdct.existingScheduledClusterGroups[groupKey]; found {
			resultClusterGroups[groupKey] = clusters
			resultGroupKeys = append(resultGroupKeys, groupKey)
		}
	}

	return resultGroupKeys, resultClusterGroups
}

func (pdct *PlacementDecisionClustersTracker) ClusterToGroupKey(groupToCluster map[GroupKey]sets.Set[string]) map[string]GroupKey {
	clusterToGroupKey := map[string]GroupKey{}

	for groupKey, clustersSet := range groupToCluster {
		for c := range clustersSet {
			clusterToGroupKey[c] = groupKey
		}
	}

	return clusterToGroupKey
}

// Fulfill the expect groupkeys with group name or group index, the returned groupkeys are ordered by input group name then group index.
// For example, the input is []GroupKey{{GroupName: "group1"}, {GroupIndex: 2}},
// the returned is []GroupKey{{GroupName: "group1", GroupIndex: 0}, {GroupName: "group1", GroupIndex: 1}, {GroupName: "group2", GroupIndex: 2}}
func (pdct *PlacementDecisionClustersTracker) fulfillGroupKeys(groupKeys []GroupKey) []GroupKey {
	fulfilledGroupKeys := []GroupKey{}
	for _, gk := range groupKeys {
		if gk.GroupName != "" {
			if indexes, exist := pdct.clustesGroupsNameToIndex[gk.GroupName]; exist {
				for _, groupIndex := range indexes {
					fulfilledGroupKeys = append(fulfilledGroupKeys, GroupKey{GroupName: gk.GroupName, GroupIndex: groupIndex})
				}
			}
		} else {
			if groupName, exist := pdct.clustesGroupsIndexToName[gk.GroupIndex]; exist {
				fulfilledGroupKeys = append(fulfilledGroupKeys, GroupKey{GroupName: groupName, GroupIndex: gk.GroupIndex})
			}
		}
	}
	return fulfilledGroupKeys
}

func (pdct *PlacementDecisionClustersTracker) getOrderedGroupKeysBesides(orderedGroupKeyToExclude []GroupKey) []GroupKey {
	orderedGroupKey := []GroupKey{}
	for i := 0; i < len(pdct.clustesGroupsIndexToName); i++ {
		groupKey := GroupKey{GroupName: pdct.clustesGroupsIndexToName[int32(i)], GroupIndex: int32(i)}
		if !containsGroupKey(orderedGroupKeyToExclude, groupKey) {
			orderedGroupKey = append(orderedGroupKey, groupKey)
		}
	}

	return orderedGroupKey
}

// Helper function to check if a groupKey is present in the groupKeys slice.
func containsGroupKey(groupKeys []GroupKey, groupKey GroupKey) bool {
	for _, gk := range groupKeys {
		if gk == groupKey {
			return true
		}
	}
	return false
}

func parseGroupKeyFromDecision(d *PlacementDecision) (GroupKey, error) {
	groupName := d.Labels[DecisionGroupNameLabel]
	groupIndex := d.Labels[DecisionGroupIndexLabel]
	groupIndexNum, err := strconv.Atoi(groupIndex)
	if err != nil {
		return GroupKey{}, fmt.Errorf("incorrect group index: %w", err)
	}
	return GroupKey{GroupName: groupName, GroupIndex: int32(groupIndexNum)}, nil
}
