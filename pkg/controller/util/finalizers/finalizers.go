package finalizers

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/util/sets"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func HasFinalizer(obj runtimeclient.Object, finalizer string) (bool, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return false, err
	}
	finalizers := sets.NewString(accessor.GetFinalizers()...)
	return finalizers.Has(finalizer), nil
}

func AddFinalizers(obj runtimeclient.Object, newFinalizers sets.String) (bool, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return false, err
	}
	oldFinalizers := sets.NewString(accessor.GetFinalizers()...)
	if oldFinalizers.IsSuperset(newFinalizers) {
		return false, nil
	}
	allFinalizers := oldFinalizers.Union(newFinalizers)
	accessor.SetFinalizers(allFinalizers.List())
	return true, nil
}

func RemoveFinalizers(obj runtimeclient.Object, finalizers sets.String) (bool, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return false, err
	}
	oldFinalizers := sets.NewString(accessor.GetFinalizers()...)
	if oldFinalizers.Intersection(finalizers).Len() == 0 {
		return false, nil
	}
	newFinalizers := oldFinalizers.Difference(finalizers)
	accessor.SetFinalizers(newFinalizers.List())
	return true, nil
}
