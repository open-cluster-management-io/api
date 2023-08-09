package v1alpha1

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	testingclock "k8s.io/utils/clock/testing"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
)

var fakeTime = metav1.NewTime(time.Date(2022, time.January, 01, 0, 0, 0, 0, time.UTC))
var fakeTimeMax = metav1.NewTime(fakeTime.Add(maxTimeDuration))
var fakeTimeMax_60s = metav1.NewTime(fakeTime.Add(maxTimeDuration - time.Minute))
var fakeTimeMax_120s = metav1.NewTime(fakeTime.Add(maxTimeDuration - 2*time.Minute))
var fakeTime30s = metav1.NewTime(fakeTime.Add(30 * time.Second))
var fakeTime_30s = metav1.NewTime(fakeTime.Add(-30 * time.Second))
var fakeTime_60s = metav1.NewTime(fakeTime.Add(-time.Minute))
var fakeTime_120s = metav1.NewTime(fakeTime.Add(-2 * time.Minute))

type FakePlacementDecisionGetter struct {
	FakeDecisions []*clusterv1beta1.PlacementDecision
}

func (f *FakePlacementDecisionGetter) List(selector labels.Selector, namespace string) (ret []*clusterv1beta1.PlacementDecision, err error) {
	return f.FakeDecisions, nil
}

func TestGetRolloutCluster_All(t *testing.T) {
	tests := []struct {
		name                           string
		rolloutStrategy                RolloutStrategy
		existingScheduledClusters      sets.Set[string]
		existingScheduledClusterGroups map[clusterv1beta1.GroupKey]sets.Set[string]
		clusterRolloutStatusFunc       ClusterRolloutStatusFunc
		expectRolloutStrategy          *RolloutStrategy
		expectRolloutClusters          map[string]ClusterRolloutStatus
		expectTimeOutClusters          map[string]ClusterRolloutStatus
	}{
		{
			name:                      "test rollout all with timeout 90s",
			rolloutStrategy:           RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"90s"}}},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5", "cluster6"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: ToApply, LastTransitionTime: &fakeTime_60s},
					"cluster2": {Status: Progressing, LastTransitionTime: &fakeTime_60s},
					"cluster3": {Status: Succeed, LastTransitionTime: &fakeTime_60s},
					"cluster4": {Status: Failed, LastTransitionTime: &fakeTime_60s},
					"cluster5": {Status: Failed, LastTransitionTime: &fakeTime_120s},
					"cluster6": {},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"90s"}}},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply, LastTransitionTime: &fakeTime_60s},
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				"cluster4": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Failed, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				"cluster6": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{
				"cluster5": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
		},
		{
			name:                      "test rollout all (default timeout None)",
			rolloutStrategy:           RolloutStrategy{Type: All},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: ToApply},
					"cluster2": {Status: Progressing},
					"cluster3": {Status: Succeed},
					"cluster4": {Status: Failed},
					"cluster5": {},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{""}}},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, TimeOutTime: &fakeTimeMax},
				"cluster4": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Failed, TimeOutTime: &fakeTimeMax},
				"cluster5": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{},
		},
		{
			name:                      "test rollout all with timeout 0s",
			rolloutStrategy:           RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"0s"}}},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: ToApply},
					"cluster2": {Status: Progressing},
					"cluster3": {Status: Succeed},
					"cluster4": {Status: Failed},
					"cluster5": {},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"0s"}}},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				"cluster5": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, TimeOutTime: &fakeTime},
				"cluster4": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: TimeOut, TimeOutTime: &fakeTime},
			},
		},
	}

	RolloutClock = testingclock.NewFakeClock(fakeTime.Time)
	for _, test := range tests {
		// init fake placement decision tracker
		fakeGetter := FakePlacementDecisionGetter{}
		tracker := clusterv1beta1.NewPlacementDecisionClustersTracker(nil, &fakeGetter, test.existingScheduledClusters, test.existingScheduledClusterGroups)

		rolloutHandler, _ := NewRolloutHandler(tracker)
		actualRolloutStrategy, actualRolloutResult, _ := rolloutHandler.GetRolloutCluster(test.rolloutStrategy, test.clusterRolloutStatusFunc)

		if !reflect.DeepEqual(actualRolloutStrategy.All, test.expectRolloutStrategy.All) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect strategy : %v, actual : %v", test.name, test.expectRolloutStrategy, actualRolloutStrategy)
			return
		}
		if !reflect.DeepEqual(actualRolloutResult.ClustersToRollout, test.expectRolloutClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect rollout clusters: %v, actual : %v", test.name, test.expectRolloutClusters, actualRolloutResult.ClustersToRollout)
			return
		}
		if !reflect.DeepEqual(actualRolloutResult.ClustersTimeOut, test.expectTimeOutClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect timeout clusters: %v, actual : %v", test.name, test.expectTimeOutClusters, actualRolloutResult.ClustersTimeOut)
			return
		}
	}
}

