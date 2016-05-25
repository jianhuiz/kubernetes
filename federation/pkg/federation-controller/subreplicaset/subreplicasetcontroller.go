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

package subreplicaset

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/kubernetes/federation/apis/federation"
	federationcache "k8s.io/kubernetes/federation/client/cache"
	federationclientset "k8s.io/kubernetes/federation/client/clientset_generated/federation_internalclientset"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/cache"
	k8sclientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/conversion"
	"k8s.io/kubernetes/pkg/runtime"
	utilruntime "k8s.io/kubernetes/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/util/sets"
	"k8s.io/kubernetes/pkg/util/wait"
	"k8s.io/kubernetes/pkg/util/workqueue"
	"k8s.io/kubernetes/pkg/watch"
	"strings"
	"time"
)

const (
	AnnotationKeyOfTargetCluster        = "kubernetes.io/target-cluster"
	AnnotationKeyOfFederationReplicaSet = "kubernetes.io/created-by"
	UserAgentName                       = "Federation-replicaset-Controller"
	KubeAPIQPS                          = 20.0
	KubeAPIBurst                        = 30
)

type SubReplicaSetController struct {
	knownClusterSet sets.String

	//federationClient used to operate cluster, rs and subrs
	federationClient federationclientset.Interface

	// To allow injection of syncSubRC for testing.
	syncHandler func(subRcKey string) error

	// clusterMonitorPeriod is the period for updating status of cluster
	clusterMonitorPeriod time.Duration

	// clusterKubeClientMap is a mapping of clusterName and restclient
	clusterKubeClientMap map[string]*k8sclientset.Clientset

	// subRc framework and store
	subReplicaSetController *framework.Controller
	subReplicaSetStore      federationcache.StoreToSubReplicaSetLister

	// cluster framework and store
	clusterController *framework.Controller
	clusterStore      federationcache.StoreToClusterLister

	// UberRc framework and store
	replicaSetController *framework.Controller
	replicaSetStore      cache.StoreToReplicaSetLister

	// subRC that have been queued up for processing by workers
	queue *workqueue.Type
}

// NewclusterController returns a new cluster controller
func NewSubReplicaSetController(federationClient federationclientset.Interface, clusterMonitorPeriod time.Duration) *SubReplicaSetController {
	cc := &SubReplicaSetController{
		knownClusterSet:      make(sets.String),
		federationClient:     federationClient,
		clusterMonitorPeriod: clusterMonitorPeriod,
		clusterKubeClientMap: make(map[string]*k8sclientset.Clientset),
		queue:                workqueue.New(),
	}

	cc.subReplicaSetStore.Store, cc.subReplicaSetController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options api.ListOptions) (runtime.Object, error) {
				return cc.federationClient.Federation().SubReplicaSets(api.NamespaceAll).List(options)
			},
			WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
				return cc.federationClient.Federation().SubReplicaSets(api.NamespaceAll).Watch(options)
			},
		},
		&federation.SubReplicaSet{},
		controller.NoResyncPeriodFunc(),
		framework.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				subRc := obj.(*federation.SubReplicaSet)
				cc.enqueueSubRc(subRc)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				subRc := newObj.(*federation.SubReplicaSet)
				cc.enqueueSubRc(subRc)
			},
		},
	)

	cc.clusterStore.Store, cc.clusterController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options api.ListOptions) (runtime.Object, error) {
				return cc.federationClient.Federation().Clusters().List(options)
			},
			WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
				return cc.federationClient.Federation().Clusters().Watch(options)
			},
		},
		&federation.Cluster{},
		controller.NoResyncPeriodFunc(),
		framework.ResourceEventHandlerFuncs{
			DeleteFunc: cc.delFromClusterSet,
			AddFunc:    cc.addToClusterSet,
		},
	)

	cc.replicaSetStore.Store, cc.replicaSetController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options api.ListOptions) (runtime.Object, error) {
				return cc.federationClient.Extensions().ReplicaSets(api.NamespaceAll).List(options)
			},
			WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
				return cc.federationClient.Extensions().ReplicaSets(api.NamespaceAll).Watch(options)
			},
		},
		&extensions.ReplicaSet{},
		controller.NoResyncPeriodFunc(),
		framework.ResourceEventHandlerFuncs{
			DeleteFunc: cc.deleteSubRs,
		},
	)
	cc.syncHandler = cc.syncSubReplicaSet
	return cc
}

// delFromClusterSet delete a cluster from clusterSet and
// delete the corresponding restclient from the map clusterKubeClientMap
func (cc *SubReplicaSetController) delFromClusterSet(obj interface{}) {
	cluster := obj.(*federation.Cluster)
	cc.knownClusterSet.Delete(cluster.Name)
	delete(cc.clusterKubeClientMap, cluster.Name)
}

