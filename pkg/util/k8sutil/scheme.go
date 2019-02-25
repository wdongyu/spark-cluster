package k8sutil

import "k8s.io/apimachinery/pkg/runtime"

func NewScheme(AddToSchemeFuncs []func(scheme *runtime.Scheme)) *runtime.Scheme {
	scheme := runtime.NewScheme()

	for _, addToSchemeFunc := range AddToSchemeFuncs {
		addToSchemeFunc(scheme)
	}

	return scheme
}
