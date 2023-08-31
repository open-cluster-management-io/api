package utils

import (
	"encoding/json"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	workv1 "open-cluster-management.io/api/work/v1"
)

func TestPatch(t *testing.T) {
	cases := []struct {
		name      string
		patchType types.PatchType
		work      *workv1.ManifestWork
		patch     []byte
		validate  func(t *testing.T, work *workv1.ManifestWork)
	}{
		{
			name:      "json patch",
			patchType: types.JSONPatchType,
			work: &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			},
			patch: []byte("[{\"op\":\"replace\",\"path\":\"/metadata/name\",\"value\":\"test1\"}]"),
			validate: func(t *testing.T, work *workv1.ManifestWork) {
				if work.Name != "test1" {
					t.Errorf("unexpected work %v", work)
				}
			},
		},
		{
			name:      "merge patch",
			patchType: types.MergePatchType,
			work: &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			},
			patch: func() []byte {
				newWork := &workv1.ManifestWork{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test2",
						Namespace: "test2",
					},
				}
				data, err := json.Marshal(newWork)
				if err != nil {
					t.Fatal(err)
				}
				return data
			}(),
			validate: func(t *testing.T, work *workv1.ManifestWork) {
				if work.Name != "test2" {
					t.Errorf("unexpected work %v", work)
				}
				if work.Namespace != "test2" {
					t.Errorf("unexpected work %v", work)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			work, err := Patch(c.patchType, c.work, c.patch)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			c.validate(t, work)
		})
	}
}
