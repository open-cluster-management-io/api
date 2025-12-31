// Copyright Contributors to the Open Cluster Management project
package v1beta1

import (
	"testing"

	certificates "k8s.io/api/certificates/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"open-cluster-management.io/api/addon/v1alpha1"
)

func TestConvert_v1alpha1_ClusterManagementAddOn_To_v1beta1_ClusterManagementAddOn(t *testing.T) {
	tests := []struct {
		name    string
		in      *v1alpha1.ClusterManagementAddOn
		want    *ClusterManagementAddOn
		wantErr bool
	}{
		{
			name: "basic conversion with default config",
			in: &v1alpha1.ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: v1alpha1.ClusterManagementAddOnSpec{
					AddOnMeta: v1alpha1.AddOnMeta{
						DisplayName: "Test Addon",
						Description: "Test Description",
					},
					SupportedConfigs: []v1alpha1.ConfigMeta{
						{
							ConfigGroupResource: v1alpha1.ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							DefaultConfig: &v1alpha1.ConfigReferent{
								Name:      "test-config",
								Namespace: "test-namespace",
							},
						},
					},
					InstallStrategy: v1alpha1.InstallStrategy{
						Type: v1alpha1.AddonInstallStrategyManual,
					},
				},
			},
			want: &ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: ClusterManagementAddOnSpec{
					AddOnMeta: AddOnMeta{
						DisplayName: "Test Addon",
						Description: "Test Description",
					},
					DefaultConfigs: []AddOnConfig{
						{
							ConfigGroupResource: ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							ConfigReferent: ConfigReferent{
								Name:      "test-config",
								Namespace: "test-namespace",
							},
						},
					},
					InstallStrategy: InstallStrategy{
						Type: AddonInstallStrategyManual,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with nil default config",
			in: &v1alpha1.ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: v1alpha1.ClusterManagementAddOnSpec{
					SupportedConfigs: []v1alpha1.ConfigMeta{
						{
							ConfigGroupResource: v1alpha1.ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							DefaultConfig: nil,
						},
					},
				},
			},
			want: &ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: ClusterManagementAddOnSpec{
					DefaultConfigs: []AddOnConfig{
						{
							ConfigGroupResource: ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							ConfigReferent: ConfigReferent{
								Name: ReservedNoDefaultConfigName,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with empty supported configs",
			in: &v1alpha1.ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: v1alpha1.ClusterManagementAddOnSpec{
					SupportedConfigs: []v1alpha1.ConfigMeta{},
				},
			},
			want: &ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: ClusterManagementAddOnSpec{
					DefaultConfigs: []AddOnConfig{},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &ClusterManagementAddOn{}
			err := Convert_v1alpha1_ClusterManagementAddOn_To_v1beta1_ClusterManagementAddOn(tt.in, got, conversion.Scope(nil))
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert_v1alpha1_ClusterManagementAddOn_To_v1beta1_ClusterManagementAddOn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Name != tt.want.Name {
					t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
				}
				if got.Spec.AddOnMeta.DisplayName != tt.want.Spec.AddOnMeta.DisplayName {
					t.Errorf("DisplayName = %v, want %v", got.Spec.AddOnMeta.DisplayName, tt.want.Spec.AddOnMeta.DisplayName)
				}
				if len(got.Spec.DefaultConfigs) != len(tt.want.Spec.DefaultConfigs) {
					t.Errorf("DefaultConfigs length = %v, want %v", len(got.Spec.DefaultConfigs), len(tt.want.Spec.DefaultConfigs))
				}
				for i := range got.Spec.DefaultConfigs {
					if got.Spec.DefaultConfigs[i].Name != tt.want.Spec.DefaultConfigs[i].Name {
						t.Errorf("DefaultConfigs[%d].Name = %v, want %v", i, got.Spec.DefaultConfigs[i].Name, tt.want.Spec.DefaultConfigs[i].Name)
					}
					if got.Spec.DefaultConfigs[i].Namespace != tt.want.Spec.DefaultConfigs[i].Namespace {
						t.Errorf("DefaultConfigs[%d].Namespace = %v, want %v", i, got.Spec.DefaultConfigs[i].Namespace, tt.want.Spec.DefaultConfigs[i].Namespace)
					}
					if got.Spec.DefaultConfigs[i].Resource != tt.want.Spec.DefaultConfigs[i].Resource {
						t.Errorf("DefaultConfigs[%d].Resource = %v, want %v", i, got.Spec.DefaultConfigs[i].Resource, tt.want.Spec.DefaultConfigs[i].Resource)
					}
				}
			}
		})
	}
}

// nolint:staticcheck
func TestConvert_v1beta1_ClusterManagementAddOn_To_v1alpha1_ClusterManagementAddOn(t *testing.T) {
	tests := []struct {
		name    string
		in      *ClusterManagementAddOn
		want    *v1alpha1.ClusterManagementAddOn
		wantErr bool
	}{
		{
			name: "basic conversion with default config",
			in: &ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: ClusterManagementAddOnSpec{
					AddOnMeta: AddOnMeta{
						DisplayName: "Test Addon",
						Description: "Test Description",
					},
					DefaultConfigs: []AddOnConfig{
						{
							ConfigGroupResource: ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							ConfigReferent: ConfigReferent{
								Name:      "test-config",
								Namespace: "test-namespace",
							},
						},
					},
					InstallStrategy: InstallStrategy{
						Type: AddonInstallStrategyManual,
					},
				},
			},
			want: &v1alpha1.ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: v1alpha1.ClusterManagementAddOnSpec{
					AddOnMeta: v1alpha1.AddOnMeta{
						DisplayName: "Test Addon",
						Description: "Test Description",
					},
					SupportedConfigs: []v1alpha1.ConfigMeta{
						{
							ConfigGroupResource: v1alpha1.ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							DefaultConfig: &v1alpha1.ConfigReferent{
								Name:      "test-config",
								Namespace: "test-namespace",
							},
						},
					},
					InstallStrategy: v1alpha1.InstallStrategy{
						Type: v1alpha1.AddonInstallStrategyManual,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with reserved sentinel value",
			in: &ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: ClusterManagementAddOnSpec{
					DefaultConfigs: []AddOnConfig{
						{
							ConfigGroupResource: ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							ConfigReferent: ConfigReferent{
								Name: ReservedNoDefaultConfigName,
							},
						},
					},
				},
			},
			want: &v1alpha1.ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: v1alpha1.ClusterManagementAddOnSpec{
					SupportedConfigs: []v1alpha1.ConfigMeta{
						{
							ConfigGroupResource: v1alpha1.ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							DefaultConfig: nil,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with empty default configs",
			in: &ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: ClusterManagementAddOnSpec{
					DefaultConfigs: []AddOnConfig{},
				},
			},
			want: &v1alpha1.ClusterManagementAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-addon",
				},
				Spec: v1alpha1.ClusterManagementAddOnSpec{
					SupportedConfigs: []v1alpha1.ConfigMeta{},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &v1alpha1.ClusterManagementAddOn{}
			err := Convert_v1beta1_ClusterManagementAddOn_To_v1alpha1_ClusterManagementAddOn(tt.in, got, conversion.Scope(nil))
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert_v1beta1_ClusterManagementAddOn_To_v1alpha1_ClusterManagementAddOn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Name != tt.want.Name {
					t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
				}
				if got.Spec.AddOnMeta.DisplayName != tt.want.Spec.AddOnMeta.DisplayName {
					t.Errorf("DisplayName = %v, want %v", got.Spec.AddOnMeta.DisplayName, tt.want.Spec.AddOnMeta.DisplayName)
				}
				if len(got.Spec.SupportedConfigs) != len(tt.want.Spec.SupportedConfigs) {
					t.Errorf("SupportedConfigs length = %v, want %v", len(got.Spec.SupportedConfigs), len(tt.want.Spec.SupportedConfigs))
				}
				for i := range got.Spec.SupportedConfigs {
					if tt.want.Spec.SupportedConfigs[i].DefaultConfig == nil {
						if got.Spec.SupportedConfigs[i].DefaultConfig != nil {
							t.Errorf("SupportedConfigs[%d].DefaultConfig = %v, want nil", i, got.Spec.SupportedConfigs[i].DefaultConfig)
						}
					} else {
						if got.Spec.SupportedConfigs[i].DefaultConfig == nil {
							t.Errorf("SupportedConfigs[%d].DefaultConfig = nil, want non-nil", i)
						} else {
							if got.Spec.SupportedConfigs[i].DefaultConfig.Name != tt.want.Spec.SupportedConfigs[i].DefaultConfig.Name {
								t.Errorf("SupportedConfigs[%d].DefaultConfig.Name = %v, want %v", i, got.Spec.SupportedConfigs[i].DefaultConfig.Name, tt.want.Spec.SupportedConfigs[i].DefaultConfig.Name)
							}
							if got.Spec.SupportedConfigs[i].DefaultConfig.Namespace != tt.want.Spec.SupportedConfigs[i].DefaultConfig.Namespace {
								t.Errorf("SupportedConfigs[%d].DefaultConfig.Namespace = %v, want %v", i, got.Spec.SupportedConfigs[i].DefaultConfig.Namespace, tt.want.Spec.SupportedConfigs[i].DefaultConfig.Namespace)
							}
						}
					}
					if got.Spec.SupportedConfigs[i].Resource != tt.want.Spec.SupportedConfigs[i].Resource {
						t.Errorf("SupportedConfigs[%d].Resource = %v, want %v", i, got.Spec.SupportedConfigs[i].Resource, tt.want.Spec.SupportedConfigs[i].Resource)
					}
				}
			}
		})
	}
}

func TestConvert_v1alpha1_ManagedClusterAddOn_To_v1beta1_ManagedClusterAddOn(t *testing.T) {
	tests := []struct {
		name    string
		in      *v1alpha1.ManagedClusterAddOn
		want    *ManagedClusterAddOn
		wantErr bool
	}{
		{
			name: "basic conversion with configs",
			in: &v1alpha1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: v1alpha1.ManagedClusterAddOnSpec{
					InstallNamespace: "test-install-ns",
					Configs: []v1alpha1.AddOnConfig{
						{
							ConfigGroupResource: v1alpha1.ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							ConfigReferent: v1alpha1.ConfigReferent{
								Name:      "test-config",
								Namespace: "test-namespace",
							},
						},
					},
				},
				Status: v1alpha1.ManagedClusterAddOnStatus{
					Namespace: "test-ns",
					Registrations: []v1alpha1.RegistrationConfig{
						{
							SignerName: certificates.KubeAPIServerClientSignerName,
							Subject: v1alpha1.Subject{
								User:   "test-user",
								Groups: []string{"test-group"},
							},
						},
					},
				},
			},
			want: &ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: ManagedClusterAddOnSpec{
					Configs: []AddOnConfig{
						{
							ConfigGroupResource: ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							ConfigReferent: ConfigReferent{
								Name:      "test-config",
								Namespace: "test-namespace",
							},
						},
					},
				},
				Status: ManagedClusterAddOnStatus{
					Namespace: "test-ns",
					Registrations: []RegistrationConfig{
						{
							Type: KubeClient,
							KubeClient: &KubeClientConfig{
								Subject: KubeClientSubject{
									BaseSubject: BaseSubject{
										User:   "test-user",
										Groups: []string{"test-group"},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with CSR registration",
			in: &v1alpha1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: v1alpha1.ManagedClusterAddOnSpec{
					Configs: []v1alpha1.AddOnConfig{},
				},
				Status: v1alpha1.ManagedClusterAddOnStatus{
					Registrations: []v1alpha1.RegistrationConfig{
						{
							SignerName: "custom.signer.io/custom",
							Subject: v1alpha1.Subject{
								User:              "test-user",
								Groups:            []string{"test-group"},
								OrganizationUnits: []string{"test-ou"},
							},
						},
					},
				},
			},
			want: &ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: ManagedClusterAddOnSpec{
					Configs: []AddOnConfig{},
				},
				Status: ManagedClusterAddOnStatus{
					Registrations: []RegistrationConfig{
						{
							Type: CustomSigner,
							CustomSigner: &CustomSignerConfig{
								SignerName: "custom.signer.io/custom",
								Subject: Subject{
									BaseSubject: BaseSubject{
										User:   "test-user",
										Groups: []string{"test-group"},
									},
									OrganizationUnits: []string{"test-ou"},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with KubeClient registration and csr driver",
			in: &v1alpha1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: v1alpha1.ManagedClusterAddOnSpec{
					Configs: []v1alpha1.AddOnConfig{},
				},
				Status: v1alpha1.ManagedClusterAddOnStatus{
					Registrations: []v1alpha1.RegistrationConfig{
						{
							SignerName: certificates.KubeAPIServerClientSignerName,
							Subject: v1alpha1.Subject{
								User:   "test-user",
								Groups: []string{"test-group"},
							},
							Driver: "csr",
						},
					},
				},
			},
			want: &ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: ManagedClusterAddOnSpec{
					Configs: []AddOnConfig{},
				},
				Status: ManagedClusterAddOnStatus{
					Registrations: []RegistrationConfig{
						{
							Type: KubeClient,
							KubeClient: &KubeClientConfig{
								Subject: KubeClientSubject{
									BaseSubject: BaseSubject{
										User:   "test-user",
										Groups: []string{"test-group"},
									},
								},
								Driver: "csr",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with KubeClient registration and token driver",
			in: &v1alpha1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: v1alpha1.ManagedClusterAddOnSpec{
					Configs: []v1alpha1.AddOnConfig{},
				},
				Status: v1alpha1.ManagedClusterAddOnStatus{
					Registrations: []v1alpha1.RegistrationConfig{
						{
							SignerName: certificates.KubeAPIServerClientSignerName,
							Subject: v1alpha1.Subject{
								User:   "test-user",
								Groups: []string{"test-group"},
							},
							Driver: "token",
						},
					},
				},
			},
			want: &ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: ManagedClusterAddOnSpec{
					Configs: []AddOnConfig{},
				},
				Status: ManagedClusterAddOnStatus{
					Registrations: []RegistrationConfig{
						{
							Type: KubeClient,
							KubeClient: &KubeClientConfig{
								Subject: KubeClientSubject{
									BaseSubject: BaseSubject{
										User:   "test-user",
										Groups: []string{"test-group"},
									},
								},
								Driver: "token",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with empty configs",
			in: &v1alpha1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: v1alpha1.ManagedClusterAddOnSpec{
					Configs: []v1alpha1.AddOnConfig{},
				},
			},
			want: &ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: ManagedClusterAddOnSpec{
					Configs: []AddOnConfig{},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &ManagedClusterAddOn{}
			err := Convert_v1alpha1_ManagedClusterAddOn_To_v1beta1_ManagedClusterAddOn(tt.in, got, conversion.Scope(nil))
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert_v1alpha1_ManagedClusterAddOn_To_v1beta1_ManagedClusterAddOn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Name != tt.want.Name {
					t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
				}
				if got.Namespace != tt.want.Namespace {
					t.Errorf("Namespace = %v, want %v", got.Namespace, tt.want.Namespace)
				}
				if len(got.Spec.Configs) != len(tt.want.Spec.Configs) {
					t.Errorf("Configs length = %v, want %v", len(got.Spec.Configs), len(tt.want.Spec.Configs))
				}
				for i := range got.Spec.Configs {
					if got.Spec.Configs[i].Name != tt.want.Spec.Configs[i].Name {
						t.Errorf("Configs[%d].Name = %v, want %v", i, got.Spec.Configs[i].Name, tt.want.Spec.Configs[i].Name)
					}
				}
				if len(got.Status.Registrations) != len(tt.want.Status.Registrations) {
					t.Errorf("Registrations length = %v, want %v", len(got.Status.Registrations), len(tt.want.Status.Registrations))
				}
				for i := range got.Status.Registrations {
					if got.Status.Registrations[i].Type != tt.want.Status.Registrations[i].Type {
						t.Errorf("Registrations[%d].Type = %v, want %v", i, got.Status.Registrations[i].Type, tt.want.Status.Registrations[i].Type)
					}
					if got.Status.Registrations[i].Type == KubeClient {
						if got.Status.Registrations[i].KubeClient == nil {
							t.Errorf("Registrations[%d].KubeClient is nil", i)
						} else {
							if got.Status.Registrations[i].KubeClient.Subject.User != tt.want.Status.Registrations[i].KubeClient.Subject.User {
								t.Errorf("Registrations[%d].KubeClient.Subject.User = %v, want %v", i, got.Status.Registrations[i].KubeClient.Subject.User, tt.want.Status.Registrations[i].KubeClient.Subject.User)
							}
							if got.Status.Registrations[i].KubeClient.Driver != tt.want.Status.Registrations[i].KubeClient.Driver {
								t.Errorf("Registrations[%d].KubeClient.Driver = %v, want %v", i, got.Status.Registrations[i].KubeClient.Driver, tt.want.Status.Registrations[i].KubeClient.Driver)
							}
						}
					}
					if got.Status.Registrations[i].Type == CustomSigner {
						if got.Status.Registrations[i].CustomSigner == nil {
							t.Errorf("Registrations[%d].CustomSigner is nil", i)
						} else {
							if got.Status.Registrations[i].CustomSigner.SignerName != tt.want.Status.Registrations[i].CustomSigner.SignerName {
								t.Errorf("Registrations[%d].CustomSigner.SignerName = %v, want %v", i, got.Status.Registrations[i].CustomSigner.SignerName, tt.want.Status.Registrations[i].CustomSigner.SignerName)
							}
							if got.Status.Registrations[i].CustomSigner.Subject.User != tt.want.Status.Registrations[i].CustomSigner.Subject.User {
								t.Errorf("Registrations[%d].CustomSigner.Subject.User = %v, want %v", i, got.Status.Registrations[i].CustomSigner.Subject.User, tt.want.Status.Registrations[i].CustomSigner.Subject.User)
							}
						}
					}
				}
			}
		})
	}
}

// nolint:staticcheck
func TestConvert_v1beta1_ManagedClusterAddOn_To_v1alpha1_ManagedClusterAddOn(t *testing.T) {
	tests := []struct {
		name    string
		in      *ManagedClusterAddOn
		want    *v1alpha1.ManagedClusterAddOn
		wantErr bool
	}{
		{
			name: "basic conversion with KubeClient registration",
			in: &ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: ManagedClusterAddOnSpec{
					Configs: []AddOnConfig{
						{
							ConfigGroupResource: ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							ConfigReferent: ConfigReferent{
								Name:      "test-config",
								Namespace: "test-namespace",
							},
						},
					},
				},
				Status: ManagedClusterAddOnStatus{
					Namespace: "test-ns",
					Registrations: []RegistrationConfig{
						{
							Type: KubeClient,
							KubeClient: &KubeClientConfig{
								Subject: KubeClientSubject{
									BaseSubject: BaseSubject{
										User:   "test-user",
										Groups: []string{"test-group"},
									},
								},
							},
						},
					},
				},
			},
			want: &v1alpha1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: v1alpha1.ManagedClusterAddOnSpec{
					Configs: []v1alpha1.AddOnConfig{
						{
							ConfigGroupResource: v1alpha1.ConfigGroupResource{
								Group:    "test.group",
								Resource: "testconfigs",
							},
							ConfigReferent: v1alpha1.ConfigReferent{
								Name:      "test-config",
								Namespace: "test-namespace",
							},
						},
					},
				},
				Status: v1alpha1.ManagedClusterAddOnStatus{
					Namespace: "test-ns",
					Registrations: []v1alpha1.RegistrationConfig{
						{
							SignerName: certificates.KubeAPIServerClientSignerName,
							Subject: v1alpha1.Subject{
								User:   "test-user",
								Groups: []string{"test-group"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with CSR registration",
			in: &ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: ManagedClusterAddOnSpec{
					Configs: []AddOnConfig{},
				},
				Status: ManagedClusterAddOnStatus{
					Registrations: []RegistrationConfig{
						{
							Type: CustomSigner,
							CustomSigner: &CustomSignerConfig{
								SignerName: "custom.signer.io/custom",
								Subject: Subject{
									BaseSubject: BaseSubject{
										User:   "test-user",
										Groups: []string{"test-group"},
									},
									OrganizationUnits: []string{"test-ou"},
								},
							},
						},
					},
				},
			},
			want: &v1alpha1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: v1alpha1.ManagedClusterAddOnSpec{
					Configs: []v1alpha1.AddOnConfig{},
				},
				Status: v1alpha1.ManagedClusterAddOnStatus{
					Registrations: []v1alpha1.RegistrationConfig{
						{
							SignerName: "custom.signer.io/custom",
							Subject: v1alpha1.Subject{
								User:              "test-user",
								Groups:            []string{"test-group"},
								OrganizationUnits: []string{"test-ou"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with KubeClient registration and csr driver",
			in: &ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: ManagedClusterAddOnSpec{
					Configs: []AddOnConfig{},
				},
				Status: ManagedClusterAddOnStatus{
					Registrations: []RegistrationConfig{
						{
							Type: KubeClient,
							KubeClient: &KubeClientConfig{
								Subject: KubeClientSubject{
									BaseSubject: BaseSubject{
										User:   "test-user",
										Groups: []string{"test-group"},
									},
								},
								Driver: "csr",
							},
						},
					},
				},
			},
			want: &v1alpha1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: v1alpha1.ManagedClusterAddOnSpec{
					Configs: []v1alpha1.AddOnConfig{},
				},
				Status: v1alpha1.ManagedClusterAddOnStatus{
					Registrations: []v1alpha1.RegistrationConfig{
						{
							SignerName: certificates.KubeAPIServerClientSignerName,
							Subject: v1alpha1.Subject{
								User:   "test-user",
								Groups: []string{"test-group"},
							},
							Driver: "csr",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with KubeClient registration and token driver",
			in: &ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: ManagedClusterAddOnSpec{
					Configs: []AddOnConfig{},
				},
				Status: ManagedClusterAddOnStatus{
					Registrations: []RegistrationConfig{
						{
							Type: KubeClient,
							KubeClient: &KubeClientConfig{
								Subject: KubeClientSubject{
									BaseSubject: BaseSubject{
										User:   "test-user",
										Groups: []string{"test-group"},
									},
								},
								Driver: "token",
							},
						},
					},
				},
			},
			want: &v1alpha1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: v1alpha1.ManagedClusterAddOnSpec{
					Configs: []v1alpha1.AddOnConfig{},
				},
				Status: v1alpha1.ManagedClusterAddOnStatus{
					Registrations: []v1alpha1.RegistrationConfig{
						{
							SignerName: certificates.KubeAPIServerClientSignerName,
							Subject: v1alpha1.Subject{
								User:   "test-user",
								Groups: []string{"test-group"},
							},
							Driver: "token",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conversion with empty configs",
			in: &ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: ManagedClusterAddOnSpec{
					Configs: []AddOnConfig{},
				},
			},
			want: &v1alpha1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-addon",
					Namespace: "test-cluster",
				},
				Spec: v1alpha1.ManagedClusterAddOnSpec{
					Configs: []v1alpha1.AddOnConfig{},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &v1alpha1.ManagedClusterAddOn{}
			err := Convert_v1beta1_ManagedClusterAddOn_To_v1alpha1_ManagedClusterAddOn(tt.in, got, conversion.Scope(nil))
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert_v1beta1_ManagedClusterAddOn_To_v1alpha1_ManagedClusterAddOn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Name != tt.want.Name {
					t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
				}
				if got.Namespace != tt.want.Namespace {
					t.Errorf("Namespace = %v, want %v", got.Namespace, tt.want.Namespace)
				}
				if len(got.Spec.Configs) != len(tt.want.Spec.Configs) {
					t.Errorf("Configs length = %v, want %v", len(got.Spec.Configs), len(tt.want.Spec.Configs))
				}
				for i := range got.Spec.Configs {
					if got.Spec.Configs[i].Name != tt.want.Spec.Configs[i].Name {
						t.Errorf("Configs[%d].Name = %v, want %v", i, got.Spec.Configs[i].Name, tt.want.Spec.Configs[i].Name)
					}
				}
				if len(got.Status.Registrations) != len(tt.want.Status.Registrations) {
					t.Errorf("Registrations length = %v, want %v", len(got.Status.Registrations), len(tt.want.Status.Registrations))
				}
				for i := range got.Status.Registrations {
					if got.Status.Registrations[i].SignerName != tt.want.Status.Registrations[i].SignerName {
						t.Errorf("Registrations[%d].SignerName = %v, want %v", i, got.Status.Registrations[i].SignerName, tt.want.Status.Registrations[i].SignerName)
					}
					if got.Status.Registrations[i].Subject.User != tt.want.Status.Registrations[i].Subject.User {
						t.Errorf("Registrations[%d].Subject.User = %v, want %v", i, got.Status.Registrations[i].Subject.User, tt.want.Status.Registrations[i].Subject.User)
					}
					if got.Status.Registrations[i].Driver != tt.want.Status.Registrations[i].Driver {
						t.Errorf("Registrations[%d].Driver = %v, want %v", i, got.Status.Registrations[i].Driver, tt.want.Status.Registrations[i].Driver)
					}
					if len(got.Status.Registrations[i].Subject.OrganizationUnits) != len(tt.want.Status.Registrations[i].Subject.OrganizationUnits) {
						t.Errorf("Registrations[%d].Subject.OrganizationUnits length = %v, want %v", i, len(got.Status.Registrations[i].Subject.OrganizationUnits), len(tt.want.Status.Registrations[i].Subject.OrganizationUnits))
					}
				}
			}
		})
	}
}
