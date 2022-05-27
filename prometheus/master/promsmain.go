package master

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func prometheusMasterMain() {
	var (
		MyTestCounter = prometheus.NewCounter(prometheus.CounterOpts{
			//因为Name不可以重复，所以建议规则为："部门名_业务名_模块名_标量名_类型"
			Name: "demo3_counter",     //唯一id，不可重复Register()，可以Unregister()
			Help: "this is a counter", //对此Counter的描述
		})
		MyTestGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "demo3_gauge",
			Help: "this is a gauge",
		})
		memGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "demo3_mem_gauge_vec",
			Help: "this is a mem gauge_vec",
		}, []string{
			"percent",
		})
	)
	err := prometheus.Register(MyTestGauge)
	if err != nil {
		fmt.Println("register test gauge error")
	}
	err = prometheus.Register(memGauge)
	if err != nil {
		fmt.Println("register mem gauge error")
	}
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			MyTestCounter.Add(100)

			totalPercent, _ := cpu.Percent(time.Second*1, false)
			MyTestGauge.Set(totalPercent[0])
			fmt.Printf("CPU use precent:[%v]", totalPercent[0])

			memInfo, _ := mem.VirtualMemory()
			memGauge.WithLabelValues("total").Set(float64(memInfo.Total))
			memGauge.WithLabelValues("used").Set(float64(memInfo.Used))
			memGauge.WithLabelValues("Available").Set(float64(memInfo.Available))
			fmt.Printf("virtual memory total:[%v], used:[%v], free:[%v]", memInfo.Total, memInfo.Used, memInfo.Available)
		}
	}

}

func listen() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe("0.0.0.0:5050", nil)
}
