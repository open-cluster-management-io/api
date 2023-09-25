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
var fakeTimeMax_60s = metav1.NewTime(fakeTime.Add(maxTimeDuration - time.Minute))
var fakeTimeMax_120s = metav1.NewTime(fakeTime.Add(maxTimeDuration - 2*time.Minute))
var fakeTime30s = metav1.NewTime(fakeTime.Add(30 * time.Second))
var fakeTime_30s = metav1.NewTime(fakeTime.Add(-30 * time.Second))
var fakeTime_60s = metav1.NewTime(fakeTime.Add(-time.Minute))
var fakeTime_120s = metav1.NewTime(fakeTime.Add(-2 * time.Minute))

type FakePlacementDecisionGetter struct {
	FakeDecisions []*clusterv1beta1.PlacementDecision
}

// Dummy workload type that will be used to create a RolloutHandler.
type dummyWorkload struct {
	ClusterGroup       clusterv1beta1.GroupKey
	ClusterName        string
	State              string
	LastTransitionTime *metav1.Time
}

// Dummy Workload status
const (
	valid    = "valid"
	applying = "applying"
	done     = "done"
	missing  = "missing"
)

// Dummy ClusterRolloutStatusFunc implementation that will be used to create a RolloutHandler.
func dummyWorkloadClusterRolloutStatusFunc(clusterName string, workload dummyWorkload) (ClusterRolloutStatus, error) {
	// workload obj should be used to determine the clusterRolloutStatus.
	switch workload.State {
	case valid:
		return ClusterRolloutStatus{GroupKey: workload.ClusterGroup, ClusterName: clusterName, Status: ToApply, LastTransitionTime: workload.LastTransitionTime}, nil
	case applying:
		return ClusterRolloutStatus{GroupKey: workload.ClusterGroup, ClusterName: clusterName, Status: Progressing, LastTransitionTime: workload.LastTransitionTime}, nil
	case done:
		return ClusterRolloutStatus{GroupKey: workload.ClusterGroup, ClusterName: clusterName, Status: Succeeded, LastTransitionTime: workload.LastTransitionTime}, nil
	case missing:
		return ClusterRolloutStatus{GroupKey: workload.ClusterGroup, ClusterName: clusterName, Status: Failed, LastTransitionTime: workload.LastTransitionTime}, nil
	default:
		return ClusterRolloutStatus{GroupKey: workload.ClusterGroup, ClusterName: clusterName, Status: ToApply, LastTransitionTime: workload.LastTransitionTime}, nil
	}
}

type testCase struct {
	name                           string
	rolloutStrategy                RolloutStrategy
	existingScheduledClusterGroups map[clusterv1beta1.GroupKey]sets.Set[string]
	clusterRolloutStatusFunc       ClusterRolloutStatusFunc[dummyWorkload] // Using type dummy workload obj
	expectRolloutStrategy          *RolloutStrategy
	existingWorkloads              []dummyWorkload
	expectRolloutClusters          []ClusterRolloutStatus
	expectTimeOutClusters          []ClusterRolloutStatus
	expectRemovedClusters          []ClusterRolloutStatus
}

func (f *FakePlacementDecisionGetter) List(selector labels.Selector, namespace string) (ret []*clusterv1beta1.PlacementDecision, err error) {
	return f.FakeDecisions, nil
}

func TestGetRolloutCluster_All(t *testing.T) {
	tests := []testCase{
		{
			name:            "test rollout all with timeout 90s witout workload created",
			rolloutStrategy: RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"90s"}}},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			existingWorkloads:        []dummyWorkload{},
			expectRolloutStrategy:    &RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"90s"}}},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster1", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster6", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
			},
		},
		{
			name:            "test rollout all with timeout 90s",
			rolloutStrategy: RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"90s"}}},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1},
					ClusterName:        "cluster3",
					State:              applying,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1},
					ClusterName:        "cluster4",
					State:              missing,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1},
					ClusterName:        "cluster5",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutStrategy: &RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"90s"}}},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Failed, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster6", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
		},
	}

	RolloutClock = testingclock.NewFakeClock(fakeTime.Time)
	for _, test := range tests {
		// init fake placement decision tracker
		fakeGetter := FakePlacementDecisionGetter{}
		tracker := clusterv1beta1.NewPlacementDecisionClustersTrackerWithGroups(nil, &fakeGetter, test.existingScheduledClusterGroups)

		rolloutHandler, _ := NewRolloutHandler(tracker, test.clusterRolloutStatusFunc)
		existingRolloutClusters := []ClusterRolloutStatus{}
		for _, workload := range test.existingWorkloads {
			clsRolloutStatus, _ := test.clusterRolloutStatusFunc(workload.ClusterName, workload)
			existingRolloutClusters = append(existingRolloutClusters, clsRolloutStatus)
		}

		actualRolloutStrategy, actualRolloutResult, _ := rolloutHandler.GetRolloutCluster(test.rolloutStrategy, existingRolloutClusters)

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
		if !reflect.DeepEqual(actualRolloutResult.ClustersRemoved, test.expectRemovedClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect removed clusters: %v, actual : %v", test.name, test.expectRemovedClusters, actualRolloutResult.ClustersRemoved)
			return
		}
	}
}

