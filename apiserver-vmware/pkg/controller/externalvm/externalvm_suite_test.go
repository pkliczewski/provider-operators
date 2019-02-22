
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


package externalvm_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"github.com/kubernetes-incubator/apiserver-builder-alpha/pkg/test"

	"github.com/pkliczewski/provider-operators/apiserver-vmware/pkg/apis"
	"github.com/pkliczewski/provider-operators/apiserver-vmware/pkg/client/clientset_generated/clientset"
	"github.com/pkliczewski/provider-operators/apiserver-vmware/pkg/openapi"
	"github.com/pkliczewski/provider-operators/apiserver-vmware/pkg/controller/sharedinformers"
	"github.com/pkliczewski/provider-operators/apiserver-vmware/pkg/controller/externalvm"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *clientset.Clientset
var shutdown chan struct{}
var controller *externalvm.ExternalVmController
var si *sharedinformers.SharedInformers

func TestExternalVm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "ExternalVm Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	testenv = test.NewTestEnvironment()
	config = testenv.Start(apis.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions)
	cs = clientset.NewForConfigOrDie(config)

	shutdown = make(chan struct{})
	si = sharedinformers.NewSharedInformers(config, shutdown)
	controller = externalvm.NewExternalVmController(config, si)
	controller.Run(shutdown)
})

var _ = AfterSuite(func() {
	close(shutdown)
	testenv.Stop()
})
