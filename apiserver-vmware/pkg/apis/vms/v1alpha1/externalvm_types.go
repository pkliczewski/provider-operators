
/*
Copyright 2019 The Kubernetes Authors.

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
	"log"
	"context"

	"k8s.io/apimachinery/pkg/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/pkliczewski/provider-operators/apiserver-vmware/pkg/apis/vms"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExternalVm
// +k8s:openapi-gen=true
// +resource:path=externalvms,strategy=ExternalVmStrategy
type ExternalVm struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalVmSpec   `json:"spec,omitempty"`
	Status ExternalVmStatus `json:"status,omitempty"`
}

// ExternalVmSpec defines the desired state of ExternalVm
type ExternalVmSpec struct {
}

// ExternalVmStatus defines the observed state of ExternalVm
type ExternalVmStatus struct {
}

// Validate checks that an instance of ExternalVm is well formed
func (ExternalVmStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*vms.ExternalVm)
	log.Printf("Validating fields for ExternalVm %s\n", o.Name)
	errors := field.ErrorList{}
	// perform validation here and add to errors using field.Invalid
	return errors
}

// DefaultingFunction sets default ExternalVm field values
func (ExternalVmSchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*ExternalVm)
	// set default field values here
	log.Printf("Defaulting fields for ExternalVm %s\n", obj.Name)
}
