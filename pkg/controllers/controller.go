package controllers

import (
	"github.com/hliangzhao/balancer/pkg/controllers/balancer"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// All controllers will be added in this file.

var AddToManagerFuncs []func(manager.Manager) error

// AddToManager is the factory to add controllers to the controller-manager.
func AddToManager(manager manager.Manager) error {
	for _, f := range AddToManagerFuncs {
		if err := f(manager); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	// actually we have only one controller to add
	AddToManagerFuncs = append(AddToManagerFuncs, balancer.Add)
}
