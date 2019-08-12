package crd

import (
	"log"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

const (
	CRDKind      string = "RDSDatabase"
	CRDSingular  string = "rdsdatabase"
	CRDPlural    string = CRDSingular + "s"
	CRDGroup     string = "aws.pmacik.dev"
	CRDVersion   string = "v1alpha1"
	ShortCRDName string = "rdsdb"
	FullCRDName  string = CRDPlural + "." + CRDGroup
)

// EnsureCRD creates the CRD resource if it does not exist, or returns the one that it's found
func EnsureCRD(clientSet apiextcs.Interface) (*apiextv1beta1.CustomResourceDefinition, error) {
	log.Printf("Ensuring CRD is created...")

	crd, err := findCRD(clientSet)

	if err != nil {
		log.Printf("CRD not found: %v", err)
		log.Printf("Creating a new CRD...")
		crd := &apiextv1beta1.CustomResourceDefinition{
			ObjectMeta: meta_v1.ObjectMeta{Name: FullCRDName},
			Spec: apiextv1beta1.CustomResourceDefinitionSpec{
				Group:   CRDGroup,
				Version: CRDVersion,
				Scope:   apiextv1beta1.NamespaceScoped,
				Names: apiextv1beta1.CustomResourceDefinitionNames{
					Singular:   CRDSingular,
					Plural:     CRDPlural,
					Kind:       CRDKind,
					ShortNames: []string{ShortCRDName},
				},
				Subresources: &apiextv1beta1.CustomResourceSubresources{
					Status: &apiextv1beta1.CustomResourceSubresourceStatus{},
				},
				Versions: []apiextv1beta1.CustomResourceDefinitionVersion{
					apiextv1beta1.CustomResourceDefinitionVersion{
						Name:    CRDVersion,
						Served:  true,
						Storage: true,
					},
				},
			},
		}
		crd, err := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
		if err != nil {
			log.Printf("Failed to create CRD: %v", err)
			return nil, err
		}
	}
	log.Printf("CRD Found, waiting for it to be available...")
	for {
		c, wErr := findCRD(clientSet)
		if wErr == nil {
			crd = c
			break
		}
		log.Printf("CRD not available, yet - trying again.")
	}
	log.Println("CRD Available.")
	return crd, nil
}

func findCRD(clientSet apiextcs.Interface) (*apiextv1beta1.CustomResourceDefinition, error) {
	crd, err := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Get(FullCRDName, meta_v1.GetOptions{
		/*TypeMeta: meta_v1.TypeMeta{
			Kind:       CRDKind,
			APIVersion: CRDGroup + "/" + CRDVersion,
		},*/
	})
	if err != nil {
		return nil, err
	}
	return crd, nil
}

// RDSDatabase is the definition of our CRD RDSDatabase
type RDSDatabase struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               RDSDatabaseSpec   `json:"spec"`
	Status             RDSDatabaseStatus `json:"status"`
}

// RDSDatabaseSpec main structure describing the database instance
type RDSDatabaseSpec struct {
	Username              string         `json:"username"`
	Password              PasswordSecret `json:"password"`
	DBName                string         `json:"dbName"`
	Engine                string         `json:"engine"` // "postgres"
	Class                 string         `json:"class"`  // like "db.t2.micro"
	Size                  int64          `json:"size"`   // size in gb
	MultiAZ               bool           `json:"multiAZ,omitempty"`
	PubliclyAccessible    bool           `json:"publiclyAccessible,omitempty"`
	StorageEncrypted      bool           `json:"storageEncrypted,omitempty"`
	StorageType           string         `json:"storageType,omitempty"`
	Iops                  int64          `json:"iops,omitempty"`
	BackupRetentionPeriod int64          `json:"backupRetentionPeriod,omitempty"` // between 0 and 35, zero means disable
	DeleteProtection      bool           `json:"deleteProtection,omitempty"`
}

type PasswordSecret struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type RDSDatabaseStatus struct {
	State              string `json:"state" description:"State of the deploy"`
	Message            string `json:"message" description:"Detailed message around the state"`
	DBConnectionConfig string `json:"dbConnectionConfig" description:"Name of a Config Map with DB Connection Configuration"`
	DBCredentials      string `json:"dbCredentials" description:"Name of the secret to hold DB Credentials"`
}

type RDSDatabaseList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []RDSDatabase `json:"items"`
}

var SchemeGroupVersion = schema.GroupVersion{Group: CRDGroup, Version: CRDVersion}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&RDSDatabase{},
		&RDSDatabaseList{},
	)
	meta_v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

// Create a Rest client with the new CRD Schema
func NewRESTClient(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	config := *cfg
	config.GroupVersion = &SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, nil, err
	}
	return client, scheme, nil
}