func TestGetRolloutCluster_Progressive(t *testing.T) {
	tests := []struct {
		name                           string
		rolloutStrategy                RolloutStrategy
		existingScheduledClusters      sets.Set[string]
		existingScheduledClusterGroups map[clusterv1beta1.GroupKey]sets.Set[string]
		clusterRolloutStatusFunc       ClusterRolloutStatusFunc
		expectRolloutStrategy          *RolloutStrategy
		expectRolloutClusters          map[string]ClusterRolloutStatus
		expectTimeOutClusters          map[string]ClusterRolloutStatus
	}{
		{
			name: "test progressive rollout with timeout 90s",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5", "cluster6"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: ToApply, LastTransitionTime: &fakeTime_60s},
					"cluster2": {Status: Progressing, LastTransitionTime: &fakeTime_60s},
					"cluster3": {Status: Succeed, LastTransitionTime: &fakeTime_60s},
					"cluster4": {Status: Failed, LastTransitionTime: &fakeTime_60s},
					"cluster5": {Status: Failed, LastTransitionTime: &fakeTime_120s},
					"cluster6": {},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout: Timeout{"90s"},
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply, LastTransitionTime: &fakeTime_60s},
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				"cluster4": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Failed, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				"cluster6": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{
				"cluster5": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
		},
		{
			name: "test progressive rollout with timeout None and MaxConcurrency 50%",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{""},
					MaxConcurrency: intstr.FromString("50%"), // 50% of total clusters
				},
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: ToApply},
					"cluster2": {Status: Progressing},
					"cluster3": {Status: Succeed},
					"cluster4": {Status: Failed},
					"cluster5": {},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{""},
					MaxConcurrency: intstr.FromString("50%"), // 50% of total clusters
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, TimeOutTime: &fakeTimeMax},
				"cluster4": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Failed, TimeOutTime: &fakeTimeMax},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{},
		},
		{
			name: "test progressive rollout with timeout 0s and MaxConcurrency 3",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"0s"},
					MaxConcurrency: intstr.FromInt(3), // Maximum 3 clusters concurrently
				},
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: ToApply},
					"cluster2": {Status: Progressing},
					"cluster3": {Status: Succeed},
					"cluster4": {Status: Failed},
					"cluster5": {},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"0s"},
					MaxConcurrency: intstr.FromInt(3), // Maximum 3 clusters concurrently
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				"cluster5": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, TimeOutTime: &fakeTime},
				"cluster4": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: TimeOut, TimeOutTime: &fakeTime},
			},
		},
		{
			name: "test progressive rollout with mandatory decision groups",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					MaxConcurrency: intstr.FromString("50%"),
				},
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: ToApply},
					"cluster2": {Status: ToApply},
					"cluster3": {Status: ToApply},
					"cluster4": {Status: ToApply},
					"cluster5": {Status: ToApply},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					MaxConcurrency: intstr.FromString("50%"),
					Timeout:        Timeout{""},
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{},
		},
		{
			name: "test progressive rollout with mandatory decision groups Succeed",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					MaxConcurrency: intstr.FromInt(2),
				},
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5", "cluster6"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: Succeed},
					"cluster2": {Status: Succeed},
					"cluster3": {Status: ToApply},
					"cluster4": {Status: Succeed},
					"cluster5": {Status: ToApply},
					"cluster6": {Status: ToApply},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					MaxConcurrency: intstr.FromInt(2),
					Timeout:        Timeout{""},
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster3": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
				"cluster5": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{},
		},
		{
			name: "test progressive rollout with mandatory decision groups failed",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					MaxConcurrency: intstr.FromString("50%"), // 50% of total clusters
					Timeout:        Timeout{"0s"},
				},
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: Failed},
					"cluster2": {Status: Failed},
					"cluster3": {Status: ToApply},
					"cluster4": {Status: ToApply},
					"cluster5": {Status: ToApply},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					MaxConcurrency: intstr.FromString("50%"), // 50% of total clusters
					Timeout:        Timeout{"0s"},
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, TimeOutTime: &fakeTimeMax},
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, TimeOutTime: &fakeTimeMax},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{},
		},
	}

	// Set the fake time for testing
	RolloutClock = testingclock.NewFakeClock(fakeTime.Time)

	for _, test := range tests {
		// Init fake placement decision tracker
		fakeGetter := FakePlacementDecisionGetter{}
		tracker := clusterv1beta1.NewPlacementDecisionClustersTracker(nil, &fakeGetter, test.existingScheduledClusters, test.existingScheduledClusterGroups)

		rolloutHandler, _ := NewRolloutHandler(tracker)
		actualRolloutStrategy, actualRolloutResult, _ := rolloutHandler.GetRolloutCluster(test.rolloutStrategy, test.clusterRolloutStatusFunc)

		if !reflect.DeepEqual(actualRolloutStrategy.Progressive, test.expectRolloutStrategy.Progressive) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect strategy : %v, actual : %v", test.name, test.expectRolloutStrategy, actualRolloutStrategy)
			return
		}
		if !reflect.DeepEqual(actualRolloutResult.ClustersToRollout, test.expectRolloutClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect rollout clusters: %v, actual : %v", test.name, test.expectRolloutClusters, actualRolloutResult.ClustersToRollout)
			return
		}
		if !reflect.DeepEqual(actualRolloutResult.ClustersTimeOut, test.expectTimeOutClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect timeout clusters: %v, actual : %v", test.name, test.expectTimeOutClusters, actualRolloutResult.ClustersTimeOut)
			return
		}
	}
}

