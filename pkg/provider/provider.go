package provider

import (
	"github.com/operator-backing-service-samples/aws-rds/pkg/crd"
	v1 "k8s.io/api/core/v1"
)

// DBEndpoint represent DB hostname and port
type DBEndpoint struct {
	Hostname string
	Port     int64
}

// RDSDatabaseProvider is the interface for creating and deleting databases
// this is the main interface that should be implemented if a new provider is created
type RDSDatabaseProvider interface {
	CreateRDSDatabase(*crd.RDSDatabase) (*DBEndpoint, error)
	DeleteRDSDatabase(*crd.RDSDatabase) error
	GetSecret(namepspace string, pwname string) (*v1.Secret, error)
}
