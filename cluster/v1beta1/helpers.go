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
	existingScheduledClusterGroups ClusterGroupsMap
	clusterGroupsIndexToName       map[int32]string
	clusterGroupsNameToIndex       map[string][]int32
	lock                           sync.RWMutex
}

// +k8s:deepcopy-gen=false
type GroupKey struct {
	GroupName  string `json:"groupName,omitempty"`
	GroupIndex int32  `json:"groupIndex,omitempty"`
}

// NewPlacementDecisionClustersTracker initializes a PlacementDecisionClustersTracker
// using existing clusters. Clusters are added to the default cluster group with index 0.
// Set existingScheduledClusters to nil if there are no existing clusters.
func NewPlacementDecisionClustersTracker(placement *Placement, pdl PlacementDecisionGetter, existingScheduledClusters sets.Set[string]) *PlacementDecisionClustersTracker {
	pdct := &PlacementDecisionClustersTracker{
		placement:                      placement,
		placementDecisionGetter:        pdl,
		existingScheduledClusterGroups: ClusterGroupsMap{{GroupIndex: 0}: existingScheduledClusters},
	}

	// Generate group name indices for the tracker.
	pdct.generateGroupsNameIndex()
	return pdct
}

// NewPlacementDecisionClustersTrackerWithGroups initializes a PlacementDecisionClustersTracker
// using existing cluster groups. Set existingScheduledClusterGroups to nil if no groups exist.
func NewPlacementDecisionClustersTrackerWithGroups(placement *Placement, pdl PlacementDecisionGetter, existingScheduledClusterGroups ClusterGroupsMap) *PlacementDecisionClustersTracker {
	pdct := &PlacementDecisionClustersTracker{
		placement:                      placement,
		placementDecisionGetter:        pdl,
		existingScheduledClusterGroups: existingScheduledClusterGroups,
	}

	// Generate group name indices for the tracker.
	pdct.generateGroupsNameIndex()
	return pdct
}

// Refresh refreshes the tracker's decisionClusters.
func (pdct *PlacementDecisionClustersTracker) Refresh() error {
	pdct.lock.Lock()
	defer pdct.lock.Unlock()

	if pdct.placement == nil || pdct.placementDecisionGetter == nil {
		return nil
	}

	// Get the generated PlacementDecisions
	decisionSelector := labels.SelectorFromSet(labels.Set{
		PlacementLabel: pdct.placement.Name,
	})
	decisions, err := pdct.placementDecisionGetter.List(decisionSelector, pdct.placement.Namespace)
	if err != nil {
		return fmt.Errorf("failed to list PlacementDecisions: %w", err)
	}

	// Get the decision cluster names and groups
	newScheduledClusterGroups := map[GroupKey]sets.Set[string]{}
	for _, d := range decisions {
		groupKey, err := parseGroupKeyFromDecision(d)
		if err != nil {
			return err
		}

		if _, exist := newScheduledClusterGroups[groupKey]; !exist {
			newScheduledClusterGroups[groupKey] = sets.New[string]()
		}

		for _, sd := range d.Status.Decisions {
			newScheduledClusterGroups[groupKey].Insert(sd.ClusterName)
		}
	}

	// Update the existing decision cluster groups
	pdct.existingScheduledClusterGroups = newScheduledClusterGroups
	pdct.generateGroupsNameIndex()

	return nil
}

// GetClusterChanges updates the tracker's decisionClusters and returns added and deleted cluster names.
func (pdct *PlacementDecisionClustersTracker) GetClusterChanges() (sets.Set[string], sets.Set[string], error) {
	// Get existing clusters
	existingScheduledClusters := pdct.existingScheduledClusterGroups.GetClusters()

	// Refresh clusters
	err := pdct.Refresh()
	if err != nil {
		return nil, nil, err
	}
	newScheduledClusters := pdct.existingScheduledClusterGroups.GetClusters()

	// Compare the difference
	added := newScheduledClusters.Difference(existingScheduledClusters)
	deleted := existingScheduledClusters.Difference(newScheduledClusters)

	return added, deleted, nil
}

func (pdct *PlacementDecisionClustersTracker) generateGroupsNameIndex() {
	pdct.clusterGroupsIndexToName = map[int32]string{}
	pdct.clusterGroupsNameToIndex = map[string][]int32{}

	for groupkey := range pdct.existingScheduledClusterGroups {
		// index to name
		pdct.clusterGroupsIndexToName[groupkey.GroupIndex] = groupkey.GroupName
		// name to index
		if index, exist := pdct.clusterGroupsNameToIndex[groupkey.GroupName]; exist {
			pdct.clusterGroupsNameToIndex[groupkey.GroupName] = append(index, groupkey.GroupIndex)
		} else {
			pdct.clusterGroupsNameToIndex[groupkey.GroupName] = []int32{groupkey.GroupIndex}
		}
	}

	// sort index order
	for _, index := range pdct.clusterGroupsNameToIndex {
		sort.Slice(index, func(i, j int) bool {
			return index[i] < index[j]
		})
	}
}

