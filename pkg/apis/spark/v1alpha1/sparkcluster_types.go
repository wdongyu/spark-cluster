/*
Copyright 2019 wdongyu.

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ListSpec defines the list of memory/cpu...
type ListSpec struct {
	Memory string `json:"memory,omitempty"`
	CPU    string `json:"cpu,omitempty"`
}

// ResourcesSpec defines the Resource request/limit
type ResourcesSpec struct {
	Limits   ListSpec `json:"limits,omitempty"`
	Requests ListSpec `json:"requests,omitempty"`
}

// SparkClusterSpec defines the desired state of SparkCluster
type SparkClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// PodName  string `json:"podName"`
	// Replicas *int32 `json:"replicas"`
	SlaveNum  int                         `json:"slaveNum"`
	Ports     []corev1.ServicePort        `json:"ports,omitempty"`
	PvcEnable bool                        `json:"pvcEnable,omitempty"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	NFS       corev1.NFSVolumeSource      `json:"nfs,omitempty"`
}

// SparkClusterStatus defines the observed state of SparkCluster
type SparkClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	AvailableReplicas int32 `json:"availableReplicas"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SparkCluster is the Schema for the sparkclusters API
// +k8s:openapi-gen=true
type SparkCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SparkClusterSpec   `json:"spec,omitempty"`
	Status SparkClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SparkClusterList contains a list of SparkCluster
type SparkClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SparkCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SparkCluster{}, &SparkClusterList{})
}