func TestGetRolloutCluster_Progressive(t *testing.T) {
	tests := []testCase{
		{
			name: "test progressive rollout with timeout 90s witout workload created",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"90s"},
					MaxConcurrency: intstr.FromInt(2),
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"90s"},
					MaxConcurrency: intstr.FromInt(2),
				},
			},
			existingWorkloads: []dummyWorkload{},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster1", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
			},
		},
		{
			name: "test progressive rollout with timeout 90s and workload clusterRollOutStatus are in ToApply status",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MaxConcurrency: intstr.FromInt(4),
					Timeout:        Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
				{GroupName: "", GroupIndex: 2}:       sets.New[string]("cluster7", "cluster8", "cluster9"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MaxConcurrency: intstr.FromInt(4),
					Timeout:        Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:  "cluster2",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:  "cluster3",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:  "cluster4",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:  "cluster5",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:  "cluster6",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 2},
					ClusterName:  "cluster7",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 2},
					ClusterName:  "cluster8",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 2},
					ClusterName:  "cluster9",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 2},
					ClusterName:  "cluster10",
					State:        valid,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster1", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
			},
			expectRemovedClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster10", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 2}, Status: ToApply},
			},
		},
		{
			name: "test progressive rollout with timeout 90s and MaxConcurrency not set",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
				{GroupName: "", GroupIndex: 2}:       sets.New[string]("cluster7", "cluster8", "cluster9"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:  "cluster2",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:  "cluster3",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:  "cluster4",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:  "cluster5",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:  "cluster6",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 2},
					ClusterName:  "cluster7",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 2},
					ClusterName:  "cluster8",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 2},
					ClusterName:  "cluster9",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 2},
					ClusterName:  "cluster10",
					State:        valid,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster1", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster6", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster7", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 2}, Status: ToApply},
				{ClusterName: "cluster8", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 2}, Status: ToApply},
				{ClusterName: "cluster9", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 2}, Status: ToApply},
			},
			expectRemovedClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster10", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 2}, Status: ToApply},
			},
		},
		{
			name: "test progressive rollout with timeout 90s",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"90s"},
					MaxConcurrency: intstr.FromInt(2),
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"90s"},
					MaxConcurrency: intstr.FromInt(2),
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
		},
		{
			name: "test progressive rollout with timeout 0s",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"0s"},
					MaxConcurrency: intstr.FromInt(2),
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"0s"},
					MaxConcurrency: intstr.FromInt(2),
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_30s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              applying,
					LastTransitionTime: &fakeTime_30s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              applying,
					LastTransitionTime: &fakeTime_30s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_30s, TimeOutTime: &fakeTime_30s},
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_30s, TimeOutTime: &fakeTime_30s},
			},
		},
		{
			name: "test progressive rollout with mandatroyDecisionGroup and timeout 90s ",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout:        Timeout{"90s"},
					MaxConcurrency: intstr.FromInt(3),
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout:        Timeout{"90s"},
					MaxConcurrency: intstr.FromInt(3),
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
		},
		{
			name: "test progressive rollout with timeout None and MaxConcurrency 50%",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"None"},
					MaxConcurrency: intstr.FromString("50%"), // 50% of total clusters
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"None"},
					MaxConcurrency: intstr.FromString("50%"), // 50% of total clusters
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:        "cluster3",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTimeMax_120s},
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTimeMax_60s},
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
			},
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
					MaxConcurrency: intstr.FromInt(3),
					Timeout:        Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4"),
				{GroupName: "", GroupIndex: 2}:       sets.New[string]("cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					MaxConcurrency: intstr.FromInt(3),
					Timeout:        Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster1", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
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
					MaxConcurrency: intstr.FromInt(3),
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4"),
				{GroupName: "", GroupIndex: 2}:       sets.New[string]("cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					MaxConcurrency: intstr.FromInt(3),
					Timeout:        Timeout{""},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              done,
					LastTransitionTime: &fakeTime_120s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 2}, Status: ToApply},
			},
		},
	}

	// Set the fake time for testing
	RolloutClock = testingclock.NewFakeClock(fakeTime.Time)

	for _, test := range tests {
		// Init fake placement decision tracker
		fakeGetter := FakePlacementDecisionGetter{}
		tracker := clusterv1beta1.NewPlacementDecisionClustersTrackerWithGroups(nil, &fakeGetter, test.existingScheduledClusterGroups)

		rolloutHandler, _ := NewRolloutHandler(tracker, test.clusterRolloutStatusFunc)
		existingRolloutClusters := []ClusterRolloutStatus{}
		for _, workload := range test.existingWorkloads {
			clsRolloutStatus, _ := test.clusterRolloutStatusFunc(workload.ClusterName, workload)
			existingRolloutClusters = append(existingRolloutClusters, clsRolloutStatus)
		}

		actualRolloutStrategy, actualRolloutResult, _ := rolloutHandler.GetRolloutCluster(test.rolloutStrategy, existingRolloutClusters)

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
		if !reflect.DeepEqual(actualRolloutResult.ClustersRemoved, test.expectRemovedClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect removed clusters: %v, actual : %v", test.name, test.expectRemovedClusters, actualRolloutResult.ClustersRemoved)
			return
		}
	}
}

