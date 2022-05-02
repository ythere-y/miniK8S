package service

import "fmt"

type ServiceController struct {
}

func (s ServiceController) CreateService() {
	fmt.Println("create Service")
}

func SerMain() {

	var serviceController1 ServiceController

	serviceController1.CreateService()

}
