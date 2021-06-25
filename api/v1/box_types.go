/*
Copyright 2021.

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

package v1

import (
	"github.com/haisum/k99s/backends"
	"github.com/haisum/k99s/runtimes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BoxSpec defines the desired state of Box
type BoxSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Runtime runtimes.RuntimeType `json:"runtime"`
	Backend backends.BackendType `json:"backend"`
	GitURL  string               `json:"gitURL"`

	// +optional
	GitSubPath string `json:"gitSubPath"`

	// Executed on fresh database at creation time
	// +optional
	BootstrapSQL string `json:"bootstrapSQL,omitempty"`
}

// BoxStatus defines the observed state of Box
type BoxStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Status    string      `json:"status,omitempty"`
	StartedAt metav1.Time `json:"startedAt,omitempty"`
	URL       string      `json:"url,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// Box is the Schema for the boxes API
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="URL",type=string,JSONPath=`.status.url`
// +kubebuilder:printcolumn:name="GitURL",type=string,JSONPath=`.spec.gitURL`
// +kubebuilder:printcolumn:name="Runtime",type=string,JSONPath=`.spec.runtime`
// +kubebuilder:printcolumn:name="Backend",type=string,JSONPath=`.spec.backend`
type Box struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BoxSpec   `json:"spec,omitempty"`
	Status BoxStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BoxList contains a list of Box
type BoxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Box `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Box{}, &BoxList{})
}
