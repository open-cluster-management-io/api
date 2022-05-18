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
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	v1 "open-cluster-management.io/api/cluster/v1"
)

var (
	scheme = runtime.NewScheme()
)

type clusterGetter struct {
	client client.Client
}
type clusterSetGetter struct {
	client client.Client
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

func (mcl clusterGetter) List(selector labels.Selector) ([]*v1.ManagedCluster, error) {
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

func (msl clusterSetGetter) List(selector labels.Selector) ([]*ManagedClusterSet, error) {
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
	for _, clusters := range existingClusters {
		existingObjs = append(existingObjs, clusters)
	}
	mcl := clusterGetter{
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
	for _, clusters := range existingClusters {
		existingObjs = append(existingObjs, clusters)
	}
	for _, clusterset := range existingClusterSets {
		existingObjs = append(existingObjs, clusterset)
	}

	msl := clusterSetGetter{
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