func TestGetRolloutCluster_ProgressivePerGroup(t *testing.T) {
	tests := []struct {
		name                           string
		rolloutStrategy                RolloutStrategy
		existingScheduledClusters      sets.Set[string]
		existingScheduledClusterGroups map[clusterv1beta1.GroupKey]sets.Set[string]
		clusterRolloutStatusFunc       ClusterRolloutStatusFunc
		expectRolloutStrategy          *RolloutStrategy
		expectRolloutClusters          map[string]ClusterRolloutStatus
		expectTimeOutClusters          map[string]ClusterRolloutStatus
	}{
		{
			name: "test progressive per group rollout with timeout 90s",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5", "cluster6"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: Failed, LastTransitionTime: &fakeTime_60s},
					"cluster2": {Status: Failed, LastTransitionTime: &fakeTime_120s},
					"cluster3": {Status: ToApply, LastTransitionTime: &fakeTime_60s},
					"cluster4": {Status: ToApply, LastTransitionTime: &fakeTime_60s},
					"cluster5": {Status: ToApply, LastTransitionTime: &fakeTime_60s},
					"cluster6": {},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
		},
		{
			name: "test progressive per group rollout with timeout None",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5", "cluster6"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: Failed, LastTransitionTime: &fakeTime_60s},
					"cluster2": {Status: Failed, LastTransitionTime: &fakeTime_120s},
					"cluster3": {Status: ToApply, LastTransitionTime: &fakeTime_60s},
					"cluster4": {Status: ToApply, LastTransitionTime: &fakeTime_60s},
					"cluster5": {Status: ToApply, LastTransitionTime: &fakeTime_60s},
					"cluster6": {},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{""},
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTimeMax_60s},
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTimeMax_120s},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{},
		},
		{
			name: "test progressive per group rollout with timeout 0s",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"0s"},
				},
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5", "cluster6"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: Failed, LastTransitionTime: &fakeTime_60s},
					"cluster2": {Status: Failed, LastTransitionTime: &fakeTime_120s},
					"cluster3": {Status: ToApply, LastTransitionTime: &fakeTime_60s},
					"cluster4": {Status: ToApply, LastTransitionTime: &fakeTime_60s},
					"cluster5": {Status: ToApply, LastTransitionTime: &fakeTime_60s},
					"cluster6": {},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"0s"},
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster3": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply, LastTransitionTime: &fakeTime_60s},
				"cluster4": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply, LastTransitionTime: &fakeTime_60s},
				"cluster5": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply, LastTransitionTime: &fakeTime_60s},
				"cluster6": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime_60s},
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_120s},
			},
		},
		{
			name: "test progressive per group rollout with mandatory decision groups",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
				},
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: ToApply},
					"cluster2": {Status: ToApply},
					"cluster3": {Status: ToApply},
					"cluster4": {Status: ToApply},
					"cluster5": {Status: ToApply},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{},
		},
		{
			name: "test progressive per group rollout with mandatory decision groups Succeed",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
				},
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4"),
				{GroupName: "", GroupIndex: 2}:       sets.New[string]("cluster5"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: Succeed},
					"cluster2": {Status: Succeed},
					"cluster3": {Status: ToApply},
					"cluster4": {Status: Succeed},
					"cluster5": {Status: ToApply},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster3": {GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{},
		},
		{
			name: "test progressive per group rollout with mandatory decision groups failed",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout: Timeout{"0s"},
				},
			},
			existingScheduledClusters: sets.New[string]("cluster1", "cluster2", "cluster3", "cluster4", "cluster5"),
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5"),
			},
			clusterRolloutStatusFunc: func(clusterName string) ClusterRolloutStatus {
				clustersRolloutStatus := map[string]ClusterRolloutStatus{
					"cluster1": {Status: Failed},
					"cluster2": {Status: Failed},
					"cluster3": {Status: ToApply},
					"cluster4": {Status: ToApply},
					"cluster5": {Status: ToApply},
				}
				return clustersRolloutStatus[clusterName]
			},
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout: Timeout{"0s"},
				},
			},
			expectRolloutClusters: map[string]ClusterRolloutStatus{
				"cluster1": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, TimeOutTime: &fakeTimeMax},
				"cluster2": {GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, TimeOutTime: &fakeTimeMax},
			},
			expectTimeOutClusters: map[string]ClusterRolloutStatus{},
		},
	}

	// Set the fake time for testing
	RolloutClock = testingclock.NewFakeClock(fakeTime.Time)

	for _, test := range tests {
		// Init fake placement decision tracker
		fakeGetter := FakePlacementDecisionGetter{}
		tracker := clusterv1beta1.NewPlacementDecisionClustersTracker(nil, &fakeGetter, test.existingScheduledClusters, test.existingScheduledClusterGroups)

		rolloutHandler, _ := NewRolloutHandler(tracker)
		actualRolloutStrategy, actualRolloutResult, _ := rolloutHandler.GetRolloutCluster(test.rolloutStrategy, test.clusterRolloutStatusFunc)

		if !reflect.DeepEqual(actualRolloutStrategy.ProgressivePerGroup, test.expectRolloutStrategy.ProgressivePerGroup) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect strategy : %v, actual : %v", test.name, test.expectRolloutStrategy, actualRolloutStrategy)
			return
		}
		if !reflect.DeepEqual(actualRolloutResult.ClustersToRollout, test.expectRolloutClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect rollout clusters: %v, actual : %v", test.name, test.expectRolloutClusters, actualRolloutResult.ClustersToRollout)
			return
		}
		if !reflect.DeepEqual(actualRolloutResult.ClustersTimeOut, test.expectTimeOutClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect timeout clusters: %v, actual : %v", test.name, test.expectTimeOutClusters, actualRolloutResult.ClustersTimeOut)
			return
		}
	}
}

