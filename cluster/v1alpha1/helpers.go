package v1alpha1

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/clock"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
)

var RolloutClock = clock.Clock(clock.RealClock{})
var maxTimeDuration = time.Duration(math.MaxInt64)

type RolloutStatus int

const (
	// resource desired status is not applied yet
	ToApply RolloutStatus = iota
	// resource desired status is applied and last applied status is not updated
	Progressing
	// resource desired status is applied and last applied status is succeed
	Succeed
	// resource desired status is applied and last applied status is failed
	Failed
	// when the rollout status is progressing or failed and the status remains for longer than the timeout
	// then the status will be set to timeout
	TimeOut
	// skip rollout on this cluster
	Skip
)

// Return the rollout status on each managed cluster
type ClusterRolloutStatusFunc func(clusterName string) ClusterRolloutStatus

type ClusterRolloutStatus struct {
	// cluster group key
	// optional field
	GroupKey clusterv1beta1.GroupKey
	// rollout status
	// required field
	Status RolloutStatus
	// the last transition time of rollout status
	// this is used to calculate timeout for progressing and failed status
	// optional field
	LastTransitionTime *metav1.Time
}

type RolloutResult struct {
	// clusters to rollout
	ClustersToRollout map[string]ClusterRolloutStatus
	// clusters that are timeout
	ClustersTimeOut map[string]ClusterRolloutStatus
}

// +k8s:deepcopy-gen=false
type RolloutHandler struct {
	// placement decision tracker
	pdTracker *clusterv1beta1.PlacementDecisionClustersTracker
}

func NewRolloutHandler(pdTracker *clusterv1beta1.PlacementDecisionClustersTracker) (*RolloutHandler, error) {
	if pdTracker == nil {
		return nil, fmt.Errorf("invalid placement decision tracker %v", pdTracker)
	}

	return &RolloutHandler{
		pdTracker: pdTracker,
	}, nil
}

// The input is a duck type RolloutStrategy and a ClusterRolloutStatusFunc to return the rollout status on each managed cluster.
// Return the strategy actual take effect and a list of clusters that need to rollout and that are timeout.
//
// ClustersToRollout: If mandatory decision groups are defined in strategy, will return the clusters to rollout in mandatory decision groups first.
// When all the mandatory decision groups rollout successfully, will return the rest of the clusters that need to rollout.
//
// ClustersTimeOut: If the cluster status is Progressing or Failed, and the status lasts longer than timeout defined in strategy,
// will list them RolloutResult.ClustersTimeOut with status TimeOut.
func (r *RolloutHandler) GetRolloutCluster(rolloutStrategy RolloutStrategy, statusFunc ClusterRolloutStatusFunc) (*RolloutStrategy, RolloutResult, error) {
	switch rolloutStrategy.Type {
	case All:
		return r.getRolloutAllClusters(rolloutStrategy, statusFunc)
	case Progressive:
		return r.getProgressiveClusters(rolloutStrategy, statusFunc)
	case ProgressivePerGroup:
		return r.getProgressivePerGroupClusters(rolloutStrategy, statusFunc)
	default:
		return nil, RolloutResult{}, fmt.Errorf("incorrect rollout strategy type %v", rolloutStrategy.Type)
	}
}

func (r *RolloutHandler) getRolloutAllClusters(rolloutStrategy RolloutStrategy, statusFunc ClusterRolloutStatusFunc) (*RolloutStrategy, RolloutResult, error) {
	strategy := RolloutStrategy{Type: All}
	strategy.All = rolloutStrategy.All.DeepCopy()
	if strategy.All == nil {
		strategy.All = &RolloutAll{}
	}

	// parse timeout
	failureTimeout, err := parseTimeout(strategy.All.Timeout.Timeout)
	if err != nil {
		return &strategy, RolloutResult{}, err
	}

	// get all the clusters
	totalClusters, _, totalClusterGroups := r.pdTracker.ExistingBesides([]clusterv1beta1.GroupKey{})
	clusterToGroupKey := r.pdTracker.ClusterToGroupKey(totalClusterGroups)
	rolloutresult := progressivePerCluster(totalClusters.UnsortedList(), clusterToGroupKey, len(totalClusters), failureTimeout, statusFunc)

	return &strategy, rolloutresult, nil
}

