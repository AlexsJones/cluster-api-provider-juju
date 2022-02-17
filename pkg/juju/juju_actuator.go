package juju

import (
	"cluster-api-provider-juju/api/v1alpha3"
	"fmt"
	"os/exec"

	"github.com/fatih/color"
)

type JujuActuator struct {
}

func NewJujuActuator() *JujuActuator {

	return &JujuActuator{}
}

func (j *JujuActuator) exec(args ...string) error {
	cmd, err := exec.Command("juju", args...).Output()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Print(color.BlueString(string(cmd)))
	return nil
}

func (j *JujuActuator) switchControllerContext(controllerName string) error {

	return j.exec("switch", controllerName)

}

func (j *JujuActuator) GetClusterStatus(jujuConfiguration *v1alpha3.JujuConfiguration) (E_JUJU_CLUSTER_STATUS, error) {

	return E_JUJU_CLUSTER_STATUS_UNKNOWN, nil
}
func (j *JujuActuator) CreateControllerIfNotExists(jujuConfiguration *v1alpha3.JujuConfiguration) error {
	return nil
}
func (j *JujuActuator) CreateModelIfNotExists(jujuConfiguration *v1alpha3.JujuConfiguration) error {
	return j.exec("add-model", jujuConfiguration.Spec.ModelName)
}
func (j *JujuActuator) CreateCluster(jujuConfiguration *v1alpha3.JujuConfiguration, cluster *v1alpha3.JujuCluster) error {

	if err := j.switchControllerContext(jujuConfiguration.Spec.ControllerName); err != nil {
		return err
	}

	j.CreateModelIfNotExists(jujuConfiguration)

	return j.exec("deploy", "charmed-kubernetes", "-m", jujuConfiguration.Spec.ModelName)
}

func (j *JujuActuator) DestroyCluster(jujuConfiguration *v1alpha3.JujuConfiguration, cluster *v1alpha3.JujuCluster) error {

	if err := j.switchControllerContext(jujuConfiguration.Spec.ControllerName); err != nil {
		return err
	}

	if err := j.exec("destroy-model", jujuConfiguration.Spec.ModelName); err != nil {
		return err
	}

	return nil
}
