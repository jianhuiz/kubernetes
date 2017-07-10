package proxy

import (
	"fmt"
	"github.com/golang/glog"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/rest"
	fedv1 "k8s.io/kubernetes/federation/apis/federation/v1beta1"
	fedclient "k8s.io/kubernetes/federation/client/clientset_generated/federation_clientset"
	fedutil "k8s.io/kubernetes/federation/pkg/federation-controller/util"
	"k8s.io/kubernetes/pkg/api"
	kubeclient "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

const (
	userAgentName = "federation-apiserver"
)

type ClusterClient struct {
	cluster   *fedv1.Cluster
	clientset *kubeclient.Clientset
	err       error
}

type Store struct {
	Copier            runtime.ObjectCopier
	NewFunc           func() runtime.Object
	NewListFunc       func() runtime.Object
	IterateListFunc   func(objs runtime.Object, cb func(obj runtime.Object))
	MergeListFunc     func(dest, src runtime.Object)
	RESTClientFunc    func(kubeClientset kubeclient.Interface) rest.Interface
	QualifiedResource schema.GroupResource

	FedClient fedclient.Interface
}

func (s *Store) New() runtime.Object {
	return s.NewFunc()
}

func (s *Store) NewList() runtime.Object {
	return s.NewListFunc()
}

func (s *Store) List(ctx genericapirequest.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	clusterSelector := labels.Everything()
	if options != nil && options.ClusterSelector != nil {
		clusterSelector = options.ClusterSelector
	}
	clusterClients, err := s.listClusters(clusterSelector)
	if err != nil {
		return nil, err
	}

	results := s.NewList()
	errs := []error{}
	for _, clusterClient := range clusterClients {
		if clusterClient.err != nil {
			errs = append(errs, clusterClient.err)
		} else {
			optsv1 := metav1.ListOptions{}
			api.Scheme.Convert(options, &optsv1, nil)
			objs := s.NewList()
			err := s.RESTClientFunc(clusterClient.clientset).Get().
				Namespace(genericapirequest.NamespaceValue(ctx)).
				Resource(s.QualifiedResource.Resource).
				VersionedParams(&optsv1, scheme.ParameterCodec).
				Do().Into(objs)
			if err != nil {
				errs = append(errs, err)
			} else {
				s.IterateListFunc(objs, func(obj runtime.Object) {
					accessor, _ := meta.Accessor(obj)
					accessor.SetClusterName(clusterClient.cluster.Name)
				})
				s.MergeListFunc(results, objs)
			}
		}
	}

	if len(errs) != 0 {
		return nil, utilerrors.NewAggregate(errs)
	}

	return results, nil
}

func (s *Store) Get(ctx genericapirequest.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	clusterClient, err := s.getClusterClientByName(options.ClusterName)
	if err != nil {
		glog.Warningf("error creating cluster client for cluster %s, err: %v", options.ClusterName, err)
		return nil, err
	}

	obj := s.New()
	err = s.RESTClientFunc(clusterClient).Get().
		Namespace(genericapirequest.NamespaceValue(ctx)).Name(name).
		Resource(s.QualifiedResource.Resource).
		VersionedParams(options, scheme.ParameterCodec).
		Do().
		Into(obj)
	if err != nil {
		return nil, err
	}

	accessor, _ := meta.Accessor(obj)
	accessor.SetClusterName(options.ClusterName)
	return obj, nil
}

func (s *Store) listClusters(clusterSelector labels.Selector) ([]ClusterClient, error) {
	clusterClients := []ClusterClient{}
	clusters, err := s.FedClient.Federation().Clusters().List(metav1.ListOptions{LabelSelector: clusterSelector.String()})
	if err != nil {
		glog.Warningf("error listing clusters, err: %v", err)
		return nil, err
	}
	for i := range clusters.Items {
		clientset, err := s.getClusterClient(&clusters.Items[i])
		clusterClients = append(clusterClients, ClusterClient{&clusters.Items[i], clientset, err})
	}

	return clusterClients, nil
}

func (s *Store) getClusterClientByName(clusterName string) (*kubeclient.Clientset, error) {
	cluster, err := s.FedClient.Federation().Clusters().Get(clusterName, metav1.GetOptions{})
	if err != nil {
		glog.Warningf("error getting cluster %s, err: %v", clusterName, err)
		return nil, err
	}
	return s.getClusterClient(cluster)
}

func (s *Store) getClusterClient(cluster *fedv1.Cluster) (*kubeclient.Clientset, error) {
	if !isClusterReady(cluster) {
		glog.Warningf("cluster %s is not ready, conditions: %v", cluster.Name, cluster.Status.Conditions)
		return nil, fmt.Errorf("cluster %s is not ready", cluster.Name)
	}
	clusterConfig, err := fedutil.BuildClusterConfig(cluster)
	if err != nil {
		glog.Warningf("error build cluster config for cluster %s, err: %v", cluster.Name, err)
		return nil, err
	}
	return kubeclient.NewForConfig(rest.AddUserAgent(clusterConfig, userAgentName))
}

func isClusterReady(cluster *fedv1.Cluster) bool {
	for _, condition := range cluster.Status.Conditions {
		if condition.Type == fedv1.ClusterReady {
			if condition.Status == apiv1.ConditionTrue {
				return true
			}
		}
	}
	return false
}
