package service

import (
	"fmt"
	"minik8s/cmd"
)

type ServiceController struct {
}

func (s ServiceController) CreateService() {
	fmt.Println("create Service")
	cmd.CreatServerTest()
}

func SerMain() {

	var serviceController1 ServiceController

	serviceController1.CreateService()

}