func TestGetRolloutCluster_ProgressivePerGroup(t *testing.T) {
	tests := []testCase{
		{
			name: "test progressivePerGroup rollout with timeout 90s witout workload created",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster1", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
			},
		},
		{
			name: "test progressivePerGroup rollout with timeout 90s and all workload clusterRollOutStatus are in ToApply status",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:  "cluster2",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:  "cluster3",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:  "cluster4",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:  "cluster5",
					State:        valid,
				},
				{
					ClusterGroup: clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:  "cluster6",
					State:        valid,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster1", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
			},
		},
		{
			name: "test progressivePerGroup rollout with timeout 90s",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
		},
		{
			name: "test progressivePerGroup rollout with timeout 90s and first group timeOut",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
				{GroupName: "", GroupIndex: 2}:       sets.New[string]("cluster7", "cluster8", "cluster9"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster6", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
		},
		{
			name: "test progressivePerGroup rollout with timeout 90s and first group timeOut, second group successed",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
				{GroupName: "", GroupIndex: 2}:       sets.New[string]("cluster7", "cluster8", "cluster9"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:        "cluster4",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:        "cluster5",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:        "cluster6",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster7", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 2}, Status: ToApply},
				{ClusterName: "cluster8", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 2}, Status: ToApply},
				{ClusterName: "cluster9", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 2}, Status: ToApply},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
		},
		{
			name: "test progressivePerGroup rollout with timeout None and first group failing",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"None"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
				{GroupName: "", GroupIndex: 2}:       sets.New[string]("cluster7", "cluster8", "cluster9"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"None"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTimeMax_120s},
			},
		},
		{
			name: "test ProgressivePerGroup rollout with mandatroyDecisionGroup failing and timeout 90s ",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
		},
		{
			name: "test ProgressivePerGroup rollout with mandatroyDecisionGroup Succeeded and timeout 90s ",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
				{GroupName: "", GroupIndex: 2}:       sets.New[string]("cluster7", "cluster8", "cluster9"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster6", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
			},
		},
	}

	// Set the fake time for testing
	RolloutClock = testingclock.NewFakeClock(fakeTime.Time)

	for _, test := range tests {
		// Init fake placement decision tracker
		fakeGetter := FakePlacementDecisionGetter{}
		tracker := clusterv1beta1.NewPlacementDecisionClustersTrackerWithGroups(nil, &fakeGetter, test.existingScheduledClusterGroups)

		rolloutHandler, _ := NewRolloutHandler(tracker, test.clusterRolloutStatusFunc)
		existingRolloutClusters := []ClusterRolloutStatus{}
		for _, workload := range test.existingWorkloads {
			clsRolloutStatus, _ := test.clusterRolloutStatusFunc(workload.ClusterName, workload)
			existingRolloutClusters = append(existingRolloutClusters, clsRolloutStatus)
		}

		actualRolloutStrategy, actualRolloutResult, _ := rolloutHandler.GetRolloutCluster(test.rolloutStrategy, existingRolloutClusters)

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
		if !reflect.DeepEqual(actualRolloutResult.ClustersRemoved, test.expectRemovedClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect removed clusters: %v, actual : %v", test.name, test.expectRemovedClusters, actualRolloutResult.ClustersRemoved)
			return
		}
	}
}

