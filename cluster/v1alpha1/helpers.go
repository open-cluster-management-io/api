package v1alpha1

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/clock"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
)

type RolloutStatus int

var RolloutClock = clock.Clock(clock.RealClock{})
var maxTimeDuration = time.Duration(math.MaxInt64)

const (
	// resource desired status is not applied yet
	ToApply RolloutStatus = iota
	// resource desired status is applied and last applied status is not updated
	Progressing
	// resource desired status is applied and last applied status is succeed
	Succeed
	// resource desired status is applied and last applied status is failed
	Failed
)

// Return the rollout status on each managed cluster
type ClusterRolloutStatusFunc func(clusterName string) ClusterRolloutStatus

type ClusterRolloutStatus struct {
	// rollout status
	status RolloutStatus
	// the last transition time of rollout status
	lastTransitionTime *time.Time
}

// +k8s:deepcopy-gen=false
type RolloutHandler struct {
	// duck type rollout strategy
	rolloutStrategy RolloutStrategy
	// placement decision tracker
	pdTracker *clusterv1beta1.PlacementDecisionClustersTracker
	// Return the rollout status on each managed cluster
	clusterRolloutStatusFunc ClusterRolloutStatusFunc
}

func NewRolloutHandler(rolloutStrategy RolloutStrategy, pdTracker *clusterv1beta1.PlacementDecisionClustersTracker, clusterRolloutStatusFunc ClusterRolloutStatusFunc) *RolloutHandler {
	if pdTracker == nil {
		return nil
	}

	return &RolloutHandler{
		rolloutStrategy:          rolloutStrategy,
		pdTracker:                pdTracker,
		clusterRolloutStatusFunc: clusterRolloutStatusFunc,
	}
}

// Return the strategy actual take effect and a list of cluster that needs to update with current state
func (r *RolloutHandler) GetRolloutCluster() (*RolloutStrategy, map[string]ClusterRolloutStatus, error) {
	switch r.rolloutStrategy.Type {
	case All:
		return r.getRolloutAllClusters()
	case Progressive:
		return r.getProgressiveClusters()
	case ProgressivePerGroup:
		return r.getProgressivePerGroupClusters()
	default:
		return nil, map[string]ClusterRolloutStatus{}, fmt.Errorf("incorrect rollout strategy type %v", r.rolloutStrategy.Type)
	}
}

func (r *RolloutHandler) getRolloutAllClusters() (*RolloutStrategy, map[string]ClusterRolloutStatus, error) {
	strategy := RolloutStrategy{Type: All}
	strategy.All = r.rolloutStrategy.All.DeepCopy()
	if strategy.All == nil {
		strategy.All = &RolloutAll{}
	}

	// parse timeout
	failureTimeout, err := parseTimeout(strategy.All.Timeout.Timeout)
	if err != nil {
		return &strategy, map[string]ClusterRolloutStatus{}, err
	}

	// get all the clusters
	totalclusters := r.pdTracker.ExistingBesides([]clusterv1beta1.GroupKey{}).UnsortedList()

	result := r.progressivePerCluster(totalclusters, len(totalclusters), failureTimeout)
	return &strategy, result, nil
}

func (r *RolloutHandler) getProgressiveClusters() (*RolloutStrategy, map[string]ClusterRolloutStatus, error) {
	strategy := RolloutStrategy{Type: Progressive}
	strategy.Progressive = r.rolloutStrategy.Progressive.DeepCopy()
	if strategy.Progressive == nil {
		strategy.Progressive = &RolloutProgressive{}
	}

	// upgrade madatory decision groups first
	groupKeys := decisionGroupsToGroupKeys(strategy.Progressive.MandatoryDecisionGroups.MandatoryDecisionGroups)
	clusterGroupKeys, clusterGroups := r.pdTracker.ExistingClusterGroups(groupKeys)

	result := r.progressivePerGroup(clusterGroupKeys, clusterGroups, maxTimeDuration)
	if len(result) > 0 {
		return &strategy, result, nil
	}

	// parse timeout
	failureTimeout, err := parseTimeout(strategy.Progressive.Timeout.Timeout)
	if err != nil {
		return &strategy, result, err
	}

	// upgrade the rest clusters
	totalclusters := r.pdTracker.ExistingBesides([]clusterv1beta1.GroupKey{})
	length, err := calculateLength(strategy.Progressive.MaxConcurrency, len(totalclusters))
	if err != nil {
		return &strategy, result, err
	}
	restclusters := r.pdTracker.ExistingBesides(clusterGroupKeys).UnsortedList()

	result = r.progressivePerCluster(restclusters, length, failureTimeout)
	return &strategy, result, nil
}

