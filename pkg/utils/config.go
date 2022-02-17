package utils

import (
	"cluster-api-provider-juju/api/v1alpha3"
	"context"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func FetchJujuConfigurationObject(cluster *v1alpha3.JujuCluster, cli client.Client, ctx context.Context) (*v1alpha3.JujuConfiguration, error) {

	conf := v1alpha3.JujuConfiguration{}
	if err := cli.Get(ctx, types.NamespacedName{
		Name:      cluster.Spec.JujuConfiguration,
		Namespace: cluster.Namespace,
	}, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
