package kubernetes

import (
	"fmt"
)

func Test() {
	fmt.Println("master start test")
	StartUpMaster()
	// JoinAsWorker([]string{"localhost:2379", "./config/slaveNode.yaml"})

}

func Main() {
	//Test func
	fmt.Println("kubernetes start up~!")

	//StartUpMaster()

	JoinAsWorker([]string{"localhost:2379", "./config/slaveNode.yaml"})
}
