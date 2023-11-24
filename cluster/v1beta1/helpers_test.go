package v1beta1

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"testing"

	cliScheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	v1 "open-cluster-management.io/api/cluster/v1"
)

var (
	scheme = runtime.NewScheme()
)

type clustersGetter struct {
	client client.Client
}
type clusterSetsGetter struct {
	client client.Client
}
type clusterSetBindingsGetter struct {
	client client.Client
}

var existingClusterSetBindings = []*ManagedClusterSetBinding{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dev",
			Namespace: "default",
		},
		Spec: ManagedClusterSetBindingSpec{
			ClusterSet: "dev",
		},
		Status: ManagedClusterSetBindingStatus{
			Conditions: []metav1.Condition{
				{
					Type:   ClusterSetBindingBoundType,
					Status: metav1.ConditionTrue,
				},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "global",
			Namespace: "default",
		},
		Spec: ManagedClusterSetBindingSpec{
			ClusterSet: "global",
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "no-such-cluster-set",
			Namespace: "kube-system",
		},
		Spec: ManagedClusterSetBindingSpec{
			ClusterSet: "no-such-cluster-set",
		},
	},
}
var existingClusterSets = []*ManagedClusterSet{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "dev",
		},
		Spec: ManagedClusterSetSpec{},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "global",
		},
		Spec: ManagedClusterSetSpec{
			ClusterSelector: ManagedClusterSelector{
				SelectorType:  LabelSelector,
				LabelSelector: &metav1.LabelSelector{},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "openshift",
		},
		Spec: ManagedClusterSetSpec{
			ClusterSelector: ManagedClusterSelector{
				SelectorType: LabelSelector,
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"vendor": "openshift",
					},
				},
			},
		},
	},
}
var existingClusters = []*v1.ManagedCluster{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "c1",
			Labels: map[string]string{
				"vendor":        "openshift",
				ClusterSetLabel: "dev",
			},
		},
		Spec: v1.ManagedClusterSpec{},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "c2",
			Labels: map[string]string{
				"cloud":         "aws",
				"vendor":        "openshift",
				ClusterSetLabel: "dev",
			},
		},
		Spec: v1.ManagedClusterSpec{},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "c3",
			Labels: map[string]string{
				"cloud": "aws",
			},
		},
		Spec: v1.ManagedClusterSpec{},
	},
}

