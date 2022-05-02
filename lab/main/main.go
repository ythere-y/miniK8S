package main

import (
	"minik/lab/dksdk"
)

func main() {
	//cl, err := client.NewEnvClient()
	//if err != nil {
	//	fmt.Println("Unable to create docker client")
	//	panic(err)
	//}
	//
	//fmt.Println(cl.ImageList(context.Background(), types.ImageListOptions{}))


	dksdk.CreateContainerInBackground()

}
