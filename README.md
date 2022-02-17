# cluster-api-provider-juju

<img src="images/juju.svg" width="250">

This is the [juju](https://juju.is/) provider for cluster API.
It enables you to provision [Charmed Kubernetes](https://ubuntu.com/kubernetes) infrastructure.

<img src="images/provider.png" width="500">

### Development

- `make cluster-up`
- Install manifests `make install`
- Run the `main.go` locally and connect to the cluster


#### Dependencies

- Either with MacOS or LinuxBrew

```
brew install kubebuilder kustomize clusterctl kind
```
