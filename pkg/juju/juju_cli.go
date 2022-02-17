package juju

import (
	"bufio"
	"context"
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
	cmd := exec.Command("juju", args...)
	stdout, err := cmd.Output()
	stderr, _ := cmd.StderrPipe()
	fmt.Print(color.BlueString(string(stdout)))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(color.RedString(scanner.Text()))
	}

	return nil
}

func (j *JujuActuator) switchControllerContext(controllerName string) error {

	j.exec("switch", "controllerName")

	return nil
}

func (j *JujuActuator) CreateControllerIfNotExists(ctx context.Context, controllerName string) error {
	return j.exec("controller", "add", controllerName)
}

func (j *JujuActuator) CreateModelIfNotExists(ctx context.Context, modelName string, controllerName string) error {

	if err := j.switchControllerContext(controllerName); err != nil {
		return err
	}
	return j.exec("add-model", "add", modelName)
}

func (j *JujuActuator) CreateCluster(ctx context.Context, modelName string, controllerName string) error {

	if err := j.switchControllerContext(controllerName); err != nil {
		return err
	}
	return j.exec("deploy", "charmed-kubernetes", "-m", modelName)
}

func (j *JujuActuator) GetClusterStatus(modelName string, controllerName string) (E_JUJU_CLUSTER_STATUS, error) {
	return E_JUJU_CLUSTER_STATUS_UNKNOWN, nil
}
