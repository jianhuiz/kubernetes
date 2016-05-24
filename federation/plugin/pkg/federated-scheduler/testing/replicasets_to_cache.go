/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package schedulercache

import (
	"k8s.io/kubernetes/federation/plugin/pkg/federated-scheduler/schedulercache"
	extensions "k8s.io/kubernetes/pkg/apis/extensions"
)

// ReplicaSetsToCache is used for testing
type ReplicaSetsToCache []*extensions.ReplicaSet

func (r ReplicaSetsToCache) AssumeReplicaSet(replicaSet *extensions.ReplicaSet) error {
	return nil
}

func (r ReplicaSetsToCache) AddReplicaSet(replicaSet *extensions.ReplicaSet) error { return nil }

func (r ReplicaSetsToCache) UpdateReplicaSet(oldReplicaSet, newReplicaSet *extensions.ReplicaSet) error {
	return nil
}

func (r ReplicaSetsToCache) RemoveReplicaSet(replicaSet *extensions.ReplicaSet) error { return nil }

func (r ReplicaSetsToCache) GetClusterNameToInfoMap() (map[string]*schedulercache.ClusterInfo, error) {
	return nil, nil
}

func (r ReplicaSetsToCache) List() (selected []*extensions.ReplicaSet, err error) {
	return nil, nil
}
