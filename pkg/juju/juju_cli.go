package juju

import (
	"context"
	"fmt"
	"os/exec"
)

type JujuCLI struct {
}

func NewJujuCLI() *JujuCLI {

	return &JujuCLI{}
}

func (j *JujuCLI) exec(args ...string) error {
	cmd := exec.Command("juju", args...)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Print(string(stdout))

	return nil
}

func (j *JujuCLI) CreateControllerIfNotExists(ctx context.Context, controllerName string) error {
	return j.exec("controller", "add", controllerName)
}

func (j *JujuCLI) CreateModelIfNotExists(ctx context.Context, modelName string) error {
	// TODO: Make this better
	return j.exec("add-model", modelName)

}

func (j *JujuCLI) CreateCluster(ctx context.Context) error {
	return j.exec("deploy", "charmed-kubernetes")
}
