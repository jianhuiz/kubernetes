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

package cache

import (
	"github.com/golang/glog"
	"k8s.io/kubernetes/federation/apis/federation"
	"k8s.io/kubernetes/pkg/api"
	kubeCache "k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/labels"
)

// StoreToClusterLister makes a Store have the List method of the unversioned.ClusterInterface
// The Store must contain (only) clusters.
type StoreToClusterLister struct {
	kubeCache.Store
}

func (s *StoreToClusterLister) List() (clusters federation.ClusterList, err error) {
	for _, m := range s.Store.List() {
		clusters.Items = append(clusters.Items, *(m.(*federation.Cluster)))
	}
	return clusters, nil
}

// ClusterConditionPredicate is a function that indicates whether the given cluster's conditions meet
// some set of criteria defined by the function.
type ClusterConditionPredicate func(cluster federation.Cluster) bool

// storeToClusterConditionLister filters and returns nodes matching the given type and status from the store.
type storeToClusterConditionLister struct {
	store     kubeCache.Store
	predicate ClusterConditionPredicate
}

// ClusterCondition returns a storeToClusterConditionLister
func (s *StoreToClusterLister) ClusterCondition(predicate ClusterConditionPredicate) storeToClusterConditionLister {
	return storeToClusterConditionLister{s.Store, predicate}
}

// List returns a list of clusters that match the conditions defined by the predicate functions in the storeToClusterConditionLister.
func (s storeToClusterConditionLister) List() (clusters federation.ClusterList, err error) {
	for _, m := range s.store.List() {
		cluster := *m.(*federation.Cluster)
		if s.predicate(cluster) {
			clusters.Items = append(clusters.Items, cluster)
		} else {
			glog.V(5).Infof("Cluster %s matches none of the conditions", cluster.Name)
		}
	}
	return
}

// StoreToReplicationControllerLister gives a store List and Exists methods. The store must contain only ReplicationControllers.
type StoreToSubReplicaSetLister struct {
	kubeCache.Store
}

// StoreToSubReplicaSetLister lists all replicaSets in the store.
func (s *StoreToSubReplicaSetLister) List() (replicaSets []federation.SubReplicaSet, err error) {
	for _, c := range s.Store.List() {
		replicaSets = append(replicaSets, *(c.(*federation.SubReplicaSet)))
	}
	return replicaSets, nil
}

// Exists checks if the given SubRS exists in the store.
func (s *StoreToSubReplicaSetLister) Exists(rs *federation.SubReplicaSet) (bool, error) {
	_, exists, err := s.Store.Get(rs)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type storeSubRSsNamespacer struct {
	store     kubeCache.Store
	namespace string
}

func (s storeSubRSsNamespacer) List(selector labels.Selector) (rss []federation.SubReplicaSet, err error) {
	for _, c := range s.store.List() {
		rs := *(c.(*federation.SubReplicaSet))
		if s.namespace == api.NamespaceAll || s.namespace == rs.Namespace {
			if selector.Matches(labels.Set(rs.Labels)) {
				rss = append(rss, rs)
			}
		}
	}
	return
}

func (s *StoreToSubReplicaSetLister) SubRSs(namespace string) storeSubRSsNamespacer {
	return storeSubRSsNamespacer{s.Store, namespace}
}
