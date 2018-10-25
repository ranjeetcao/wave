/*
Copyright 2018 Pusher Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package deployment

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pusher/wave/test/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Wave hash Suite", func() {
	// Waiting for calculateConfigHash to be implemented
	PContext("calculateConfigHash", func() {
		var cm1 *corev1.ConfigMap
		var cm2 *corev1.ConfigMap
		var s1 *corev1.Secret
		var s2 *corev1.Secret

		BeforeEach(func() {
			cm1 = utils.ExampleConfigMap1.DeepCopy()
			cm2 = utils.ExampleConfigMap2.DeepCopy()
			s1 = utils.ExampleSecret1.DeepCopy()
			s2 = utils.ExampleSecret2.DeepCopy()
		})

		It("returns a different hash when a child's data is updated", func() {
			c := []metav1.Object{cm1, cm2, s1, s2}

			h1, err := calculateConfigHash(c)
			Expect(err).NotTo(HaveOccurred())

			cm1.Data["key1"] = "modified"
			h2, err := calculateConfigHash(c)
			Expect(err).NotTo(HaveOccurred())

			Expect(h2).NotTo(Equal(h1))
		})

		It("returns the same hash when a child's metadata is updated", func() {
			c := []metav1.Object{cm1, cm2, s1, s2}

			h1, err := calculateConfigHash(c)
			Expect(err).NotTo(HaveOccurred())

			s1.Annotations = map[string]string{"new": "annotations"}
			h2, err := calculateConfigHash(c)
			Expect(err).NotTo(HaveOccurred())

			Expect(h2).To(Equal(h1))
		})

		It("returns the same hash independent of child ordering", func() {
			c1 := []metav1.Object{cm1, cm2, s1, s2}
			c2 := []metav1.Object{cm1, s2, cm2, s1}

			h1, err := calculateConfigHash(c1)
			Expect(err).NotTo(HaveOccurred())
			h2, err := calculateConfigHash(c2)
			Expect(err).NotTo(HaveOccurred())

			Expect(h2).To(Equal(h1))
		})
	})

	Context("updateHash", func() {
		var deployment *appsv1.Deployment

		BeforeEach(func() {
			deployment = utils.ExampleDeployment.DeepCopy()
		})

		It("sets the hash annotation to the provided value", func() {
			updateHash(deployment, "1234")

			podAnnotations := deployment.Spec.Template.GetAnnotations()
			Expect(podAnnotations).NotTo(BeNil())

			hash, ok := podAnnotations[configHashAnnotation]
			Expect(ok).To(BeTrue())
			Expect(hash).To(Equal("1234"))
		})
	})
})
