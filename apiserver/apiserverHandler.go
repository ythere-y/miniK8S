package apiserver

import "fmt"

func apiserverHandler(key string, value string) {
	fmt.Printf("apiserver handling key = %v, value = %v\n", key, value)
}