func TestNeedToUpdate(t *testing.T) {
	testCases := []struct {
		name           string
		status         RolloutStatus
		lastTransition *metav1.Time
		timeout        time.Duration
		expectedResult bool
	}{
		{
			name:           "ToApply status",
			status:         ToApply,
			lastTransition: nil,
			timeout:        time.Minute,
			expectedResult: true,
		},
		{
			name:           "Progressing status",
			status:         Progressing,
			lastTransition: nil,
			timeout:        time.Minute,
			expectedResult: true,
		},
		{
			name:           "Succeed status",
			status:         Succeed,
			lastTransition: nil,
			timeout:        time.Minute,
			expectedResult: false,
		},
		{
			name:           "Failed status, timeout is None",
			status:         Failed,
			lastTransition: &fakeTime,
			timeout:        maxTimeDuration,
			expectedResult: true,
		},
		{
			name:           "Failed status, timeout is 0",
			status:         Failed,
			lastTransition: &fakeTime,
			timeout:        0,
			expectedResult: false,
		},
		{
			name:           "Failed status, within the timeout duration",
			status:         Failed,
			lastTransition: &fakeTime_60s,
			timeout:        2 * time.Minute,
			expectedResult: true,
		},
		{
			name:           "Failed status, outside the timeout duration",
			status:         Failed,
			lastTransition: &fakeTime_120s,
			timeout:        time.Minute,
			expectedResult: false,
		},
	}

	RolloutClock = testingclock.NewFakeClock(fakeTime.Time)
	// Run the tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a ClusterRolloutStatus instance
			status := ClusterRolloutStatus{
				Status:             tc.status,
				LastTransitionTime: tc.lastTransition,
			}

			// Call the determineRolloutStatusAndContinue function
			_, result := determineRolloutStatusAndContinue(status, tc.timeout)

			// Compare the result with the expected result
			if result != tc.expectedResult {
				t.Errorf("Expected result: %v, got: %v", tc.expectedResult, result)
			}
		})
	}
}

