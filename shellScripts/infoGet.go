package shellScripts

import (
	"log"
	"net"
)

func GetEth0() (string, error) {
	var (
		res string
		err error
	)
	res = ""
	err = nil

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	boliAddr := localAddr.String()

	for i := 0; i < len(boliAddr); i++ {
		if boliAddr[i] == ':' {
			break
		}
		res += string(boliAddr[i])
	}

	//fmt.Printf("localAddr = %v\n", res)
	return res, err
}
