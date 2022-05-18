package scheduler

import (
	"fmt"
	"minik8s/constant"
	. "minik8s/lab/etcd"
	"minik8s/utils"
)

func CreateSchedulerManager(listen string) {
	SyncWatch(listen, schedulerHandler)
}

func Main() {
	fmt.Println("Scheduler Main started!")

	watchName := SetKey(
		SetPrefix(constant.SchedulerPrefix),
		SetSourceType("pods"),
	)

	CreateSchedulerManager(watchName)

	utils.HoldPro()
}
