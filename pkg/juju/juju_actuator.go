package juju

import (
	"github.com/AlexsJones/cluster-api-provider-juju/api/v1alpha3"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

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

	if err := j.switchControllerContext(jujuConfiguration.Spec.ControllerName); err != nil {
		return E_JUJU_CLUSTER_STATUS_UNKNOWN, err
	}

	cmd, err := exec.Command("bash", "-c", "juju status | grep active | wc -l").Output()
	cmdOutput := strings.TrimSuffix(strings.TrimSpace(string(cmd)), "\n")
	if err != nil {
		fmt.Println(err.Error())
		return E_JUJU_CLUSTER_STATUS_UNKNOWN, err
	}
	i, err := strconv.Atoi(cmdOutput)
	if err != nil {
		return E_JUJU_CLUSTER_STATUS_UNKNOWN, err
	}
	if i > 25 {
		return E_JUJU_CLUSTER_STATUS_RUNNING, nil
	}

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

	if err := j.exec("destroy-model", "-y", jujuConfiguration.Spec.ModelName, "--force", "--destroy-storage", "--no-wait"); err != nil {
		return err
	}

	return nil
}