// ExistingClusterGroups returns the tracker's existing decision cluster groups for groups listed in groupKeys.
// Return empty set when groupKeys is empty.
func (pdct *PlacementDecisionClustersTracker) ExistingClusterGroups(groupKeys ...GroupKey) ClusterGroupsMap {
	pdct.lock.RLock()
	defer pdct.lock.RUnlock()

	resultClusterGroups := make(map[GroupKey]sets.Set[string])

	includeGroupKeys := pdct.fulfillGroupKeys(groupKeys)
	for _, groupKey := range includeGroupKeys {
		if clusters, found := pdct.existingScheduledClusterGroups[groupKey]; found {
			resultClusterGroups[groupKey] = clusters
		}
	}

	return resultClusterGroups
}

// ExistingClusterGroupsBesides returns the tracker's existing decision cluster groups except cluster groups listed in groupKeys.
// Return all the clusters when groupKeys is empty.
func (pdct *PlacementDecisionClustersTracker) ExistingClusterGroupsBesides(groupKeys ...GroupKey) ClusterGroupsMap {
	pdct.lock.RLock()
	defer pdct.lock.RUnlock()

	resultClusterGroups := make(map[GroupKey]sets.Set[string])

	excludeGroupKeys := pdct.fulfillGroupKeys(groupKeys)
	includeGroupKeys := pdct.getGroupKeysBesides(excludeGroupKeys)
	for _, groupKey := range includeGroupKeys {
		if clusters, found := pdct.existingScheduledClusterGroups[groupKey]; found {
			resultClusterGroups[groupKey] = clusters
		}
	}

	return resultClusterGroups
}

// Fulfill the expect groupkeys with group name or group index, the returned groupkeys are ordered by input group name then group index.
// For example, the input is []GroupKey{{GroupName: "group1"}, {GroupIndex: 2}},
// the returned is []GroupKey{{GroupName: "group1", GroupIndex: 0}, {GroupName: "group1", GroupIndex: 1}, {GroupName: "group2", GroupIndex: 2}}
func (pdct *PlacementDecisionClustersTracker) fulfillGroupKeys(groupKeys []GroupKey) []GroupKey {
	fulfilledGroupKeys := []GroupKey{}
	for _, gk := range groupKeys {
		if gk.GroupName != "" {
			if indexes, exist := pdct.clusterGroupsNameToIndex[gk.GroupName]; exist {
				for _, groupIndex := range indexes {
					fulfilledGroupKeys = append(fulfilledGroupKeys, GroupKey{GroupName: gk.GroupName, GroupIndex: groupIndex})
				}
			}
		} else {
			if groupName, exist := pdct.clusterGroupsIndexToName[gk.GroupIndex]; exist {
				fulfilledGroupKeys = append(fulfilledGroupKeys, GroupKey{GroupName: groupName, GroupIndex: gk.GroupIndex})
			}
		}
	}
	return fulfilledGroupKeys
}

func (pdct *PlacementDecisionClustersTracker) getGroupKeysBesides(groupKeyToExclude []GroupKey) []GroupKey {
	groupKey := []GroupKey{}
	for i := 0; i < len(pdct.clusterGroupsIndexToName); i++ {
		gKey := GroupKey{GroupName: pdct.clusterGroupsIndexToName[int32(i)], GroupIndex: int32(i)}
		if !containsGroupKey(groupKeyToExclude, gKey) {
			groupKey = append(groupKey, gKey)
		}
	}

	return groupKey
}

// ClusterGroupsMap is a custom type representing a map of group keys to sets of cluster names.
type ClusterGroupsMap map[GroupKey]sets.Set[string]

// GetOrderedGroupKeys returns an ordered slice of GroupKeys, sorted by group index.
func (g ClusterGroupsMap) GetOrderedGroupKeys() []GroupKey {
	groupKeys := []GroupKey{}
	for groupKey := range g {
		groupKeys = append(groupKeys, groupKey)
	}

	// sort by group index index
	sort.Slice(groupKeys, func(i, j int) bool {
		return groupKeys[i].GroupIndex < groupKeys[j].GroupIndex
	})

	return groupKeys
}

// GetClusters returns a set containing all clusters from all group sets.
func (g ClusterGroupsMap) GetClusters() sets.Set[string] {
	clusterSet := sets.New[string]()
	for _, clusterGroup := range g {
		clusterSet = clusterSet.Union(clusterGroup)
	}
	return clusterSet
}

// ClusterToGroupKey returns a mapping of cluster names to their respective group keys.
func (g ClusterGroupsMap) ClusterToGroupKey() map[string]GroupKey {
	clusterToGroupKey := map[string]GroupKey{}

	for groupKey, clusterGroup := range g {
		for c := range clusterGroup {
			clusterToGroupKey[c] = groupKey
		}
	}

	return clusterToGroupKey
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