// addToClusterSet insert the new cluster to clusterSet and create a corresponding
// restclient to map clusterKubeClientMap
func (cc *SubReplicaSetController) addToClusterSet(obj interface{}) {
	cluster := obj.(*federation.Cluster)
	cc.knownClusterSet.Insert(cluster.Name)
	// create the restclient of cluster
	restClient, err := newClusterClientset(cluster)
	if err != nil || restClient == nil {
		glog.Errorf("Failed to create corresponding restclient of kubernetes cluster: %v", err)
		return
	}
	cc.clusterKubeClientMap[cluster.Name] = restClient
}

// Run begins watching and syncing.
func (cc *SubReplicaSetController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	go cc.clusterController.Run(wait.NeverStop)
	go cc.replicaSetController.Run(stopCh)
	go cc.subReplicaSetController.Run(stopCh)

	for i := 0; i < workers; i++ {
		go wait.Until(cc.worker, time.Second, stopCh)
	}
	<-stopCh
	glog.Infof("Shutting down ClusterController")
	cc.queue.ShutDown()

}

// enqueueSubRc adds an object to the controller work queue
// obj could be an *federation_v1alpha1.SubReplicaSet, or a DeletionFinalStateUnknown item.
func (cc *SubReplicaSetController) enqueueSubRc(obj interface{}) {
	key, err := controller.KeyFunc(obj)
	if err != nil {
		glog.Errorf("Couldn't get key for object %+v: %v", obj, err)
		return
	}
	cc.queue.Add(key)
}

func (cc *SubReplicaSetController) worker() {
	for {
		func() {
			key, quit := cc.queue.Get()
			if quit {
				return
			}
			defer cc.queue.Done(key)
			err := cc.syncHandler(key.(string))
			if err != nil {
				glog.Errorf("Error syncing cluster controller: %v", err)
			}
		}()
	}
}

// syncSubReplicaSet will sync the subrc with the given key,
func (cc *SubReplicaSetController) syncSubReplicaSet(key string) error {
	startTime := time.Now()
	defer func() {
		glog.V(4).Infof("Finished syncing controller %q (%v)", key, time.Now().Sub(startTime))
	}()
	obj, exists, err := cc.subReplicaSetStore.Store.GetByKey(key)
	if !exists {
		glog.Infof("sub replicaset: %v has been deleted", key)
		return nil
	}
	if err != nil {
		glog.Infof("Unable to retrieve sub replicaset %v from store: %v", key, err)
		cc.queue.Add(key)
		return err
	}
	subRs := obj.(*federation.SubReplicaSet)
	err = cc.manageSubReplicaSet(subRs)
	if err != nil {
		glog.Infof("Unable to manage subRs in kubernetes cluster: %v", key, err)
		//TODO: if manageSubReplicaSet fails multi times, the subrs need to be rescheduled
		//TODO: Now we just drop this subreplicaset
		//TODO: cc.queue.Add(key)
		return err
	}
	return nil
}

// getBindingClusterOfSubRS get the target cluster(scheduled by federation scheduler) of subRS
// return the targetCluster name
func (cc *SubReplicaSetController) getBindingClusterOfSubRS(subRs *federation.SubReplicaSet) (*federation.Cluster, error) {
	accessor, err := meta.Accessor(subRs)
	if err != nil {
		return nil, err
	}
	annotations := accessor.GetAnnotations()
	if annotations == nil {
		return nil, fmt.Errorf("Failed to get target cluster from the annotation of subreplicaset")
	}
	targetClusterName, found := annotations[AnnotationKeyOfTargetCluster]
	if !found {
		return nil, fmt.Errorf("Failed to get target cluster from the annotation of subreplicaset")
	}
	return cc.federationClient.Federation().Clusters().Get(targetClusterName)
}

// getFederateRsCreateBy get the federation ReplicaSet created by of subRS
// return the replica set name
func (cc *SubReplicaSetController) getFederateRsCreateBy(subRs *federation.SubReplicaSet) (string, error) {
	accessor, err := meta.Accessor(subRs)
	if err != nil {
		return "", err
	}
	annotations := accessor.GetAnnotations()
	if annotations == nil {
		return "", fmt.Errorf("Failed to get Federate Rs Create By from the annotation of subreplicaset")
	}
	rsCreateBy, found := annotations[AnnotationKeyOfFederationReplicaSet]
	if !found {
		return "", fmt.Errorf("Failed to get Federate Rs Create By from the annotation of subreplicaset")
	}
	return rsCreateBy, nil
}

