package juju

import (
	"github.com/AlexsJones/cluster-api-provider-juju/api/v1alpha3"
)

type E_JUJU_CLUSTER_STATUS int

const (
	E_JUJU_CLUSTER_STATUS_UNKNOWN E_JUJU_CLUSTER_STATUS = iota
	E_JUJU_CLUSTER_STATUS_RUNNING
	E_JUJU_CLUSTER_STATUS_STOPPED
	E_JUJU_CLUSTER_STATUS_PROVISIONING
	E_JUJU_CLUSTER_STATUS_DESTROYING
)

type IJuju interface {
	// Within this initial implementation, a Kubernetes cluster name is equal to a model name
	// Also Juju CLI assumes stateful context, however for future implementations we always require the
	// controller and model for working with the cluster
	GetClusterStatus(jujuConfiguration *v1alpha3.JujuConfiguration) (E_JUJU_CLUSTER_STATUS, error)
	CreateControllerIfNotExists(jujuConfiguration *v1alpha3.JujuConfiguration) error
	CreateModelIfNotExists(jujuConfiguration *v1alpha3.JujuConfiguration) error
	CreateCluster(jujuConfiguration *v1alpha3.JujuConfiguration, cluster *v1alpha3.JujuCluster) error
	DestroyCluster(jujuConfiguration *v1alpha3.JujuConfiguration, cluster *v1alpha3.JujuCluster) error
}
