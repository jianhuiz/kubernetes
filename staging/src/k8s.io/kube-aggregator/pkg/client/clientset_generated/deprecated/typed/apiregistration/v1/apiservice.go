/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	scheme "k8s.io/kube-aggregator/pkg/client/clientset_generated/deprecated/scheme"
)

// APIServicesGetter has a method to return a APIServiceInterface.
// A group's client should implement this interface.
type APIServicesGetter interface {
	APIServices() APIServiceInterface
}

// APIServiceInterface has methods to work with APIService resources.
type APIServiceInterface interface {
	Create(*v1.APIService) (*v1.APIService, error)
	Update(*v1.APIService) (*v1.APIService, error)
	UpdateStatus(*v1.APIService) (*v1.APIService, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.APIService, error)
	List(opts metav1.ListOptions) (*v1.APIServiceList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.APIService, err error)
	APIServiceExpansion
}

// aPIServices implements APIServiceInterface
type aPIServices struct {
	client rest.Interface
}

// newAPIServices returns a APIServices
func newAPIServices(c *ApiregistrationV1Client) *aPIServices {
	return &aPIServices{
		client: c.RESTClient(),
	}
}

// Get takes name of the aPIService, and returns the corresponding aPIService object, and an error if there is any.
func (c *aPIServices) Get(name string, options metav1.GetOptions) (result *v1.APIService, err error) {
	result = &v1.APIService{}
	err = c.client.Get().
		Resource("apiservices").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of APIServices that match those selectors.
func (c *aPIServices) List(opts metav1.ListOptions) (result *v1.APIServiceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.APIServiceList{}
	err = c.client.Get().
		Resource("apiservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(context.TODO()).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested aPIServices.
func (c *aPIServices) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("apiservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(context.TODO())
}

// Create takes the representation of a aPIService and creates it.  Returns the server's representation of the aPIService, and an error, if there is any.
func (c *aPIServices) Create(aPIService *v1.APIService) (result *v1.APIService, err error) {
	result = &v1.APIService{}
	err = c.client.Post().
		Resource("apiservices").
		Body(aPIService).
		Do(context.TODO()).
		Into(result)
	return
}

// Update takes the representation of a aPIService and updates it. Returns the server's representation of the aPIService, and an error, if there is any.
func (c *aPIServices) Update(aPIService *v1.APIService) (result *v1.APIService, err error) {
	result = &v1.APIService{}
	err = c.client.Put().
		Resource("apiservices").
		Name(aPIService.Name).
		Body(aPIService).
		Do(context.TODO()).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *aPIServices) UpdateStatus(aPIService *v1.APIService) (result *v1.APIService, err error) {
	result = &v1.APIService{}
	err = c.client.Put().
		Resource("apiservices").
		Name(aPIService.Name).
		SubResource("status").
		Body(aPIService).
		Do(context.TODO()).
		Into(result)
	return
}

// Delete takes name of the aPIService and deletes it. Returns an error if one occurs.
func (c *aPIServices) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("apiservices").
		Name(name).
		Body(options).
		Do(context.TODO()).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *aPIServices) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("apiservices").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do(context.TODO()).
		Error()
}

// Patch applies the patch and returns the patched aPIService.
func (c *aPIServices) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.APIService, err error) {
	result = &v1.APIService{}
	err = c.client.Patch(pt).
		Resource("apiservices").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do(context.TODO()).
		Into(result)
	return
}