func TestMain(m *testing.M) {
	if err := v1.AddToScheme(cliScheme.Scheme); err != nil {
		klog.Errorf("Failed adding cluster to scheme, %v", err)
		os.Exit(1)
	}
	if err := AddToScheme(cliScheme.Scheme); err != nil {
		klog.Errorf("Failed adding set to scheme, %v", err)
		os.Exit(1)
	}

	if err := v1.Install(scheme); err != nil {
		klog.Errorf("Failed adding cluster to scheme, %v", err)
		os.Exit(1)
	}
	if err := AddToScheme(scheme); err != nil {
		klog.Errorf("Failed adding set to scheme, %v", err)
		os.Exit(1)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func (mcl clustersGetter) List(selector labels.Selector) ([]*v1.ManagedCluster, error) {
	clusterList := v1.ManagedClusterList{}
	err := mcl.client.List(context.Background(), &clusterList, &client.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	var retClusters []*v1.ManagedCluster
	for i := range clusterList.Items {
		retClusters = append(retClusters, &clusterList.Items[i])
	}
	return retClusters, nil
}

func (msl clusterSetsGetter) List(selector labels.Selector) ([]*ManagedClusterSet, error) {
	clusterSetList := ManagedClusterSetList{}
	err := msl.client.List(context.Background(), &clusterSetList, &client.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	var retClusterSets []*ManagedClusterSet
	for i := range clusterSetList.Items {
		retClusterSets = append(retClusterSets, &clusterSetList.Items[i])
	}
	return retClusterSets, nil
}

func (mbl clusterSetBindingsGetter) List(namespace string,
	selector labels.Selector) ([]*ManagedClusterSetBinding, error) {
	clusterSetBindingList := ManagedClusterSetBindingList{}
	err := mbl.client.List(context.Background(), &clusterSetBindingList,
		client.InNamespace(namespace), &client.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	var retClusterSetBindings []*ManagedClusterSetBinding
	for i := range clusterSetBindingList.Items {
		retClusterSetBindings = append(retClusterSetBindings, &clusterSetBindingList.Items[i])
	}
	return retClusterSetBindings, nil
}

func TestGetClustersFromClusterSet(t *testing.T) {
	tests := []struct {
		name               string
		clusterset         *ManagedClusterSet
		expectClustersName sets.Set[string]
		expectError        bool
	}{
		{
			name: "test legacy cluster set",
			clusterset: &ManagedClusterSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dev",
				},
				Spec: ManagedClusterSetSpec{},
			},
			expectClustersName: sets.New[string]("c1", "c2"),
		},
		{
			name: "test label selector(openshift) cluster set",
			clusterset: &ManagedClusterSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "openshift",
				},
				Spec: ManagedClusterSetSpec{
					ClusterSelector: ManagedClusterSelector{
						SelectorType: LabelSelector,
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"vendor": "openshift",
							},
						},
					},
				},
			},
			expectClustersName: sets.New[string]("c1", "c2"),
		},
		{
			name: "test global cluster set",
			clusterset: &ManagedClusterSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "global",
				},
				Spec: ManagedClusterSetSpec{
					ClusterSelector: ManagedClusterSelector{
						SelectorType:  LabelSelector,
						LabelSelector: &metav1.LabelSelector{},
					},
				},
			},
			expectClustersName: sets.New[string]("c1", "c2", "c3"),
		},
		{
			name: "test label selector cluster set",
			clusterset: &ManagedClusterSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "invalidset",
				},
				Spec: ManagedClusterSetSpec{
					ClusterSelector: ManagedClusterSelector{
						SelectorType: "invalidType",
					},
				},
			},
			expectError: true,
		},
	}

	var existingObjs []client.Object
	for _, cluster := range existingClusters {
		existingObjs = append(existingObjs, cluster)
	}
	mcl := clustersGetter{
		client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(existingObjs...).Build(),
	}

	for _, test := range tests {
		clusters, err := GetClustersFromClusterSet(test.clusterset, mcl)
		if err != nil {
			if test.expectError {
				continue
			}
			t.Errorf("Case: %v, Failed to run GetClustersFromClusterSet with clusterset: %v", test.name, test.clusterset)
			return
		}
		returnClusters := convertClusterToSet(clusters)
		if !reflect.DeepEqual(returnClusters, test.expectClustersName) {
			t.Errorf("Case: %v, Failed to run GetClustersFromClusterSet. Expect clusters: %v, return cluster: %v", test.name, test.expectClustersName, returnClusters)
			return
		}
	}
}

func TestGetClusterSetsOfCluster(t *testing.T) {
	tests := []struct {
		name                 string
		cluster              v1.ManagedCluster
		expectClusterSetName sets.Set[string]
		expectError          bool
	}{
		{
			name: "test c1 cluster",
			cluster: v1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "c1",
					Labels: map[string]string{
						"vendor":        "openshift",
						ClusterSetLabel: "dev",
					},
				},
				Spec: v1.ManagedClusterSpec{},
			},
			expectClusterSetName: sets.New[string]("dev", "openshift", "global"),
		},
		{
			name: "test c2 cluster",
			cluster: v1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "c2",
					Labels: map[string]string{
						"cloud":         "aws",
						"vendor":        "openshift",
						ClusterSetLabel: "dev",
					},
				},
				Spec: v1.ManagedClusterSpec{},
			},
			expectClusterSetName: sets.New[string]("dev", "openshift", "global"),
		},
		{
			name: "test c3 cluster",
			cluster: v1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "c2",
					Labels: map[string]string{
						"cloud": "aws",
					},
				},
				Spec: v1.ManagedClusterSpec{},
			},
			expectClusterSetName: sets.New[string]("global"),
		},
		{
			name: "test nonexist cluster in client",
			cluster: v1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "doNotExistCluster",
					Labels: map[string]string{
						"cloud":  "aws",
						"vendor": "openshift",
					},
				},
				Spec: v1.ManagedClusterSpec{},
			},
			expectClusterSetName: sets.New[string]("openshift", "global"),
		},
	}

	var existingObjs []client.Object
	for _, cluster := range existingClusters {
		existingObjs = append(existingObjs, cluster)
	}
	for _, clusterset := range existingClusterSets {
		existingObjs = append(existingObjs, clusterset)
	}

	msl := clusterSetsGetter{
		client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(existingObjs...).Build(),
	}

	for _, test := range tests {
		returnSets, err := GetClusterSetsOfCluster(&test.cluster, msl)

		if err != nil {
			if test.expectError {
				continue
			}
			t.Errorf("Case: %v, Failed to run GetClusterSetsOfCluster with cluster: %v", test.name, test.cluster)
			return
		}
		returnClusters := convertClusterSetToSet(returnSets)
		if !reflect.DeepEqual(returnClusters, test.expectClusterSetName) {
			t.Errorf("Case: %v, Failed to run GetClusterSetsOfCluster. Expect clusters: %v, return cluster: %v", test.name, test.expectClusterSetName, returnClusters)
			return
		}
	}
}

