package juju

import "context"

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
	GetClusterStatus(modelname string, controllerName string) (E_JUJU_CLUSTER_STATUS, error)
	CreateControllerIfNotExists(ctx context.Context, controllerName string) error
	CreateModelIfNotExists(ctx context.Context, modelName string, controllerName string) error
	CreateCluster(ctx context.Context, modelName string, controllerName string) error
}