func (r *RolloutHandler) getProgressivePerGroupClusters() (*RolloutStrategy, map[string]ClusterRolloutStatus, error) {
	strategy := RolloutStrategy{Type: ProgressivePerGroup}
	strategy.ProgressivePerGroup = r.rolloutStrategy.ProgressivePerGroup.DeepCopy()
	if strategy.ProgressivePerGroup == nil {
		strategy.ProgressivePerGroup = &RolloutProgressivePerGroup{}
	}

	// upgrade madatory decision groups first
	groupKeys := decisionGroupsToGroupKeys(strategy.ProgressivePerGroup.MandatoryDecisionGroups.MandatoryDecisionGroups)
	clusterGroupKeys, clusterGroups := r.pdTracker.ExistingClusterGroups(groupKeys)

	result := r.progressivePerGroup(clusterGroupKeys, clusterGroups, maxTimeDuration)
	if len(result) > 0 {
		return &strategy, result, nil
	}

	// parse timeout
	failureTimeout, err := parseTimeout(strategy.ProgressivePerGroup.Timeout.Timeout)
	if err != nil {
		return &strategy, result, err
	}

	// upgrade the rest decision groups
	restClusterGroupKeys, restClusterGroups := r.pdTracker.ExistingClusterGroupsBesides(clusterGroupKeys)

	result = r.progressivePerGroup(restClusterGroupKeys, restClusterGroups, failureTimeout)
	return &strategy, result, nil
}

func (r *RolloutHandler) progressivePerCluster(clusters []string, length int, timeout time.Duration) map[string]ClusterRolloutStatus {
	result := map[string]ClusterRolloutStatus{}
	if length == 0 {
		return result
	}

	// sort the clusters
	sort.Strings(clusters)
	for _, cluster := range clusters {
		status := r.clusterRolloutStatusFunc(cluster)
		if needToUpdate(status, timeout) {
			result[cluster] = status
		}

		if (len(result)%length == 0) && len(result) > 0 {
			return result
		}
	}

	return result
}

func (r *RolloutHandler) progressivePerGroup(clusterGroupKeys []clusterv1beta1.GroupKey, clusterGroups map[clusterv1beta1.GroupKey]sets.Set[string], timeout time.Duration) map[string]ClusterRolloutStatus {
	result := map[string]ClusterRolloutStatus{}

	for _, key := range clusterGroupKeys {
		if subclusters, ok := clusterGroups[key]; ok {
			for _, cluster := range subclusters.UnsortedList() {
				status := r.clusterRolloutStatusFunc(cluster)
				if needToUpdate(status, timeout) {
					result[cluster] = status
				}
			}

			if len(result) > 0 {
				return result
			}
		}
	}

	return result
}

func needToUpdate(status ClusterRolloutStatus, timeout time.Duration) bool {
	switch status.status {
	case ToApply:
		return true
	case Progressing:
		return true
	case Succeed:
		return false
	case Failed:
		// 1) When timeout is None, it means stop upgrade when meet failure, default is None.
		// Return true to append it to the result and stop upgrading other clusters or groups.
		// 2) Timeout is 0 means continue upgrade others when meet failure.
		// Return false to skip updating it and continue upgrading other clusters or groups.
		var checkFailureUntil time.Time
		if status.lastTransitionTime == nil {
			checkFailureUntil = RolloutClock.Now().Add(timeout)
		} else {
			checkFailureUntil = status.lastTransitionTime.Add(timeout)
		}
		if RolloutClock.Now().Before(checkFailureUntil) {
			return true
		}
		return false
	default:
		return true
	}
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

	// Extract the numerical part and unit from the string
	numStr := timeoutStr[:len(timeoutStr)-1]
	unit := timeoutStr[len(timeoutStr)-1]

	// Convert the numerical part to an integer
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return maxTimeDuration, fmt.Errorf("invalid numeric value in timeout string")
	}

	// Calculate the duration based on the unit
	var duration time.Duration
	switch unit {
	case 'h':
		duration = time.Duration(num) * time.Hour
	case 'm':
		duration = time.Duration(num) * time.Minute
	case 's':
		duration = time.Duration(num) * time.Second
	default:
		return maxTimeDuration, fmt.Errorf("invalid unit in timeout string")
	}

	return duration, nil
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
