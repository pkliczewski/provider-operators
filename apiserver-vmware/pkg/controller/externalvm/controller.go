
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


package externalvm

import (
	"log"

	"github.com/kubernetes-incubator/apiserver-builder-alpha/pkg/builders"

	"github.com/pkliczewski/provider-operators/apiserver-vmware/pkg/apis/vms/v1alpha1"
	"github.com/pkliczewski/provider-operators/apiserver-vmware/pkg/controller/sharedinformers"
	listers "github.com/pkliczewski/provider-operators/apiserver-vmware/pkg/client/listers_generated/vms/v1alpha1"
)

// +controller:group=vms,version=v1alpha1,kind=ExternalVm,resource=externalvms
type ExternalVmControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about ExternalVm
	lister listers.ExternalVmLister
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *ExternalVmControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// Use the lister for indexing externalvms labels
	c.lister = arguments.GetSharedInformers().Factory.Vms().V1alpha1().ExternalVms().Lister()
}

// Reconcile handles enqueued messages
func (c *ExternalVmControllerImpl) Reconcile(u *v1alpha1.ExternalVm) error {
	// Implement controller logic here
	log.Printf("Running reconcile ExternalVm for %s\n", u.Name)
	return nil
}

func (c *ExternalVmControllerImpl) Get(namespace, name string) (*v1alpha1.ExternalVm, error) {
	return c.lister.ExternalVms(namespace).Get(name)
}
