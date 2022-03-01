# cluster-api-provider-juju

<img src="images/juju.svg" width="250">

This is the [juju](https://juju.is/) provider for cluster API.
It enables you to provision [Charmed Kubernetes](https://ubuntu.com/kubernetes) infrastructure.

<img src="images/provider.png" width="500">

### Running from local

- `microk8s config > ~/config && export KUBECONFIG=~/config`
- `clusterctl init`
- Install local CRD's with `make install`
- Run the `main.go` locally and connect to the cluster


### Running from cluster

- `microk8s config > ~/config && export KUBECONFIG=~/config`
- `clusterctl init`
- `kubectl cp /home/<NAME>/.local/share/juju <pod-name>:/root/.local/share/ -n cluster-api-provider-juju-system`
- Install local CRD's with `make install`
- `make deploy`

#### Dependencies

_Either with MacOS or LinuxBrew_

```
brew install kubebuilder kustomize clusterctl
```

- When running locally you will need `juju` installed and configured `sudo snap install juju --classic`

- Optionally using `microK8s` or `kind`
  - `snap install microk8s --classic`
  - `brew install kind`
