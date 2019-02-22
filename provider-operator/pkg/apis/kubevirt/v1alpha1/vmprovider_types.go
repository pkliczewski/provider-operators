package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VmProviderSpec defines the desired state of VmProvider
// +k8s:openapi-gen=true
type VmProviderSpec struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// VmProviderStatus defines the observed state of VmProvider
// +k8s:openapi-gen=true
type VmProviderStatus struct {
	Validated bool   `json:"validated"`
	Message   string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VmProvider is the Schema for the vmproviders API
// +k8s:openapi-gen=true
type VmProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VmProviderSpec   `json:"spec,omitempty"`
	Status VmProviderStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VmProviderList contains a list of VmProvider
type VmProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VmProvider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VmProvider{}, &VmProviderList{})
}
