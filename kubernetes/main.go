package kubernetes

import (
	"fmt"
	"time"
)

func Test() {
	fmt.Println("slave join test")
	StartUpMaster()
	//controllerManager.ControllerStartUp()
	//JoinAsWorker([]string{"localhost:2379", "./config/masterNode.yaml"})
	time.Sleep(4 * time.Second)

	JoinAsWorker([]string{"localhost:2379", "./config/slaveNode.yaml"})

}

func Main() {
	//Test func
	fmt.Println("kubernetes start up~!")

	//StartUpMaster()

	JoinAsWorker([]string{"localhost:2379", "./config/slaveNode.yaml"})
}
