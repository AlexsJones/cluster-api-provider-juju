package juju

import "context"

type IJuju interface {
	CreateControllerIfNotExists(ctx context.Context, controllerName string) error
	CreateModelIfNotExists(ctx context.Context, modelName string) error
	CreateCluster(ctx context.Context) error
}
