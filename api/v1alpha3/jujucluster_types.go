/*
Copyright 2022.

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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JujuClusterSpec defines the desired state of JujuCluster
type JujuClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ControllerName string `json:"controllerName,omitempty"`
	ModelName      string `json:"modelName,omitempty"`
	Overlay        string `json:"overlay,omitempty"`
}

// JujuClusterStatus defines the observed state of JujuCluster
type JujuClusterStatus struct {
	State string `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=`.status.state`
//+kubebuilder:printcolumn:name="Controller",type="string",JSONPath=`.spec.controllerName`
//+kubebuilder:printcolumn:name="Model",type="string",JSONPath=`.spec.modelName`
// JujuCluster is the Schema for the jujuclusters API
type JujuCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JujuClusterSpec   `json:"spec,omitempty"`
	Status JujuClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JujuClusterList contains a list of JujuCluster
type JujuClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JujuCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JujuCluster{}, &JujuClusterList{})
}