type FakePlacementDecisionGetter struct {
	FakeDecisions []*PlacementDecision
}

func (f *FakePlacementDecisionGetter) List(selector labels.Selector, namespace string) (ret []*PlacementDecision, err error) {
	return f.FakeDecisions, nil
}

func (f *FakePlacementDecisionGetter) Update(newPlacementDecisions []*PlacementDecision) (ret []*PlacementDecision, err error) {
	f.FakeDecisions = newPlacementDecisions
	return f.FakeDecisions, nil
}

func newFakePlacementDecision(placementName, groupName string, groupIndex int, clusterNames ...string) *PlacementDecision {
	decisions := make([]ClusterDecision, len(clusterNames))
	for i, clusterName := range clusterNames {
		decisions[i] = ClusterDecision{ClusterName: clusterName}
	}

	return &PlacementDecision{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				PlacementLabel:          placementName,
				DecisionGroupNameLabel:  groupName,
				DecisionGroupIndexLabel: strconv.Itoa(groupIndex),
			},
		},
		Status: PlacementDecisionStatus{
			Decisions: decisions,
		},
	}
}

func TestPlacementDecisionClustersTracker_GetClusterChanges(t *testing.T) {
	tests := []struct {
		name                           string
		placement                      Placement
		existingScheduledClusterGroups map[GroupKey]sets.Set[string]
		updateDecisions                []*PlacementDecision
		expectAddedScheduledClusters   sets.Set[string]
		expectDeletedScheduledClusters sets.Set[string]
	}{
		{
			name: "test placementdecisions",
			placement: Placement{
				ObjectMeta: metav1.ObjectMeta{Name: "placement1", Namespace: "default"},
				Spec:       PlacementSpec{},
			},
			existingScheduledClusterGroups: map[GroupKey]sets.Set[string]{
				{GroupName: "", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
			},
			updateDecisions: []*PlacementDecision{
				newFakePlacementDecision("placement1", "", 0, "cluster1"),
				newFakePlacementDecision("placement1", "", 0, "cluster3"),
			},
			expectAddedScheduledClusters:   sets.New[string]("cluster3"),
			expectDeletedScheduledClusters: sets.New[string]("cluster2"),
		},
		{
			name: "test empty placementdecision",
			placement: Placement{
				ObjectMeta: metav1.ObjectMeta{Name: "placement2", Namespace: "default"},
				Spec:       PlacementSpec{},
			},
			existingScheduledClusterGroups: map[GroupKey]sets.Set[string]{},
			updateDecisions: []*PlacementDecision{
				newFakePlacementDecision("placement2", "", 0, "cluster1", "cluster2"),
			},
			expectAddedScheduledClusters:   sets.New[string]("cluster1", "cluster2"),
			expectDeletedScheduledClusters: sets.New[string](),
		},
		{
			name: "test nil exist cluster groups",
			placement: Placement{
				ObjectMeta: metav1.ObjectMeta{Name: "placement2", Namespace: "default"},
				Spec:       PlacementSpec{},
			},
			existingScheduledClusterGroups: nil,
			updateDecisions: []*PlacementDecision{
				newFakePlacementDecision("placement2", "", 0, "cluster1", "cluster2"),
			},
			expectAddedScheduledClusters:   sets.New[string]("cluster1", "cluster2"),
			expectDeletedScheduledClusters: sets.New[string](),
		},
	}

	for _, test := range tests {
		// init fake placement decision getter
		fakeGetter := FakePlacementDecisionGetter{
			FakeDecisions: test.updateDecisions,
		}
		// init tracker
		tracker := NewPlacementDecisionClustersTrackerWithGroups(&test.placement, &fakeGetter, test.existingScheduledClusterGroups)

		// check changed decision clusters
		addedClusters, deletedClusters, err := tracker.GetClusterChanges()
		if err != nil {
			t.Errorf("Case: %v, Failed to run Get(): %v", test.name, err)
		}
		if !reflect.DeepEqual(addedClusters, test.expectAddedScheduledClusters) {
			t.Errorf("Case: %v, expect added decisions: %v, return decisions: %v", test.name, test.expectAddedScheduledClusters, addedClusters)
			return
		}
		if !reflect.DeepEqual(deletedClusters, test.expectDeletedScheduledClusters) {
			t.Errorf("Case: %v, expect deleted decisions: %v, return decisions: %v", test.name, test.expectDeletedScheduledClusters, deletedClusters)
			return
		}
	}
}

func TestPlacementDecisionClustersTracker_Existing(t *testing.T) {
	tests := []struct {
		name                            string
		placement                       Placement
		placementDecisions              []*PlacementDecision
		groupKeys                       []GroupKey
		expectedExistingClusters        sets.Set[string]
		expectedExistingBesidesClusters sets.Set[string]
	}{
		{
			name: "test full group key",
			placement: Placement{
				ObjectMeta: metav1.ObjectMeta{Name: "placement1", Namespace: "default"},
				Spec:       PlacementSpec{},
			},
			placementDecisions: []*PlacementDecision{
				newFakePlacementDecision("placement1", "group1", 0, "cluster1", "cluster2"),
				newFakePlacementDecision("placement1", "group2", 1, "cluster3", "cluster4"),
			},
			groupKeys: []GroupKey{
				{GroupName: "group1"},
				{GroupIndex: 1},
				{GroupName: "group3"},
			},
			expectedExistingClusters:        sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4"),
			expectedExistingBesidesClusters: sets.New[string](),
		},
		{
			name: "test part of group key",
			placement: Placement{
				ObjectMeta: metav1.ObjectMeta{Name: "placement1", Namespace: "default"},
				Spec:       PlacementSpec{},
			},
			placementDecisions: []*PlacementDecision{
				newFakePlacementDecision("placement1", "group1", 0, "cluster1", "cluster2"),
				newFakePlacementDecision("placement1", "group2", 1, "cluster3", "cluster4"),
			},
			groupKeys: []GroupKey{
				{GroupName: "group1"},
			},
			expectedExistingClusters:        sets.New[string]("cluster1", "cluster2"),
			expectedExistingBesidesClusters: sets.New[string]("cluster3", "cluster4"),
		},
		{
			name: "test empty group key",
			placement: Placement{
				ObjectMeta: metav1.ObjectMeta{Name: "placement1", Namespace: "default"},
				Spec:       PlacementSpec{},
			},
			placementDecisions: []*PlacementDecision{
				newFakePlacementDecision("placement1", "group1", 0, "cluster1", "cluster2"),
				newFakePlacementDecision("placement1", "group2", 1, "cluster3", "cluster4"),
			},
			groupKeys:                       []GroupKey{},
			expectedExistingClusters:        sets.New[string](),
			expectedExistingBesidesClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4"),
		},
	}

	for _, test := range tests {
		// init fake placement decision getter
		fakeGetter := FakePlacementDecisionGetter{
			FakeDecisions: test.placementDecisions,
		}
		// init tracker
		tracker := NewPlacementDecisionClustersTrackerWithGroups(&test.placement, &fakeGetter, nil)
		err := tracker.Refresh()
		if err != nil {
			t.Errorf("Case: %v, Failed to run Refresh(): %v", test.name, err)
		}

		// Call the Existing method with different groupKeys inputs.
		existingClusters := tracker.ExistingClusterGroups(test.groupKeys...).GetClusters()
		existingBesidesClusters := tracker.ExistingClusterGroupsBesides(test.groupKeys...).GetClusters()

		// Assert the existingClusters
		if !test.expectedExistingClusters.Equal(existingClusters) {
			t.Errorf("Expected: %v, Actual: %v", test.expectedExistingClusters.UnsortedList(), existingClusters.UnsortedList())
		}
		if !test.expectedExistingBesidesClusters.Equal(existingBesidesClusters) {
			t.Errorf("Expected: %v, Actual: %v", test.expectedExistingBesidesClusters.UnsortedList(), existingBesidesClusters.UnsortedList())
		}
	}
}

func TestPlacementDecisionClustersTracker_ExistingClusterGroups(t *testing.T) {
	tests := []struct {
		name                                 string
		placement                            Placement
		placementDecisions                   []*PlacementDecision
		groupKeys                            []GroupKey
		expectedGroupKeys                    []GroupKey
		expectedExistingClusterGroups        map[GroupKey]sets.Set[string]
		expectedBesidesGroupKeys             []GroupKey
		expectedExistingBesidesClusterGroups map[GroupKey]sets.Set[string]
	}{
		{
			name: "test full group key",
			placement: Placement{
				ObjectMeta: metav1.ObjectMeta{Name: "placement1", Namespace: "default"},
				Spec:       PlacementSpec{},
			},
			placementDecisions: []*PlacementDecision{
				newFakePlacementDecision("placement1", "group1", 0, "cluster1", "cluster2"),
				newFakePlacementDecision("placement1", "group1", 1, "cluster3", "cluster4"),
				newFakePlacementDecision("placement1", "group2", 2, "cluster5", "cluster6"),
			},
			groupKeys: []GroupKey{
				{GroupName: "group1"},
				{GroupIndex: 2},
				{GroupName: "group3"},
			},
			expectedGroupKeys: []GroupKey{
				{GroupName: "group1", GroupIndex: 0},
				{GroupName: "group1", GroupIndex: 1},
				{GroupName: "group2", GroupIndex: 2},
			},
			expectedExistingClusterGroups: map[GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "group1", GroupIndex: 1}: sets.New[string]("cluster3", "cluster4"),
				{GroupName: "group2", GroupIndex: 2}: sets.New[string]("cluster5", "cluster6"),
			},
			expectedBesidesGroupKeys:             []GroupKey{},
			expectedExistingBesidesClusterGroups: map[GroupKey]sets.Set[string]{},
		},
		{
			name: "test part of group key",
			placement: Placement{
				ObjectMeta: metav1.ObjectMeta{Name: "placement1", Namespace: "default"},
				Spec:       PlacementSpec{},
			},
			placementDecisions: []*PlacementDecision{
				newFakePlacementDecision("placement1", "group1", 0, "cluster1", "cluster2"),
				newFakePlacementDecision("placement1", "group1", 1, "cluster3", "cluster4"),
				newFakePlacementDecision("placement1", "group2", 2, "cluster5", "cluster6"),
			},
			groupKeys: []GroupKey{
				{GroupName: "group1"},
			},
			expectedGroupKeys: []GroupKey{
				{GroupName: "group1", GroupIndex: 0},
				{GroupName: "group1", GroupIndex: 1},
			},
			expectedExistingClusterGroups: map[GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "group1", GroupIndex: 1}: sets.New[string]("cluster3", "cluster4"),
			},
			expectedBesidesGroupKeys: []GroupKey{
				{GroupName: "group2", GroupIndex: 2},
			},
			expectedExistingBesidesClusterGroups: map[GroupKey]sets.Set[string]{
				{GroupName: "group2", GroupIndex: 2}: sets.New[string]("cluster5", "cluster6"),
			},
		},
		{
			name: "test empty group key",
			placement: Placement{
				ObjectMeta: metav1.ObjectMeta{Name: "placement1", Namespace: "default"},
				Spec:       PlacementSpec{},
			},
			placementDecisions: []*PlacementDecision{
				newFakePlacementDecision("placement1", "group1", 0, "cluster1", "cluster2"),
				newFakePlacementDecision("placement1", "group1", 1, "cluster3", "cluster4"),
				newFakePlacementDecision("placement1", "group2", 2, "cluster5", "cluster6"),
			},
			groupKeys:                     []GroupKey{},
			expectedGroupKeys:             []GroupKey{},
			expectedExistingClusterGroups: map[GroupKey]sets.Set[string]{},
			expectedBesidesGroupKeys: []GroupKey{
				{GroupName: "group1", GroupIndex: 0},
				{GroupName: "group1", GroupIndex: 1},
				{GroupName: "group2", GroupIndex: 2},
			},
			expectedExistingBesidesClusterGroups: map[GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "group1", GroupIndex: 1}: sets.New[string]("cluster3", "cluster4"),
				{GroupName: "group2", GroupIndex: 2}: sets.New[string]("cluster5", "cluster6"),
			},
		},
	}

	for _, test := range tests {
		// init fake placement decision getter
		fakeGetter := FakePlacementDecisionGetter{
			FakeDecisions: test.placementDecisions,
		}
		// init tracker
		tracker := NewPlacementDecisionClustersTracker(&test.placement, &fakeGetter, nil)
		err := tracker.Refresh()
		if err != nil {
			t.Errorf("Case: %v, Failed to run Refresh(): %v", test.name, err)
		}

		// Call the Existing method with different groupKeys inputs.
		existingClusterGroups := tracker.ExistingClusterGroups(test.groupKeys...)
		existingBesidesClusterGroups := tracker.ExistingClusterGroupsBesides(test.groupKeys...)
		existingGroupKeys := existingClusterGroups.GetOrderedGroupKeys()
		existingBesidesGroupKeys := existingBesidesClusterGroups.GetOrderedGroupKeys()

		// Assert the existingClustersGroups
		if !reflect.DeepEqual(existingGroupKeys, test.expectedGroupKeys) {
			t.Errorf("Expected: %v, Actual: %v", test.expectedGroupKeys, existingGroupKeys)
		}
		for _, gk := range existingGroupKeys {
			if !test.expectedExistingClusterGroups[gk].Equal(existingClusterGroups[gk]) {
				t.Errorf("Expected: %v, Actual: %v", test.expectedExistingClusterGroups[gk], existingClusterGroups[gk])
			}
		}

		if !reflect.DeepEqual(existingBesidesGroupKeys, test.expectedBesidesGroupKeys) {
			t.Errorf("Expected: %v, Actual: %v", test.expectedBesidesGroupKeys, existingBesidesGroupKeys)
		}
		for _, gk := range existingBesidesGroupKeys {
			if !test.expectedExistingBesidesClusterGroups[gk].Equal(existingBesidesClusterGroups[gk]) {
				t.Errorf("Expected: %v, Actual: %v", test.expectedExistingBesidesClusterGroups[gk], existingClusterGroups[gk])
			}
		}
	}
}