func TestGetRolloutCluster_ClusterAdded(t *testing.T) {
	tests := []testCase{
		{
			name:            "test rollout all with timeout 90s and cluster added",
			rolloutStrategy: RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"90s"}}},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster7"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1},
					ClusterName:        "cluster3",
					State:              applying,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1},
					ClusterName:        "cluster4",
					State:              missing,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1},
					ClusterName:        "cluster5",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutStrategy: &RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"90s"}}},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Failed, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster6", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster7", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
		},
		{
			name: "test progressive rollout with mandatory decision groups Succeed and clusters added after rollout",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					MaxConcurrency: intstr.FromInt(3),
					Timeout:        Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster7"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster4", "cluster8"),
				{GroupName: "", GroupIndex: 2}:       sets.New[string]("cluster5", "cluster6", "cluster9"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					MaxConcurrency: intstr.FromInt(3),
					Timeout:        Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:        "cluster3",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:        "cluster4",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupIndex: 2},
					ClusterName:        "cluster5",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				//{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 2}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster7", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
			},
		},
		{
			name: "test progressivePerGroup rollout with timeout 90s and cluster added after rollout start.",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:        "cluster4",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupIndex: 1},
					ClusterName:        "cluster5",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: ToApply},
			},
		},
	}

	// Set the fake time for testing
	RolloutClock = testingclock.NewFakeClock(fakeTime.Time)

	for _, test := range tests {
		// Init fake placement decision tracker
		fakeGetter := FakePlacementDecisionGetter{}
		tracker := clusterv1beta1.NewPlacementDecisionClustersTrackerWithGroups(nil, &fakeGetter, test.existingScheduledClusterGroups)

		rolloutHandler, _ := NewRolloutHandler(tracker, test.clusterRolloutStatusFunc)
		existingRolloutClusters := []ClusterRolloutStatus{}
		for _, workload := range test.existingWorkloads {
			clsRolloutStatus, _ := test.clusterRolloutStatusFunc(workload.ClusterName, workload)
			existingRolloutClusters = append(existingRolloutClusters, clsRolloutStatus)
		}

		actualRolloutStrategy, actualRolloutResult, _ := rolloutHandler.GetRolloutCluster(test.rolloutStrategy, existingRolloutClusters)

		if !reflect.DeepEqual(actualRolloutStrategy.Type, test.expectRolloutStrategy.Type) {
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
		if !reflect.DeepEqual(actualRolloutResult.ClustersRemoved, test.expectRemovedClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect removed clusters: %v, actual : %v", test.name, test.expectRemovedClusters, actualRolloutResult.ClustersRemoved)
			return
		}
	}
}

