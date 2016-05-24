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

package service

import (
	"fmt"
	"time"

	federationclientset "k8s.io/kubernetes/federation/client/clientset_generated/federation_internalclientset"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	cache "k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/controller"

	"github.com/golang/glog"
	"reflect"
	"sort"
)

// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (sc *ServiceController) clusterServiceWorker() {
	fedClient := sc.federationClient
	for clusterName, cache := range sc.clusterCache.clientMap {
		go func(cache *clusterCache, clusterName string) {
			for {
				key, quit := cache.serviceQueue.Get()
				// update service cache
				if quit {
					return
				}
				defer cache.serviceQueue.Done(key)
				err := sc.clusterCache.syncService(key.(string), clusterName, cache, sc.serviceCache, fedClient)
				if err != nil {
					glog.Errorf("failed to sync service: %+v", err)
				}
			}
		}(cache, clusterName)
	}
}

// Whenever there is change on service, the federation service should be updated
func (cc *clusterClientCache) syncService(key, clusterName string, clusterCache *clusterCache, serviceCache *serviceCache, fedClient federationclientset.Interface) error {
	// obj holds the latest service info from apiserver, return if there is no federation cache for the service
	cachedService, ok := serviceCache.get(key)
	if !ok {
		// here we filtered all non-federation services
		return nil
	}
	serviceInterface, exists, err := clusterCache.serviceStore.GetByKey(key)
	if err != nil {
		glog.Infof("did not successfully get %v from store: %v, will retry later", key, err)
		clusterCache.serviceQueue.Add(key)
		return err
	}
	var needUpdate bool
	if exists {
		service, ok := serviceInterface.(*api.Service)
		if ok {
			glog.V(4).Infof("Found service for federation service %s/%s from cluster %s", service.Namespace, service.Name, clusterName)
			needUpdate = cc.processServiceUpdate(cachedService, service, clusterName)
		} else {
			_, ok := serviceInterface.(cache.DeletedFinalStateUnknown)
			if !ok {
				return fmt.Errorf("object contained wasn't a service or a deleted key: %+v", serviceInterface)
			}
			glog.Infof("Found tombstone for %v", key)
			needUpdate = cc.processServiceDeletion(cachedService, clusterName)
		}
	} else {
		// service absence in store means watcher caught the deletion
		glog.Infof("service has been deleted %v", key)
		needUpdate = cc.processServiceDeletion(cachedService, clusterName)
	}
	if needUpdate {
		err := cc.persistFedServiceUpdate(cachedService, fedClient)
		if err == nil {
			cachedService.appliedState = cachedService.lastState
			cachedService.resetFedUpdateDelay()
		} else {
			if err != nil {
				glog.Errorf("failed to sync service: %+v, put back to service queue", err)
				clusterCache.serviceQueue.Add(key)
			}
		}
	}
	return nil
}

// processServiceDeletion is triggered when a service is delete from underlying k8s cluster
// the deletion function will wip out the cached ingress info of the service from federation service ingress
// the function returns a bool to indicate if actual update happend on federation service cache
// and if the federation service cache is updated, the updated info should be post to federation apiserver
func (cc *clusterClientCache) processServiceDeletion(cachedService *cachedService, clusterName string) bool {
	cachedService.rwlock.Lock()
	defer cachedService.rwlock.Unlock()
	cachedStatus, ok := cachedService.serviceStatusMap[clusterName]
	// cached status found, remove ingress info from federation service cache
	if ok {
		cachedFedServiceStatus := cachedService.lastState.Status.LoadBalancer
		removeIndexes := []int{}
		for i, fed := range cachedFedServiceStatus.Ingress {
			for _, new := range cachedStatus.Ingress {
				// remove if same ingress record found
				if new.IP == fed.IP && new.Hostname == fed.Hostname {
					removeIndexes = append(removeIndexes, i)
				}
			}
		}
		sort.Ints(removeIndexes)
		for i := len(removeIndexes) - 1; i >= 0; i-- {
			cachedFedServiceStatus.Ingress = append(cachedFedServiceStatus.Ingress[:removeIndexes[i]], cachedFedServiceStatus.Ingress[removeIndexes[i]+1:]...)
			glog.V(4).Infof("remove old ingress %d for service %s/%s", removeIndexes[i], cachedService.lastState.Namespace, cachedService.lastState.Name)
		}
		delete(cachedService.serviceStatusMap, clusterName)
		delete(cachedService.endpointMap, clusterName)
		cachedService.lastState.Status.LoadBalancer = cachedFedServiceStatus
		return true
	} else {
		glog.V(4).Infof("service removal %s/%s from cluster %s observed.", cachedService.lastState.Namespace, cachedService.lastState.Name, clusterName)
	}
	return false
}

