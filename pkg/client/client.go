package client

import (
	"github.com/operator-backing-service-samples/aws-rds/pkg/crd"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// This file implement all the (CRUD) client methods we need to access our CR objects

func NewCRClient(restClient *rest.RESTClient, scheme *runtime.Scheme, namespace string) *CRClient {
	return &CRClient{restClient: restClient, ns: namespace, resource: crd.CRDPlural,
		codec: runtime.NewParameterCodec(scheme)}
}

type CRClient struct {
	restClient *rest.RESTClient
	ns         string
	resource   string
	codec      runtime.ParameterCodec
}

func (crClient *CRClient) Create(obj *crd.RDSDatabase) (*crd.RDSDatabase, error) {
	var result crd.RDSDatabase
	err := crClient.restClient.Post().
		Namespace(crClient.ns).Resource(crClient.resource).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (crClient *CRClient) Update(obj *crd.RDSDatabase) (*crd.RDSDatabase, error) {
	var result crd.RDSDatabase
	err := crClient.restClient.Put().
		Namespace(crClient.ns).Resource(crClient.resource).Name(obj.Name).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (crClient *CRClient) UpdateStatus(obj *crd.RDSDatabase) (*crd.RDSDatabase, error) {
	var result crd.RDSDatabase
	err := crClient.restClient.Put().
		Namespace(crClient.ns).Resource(crClient.resource).SubResource("status").Name(obj.Name).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (crClient *CRClient) Delete(name string, options *meta_v1.DeleteOptions) error {
	return crClient.restClient.Delete().
		Namespace(crClient.ns).Resource(crClient.resource).
		Name(name).Body(options).Do().
		Error()
}

func (crClient *CRClient) Get(name string) (*crd.RDSDatabase, error) {
	var result crd.RDSDatabase
	err := crClient.restClient.Get().
		Namespace(crClient.ns).Resource(crClient.resource).
		Name(name).Do().Into(&result)
	return &result, err
}

func (crClient *CRClient) List(opts meta_v1.ListOptions) (*crd.RDSDatabaseList, error) {
	var result crd.RDSDatabaseList
	err := crClient.restClient.Get().
		Namespace(crClient.ns).Resource(crClient.resource).
		VersionedParams(&opts, crClient.codec).
		Do().Into(&result)
	return &result, err
}

// Create a new List watch for our CRD
func (crClient *CRClient) NewListWatch() *cache.ListWatch {
	return cache.NewListWatchFromClient(crClient.restClient, crClient.resource, crClient.ns, fields.Everything())
}
