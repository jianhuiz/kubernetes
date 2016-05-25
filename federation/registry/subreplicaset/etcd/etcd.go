/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

// If you make changes to this file, you should also make the corresponding change in ReplicationController.

package etcd

import (
	"k8s.io/kubernetes/federation/apis/federation"
	"k8s.io/kubernetes/federation/registry/subreplicaset"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/rest"
	"k8s.io/kubernetes/pkg/registry/generic"
	"k8s.io/kubernetes/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/runtime"
)

type REST struct {
	*registry.Store
}

// NewREST returns a RESTStorage object that will work against SubReplicaSet.
func NewREST(opts generic.RESTOptions) (*REST, *StatusREST) {
	prefix := "/subreplicasets"

	newListFunc := func() runtime.Object { return &federation.SubReplicaSetList{} }
	storageInterface := opts.Decorator(
		opts.Storage, 100, &federation.SubReplicaSet{}, prefix, subreplicaset.Strategy, newListFunc)

	store := &registry.Store{
		NewFunc:     func() runtime.Object { return &federation.SubReplicaSet{} },
		NewListFunc: newListFunc,
		KeyRootFunc: func(ctx api.Context) string {
			return registry.NamespaceKeyRootFunc(ctx, prefix)
		},
		KeyFunc: func(ctx api.Context, name string) (string, error) {
			return registry.NamespaceKeyFunc(ctx, prefix, name)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*federation.SubReplicaSet).Name, nil
		},
		PredicateFunc:           subreplicaset.MatchSubReplicaSet,
		QualifiedResource:       federation.Resource("subreplicasets"),
		DeleteCollectionWorkers: opts.DeleteCollectionWorkers,

		CreateStrategy: subreplicaset.Strategy,
		UpdateStrategy: subreplicaset.Strategy,
		DeleteStrategy: subreplicaset.Strategy,

		Storage: storageInterface,
	}

	statusStore := *store
	statusStore.UpdateStrategy = subreplicaset.StatusStrategy

	return &REST{store}, &StatusREST{store: &statusStore}
}

// StatusREST implements the REST endpoint for changing the status of a SubReplicaSet
type StatusREST struct {
	store *registry.Store
}

func (r *StatusREST) New() runtime.Object {
	return &federation.SubReplicaSet{}
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx api.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return r.store.Update(ctx, name, objInfo)
}
