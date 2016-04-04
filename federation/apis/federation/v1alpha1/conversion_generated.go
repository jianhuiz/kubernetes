// +build !ignore_autogenerated

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

// This file was autogenerated by conversion-gen. Do not edit it manually!

package v1alpha1

import (
	federation "k8s.io/kubernetes/federation/apis/federation"
	api "k8s.io/kubernetes/pkg/api"
	resource "k8s.io/kubernetes/pkg/api/resource"
	v1 "k8s.io/kubernetes/pkg/api/v1"
	v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	conversion "k8s.io/kubernetes/pkg/conversion"
)

func init() {
	if err := api.Scheme.AddGeneratedConversionFuncs(
		Convert_v1alpha1_Cluster_To_federation_Cluster,
		Convert_federation_Cluster_To_v1alpha1_Cluster,
		Convert_v1alpha1_ClusterCondition_To_federation_ClusterCondition,
		Convert_federation_ClusterCondition_To_v1alpha1_ClusterCondition,
		Convert_v1alpha1_ClusterList_To_federation_ClusterList,
		Convert_federation_ClusterList_To_v1alpha1_ClusterList,
		Convert_v1alpha1_ClusterMeta_To_federation_ClusterMeta,
		Convert_federation_ClusterMeta_To_v1alpha1_ClusterMeta,
		Convert_v1alpha1_ClusterSpec_To_federation_ClusterSpec,
		Convert_federation_ClusterSpec_To_v1alpha1_ClusterSpec,
		Convert_v1alpha1_ClusterStatus_To_federation_ClusterStatus,
		Convert_federation_ClusterStatus_To_v1alpha1_ClusterStatus,
		Convert_v1alpha1_ServerAddressByClientCIDR_To_federation_ServerAddressByClientCIDR,
		Convert_federation_ServerAddressByClientCIDR_To_v1alpha1_ServerAddressByClientCIDR,
		Convert_v1alpha1_SubReplicaSet_To_federation_SubReplicaSet,
		Convert_federation_SubReplicaSet_To_v1alpha1_SubReplicaSet,
		Convert_v1alpha1_SubReplicaSetList_To_federation_SubReplicaSetList,
		Convert_federation_SubReplicaSetList_To_v1alpha1_SubReplicaSetList,
	); err != nil {
		// if one of the conversion functions is malformed, detect it immediately.
		panic(err)
	}
}

