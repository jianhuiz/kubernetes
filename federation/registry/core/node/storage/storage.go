package storage

import(
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/kubernetes/pkg/api"
	regproxy "k8s.io/kubernetes/federation/registry/proxy"
	fedclient "k8s.io/kubernetes/federation/client/clientset_generated/federation_clientset"
	kubeclient "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)



func RESTClientFunc(kubeClientset kubeclient.Interface) rest.Interface {
	return kubeClientset.Core().RESTClient()
}

func IterateListFunc(objs runtime.Object, cb func(obj runtime.Object)) {
	nodList := objs.(*api.NodeList)
	for i, _ := range nodList.Items {
		cb(&nodList.Items[i])
	}
}

func MergeListFunc(dest, src runtime.Object) {
	destList := dest.(*api.NodeList)
	srcList := src.(*api.NodeList)
	destList.Items = append(destList.Items, srcList.Items...)
}

type REST struct {
	*regproxy.Store
}

func NewREST(optsGetter generic.RESTOptionsGetter, fedClient fedclient.Interface) (*REST, *StatusREST) {
	store := &regproxy.Store{
		Copier:    api.Scheme,
		NewFunc:   func() runtime.Object {return &api.Node{}},
		NewListFunc: func() runtime.Object {return &api.NodeList{}},
		IterateListFunc: IterateListFunc,
		MergeListFunc: MergeListFunc,
		RESTClientFunc: RESTClientFunc,
		QualifiedResource: api.Resource("nodes"),
		FedClient: fedClient,
	}

	statusStore := *store
	return &REST{store}, &StatusREST{store: &statusStore}
}

type StatusREST struct {
	store *regproxy.Store
}

func (r *StatusREST) New() runtime.Object {
	return r.store.New()
}

