package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

var startTime time.Time

type Metrics struct {
	reqTotal,
	evntRecvTotal, evntSentTotal, evntErrTotal *prometheus.CounterVec
	reqDuration, reqSize, respSize, reqTotalSize, respTotalSize prometheus.Summary
	uptime,
	cpuCores, cpuUsage, cpuSystemUsage, cpuUserUsage, cpuLoad1, cpuLoad5, cpuLoad15,
	memUsage, memTotal, memBuffers, memCached, memUsed, memFree,
	swapUsage, swapTotal, swapUsed, swapFree,
	diskRead, diskWrite, diskUsage, diskUsed, diskFree, diskTotal,
	networkIn, networkOut, networkInPkt, networkOutPkt prometheus.Gauge
}

func (m *Metrics) Handler() gin.HandlerFunc {

	h := promhttp.Handler()
	return func(c *gin.Context) {
		m.uptime.Set(time.Since(startTime).Seconds())

		cpuCores, err := cpu.Counts(true)
		if err != nil {
			panic(err)
		}
		m.cpuCores.Set(float64(cpuCores))

		// ....TODO: cpus
		cpuUsage, err := cpu.Percent(0, false)
		if err != nil {
			panic(err)
		}
		m.cpuUsage.Set(cpuUsage[0])

		// ....TODO: cpus
		cpuTimes, err := cpu.Times(false)
		if err != nil {
			panic(err)
		}
		m.cpuSystemUsage.Set(cpuTimes[0].System)
		m.cpuUserUsage.Set(cpuTimes[0].User)

		cpuLoad, err := load.Avg()
		if err != nil {
			panic(err)
		}
		m.cpuLoad1.Set(cpuLoad.Load1)
		m.cpuLoad5.Set(cpuLoad.Load5)
		m.cpuLoad15.Set(cpuLoad.Load15)

		memVtul, err := mem.VirtualMemory()
		if err != nil {
			panic(err)
		}
		m.memUsage.Set(memVtul.UsedPercent)
		m.memTotal.Set(float64(int(memVtul.Total) / MB))
		m.memBuffers.Set(float64(int(memVtul.Buffers) / MB))
		m.memCached.Set(float64(int(memVtul.Cached) / MB))
		m.memUsed.Set(float64(int(memVtul.Used) / MB))
		m.memFree.Set(float64(int(memVtul.Free) / MB))

		memSwap, err := mem.SwapMemory()
		if err != nil {
			panic(err)
		}
		m.swapUsage.Set(memSwap.UsedPercent)
		m.swapTotal.Set(float64(int(memSwap.Total) / MB))
		m.swapUsed.Set(float64(int(memSwap.Used) / MB))
		m.swapFree.Set(float64(int(memSwap.Free) / MB))

		//proc, err := process.NewProcess(int32(os.Getpid()))
		//if err != nil {
		//	panic(err)
		//}
		//ioCnt, err := proc.IOCounters()
		//if err != nil {
		//	panic(err)
		//}
		//m.diskRead.Set(float64(ioCnt.ReadBytes))
		//m.diskWrite.Set(float64(ioCnt.WriteBytes))

		disk, err := disk.Usage("/")
		if err != nil {
			panic(err)
		}
		m.diskUsage.Set(float64(disk.UsedPercent))
		m.diskUsed.Set(float64(int(disk.Used) / GB))
		m.diskFree.Set(float64(int(disk.Free) / GB))
		m.diskTotal.Set(float64(int(disk.Total) / GB))

		// ....TODO: interfaces
		n, err := net.IOCounters(false)
		if err != nil {
			panic(err)
		}
		m.networkIn.Set(float64(n[0].BytesRecv))
		m.networkOut.Set(float64(n[0].BytesSent))
		m.networkInPkt.Set(float64(n[0].PacketsRecv))
		m.networkOutPkt.Set(float64(n[0].PacketsSent))

		h.ServeHTTP(c.Writer, c.Request)
	}
}

