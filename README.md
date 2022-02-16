# cluster-api-provider-juju

<img src="images/juju.svg" width="250">

This is the [juju](https://juju.is/) provider for cluster API.
It enables you to provision [Charmed Kubernetes](https://ubuntu.com/kubernetes) infrastructure.


### Development

- Setup a local Kubernetes cluster with Microk8s
- `clusterctl init`
- Install manifests `make install`
- Run the `main.go` locally and connect to the cluster


#### Dependencies

- Either with MacOS or LinuxBrew

```
brew install kubebuilder kustomize clusterctl
```