// processServiceUpdate Update ingress info when service updated
// the function returns a bool to indicate if actual update happend on federation service cache
// and if the federation service cache is updated, the updated info should be post to federation apiserver
func (cc *clusterClientCache) processServiceUpdate(cachedService *cachedService, service *api.Service, clusterName string) bool {
	glog.V(4).Infof("Processing service update for %s/%s, cluster %s", service.Namespace, service.Name, clusterName)
	cachedService.rwlock.Lock()
	defer cachedService.rwlock.Unlock()
	var needUpdate bool
	newServiceLB := service.Status.LoadBalancer
	cachedFedServiceStatus := cachedService.lastState.Status.LoadBalancer
	if len(newServiceLB.Ingress) == 0 {
		// not yet get LB IP
		return false
	}

	cachedStatus, ok := cachedService.serviceStatusMap[clusterName]
	if ok {
		if reflect.DeepEqual(cachedStatus, newServiceLB) {
			glog.V(4).Infof("same ingress info observed for service %s/%s: %+v ", service.Namespace, service.Name, cachedStatus.Ingress)
		} else {
			glog.V(4).Infof("ingress info was changed for service %s/%s: cache: %+v, new: %+v ",
				service.Namespace, service.Name, cachedStatus.Ingress, newServiceLB)
			needUpdate = true
		}
	} else {
		glog.V(4).Infof("cached service status was not found for %s/%s, cluster %s, building one", service.Namespace, service.Name, clusterName)

		// cache is not always reliable(cache will be cleaned when service controller restart)
		// two cases will run into this branch:
		// 1. new service loadbalancer info received -> no info in cache, and no in federation service
		// 2. service controller being restarted -> no info in cache, but it is in federation service

		// check if the lb info is already in federation service

		cachedService.serviceStatusMap[clusterName] = newServiceLB
		needUpdate = false
		// iterate service ingress info
		for _, new := range newServiceLB.Ingress {
			var found bool
			// if it is known by federation service
			for _, fed := range cachedFedServiceStatus.Ingress {
				if new.IP == fed.IP && new.Hostname == fed.Hostname {
					found = true
				}
			}
			if !found {
				needUpdate = true
				break
			}
		}
	}

	if needUpdate {
		// new status = cached federation status - cached status + new status from k8s cluster

		removeIndexes := []int{}
		for i, fed := range cachedFedServiceStatus.Ingress {
			for _, new := range cachedStatus.Ingress {
				// remove if same ingress record found
				if new.IP == fed.IP && new.Hostname == fed.Hostname {
					removeIndexes = append(removeIndexes, i)
				}
			}
		}
		sort.Ints(removeIndexes)
		for i := len(removeIndexes) - 1; i >= 0; i-- {
			cachedFedServiceStatus.Ingress = append(cachedFedServiceStatus.Ingress[:removeIndexes[i]], cachedFedServiceStatus.Ingress[removeIndexes[i]+1:]...)
		}
		cachedFedServiceStatus.Ingress = append(cachedFedServiceStatus.Ingress, service.Status.LoadBalancer.Ingress...)
		cachedService.lastState.Status.LoadBalancer = cachedFedServiceStatus
		glog.V(4).Infof("add new ingress info %+v for service %s/%s", service.Status.LoadBalancer, service.Namespace, service.Name)
	} else {
		glog.V(4).Infof("same ingress info found for %s/%s, cluster %s", service.Namespace, service.Name, clusterName)
	}
	return needUpdate
}

func (cc *clusterClientCache) persistFedServiceUpdate(cachedService *cachedService, fedClient federationclientset.Interface) error {
	service := cachedService.lastState
	glog.V(5).Infof("Persist federation service status %s/%s", service.Namespace, service.Name)
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err := fedClient.Core().Services(service.Namespace).UpdateStatus(service)
		if err == nil {
			glog.V(2).Infof("Successfully update service %s/%s to federation apiserver", service.Namespace, service.Name)
			return nil
		}
		if errors.IsNotFound(err) {
			glog.Infof("Not persisting update to service '%s/%s' that no longer exists: %v",
				service.Namespace, service.Name, err)
			return nil
		}
		if errors.IsConflict(err) {
			glog.V(4).Infof("Not persisting update to service '%s/%s' that has been changed since we received it: %v",
				service.Namespace, service.Name, err)
			return err
		}
		time.Sleep(cachedService.nextFedUpdateDelay())
	}
	return err
}

// obj could be an *api.Service, or a DeletionFinalStateUnknown marker item.
func (cc *clusterClientCache) enqueueService(obj interface{}, clusterName string) {
	key, err := controller.KeyFunc(obj)
	if err != nil {
		glog.Errorf("Couldn't get key for object %+v: %v", obj, err)
		return
	}
	cc.rwlock.Lock()
	defer cc.rwlock.Unlock()
	_, ok := cc.clientMap[clusterName]
	if ok {
		cc.clientMap[clusterName].serviceQueue.Add(key)
	}
}
