package v1beta1

import (
	"reflect"
	"strconv"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
)

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

func TestPlacementDecisionClustersTracker_Get(t *testing.T) {
	tests := []struct {
		name                           string
		placement                      Placement
		existingScheduledClusters      sets.Set[string]
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
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2"),
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
			existingScheduledClusters:      sets.New[string](),
			existingScheduledClusterGroups: map[GroupKey]sets.Set[string]{},
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
		tracker := NewPlacementDecisionClustersTracker(&test.placement, &fakeGetter, test.existingScheduledClusters, test.existingScheduledClusterGroups)

		// check changed decision clusters
		addedClusters, deletedClusters, err := tracker.Get()
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
		tracker := NewPlacementDecisionClustersTracker(&test.placement, &fakeGetter, nil, nil)
		_, _, err := tracker.Get()
		if err != nil {
			t.Errorf("Case: %v, Failed to run Get(): %v", test.name, err)
		}

		// Call the Existing method with different groupKeys inputs.
		existingClusters, _, _ := tracker.Existing(test.groupKeys)
		existingBesidesClusters, _, _ := tracker.ExistingBesides(test.groupKeys)

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
		tracker := NewPlacementDecisionClustersTracker(&test.placement, &fakeGetter, nil, nil)
		_, _, err := tracker.Get()
		if err != nil {
			t.Errorf("Case: %v, Failed to run Get(): %v", test.name, err)
		}

		// Call the Existing method with different groupKeys inputs.
		existingGroupKeys, existingClusterGroups := tracker.ExistingClusterGroups(test.groupKeys)
		existingBesidesGroupKeys, existingBesidesClusterGroups := tracker.ExistingClusterGroupsBesides(test.groupKeys)

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