func covertSubRSToRS(subRS *federation.SubReplicaSet) (*extensions.ReplicaSet, error) {
	clone, err := conversion.NewCloner().DeepCopy(subRS)
	if err != nil {
		return nil, err
	}
	subrs, ok := clone.(*federation.SubReplicaSet)
	if !ok {
		return nil, fmt.Errorf("Unexpected subreplicaset cast error : %v\n", subrs)
	}
	result := &extensions.ReplicaSet{}
	result.Kind = "ReplicaSet"
	result.APIVersion = "extensions/v1beta1"
	result.ObjectMeta = subrs.ObjectMeta
	result.Spec = subrs.Spec
	result.Status = subrs.Status
	result.Name = subrs.Name
	return result, nil
}

// manageSubReplicaSet will sync the sub replicaset with the given key,and then create
// or update replicaset to kubernetes cluster
func (cc *SubReplicaSetController) manageSubReplicaSet(subRs *federation.SubReplicaSet) error {

	targetCluster, err := cc.getBindingClusterOfSubRS(subRs)
	if targetCluster == nil || err != nil {
		glog.Infof("Failed to get target cluster of SubRS: %v", err)
		return fmt.Errorf("Failed to get target cluster of SubRS: %v", err)
	}

	clusterClient, found := cc.clusterKubeClientMap[targetCluster.Name]
	if !found {
		glog.Infof("It's a new cluster, a cluster client will be created")
		clusterClient, err := newClusterClientset(targetCluster)
		if err != nil || clusterClient == nil {
			return err
		}
		cc.clusterKubeClientMap[targetCluster.Name] = clusterClient
	}

	rs, err := covertSubRSToRS(subRs)
	if err != nil {
		glog.Infof("Failed to convert subrs to rs: %v", err)
		return err
	}

	// check the sub replicaset already exists in kubernetes cluster or not
	replicaSet, err := clusterClient.Extensions().ReplicaSets(subRs.Namespace).Get(subRs.Name)
	// if not exist, means that this sub replicaset need to be created
	if replicaSet == nil || err != nil {
		// create the sub replicaset to kubernetes cluster
		rs.ResourceVersion = ""
		replicaSet, err := clusterClient.Extensions().ReplicaSets(rs.Namespace).Create(rs)
		if err != nil || replicaSet == nil {
			glog.Infof("Failed to create sub replicaset in kubernetes cluster: %v", err)
			return fmt.Errorf("Failed to create sub replicaset in kubernetes cluster: %v", err)
		}
	} else {
		// if exists, then update it
		replicaSet, err = clusterClient.Extensions().ReplicaSets(rs.Namespace).Update(rs)
		if err != nil || replicaSet == nil {
			glog.Infof("Failed to update sub replicaset in kubernetes cluster: %v", err)
			return fmt.Errorf("Failed to update sub replicaset in kubernetes cluster: %v", err)
		}
	}
	return nil
}

func (cc *SubReplicaSetController) deleteSubRs(cur interface{}) {
	rs := cur.(*extensions.ReplicaSet)
	// get the corresponing subrs from the cache
	subRSList, err := cc.federationClient.Federation().SubReplicaSets(api.NamespaceAll).List(api.ListOptions{})
	if err != nil || len(subRSList.Items) == 0 {
		glog.Infof("Couldn't get subRS to delete : %+v", cur)
		return
	}

	// get the related subRS created by the replicaset
	for _, subRs := range subRSList.Items {
		name, err := cc.getFederateRsCreateBy(&subRs)
		if err != nil || !strings.EqualFold(rs.Name, name) {
			continue
		}
		targetCluster, err := cc.getBindingClusterOfSubRS(&subRs)
		if targetCluster == nil || err != nil {
			continue
		}
		clusterClient, found := cc.clusterKubeClientMap[targetCluster.Name]
		if !found {
			continue
		}
		rs, err = covertSubRSToRS(&subRs)
		if err != nil {
			continue
		}
		err = clusterClient.Extensions().ReplicaSets(rs.Namespace).Delete(rs.Name, &api.DeleteOptions{})
		if err != nil {
			return
		}
	}
	return
}

func newClusterClientset(c *federation.Cluster) (*k8sclientset.Clientset, error) {
	clusterConfig, err := clientcmd.BuildConfigFromFlags(c.Spec.ServerAddressByClientCIDRs[0].ServerAddress, "")
	if err != nil {
		return nil, err
	}
	clusterConfig.QPS = KubeAPIQPS
	clusterConfig.Burst = KubeAPIBurst
	clientset := k8sclientset.NewForConfigOrDie(restclient.AddUserAgent(clusterConfig, UserAgentName))
	return clientset, nil
}
