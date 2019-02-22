package controller

import (
	"github.com/pkliczewski/provider-operators/provider-operator/pkg/controller/vmprovider"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, vmprovider.Add)
}