func (r *RolloutHandler) getProgressiveClusters(rolloutStrategy RolloutStrategy, statusFunc ClusterRolloutStatusFunc) (*RolloutStrategy, RolloutResult, error) {
	strategy := RolloutStrategy{Type: Progressive}
	strategy.Progressive = rolloutStrategy.Progressive.DeepCopy()
	if strategy.Progressive == nil {
		strategy.Progressive = &RolloutProgressive{}
	}

	// upgrade mandatory decision groups first
	groupKeys := decisionGroupsToGroupKeys(strategy.Progressive.MandatoryDecisionGroups.MandatoryDecisionGroups)
	clusterGroupKeys, clusterGroups := r.pdTracker.ExistingClusterGroups(groupKeys)

	rolloutresult := progressivePerGroup(clusterGroupKeys, clusterGroups, maxTimeDuration, statusFunc)
	if len(rolloutresult.ClustersToRollout) > 0 {
		return &strategy, rolloutresult, nil
	}

	// parse timeout
	failureTimeout, err := parseTimeout(strategy.Progressive.Timeout.Timeout)
	if err != nil {
		return &strategy, RolloutResult{}, err
	}

	// calculate length
	totalClusters, _, _ := r.pdTracker.ExistingBesides([]clusterv1beta1.GroupKey{})
	length, err := calculateLength(strategy.Progressive.MaxConcurrency, len(totalClusters))
	if err != nil {
		return &strategy, RolloutResult{}, err
	}

	// upgrade the rest clusters
	restClusters, _, restClusterGroups := r.pdTracker.ExistingBesides(clusterGroupKeys)
	restClustersToGroupKey := r.pdTracker.ClusterToGroupKey(restClusterGroups)
	rolloutresult = progressivePerCluster(restClusters.UnsortedList(), restClustersToGroupKey, length, failureTimeout, statusFunc)

	return &strategy, rolloutresult, nil
}

func (r *RolloutHandler) getProgressivePerGroupClusters(rolloutStrategy RolloutStrategy, statusFunc ClusterRolloutStatusFunc) (*RolloutStrategy, RolloutResult, error) {
	strategy := RolloutStrategy{Type: ProgressivePerGroup}
	strategy.ProgressivePerGroup = rolloutStrategy.ProgressivePerGroup.DeepCopy()
	if strategy.ProgressivePerGroup == nil {
		strategy.ProgressivePerGroup = &RolloutProgressivePerGroup{}
	}

	// upgrade mandatory decision groups first
	groupKeys := decisionGroupsToGroupKeys(strategy.ProgressivePerGroup.MandatoryDecisionGroups.MandatoryDecisionGroups)
	clusterGroupKeys, clusterGroups := r.pdTracker.ExistingClusterGroups(groupKeys)

	rolloutresult := progressivePerGroup(clusterGroupKeys, clusterGroups, maxTimeDuration, statusFunc)
	if len(rolloutresult.ClustersToRollout) > 0 {
		return &strategy, rolloutresult, nil
	}

	// parse timeout
	failureTimeout, err := parseTimeout(strategy.ProgressivePerGroup.Timeout.Timeout)
	if err != nil {
		return &strategy, RolloutResult{}, err
	}

	// upgrade the rest decision groups
	restClusterGroupKeys, restClusterGroups := r.pdTracker.ExistingClusterGroupsBesides(clusterGroupKeys)

	rolloutresult = progressivePerGroup(restClusterGroupKeys, restClusterGroups, failureTimeout, statusFunc)
	return &strategy, rolloutresult, nil
}

func progressivePerCluster(clusters []string, clusterToGroupKey map[string]clusterv1beta1.GroupKey, length int, timeout time.Duration, statusFunc ClusterRolloutStatusFunc) RolloutResult {
	rolloutclusters := map[string]ClusterRolloutStatus{}
	timeoutclusters := map[string]ClusterRolloutStatus{}

	if length == 0 {
		return RolloutResult{
			ClustersToRollout: rolloutclusters,
			ClustersTimeOut:   timeoutclusters,
		}
	}

	// sort the clusters in alphabetical order, ensure each time returns the same clusters.
	sort.Strings(clusters)
	for _, cluster := range clusters {
		status := statusFunc(cluster)
		if groupKey, exist := clusterToGroupKey[cluster]; exist {
			status.GroupKey = groupKey
		}

		newstatus, needtorollout := needToRollout(status, timeout)
		status.Status = newstatus
		if needtorollout {
			rolloutclusters[cluster] = status
		}
		if newstatus == TimeOut {
			timeoutclusters[cluster] = status
		}

		if (len(rolloutclusters)%length == 0) && len(rolloutclusters) > 0 {
			return RolloutResult{
				ClustersToRollout: rolloutclusters,
				ClustersTimeOut:   timeoutclusters,
			}
		}
	}

	return RolloutResult{
		ClustersToRollout: rolloutclusters,
		ClustersTimeOut:   timeoutclusters,
	}
}