func TestCalculateLength(t *testing.T) {
	total := 100

	tests := []struct {
		name           string
		maxConcurrency intstr.IntOrString
		expected       int
		expectedError  error
	}{
		{"maxConcurrency type is int", intstr.FromInt(50), 50, nil},
		{"maxConcurrency type is string with percentage", intstr.FromString("30%"), int(math.Ceil(0.3 * float64(total))), nil},
		{"maxConcurrency type is string without percentage", intstr.FromString("invalid"), total, fmt.Errorf("%v invalid type: string is not a percentage", intstr.FromString("invalid"))},
		{"maxConcurrency type is int 0", intstr.FromInt(0), total, nil},
		{"maxConcurrency type is int but out of range", intstr.FromInt(total + 1), total, nil},
		{"maxConcurrency type is string with percentage but out of range", intstr.FromString("200%"), total, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			length, err := calculateLength(test.maxConcurrency, total)

			// Compare the result with the expected result
			if length != test.expected {
				t.Errorf("Expected result: %v, got: %v", test.expected, length)
			}

			// Compare the error with the expected error
			if err != nil && test.expectedError != nil {
				if err.Error() != test.expectedError.Error() {
					t.Errorf("Expected error: %v, got: %v", test.expectedError, err)
				}
			} else if err != nil || test.expectedError != nil {
				t.Errorf("Expected error: %v, got: %v", test.expectedError, err)
			}
		})
	}
}

func TestParseTimeout(t *testing.T) {
	maxTimeDuration := time.Duration(int64(^uint64(0) >> 1))

	tests := []struct {
		name          string
		timeoutStr    string
		expected      time.Duration
		expectedError error
	}{
		{"Valid timeout with hours", "2h", 2 * time.Hour, nil},
		{"Valid timeout with minutes", "30m", 30 * time.Minute, nil},
		{"Valid timeout with seconds", "10s", 10 * time.Second, nil},
		{"\"None\" timeout", "None", maxTimeDuration, nil},
		{"Empty string timeout", "", maxTimeDuration, nil},
		{"Invalid format", "2d", maxTimeDuration, fmt.Errorf("invalid timeout format")},
		{"Invalid format", "0", maxTimeDuration, fmt.Errorf("invalid timeout format")},
		{"Invalid numeric value", "2d0s", maxTimeDuration, fmt.Errorf("invalid timeout format")},
		{"Invalid unit", "10g", maxTimeDuration, fmt.Errorf("invalid timeout format")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Call the parseTimeout function
			duration, err := parseTimeout(test.timeoutStr)

			// Compare the result with the expected result
			if duration != test.expected {
				t.Errorf("Expected result: %v, got: %v", test.expected, duration)
			}

			// Compare the error with the expected error
			if err != nil && test.expectedError != nil {
				if err.Error() != test.expectedError.Error() {
					t.Errorf("Expected error: %v, got: %v", test.expectedError, err)
				}
			} else if err != nil || test.expectedError != nil {
				t.Errorf("Expected error: %v, got: %v", test.expectedError, err)
			}
		})
	}
}

func TestDecisionGroupsToGroupKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    []MandatoryDecisionGroup
		expected []clusterv1beta1.GroupKey
	}{
		{
			name: "Both GroupName and GroupIndex are set",
			input: []MandatoryDecisionGroup{
				{GroupName: "group1", GroupIndex: 1},
				{GroupName: "group2", GroupIndex: 2},
			},
			expected: []clusterv1beta1.GroupKey{
				{GroupName: "group1", GroupIndex: 0},
				{GroupName: "group2", GroupIndex: 0},
			},
		},
		{
			name: "Only GroupName is set",
			input: []MandatoryDecisionGroup{
				{GroupName: "group1"},
				{GroupName: "group2"},
			},
			expected: []clusterv1beta1.GroupKey{
				{GroupName: "group1", GroupIndex: 0},
				{GroupName: "group2", GroupIndex: 0},
			},
		},
		{
			name: "Only GroupIndex is set",
			input: []MandatoryDecisionGroup{
				{GroupIndex: 1},
				{GroupIndex: 2},
			},
			expected: []clusterv1beta1.GroupKey{
				{GroupName: "", GroupIndex: 1},
				{GroupName: "", GroupIndex: 2},
			},
		},
		{
			name: "Both GroupName and GroupIndex are empty",
			input: []MandatoryDecisionGroup{
				{},
			},
			expected: []clusterv1beta1.GroupKey{
				{GroupName: "", GroupIndex: 0},
			},
		},
		{
			name:     "Empty MandatoryDecisionGroup",
			input:    []MandatoryDecisionGroup{},
			expected: []clusterv1beta1.GroupKey{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := decisionGroupsToGroupKeys(test.input)

			// Compare the result with the expected result
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected result: %v, got: %v", test.expected, result)
			}
		})
	}
}