func TestGetRolloutCluster_ClusterRemoved(t *testing.T) {
	tests := []testCase{
		{
			name:            "test rollout all with timeout 90s and clusters removed",
			rolloutStrategy: RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"90s"}}},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster3", "cluster5"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1},
					ClusterName:        "cluster3",
					State:              applying,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1},
					ClusterName:        "cluster4",
					State:              missing,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1},
					ClusterName:        "cluster5",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutStrategy: &RolloutStrategy{Type: All, All: &RolloutAll{Timeout: Timeout{"90s"}}},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
			expectRemovedClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, LastTransitionTime: &fakeTime_120s},
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: Failed, LastTransitionTime: &fakeTime_60s},
			},
		},
		{
			name: "test progressive rollout with timeout 90s and cluster removed",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"90s"},
					MaxConcurrency: intstr.FromInt(2),
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster2", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					Timeout:        Timeout{"90s"},
					MaxConcurrency: intstr.FromInt(2),
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupName: "", GroupIndex: 1}, Status: ToApply},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
			expectRemovedClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster1", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Succeeded, LastTransitionTime: &fakeTime_60s},
			},
		},
		{
			name: "test progressive rollout with mandatroyDecisionGroup, timeout 90s and cluster removed from mandatroyDecisionGroup",
			rolloutStrategy: RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout:        Timeout{"90s"},
					MaxConcurrency: intstr.FromInt(2),
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: Progressive,
				Progressive: &RolloutProgressive{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout:        Timeout{"90s"},
					MaxConcurrency: intstr.FromInt(2),
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
			},
			expectRemovedClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, LastTransitionTime: &fakeTime_120s},
			},
		},
		{
			name: "test progressivePerGroup rollout with timeout 90s and cluster removed after rollout start.",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              applying,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Progressing, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime30s},
			},
			expectRemovedClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, LastTransitionTime: &fakeTime_60s},
			},
		},
		{
			name: "test progressivePerGroup rollout with timeout 90s and cluster removed after rollout start while the group timeout.",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              applying,
					LastTransitionTime: &fakeTime_120s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster6", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
			},
			expectTimeOutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster3", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: TimeOut, LastTransitionTime: &fakeTime_120s, TimeOutTime: &fakeTime_30s},
			},
			expectRemovedClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, LastTransitionTime: &fakeTime_60s},
			},
		},
		{
			name: "test ProgressivePerGroup rollout with mandatroyDecisionGroup, timeout 90s and cluster removed from mandatroyDecisionGroup",
			rolloutStrategy: RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout: Timeout{"90s"},
				},
			},
			existingScheduledClusterGroups: map[clusterv1beta1.GroupKey]sets.Set[string]{
				{GroupName: "group1", GroupIndex: 0}: sets.New[string]("cluster1", "cluster3"),
				{GroupName: "", GroupIndex: 1}:       sets.New[string]("cluster4", "cluster5", "cluster6"),
				{GroupName: "", GroupIndex: 2}:       sets.New[string]("cluster7", "cluster8", "cluster9"),
			},
			clusterRolloutStatusFunc: dummyWorkloadClusterRolloutStatusFunc,
			expectRolloutStrategy: &RolloutStrategy{
				Type: ProgressivePerGroup,
				ProgressivePerGroup: &RolloutProgressivePerGroup{
					MandatoryDecisionGroups: MandatoryDecisionGroups{
						MandatoryDecisionGroups: []MandatoryDecisionGroup{
							{GroupName: "group1"},
						},
					},
					Timeout: Timeout{"90s"},
				},
			},
			existingWorkloads: []dummyWorkload{
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster1",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster2",
					State:              missing,
					LastTransitionTime: &fakeTime_120s,
				},
				{
					ClusterGroup:       clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0},
					ClusterName:        "cluster3",
					State:              done,
					LastTransitionTime: &fakeTime_60s,
				},
			},
			expectRolloutClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster4", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster5", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
				{ClusterName: "cluster6", GroupKey: clusterv1beta1.GroupKey{GroupIndex: 1}, Status: ToApply},
			},
			expectRemovedClusters: []ClusterRolloutStatus{
				{ClusterName: "cluster2", GroupKey: clusterv1beta1.GroupKey{GroupName: "group1", GroupIndex: 0}, Status: Failed, LastTransitionTime: &fakeTime_120s},
			},
		},
	}

	// Set the fake time for testing
	RolloutClock = testingclock.NewFakeClock(fakeTime.Time)

	for _, test := range tests {
		// Init fake placement decision tracker
		fakeGetter := FakePlacementDecisionGetter{}
		tracker := clusterv1beta1.NewPlacementDecisionClustersTrackerWithGroups(nil, &fakeGetter, test.existingScheduledClusterGroups)

		rolloutHandler, _ := NewRolloutHandler(tracker, test.clusterRolloutStatusFunc)
		existingRolloutClusters := []ClusterRolloutStatus{}
		for _, workload := range test.existingWorkloads {
			clsRolloutStatus, _ := test.clusterRolloutStatusFunc(workload.ClusterName, workload)
			existingRolloutClusters = append(existingRolloutClusters, clsRolloutStatus)
		}

		actualRolloutStrategy, actualRolloutResult, _ := rolloutHandler.GetRolloutCluster(test.rolloutStrategy, existingRolloutClusters)

		if !reflect.DeepEqual(actualRolloutStrategy.Type, test.expectRolloutStrategy.Type) {
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
		if !reflect.DeepEqual(actualRolloutResult.ClustersRemoved, test.expectRemovedClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect removed clusters: %v, actual : %v", test.name, test.expectRemovedClusters, actualRolloutResult.ClustersRemoved)
			return
		}
	}

}