func progressivePerGroup(clusterGroupKeys []clusterv1beta1.GroupKey, clusterGroups map[clusterv1beta1.GroupKey]sets.Set[string], timeout time.Duration, statusFunc ClusterRolloutStatusFunc) RolloutResult {
	rolloutclusters := map[string]ClusterRolloutStatus{}
	timeoutclusters := map[string]ClusterRolloutStatus{}

	for _, key := range clusterGroupKeys {
		if subclusters, ok := clusterGroups[key]; ok {
			// go through group by group
			for _, cluster := range subclusters.UnsortedList() {
				status := statusFunc(cluster)
				status.GroupKey = key

				newstatus, needtorollout := needToRollout(status, timeout)
				if needtorollout {
					rolloutclusters[cluster] = status
				}
				if newstatus == TimeOut {
					status.Status = newstatus
					timeoutclusters[cluster] = status
				}
			}

			if len(rolloutclusters) > 0 {
				return RolloutResult{
					ClustersToRollout: rolloutclusters,
					ClustersTimeOut:   timeoutclusters,
				}
			}
		}
	}

	return RolloutResult{
		ClustersToRollout: rolloutclusters,
		ClustersTimeOut:   timeoutclusters,
	}
}

// Check if the cluster need to update based on existing cluster status and timeout.
// Timeout is used for status progressing and failed:
// 1) When timeout is None (maxTimeDuration), it means will wait until reach success status.
// Return true to append it to the result and stop rollout other clusters or groups.
// 2) Timeout is 0 means continue upgrade others without any wait.
// Return false to skip updating it and continue rollout other clusters or groups.
func needToRollout(status ClusterRolloutStatus, timeout time.Duration) (RolloutStatus, bool) {
	switch status.Status {
	case ToApply:
		return status.Status, true
	case Progressing:
		if beforeTimeOut(status.LastTransitionTime, timeout) {
			return Progressing, true
		}
		return TimeOut, false
	case Failed:
		if beforeTimeOut(status.LastTransitionTime, timeout) {
			return Failed, true
		}
		return TimeOut, false
	case TimeOut:
		return status.Status, false
	case Succeed:
		return status.Status, false
	case Skip:
		return status.Status, false
	default:
		return status.Status, true
	}
}

// check if current time is before start time + timeout duration
func beforeTimeOut(startTime *metav1.Time, timeout time.Duration) bool {
	var timeoutTime time.Time
	if startTime == nil {
		timeoutTime = RolloutClock.Now().Add(timeout)
	} else {
		timeoutTime = startTime.Add(timeout)
	}
	if RolloutClock.Now().Before(timeoutTime) {
		return true
	}
	return false
}

func calculateLength(maxConcurrency intstr.IntOrString, total int) (int, error) {
	length := total

	switch maxConcurrency.Type {
	case intstr.Int:
		length = maxConcurrency.IntValue()
	case intstr.String:
		str := maxConcurrency.StrVal
		if strings.HasSuffix(str, "%") {
			f, err := strconv.ParseFloat(str[:len(str)-1], 64)
			if err != nil {
				return length, err
			}
			length = int(math.Ceil(f / 100 * float64(total)))
		} else {
			return length, fmt.Errorf("%v invalid type: string is not a percentage", maxConcurrency)
		}
	default:
		return length, fmt.Errorf("incorrect MaxConcurrency type %v", maxConcurrency.Type)
	}

	if length <= 0 || length > total {
		length = total
	}

	return length, nil
}

func parseTimeout(timeoutStr string) (time.Duration, error) {
	// Define the regex pattern to match the timeout string
	pattern := "^(([0-9])+[h|m|s])|None$"
	regex := regexp.MustCompile(pattern)

	if timeoutStr == "None" || timeoutStr == "" {
		// If the timeout is "None" or empty, return the maximum duration
		return maxTimeDuration, nil
	}

	// Check if the timeout string matches the pattern
	if !regex.MatchString(timeoutStr) {
		return maxTimeDuration, fmt.Errorf("invalid timeout format")
	}

	return time.ParseDuration(timeoutStr)
}

func decisionGroupsToGroupKeys(decisionsGroup []MandatoryDecisionGroup) []clusterv1beta1.GroupKey {
	result := []clusterv1beta1.GroupKey{}
	for _, d := range decisionsGroup {
		gk := clusterv1beta1.GroupKey{}
		// GroupName is considered first to select the decisionGroups then GroupIndex.
		if d.GroupName != "" {
			gk.GroupName = d.GroupName
		} else {
			gk.GroupIndex = d.GroupIndex
		}
		result = append(result, gk)
	}
	return result
}