func New() *Metrics {
	m := &Metrics{}
	startTime = time.Now()

	// Uptime
	m.uptime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "up_time_seconds",
		Help: "The uptime(seconds) of the server.",
	})
	prometheus.MustRegister(m.uptime)

	// CPU
	m.cpuCores = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_cores",
		Help: "x",
	})
	prometheus.MustRegister(m.cpuCores)
	m.cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "x",
	})
	prometheus.MustRegister(m.cpuUsage)
	m.cpuSystemUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_system_usage",
		Help: "The percentage of the system cpu usage.",
	})
	prometheus.MustRegister(m.cpuSystemUsage)
	m.cpuUserUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_user_usage",
		Help: "The percentage of the user cpu usage.",
	})
	prometheus.MustRegister(m.cpuUserUsage)
	m.cpuLoad1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_load_1",
		Help: "x",
	})
	prometheus.MustRegister(m.cpuLoad1)
	m.cpuLoad5 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_load_5",
		Help: "x",
	})
	prometheus.MustRegister(m.cpuLoad5)
	m.cpuLoad15 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_load_15",
		Help: "x",
	})
	prometheus.MustRegister(m.cpuLoad15)

	// Memory
	m.memUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage",
		Help: "The percentage of the memory usage.",
	})
	prometheus.MustRegister(m.memUsage)
	m.memTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_total",
		Help: "MB",
	})
	prometheus.MustRegister(m.memTotal)
	m.memBuffers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_buffers",
		Help: "MB",
	})
	prometheus.MustRegister(m.memBuffers)
	m.memCached = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_cached",
		Help: "MB",
	})
	prometheus.MustRegister(m.memCached)
	m.memUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_used",
		Help: "MB",
	})
	prometheus.MustRegister(m.memUsed)
	m.memFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_free",
		Help: "MB",
	})
	prometheus.MustRegister(m.memFree)
	//
	m.swapUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "swap_memory_usage",
		Help: "The percentage of the memory usage.",
	})
	prometheus.MustRegister(m.swapUsage)
	m.swapTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "swap_memory_total",
		Help: "MB",
	})
	prometheus.MustRegister(m.swapTotal)
	m.swapUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "swap_memory_used",
		Help: "MB",
	})
	prometheus.MustRegister(m.swapUsed)
	m.swapFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "swap_memory_free",
		Help: "MB",
	})
	prometheus.MustRegister(m.swapFree)

	// Disk
	m.diskRead = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_read_speed",
		Help: "The speed(MB/s) of the disk read operation.",
	})
	prometheus.MustRegister(m.diskRead)
	m.diskWrite = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_write_speed",
		Help: "The speed(MB/s) of the disk write operation.",
	})
	prometheus.MustRegister(m.diskWrite)
	m.diskUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_usage",
		Help: "The percentage of the disk usage.",
	})
	prometheus.MustRegister(m.diskUsage)
	m.diskUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_used",
		Help: "GB",
	})
	prometheus.MustRegister(m.diskUsed)
	m.diskFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_free",
		Help: "GB",
	})
	prometheus.MustRegister(m.diskFree)
	m.diskTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_total",
		Help: "GB",
	})
	prometheus.MustRegister(m.diskTotal)

	// Network
	m.networkIn = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_inbound",
		Help: "The speed(MB/s) of the network outbound.",
	})
	prometheus.MustRegister(m.networkIn)
	m.networkOut = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_outbound",
		Help: "The speed(MB/s) of the network inbound.",
	})
	prometheus.MustRegister(m.networkOut)
	m.networkInPkt = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_inbound_packets",
		Help: "x",
	})
	prometheus.MustRegister(m.networkInPkt)
	m.networkOutPkt = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_outbound_packets",
		Help: "x",
	})
	prometheus.MustRegister(m.networkOutPkt)

	// Requests & Responses
	m.reqTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_total",
			Help: "Total number of HTTP requests made.",
		},
		[]string{"method", "code"},
	)
	prometheus.MustRegister(m.reqTotal)
	m.reqDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "request_duration_seconds",
			Help: "The HTTP request latencies in seconds.",
		},
	)
	prometheus.MustRegister(m.reqDuration)
	m.reqSize = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "request_size_bytes",
			Help: "The HTTP request sizes in bytes.",
		},
	)
	prometheus.MustRegister(m.reqSize)
	m.respSize = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "response_size_bytes",
			Help: "The HTTP response sizes in bytes.",
		},
	)
	prometheus.MustRegister(m.respSize)
	m.reqTotalSize = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "request_size_bytes_total",
			Help: "x",
		},
	)
	prometheus.MustRegister(m.reqTotalSize)
	m.respTotalSize = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "response_size_bytes_total",
			Help: "x",
		},
	)
	prometheus.MustRegister(m.respTotalSize)

	// Events
	m.evntRecvTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "event_received_total",
			Help: "Total number of the receivied events.",
		},
		[]string{"method", "event"},
	)
	prometheus.MustRegister(m.evntRecvTotal)
	m.evntSentTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "event_sent_total",
			Help: "Total number of the sent events.",
		},
		[]string{"event"},
	)
	prometheus.MustRegister(m.evntSentTotal)
	m.evntErrTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "event_error_total",
			Help: "Total time of the error occurred while sending the events.",
		},
		[]string{"event"},
	)
	prometheus.MustRegister(m.evntErrTotal)

	// in
	// out
	// reset after secs
	return m
}