func TestDetermineRolloutStatus(t *testing.T) {
	testCases := []struct {
		name                  string
		timeout               time.Duration
		clusterStatus         ClusterRolloutStatus
		expectRolloutClusters []ClusterRolloutStatus
		expectTimeOutClusters []ClusterRolloutStatus
	}{
		{
			name:                  "ToApply status",
			clusterStatus:         ClusterRolloutStatus{ClusterName: "cluster1", Status: ToApply},
			timeout:               time.Minute,
			expectRolloutClusters: []ClusterRolloutStatus{{ClusterName: "cluster1", Status: ToApply}},
		},
		{
			name:          "Skip status",
			clusterStatus: ClusterRolloutStatus{ClusterName: "cluster1", Status: Skip},
			timeout:       time.Minute,
		},
		{
			name:          "Succeeded status",
			clusterStatus: ClusterRolloutStatus{ClusterName: "cluster1", Status: Succeeded},
			timeout:       time.Minute,
		},
		{
			name:          "TimeOut status",
			clusterStatus: ClusterRolloutStatus{ClusterName: "cluster1", Status: TimeOut},
			timeout:       time.Minute,
		},
		{
			name:                  "Progressing status within the timeout duration",
			clusterStatus:         ClusterRolloutStatus{ClusterName: "cluster1", Status: Progressing, LastTransitionTime: &fakeTime_30s},
			timeout:               time.Minute,
			expectRolloutClusters: []ClusterRolloutStatus{{ClusterName: "cluster1", Status: Progressing, LastTransitionTime: &fakeTime_30s, TimeOutTime: &fakeTime30s}},
		},
		{
			name:                  "Failed status out the timeout duration",
			clusterStatus:         ClusterRolloutStatus{ClusterName: "cluster1", Status: Failed, LastTransitionTime: &fakeTime_60s},
			timeout:               time.Minute,
			expectTimeOutClusters: []ClusterRolloutStatus{{ClusterName: "cluster1", Status: TimeOut, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime}},
		},
		{
			name:                  "unknown status out the timeout duration",
			clusterStatus:         ClusterRolloutStatus{ClusterName: "cluster1", Status: 8, LastTransitionTime: &fakeTime_60s},
			timeout:               time.Minute,
			expectTimeOutClusters: []ClusterRolloutStatus{{ClusterName: "cluster1", Status: TimeOut, LastTransitionTime: &fakeTime_60s, TimeOutTime: &fakeTime}},
		},
		{
			name:                  "unknown status within the timeout duration",
			clusterStatus:         ClusterRolloutStatus{ClusterName: "cluster1", Status: 9, LastTransitionTime: &fakeTime_30s},
			timeout:               time.Minute,
			expectRolloutClusters: []ClusterRolloutStatus{{ClusterName: "cluster1", Status: 9, LastTransitionTime: &fakeTime_30s, TimeOutTime: &fakeTime30s}},
		},
	}

	RolloutClock = testingclock.NewFakeClock(fakeTime.Time)
	for _, tc := range testCases {
		var rolloutClusters, timeoutClusters []ClusterRolloutStatus
		rolloutClusters, timeoutClusters = determineRolloutStatus(tc.clusterStatus, tc.timeout, rolloutClusters, timeoutClusters)
		if !reflect.DeepEqual(rolloutClusters, tc.expectRolloutClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect rollout clusters: %v, actual : %v", tc.name, tc.expectRolloutClusters, rolloutClusters)
			return
		}
		if !reflect.DeepEqual(timeoutClusters, tc.expectTimeOutClusters) {
			t.Errorf("Case: %v, Failed to run NewRolloutHandler. Expect timeout clusters: %v, actual : %v", tc.name, tc.expectTimeOutClusters, timeoutClusters)
			return
		}
	}
}

func TestCalculateRolloutSize(t *testing.T) {
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
			length, err := calculateRolloutSize(test.maxConcurrency, total)

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