func autoConvert_v1alpha1_Cluster_To_federation_Cluster(in *Cluster, out *federation.Cluster, s conversion.Scope) error {
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ClusterSpec_To_federation_ClusterSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ClusterStatus_To_federation_ClusterStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1alpha1_Cluster_To_federation_Cluster(in *Cluster, out *federation.Cluster, s conversion.Scope) error {
	return autoConvert_v1alpha1_Cluster_To_federation_Cluster(in, out, s)
}

func autoConvert_federation_Cluster_To_v1alpha1_Cluster(in *federation.Cluster, out *Cluster, s conversion.Scope) error {
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := Convert_federation_ClusterSpec_To_v1alpha1_ClusterSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_federation_ClusterStatus_To_v1alpha1_ClusterStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

func Convert_federation_Cluster_To_v1alpha1_Cluster(in *federation.Cluster, out *Cluster, s conversion.Scope) error {
	return autoConvert_federation_Cluster_To_v1alpha1_Cluster(in, out, s)
}

func autoConvert_v1alpha1_ClusterCondition_To_federation_ClusterCondition(in *ClusterCondition, out *federation.ClusterCondition, s conversion.Scope) error {
	out.Type = federation.ClusterConditionType(in.Type)
	out.Status = api.ConditionStatus(in.Status)
	if err := api.Convert_unversioned_Time_To_unversioned_Time(&in.LastProbeTime, &out.LastProbeTime, s); err != nil {
		return err
	}
	if err := api.Convert_unversioned_Time_To_unversioned_Time(&in.LastTransitionTime, &out.LastTransitionTime, s); err != nil {
		return err
	}
	out.Reason = in.Reason
	out.Message = in.Message
	return nil
}

func Convert_v1alpha1_ClusterCondition_To_federation_ClusterCondition(in *ClusterCondition, out *federation.ClusterCondition, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterCondition_To_federation_ClusterCondition(in, out, s)
}

func autoConvert_federation_ClusterCondition_To_v1alpha1_ClusterCondition(in *federation.ClusterCondition, out *ClusterCondition, s conversion.Scope) error {
	out.Type = ClusterConditionType(in.Type)
	out.Status = v1.ConditionStatus(in.Status)
	if err := api.Convert_unversioned_Time_To_unversioned_Time(&in.LastProbeTime, &out.LastProbeTime, s); err != nil {
		return err
	}
	if err := api.Convert_unversioned_Time_To_unversioned_Time(&in.LastTransitionTime, &out.LastTransitionTime, s); err != nil {
		return err
	}
	out.Reason = in.Reason
	out.Message = in.Message
	return nil
}

func Convert_federation_ClusterCondition_To_v1alpha1_ClusterCondition(in *federation.ClusterCondition, out *ClusterCondition, s conversion.Scope) error {
	return autoConvert_federation_ClusterCondition_To_v1alpha1_ClusterCondition(in, out, s)
}

func autoConvert_v1alpha1_ClusterList_To_federation_ClusterList(in *ClusterList, out *federation.ClusterList, s conversion.Scope) error {
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := api.Convert_unversioned_ListMeta_To_unversioned_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]federation.Cluster, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_Cluster_To_federation_Cluster(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func Convert_v1alpha1_ClusterList_To_federation_ClusterList(in *ClusterList, out *federation.ClusterList, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterList_To_federation_ClusterList(in, out, s)
}

func autoConvert_federation_ClusterList_To_v1alpha1_ClusterList(in *federation.ClusterList, out *ClusterList, s conversion.Scope) error {
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := api.Convert_unversioned_ListMeta_To_unversioned_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Cluster, len(*in))
		for i := range *in {
			if err := Convert_federation_Cluster_To_v1alpha1_Cluster(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func Convert_federation_ClusterList_To_v1alpha1_ClusterList(in *federation.ClusterList, out *ClusterList, s conversion.Scope) error {
	return autoConvert_federation_ClusterList_To_v1alpha1_ClusterList(in, out, s)
}

func autoConvert_v1alpha1_ClusterMeta_To_federation_ClusterMeta(in *ClusterMeta, out *federation.ClusterMeta, s conversion.Scope) error {
	out.Version = in.Version
	return nil
}

func Convert_v1alpha1_ClusterMeta_To_federation_ClusterMeta(in *ClusterMeta, out *federation.ClusterMeta, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterMeta_To_federation_ClusterMeta(in, out, s)
}

func autoConvert_federation_ClusterMeta_To_v1alpha1_ClusterMeta(in *federation.ClusterMeta, out *ClusterMeta, s conversion.Scope) error {
	out.Version = in.Version
	return nil
}

func Convert_federation_ClusterMeta_To_v1alpha1_ClusterMeta(in *federation.ClusterMeta, out *ClusterMeta, s conversion.Scope) error {
	return autoConvert_federation_ClusterMeta_To_v1alpha1_ClusterMeta(in, out, s)
}

func autoConvert_v1alpha1_ClusterSpec_To_federation_ClusterSpec(in *ClusterSpec, out *federation.ClusterSpec, s conversion.Scope) error {
	if in.ServerAddressByClientCIDRs != nil {
		in, out := &in.ServerAddressByClientCIDRs, &out.ServerAddressByClientCIDRs
		*out = make([]federation.ServerAddressByClientCIDR, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_ServerAddressByClientCIDR_To_federation_ServerAddressByClientCIDR(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.ServerAddressByClientCIDRs = nil
	}
	out.Credential = in.Credential
	return nil
}

func Convert_v1alpha1_ClusterSpec_To_federation_ClusterSpec(in *ClusterSpec, out *federation.ClusterSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterSpec_To_federation_ClusterSpec(in, out, s)
}

func autoConvert_federation_ClusterSpec_To_v1alpha1_ClusterSpec(in *federation.ClusterSpec, out *ClusterSpec, s conversion.Scope) error {
	if in.ServerAddressByClientCIDRs != nil {
		in, out := &in.ServerAddressByClientCIDRs, &out.ServerAddressByClientCIDRs
		*out = make([]ServerAddressByClientCIDR, len(*in))
		for i := range *in {
			if err := Convert_federation_ServerAddressByClientCIDR_To_v1alpha1_ServerAddressByClientCIDR(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.ServerAddressByClientCIDRs = nil
	}
	out.Credential = in.Credential
	return nil
}

func Convert_federation_ClusterSpec_To_v1alpha1_ClusterSpec(in *federation.ClusterSpec, out *ClusterSpec, s conversion.Scope) error {
	return autoConvert_federation_ClusterSpec_To_v1alpha1_ClusterSpec(in, out, s)
}

func autoConvert_v1alpha1_ClusterStatus_To_federation_ClusterStatus(in *ClusterStatus, out *federation.ClusterStatus, s conversion.Scope) error {
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]federation.ClusterCondition, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_ClusterCondition_To_federation_ClusterCondition(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Conditions = nil
	}
	if err := v1.Convert_v1_ResourceList_To_api_ResourceList(&in.Capacity, &out.Capacity, s); err != nil {
		return err
	}
	if err := v1.Convert_v1_ResourceList_To_api_ResourceList(&in.Allocatable, &out.Allocatable, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ClusterMeta_To_federation_ClusterMeta(&in.ClusterMeta, &out.ClusterMeta, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1alpha1_ClusterStatus_To_federation_ClusterStatus(in *ClusterStatus, out *federation.ClusterStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterStatus_To_federation_ClusterStatus(in, out, s)
}

func autoConvert_federation_ClusterStatus_To_v1alpha1_ClusterStatus(in *federation.ClusterStatus, out *ClusterStatus, s conversion.Scope) error {
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]ClusterCondition, len(*in))
		for i := range *in {
			if err := Convert_federation_ClusterCondition_To_v1alpha1_ClusterCondition(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Conditions = nil
	}
	if in.Capacity != nil {
		in, out := &in.Capacity, &out.Capacity
		*out = make(v1.ResourceList, len(*in))
		for key, val := range *in {
			newVal := new(resource.Quantity)
			if err := api.Convert_resource_Quantity_To_resource_Quantity(&val, newVal, s); err != nil {
				return err
			}
			(*out)[v1.ResourceName(key)] = *newVal
		}
	} else {
		out.Capacity = nil
	}
	if in.Allocatable != nil {
		in, out := &in.Allocatable, &out.Allocatable
		*out = make(v1.ResourceList, len(*in))
		for key, val := range *in {
			newVal := new(resource.Quantity)
			if err := api.Convert_resource_Quantity_To_resource_Quantity(&val, newVal, s); err != nil {
				return err
			}
			(*out)[v1.ResourceName(key)] = *newVal
		}
	} else {
		out.Allocatable = nil
	}
	if err := Convert_federation_ClusterMeta_To_v1alpha1_ClusterMeta(&in.ClusterMeta, &out.ClusterMeta, s); err != nil {
		return err
	}
	return nil
}

func Convert_federation_ClusterStatus_To_v1alpha1_ClusterStatus(in *federation.ClusterStatus, out *ClusterStatus, s conversion.Scope) error {
	return autoConvert_federation_ClusterStatus_To_v1alpha1_ClusterStatus(in, out, s)
}

func autoConvert_v1alpha1_ServerAddressByClientCIDR_To_federation_ServerAddressByClientCIDR(in *ServerAddressByClientCIDR, out *federation.ServerAddressByClientCIDR, s conversion.Scope) error {
	out.ClientCIDR = in.ClientCIDR
	out.ServerAddress = in.ServerAddress
	return nil
}

func Convert_v1alpha1_ServerAddressByClientCIDR_To_federation_ServerAddressByClientCIDR(in *ServerAddressByClientCIDR, out *federation.ServerAddressByClientCIDR, s conversion.Scope) error {
	return autoConvert_v1alpha1_ServerAddressByClientCIDR_To_federation_ServerAddressByClientCIDR(in, out, s)
}

func autoConvert_federation_ServerAddressByClientCIDR_To_v1alpha1_ServerAddressByClientCIDR(in *federation.ServerAddressByClientCIDR, out *ServerAddressByClientCIDR, s conversion.Scope) error {
	out.ClientCIDR = in.ClientCIDR
	out.ServerAddress = in.ServerAddress
	return nil
}

func Convert_federation_ServerAddressByClientCIDR_To_v1alpha1_ServerAddressByClientCIDR(in *federation.ServerAddressByClientCIDR, out *ServerAddressByClientCIDR, s conversion.Scope) error {
	return autoConvert_federation_ServerAddressByClientCIDR_To_v1alpha1_ServerAddressByClientCIDR(in, out, s)
}

func autoConvert_v1alpha1_SubReplicaSet_To_federation_SubReplicaSet(in *SubReplicaSet, out *federation.SubReplicaSet, s conversion.Scope) error {
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := v1beta1.Convert_v1beta1_ReplicaSetSpec_To_extensions_ReplicaSetSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.Status, &out.Status, 0); err != nil {
		return err
	}
	return nil
}

func Convert_v1alpha1_SubReplicaSet_To_federation_SubReplicaSet(in *SubReplicaSet, out *federation.SubReplicaSet, s conversion.Scope) error {
	return autoConvert_v1alpha1_SubReplicaSet_To_federation_SubReplicaSet(in, out, s)
}

func autoConvert_federation_SubReplicaSet_To_v1alpha1_SubReplicaSet(in *federation.SubReplicaSet, out *SubReplicaSet, s conversion.Scope) error {
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := v1beta1.Convert_extensions_ReplicaSetSpec_To_v1beta1_ReplicaSetSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.Status, &out.Status, 0); err != nil {
		return err
	}
	return nil
}

func Convert_federation_SubReplicaSet_To_v1alpha1_SubReplicaSet(in *federation.SubReplicaSet, out *SubReplicaSet, s conversion.Scope) error {
	return autoConvert_federation_SubReplicaSet_To_v1alpha1_SubReplicaSet(in, out, s)
}

func autoConvert_v1alpha1_SubReplicaSetList_To_federation_SubReplicaSetList(in *SubReplicaSetList, out *federation.SubReplicaSetList, s conversion.Scope) error {
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := api.Convert_unversioned_ListMeta_To_unversioned_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]federation.SubReplicaSet, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_SubReplicaSet_To_federation_SubReplicaSet(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func Convert_v1alpha1_SubReplicaSetList_To_federation_SubReplicaSetList(in *SubReplicaSetList, out *federation.SubReplicaSetList, s conversion.Scope) error {
	return autoConvert_v1alpha1_SubReplicaSetList_To_federation_SubReplicaSetList(in, out, s)
}

func autoConvert_federation_SubReplicaSetList_To_v1alpha1_SubReplicaSetList(in *federation.SubReplicaSetList, out *SubReplicaSetList, s conversion.Scope) error {
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := api.Convert_unversioned_ListMeta_To_unversioned_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SubReplicaSet, len(*in))
		for i := range *in {
			if err := Convert_federation_SubReplicaSet_To_v1alpha1_SubReplicaSet(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func Convert_federation_SubReplicaSetList_To_v1alpha1_SubReplicaSetList(in *federation.SubReplicaSetList, out *SubReplicaSetList, s conversion.Scope) error {
	return autoConvert_federation_SubReplicaSetList_To_v1alpha1_SubReplicaSetList(in, out, s)
}
