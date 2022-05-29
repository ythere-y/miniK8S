package cni

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	gocni "github.com/containerd/go-cni"
)

var (
	idF = flag.Int("id", 2, "")
	//confFF = flag.String("conf", "/etc/cni/net.d/11-flannel.conf", "")
	confFF = flag.String("conf", "./etc/cni/net.d/11-flannel.conf", "")
)

func init() {
	flag.Parse()
}

func Main() {
	fmt.Printf("hello world\n")

	id := *idF
	netns := fmt.Sprintf("/var/run/netns/ns-%d", id)

	l, err := gocni.New(
		gocni.WithMinNetworkCount(2),
		gocni.WithPluginDir([]string{"./etc/cni/bin"}),
		gocni.WithInterfacePrefix("eth"))
	if err != nil {
		log.Fatalf("failed to initialize cni library: %v", err)
	}

	if err := l.Load(gocni.WithLoNetwork, gocni.WithConfListFile(*confFF)); err != nil {
		log.Fatalf("failed to load cni configuration: %v", err)
	} else {
		log.Printf("success to load cni configuration !")
	}
	ctx := context.Background()
	defer func() {
		if err := l.Remove(ctx, strconv.Itoa(id), netns); err != nil {
			log.Fatalf("failed to teardown network: %v", err)
		} else {
			log.Printf("success to teardown network!\n")
		}
	}()

	result, err := l.Setup(ctx, strconv.Itoa(id), netns)
	if err != nil {
		log.Fatalf("failed to setup network for namespace: %v", err)
	} else {
		log.Printf("success to setup network for namespace !\n")
	}
	os.Exit(1)
	return

	for key, iff := range result.Interfaces {
		if len(iff.IPConfigs) > 0 {
			IP := iff.IPConfigs[0].IP.String()
			fmt.Printf("IP of the interface %s:%s\n", key, IP)
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