func convertClusterToSet(clusters []*v1.ManagedCluster) sets.Set[string] {
	if len(clusters) == 0 {
		return nil
	}
	retSet := sets.New[string]()
	for _, cluster := range clusters {
		retSet.Insert(cluster.Name)
	}
	return retSet
}

func convertClusterSetToSet(clustersets []*ManagedClusterSet) sets.Set[string] {
	if len(clustersets) == 0 {
		return nil
	}
	retSet := sets.New[string]()
	for _, clusterset := range clustersets {
		retSet.Insert(clusterset.Name)
	}
	return retSet
}

func convertClusterSetBindingsToSet(clusterSetBindings []*ManagedClusterSetBinding) sets.Set[string] {
	if len(clusterSetBindings) == 0 {
		return nil
	}
	retSet := sets.New[string]()
	for _, clusterSetBinding := range clusterSetBindings {
		retSet.Insert(clusterSetBinding.Name)
	}
	return retSet
}

func TestGetValidManagedClusterSetBindings(t *testing.T) {
	tests := []struct {
		name                          string
		namespace                     string
		expectClusterSetBindingsNames sets.Set[string]
		expectError                   bool
	}{
		{
			name:                          "test found valid cluster bindings only",
			namespace:                     "default",
			expectClusterSetBindingsNames: sets.New[string]("dev"),
		},

		{
			name:                          "test no cluster binding found",
			namespace:                     "kube-system",
			expectClusterSetBindingsNames: nil,
		},
	}

	var existingObjs []client.Object
	for _, cluster := range existingClusters {
		existingObjs = append(existingObjs, cluster)
	}
	for _, clusterSet := range existingClusterSets {
		existingObjs = append(existingObjs, clusterSet)
	}
	for _, clusterSetBinding := range existingClusterSetBindings {
		existingObjs = append(existingObjs, clusterSetBinding)
	}

	mbl := clusterSetBindingsGetter{
		client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(existingObjs...).Build(),
	}

	for _, test := range tests {
		returnSets, err := GetBoundManagedClusterSetBindings(test.namespace, mbl)

		if err != nil {
			if test.expectError {
				continue
			}
			t.Errorf("Case: %v, Failed to run GetValidManagedClusterSetBindings with namespace: %v", test.name, test.namespace)
			return
		}
		returnBindings := convertClusterSetBindingsToSet(returnSets)
		if !reflect.DeepEqual(returnBindings, test.expectClusterSetBindingsNames) {
			t.Errorf("Case: %v, Failed to run GetValidManagedClusterSetBindings. Expect bindings: %v, return bindings: %v", test.name, test.expectClusterSetBindingsNames, returnBindings)
			return
		}
	}
}
