package integration

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	operatorv1 "open-cluster-management.io/api/operator/v1"
)

var _ = Describe("Create Klusterlet API", func() {
	var klusterlet *operatorv1.Klusterlet
	BeforeEach(func() {
		suffix := rand.String(5)
		klusterManagerName := fmt.Sprintf("cm-%s", suffix)
		klusterlet = &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterManagerName,
			},
			Spec: operatorv1.KlusterletSpec{},
		}
	})

	Context("Create without nothing set", func() {
		It("should create successfully", func() {
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
	})
})
