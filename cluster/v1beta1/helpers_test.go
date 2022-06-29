package v1beta1

import (
	"context"
	"os"
	"reflect"
	"testing"

	cliScheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	v1.AddToScheme(cliScheme.Scheme)
	AddToScheme(cliScheme.Scheme)

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
		expectClustersName sets.String
		expectError        bool
	}{
		{
			name: "test legency cluster set",
			clusterset: &ManagedClusterSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dev",
				},
				Spec: ManagedClusterSetSpec{},
			},
			expectClustersName: sets.NewString("c1", "c2"),
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
			expectClustersName: sets.NewString("c1", "c2"),
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
			expectClustersName: sets.NewString("c1", "c2", "c3"),
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
		expectClusterSetName sets.String
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
			expectClusterSetName: sets.NewString("dev", "openshift", "global"),
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
			expectClusterSetName: sets.NewString("dev", "openshift", "global"),
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
			expectClusterSetName: sets.NewString("global"),
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
			expectClusterSetName: sets.NewString("openshift", "global"),
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

type placementDecisionGetter struct {
	client client.Client
}

func (pdl placementDecisionGetter) List(selector labels.Selector, namespace string) ([]*PlacementDecision, error) {
	decisionList := PlacementDecisionList{}
	err := pdl.client.List(context.Background(), &decisionList, &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	var decisions []*PlacementDecision
	for i := range decisionList.Items {
		decisions = append(decisions, &decisionList.Items[i])
	}
	return decisions, nil
}

func TestPlacementDecisionClustersTracker(t *testing.T) {
	tests := []struct {
		name                            string
		placement                       Placement
		existingDecisions               []*PlacementDecision
		expectExistingScheduledClusters sets.String
		updateDecisions                 []*PlacementDecision
		expectUpdatedScheduledClusters  sets.String
		expectAddedScheduledClusters    sets.String
		expectDeletedScheduledClusters  sets.String
	}{
		{
			name: "test placementdecisions",
			placement: Placement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "placement1",
					Namespace: "default",
				},
				Spec: PlacementSpec{},
			},
			existingDecisions: []*PlacementDecision{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "placement1-decision-1",
						Namespace: "default",
						Labels: map[string]string{
							PlacementLabel: "placement1",
						},
					},
					Status: PlacementDecisionStatus{
						Decisions: []ClusterDecision{
							{
								ClusterName: "cluster1",
								Reason:      "reason1",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "placement1-decision-2",
						Namespace: "default",
						Labels: map[string]string{
							PlacementLabel: "placement1",
						},
					},
					Status: PlacementDecisionStatus{
						Decisions: []ClusterDecision{
							{
								ClusterName: "cluster2",
								Reason:      "reason2",
							},
						},
					},
				},
			},
			updateDecisions: []*PlacementDecision{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "placement1-decision-2",
						Namespace: "default",
						Labels: map[string]string{
							PlacementLabel: "placement1",
						},
					},
					Status: PlacementDecisionStatus{
						Decisions: []ClusterDecision{
							{
								ClusterName: "cluster3",
								Reason:      "reason3",
							},
						},
					},
				},
			},
			expectExistingScheduledClusters: sets.NewString("cluster1", "cluster2"),
			expectUpdatedScheduledClusters:  sets.NewString("cluster1", "cluster3"),
			expectAddedScheduledClusters:    sets.NewString("cluster3"),
			expectDeletedScheduledClusters:  sets.NewString("cluster2"),
		},
		{
			name: "test empty placementdecision",
			placement: Placement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "placement2",
					Namespace: "default",
				},
				Spec: PlacementSpec{},
			},
			existingDecisions: []*PlacementDecision{},
			updateDecisions: []*PlacementDecision{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "placement2-decision-1",
						Namespace: "default",
						Labels: map[string]string{
							PlacementLabel: "placement2",
						},
					},
					Status: PlacementDecisionStatus{
						Decisions: []ClusterDecision{
							{
								ClusterName: "cluster1",
								Reason:      "reason1",
							},
							{
								ClusterName: "cluster2",
								Reason:      "reason2",
							},
						},
					},
				},
			},
			expectExistingScheduledClusters: sets.NewString(),
			expectUpdatedScheduledClusters:  sets.NewString("cluster1", "cluster2"),
			expectAddedScheduledClusters:    sets.NewString("cluster1", "cluster2"),
			expectDeletedScheduledClusters:  sets.NewString(),
		},
	}

	for _, test := range tests {
		// init decisions
		var existingObjs []client.Object
		for _, d := range test.existingDecisions {
			existingObjs = append(existingObjs, d)
		}

		pdl := placementDecisionGetter{
			client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(existingObjs...).Build(),
		}

		// init tracker
		tracker := NewPlacementDecisionClustersTracker(&test.placement, pdl, nil)

		// check decision clusters
		_, _, err := tracker.Get()
		if err != nil {
			t.Errorf("Case: %v, Failed to run Get(): %v", test.name, err)
		}
		if !reflect.DeepEqual(tracker.Existing(), test.expectExistingScheduledClusters) {
			t.Errorf("Case: %v, expect decisions: %v, return decisions: %v", test.name, test.expectExistingScheduledClusters, tracker.Existing())
			return
		}

		// update decisions
		for _, d := range test.updateDecisions {
			existDecision := &PlacementDecision{}
			err := pdl.client.Get(context.Background(), types.NamespacedName{Namespace: d.Namespace, Name: d.Name}, existDecision)
			if err == nil {
				existDecision.Status = d.Status
				pdl.client.Status().Update(context.Background(), existDecision)
			} else {
				pdl.client.Create(context.Background(), d)
			}
		}

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
		if !reflect.DeepEqual(tracker.Existing(), test.expectUpdatedScheduledClusters) {
			t.Errorf("Case: %v, expect updated decisions: %v, return decisions: %v", test.name, test.expectUpdatedScheduledClusters, tracker.Existing())
			return
		}
	}
}

func convertClusterToSet(clusters []*v1.ManagedCluster) sets.String {
	if len(clusters) == 0 {
		return nil
	}
	retSet := sets.NewString()
	for _, cluster := range clusters {
		retSet.Insert(cluster.Name)
	}
	return retSet
}

func convertClusterSetToSet(clustersets []*ManagedClusterSet) sets.String {
	if len(clustersets) == 0 {
		return nil
	}
	retSet := sets.NewString()
	for _, clusterset := range clustersets {
		retSet.Insert(clusterset.Name)
	}
	return retSet
}

func convertClusterSetBindingsToSet(clusterSetBindings []*ManagedClusterSetBinding) sets.String {
	if len(clusterSetBindings) == 0 {
		return nil
	}
	retSet := sets.NewString()
	for _, clusterSetBinding := range clusterSetBindings {
		retSet.Insert(clusterSetBinding.Name)
	}
	return retSet
}

func TestGetValidManagedClusterSetBindings(t *testing.T) {
	tests := []struct {
		name                          string
		namespace                     string
		expectClusterSetBindingsNames sets.String
		expectError                   bool
	}{
		{
			name:                          "test found valid cluster bindings only",
			namespace:                     "default",
			expectClusterSetBindingsNames: sets.NewString("dev"),
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
