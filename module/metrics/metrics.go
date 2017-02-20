package metrics

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/module/event/eventstore"
	"github.com/TeaMeow/KitSvc/module/logger"
	"github.com/TeaMeow/KitSvc/module/mq/mqstore"
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
var path = "/metrics"

// Metrics stores instrument data.
type Metrics struct {
	// CounterVec: Request
	reqTotal,
	// CounterVec: Message
	msgRecvTotal, msgSentTotal, msgErrTotal,
	// CounterVec: Event
	evntRecvTotal, evntSentTotal, evntErrTotal *prometheus.CounterVec

	// Summary: Request & Response
	reqDuration, reqSize, respSize prometheus.Summary

	// Gauge: Uptime
	uptime,
	// Gauge: Event
	evntUnsent, evntRecv, evntSent, evntOnline,
	// Gauge: Request
	req, recv,
	// Gauge: Message
	msgUnsent, msgRecv, msgSent, msgOnline,
	// Gauge: CPU
	cpuCores, cpuUsage, cpuSystemUsage, cpuUserUsage, cpuLoad1, cpuLoad5, cpuLoad15,
	// Gauge: Memory
	memUsage, memTotal, memBuffers, memCached, memUsed, memFree,
	// Gauge: Swap
	swapUsage, swapTotal, swapUsed, swapFree,
	// Gauge: Disk
	diskRead, diskWrite, diskUsage, diskUsed, diskFree, diskTotal,
	// Gauge: Network
	networkIn, networkOut, networkInTotal, networkOutTotal, networkInPkt, networkOutPkt prometheus.Gauge

	// Network outbound, inbout kbps.
	lastOutbound,
	lastInbound,
	// Event received, sent total in the last second.
	lastEvntRecv,
	lastEvntSent,
	// Message received, sent total in the last second.
	lastMsgRecv,
	lastMsgSent uint64
}

// information stores the system information.
type information struct {
	cpuCores int
	cpuUsage []float64
	cpuTimes []cpu.TimesStat
	cpuLoad  *load.AvgStat
	memVtul  *mem.VirtualMemoryStat
	memSwap  *mem.SwapMemoryStat
	disk     *disk.UsageStat
	network  []net.IOCountersStat
}

// systemInfo returns the system information.
func systemInfo() (information, error) {
	cpuCores, err := cpu.Counts(true)
	if err != nil {
		return information{}, err
	}
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return information{}, err
	}
	cpuTimes, err := cpu.Times(false)
	if err != nil {
		return information{}, err
	}
	cpuLoad, err := load.Avg()
	if err != nil {
		return information{}, err
	}
	memVtul, err := mem.VirtualMemory()
	if err != nil {
		return information{}, err
	}
	memSwap, err := mem.SwapMemory()
	if err != nil {
		return information{}, err
	}
	disk, err := disk.Usage("/")
	if err != nil {
		return information{}, err
	}
	network, err := net.IOCounters(false)
	if err != nil {
		return information{}, err
	}
	return information{cpuCores, cpuUsage, cpuTimes, cpuLoad, memVtul, memSwap, disk, network}, nil
}

// instrument the system, and stores the information to the `Metric` struct.
func (m *Metrics) instrument() error {
	info, err := systemInfo()
	if err != nil {
		return err
	}

	// Uptime
	m.uptime.Set(time.Since(startTime).Seconds())
	// CPU Cores
	m.cpuCores.Set(float64(info.cpuCores))
	// CPU Usages
	m.cpuUsage.Set(info.cpuUsage[0])
	m.cpuSystemUsage.Set(info.cpuTimes[0].System)
	m.cpuUserUsage.Set(info.cpuTimes[0].User)
	// CPU Load
	m.cpuLoad1.Set(info.cpuLoad.Load1)
	m.cpuLoad5.Set(info.cpuLoad.Load5)
	m.cpuLoad15.Set(info.cpuLoad.Load15)
	// Memory
	m.memUsage.Set(info.memVtul.UsedPercent)
	m.memTotal.Set(float64(int(info.memVtul.Total) / MB))
	m.memBuffers.Set(float64(int(info.memVtul.Buffers) / MB))
	m.memCached.Set(float64(int(info.memVtul.Cached) / MB))
	m.memUsed.Set(float64(int(info.memVtul.Used) / MB))
	m.memFree.Set(float64(int(info.memVtul.Free) / MB))
	// Swap
	m.swapUsage.Set(info.memSwap.UsedPercent)
	m.swapTotal.Set(float64(int(info.memSwap.Total) / MB))
	m.swapUsed.Set(float64(int(info.memSwap.Used) / MB))
	m.swapFree.Set(float64(int(info.memSwap.Free) / MB))
	// Data I/O
	//proc, err := process.NewProcess(int32(os.Getpid()))
	//if err != nil {
	//	panic(err)
	//}
	//ioCnt, err := proc.IOCounters()
	//if err != nil {
	//	panic(err)
	//}
	//m.diskRead.Set(float64(info.ioCnt.ReadBytes))
	//m.diskWrite.Set(float64(info.ioCnt.WriteBytes))
	// Disk
	m.diskUsage.Set(float64(info.disk.UsedPercent))
	m.diskUsed.Set(float64(int(info.disk.Used) / MB))
	m.diskFree.Set(float64(int(info.disk.Free) / MB))
	m.diskTotal.Set(float64(int(info.disk.Total) / MB))
	// Network
	m.networkInTotal.Set(float64(info.network[0].BytesRecv))
	m.networkOutTotal.Set(float64(info.network[0].BytesSent))
	m.networkInPkt.Set(float64(info.network[0].PacketsRecv))
	m.networkOutPkt.Set(float64(info.network[0].PacketsSent))

	return nil
}

// instruEvent instruments the event traffic.
func (m *Metrics) instruEvent() {
	for {
		<-time.After(time.Second * 1)

		queueTotal, recvTotal, sentTotal, allConnected :=
			eventstore.QueueTotal,
			eventstore.RecvTotal,
			eventstore.SentTotal,
			eventstore.AllConnected

		connected := 0
		if allConnected {
			connected = 1
		}

		m.evntOnline.Set(float64(connected))
		m.evntUnsent.Set(float64(queueTotal))
		m.evntRecv.Set(float64(recvTotal - int(m.lastEvntRecv)))
		m.evntSent.Set(float64(sentTotal - int(m.lastEvntSent)))

		m.lastEvntRecv = uint64(recvTotal)
		m.lastEvntSent = uint64(sentTotal)
	}
}

// instruMessage instruments the message traffic.
func (m *Metrics) instruMessage() {
	for {
		<-time.After(time.Second * 1)

		queueTotal, recvTotal, sentTotal, allConnected :=
			mqstore.QueueTotal,
			mqstore.RecvTotal,
			mqstore.SentTotal,
			mqstore.AllConnected

		connected := 0
		if allConnected {
			connected = 1
		}

		m.msgOnline.Set(float64(connected))
		m.msgUnsent.Set(float64(queueTotal))
		m.msgRecv.Set(float64(recvTotal - int(m.lastMsgRecv)))
		m.msgSent.Set(float64(sentTotal - int(m.lastMsgSent)))

		m.lastMsgRecv = uint64(recvTotal)
		m.lastMsgSent = uint64(sentTotal)
	}
}

// instruNetwork instruments the network traffic.
func (m *Metrics) instruNetwork() {
	for {
		<-time.After(time.Second * 1)

		n, err := net.IOCounters(false)
		if err != nil {
			panic(err)
		}

		m.networkIn.Set(float64(n[0].BytesRecv - m.lastInbound))
		m.networkOut.Set(float64(n[0].BytesSent - m.lastOutbound))

		m.lastInbound = n[0].BytesRecv
		m.lastOutbound = n[0].BytesSent
	}
}

// instruRequest instruments the http request information.
func (m *Metrics) instruRequest(c *gin.Context) {
	reqSize := make(chan int)
	go computeReqSize(c.Request, reqSize)

	start := time.Now()

	status := strconv.Itoa(c.Writer.Status())
	elapsed := float64(time.Since(start)) / float64(time.Second)
	respSize := float64(c.Writer.Size())

	m.reqDuration.Observe(elapsed)
	m.reqTotal.WithLabelValues(status, c.Request.Method, c.HandlerName()).Inc()
	m.reqSize.Observe(float64(<-reqSize))
	m.respSize.Observe(respSize)
}

// Handler handles and parses, calculates the incoming http requests.
func (m *Metrics) Handler() gin.HandlerFunc {
	// Keep instrumenting the network traffic.
	go m.instruNetwork()
	// Keep instrumenting the event traffic.
	go m.instruEvent()
	// Keep instrumenting the message traffic.
	go m.instruMessage()

	return func(c *gin.Context) {
		switch c.Request.URL.String() {
		// Ignore the health check requests.
		case "/sd/health", "/sd/ram", "/sd/cpu", "/sd/disk":
			c.Next()

			// Collect the system information when we received the metrics request.
		case path:
			if err := m.instrument(); err != nil {
				logger.WarningFields("Error occurred while instrumenting the system.", logrus.Fields{
					"err": err,
				})
			}
			c.Next()

			// Measure the bytes of the request and the response
			// if it's the common HTTP request.
		default:
			// Increase the event total if an event request was received.
			if strings.Contains(c.Request.URL.String(), "/es/") {
				m.evntRecvTotal.WithLabelValues(c.HandlerName()).Inc()
			}
			// Increase the message total if an message request was received.
			if strings.Contains(c.Request.URL.String(), "/mq/") {
				m.msgRecvTotal.WithLabelValues(c.HandlerName()).Inc()
			}
			go m.instruRequest(c)
			c.Next()
		}
	}
}

// PrometheusHandler serves the collected service information for Prometheus.
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// New defines the instrument elements.
func New() *Metrics {
	m := &Metrics{}
	startTime = time.Now()

	// Uptime
	m.uptime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "up_time_seconds",
		Help: "Server uptime.",
	})
	prometheus.MustRegister(m.uptime)

	// CPU
	m.cpuCores = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_cores",
		Help: "Number of processor cores.",
	})
	prometheus.MustRegister(m.cpuCores)
	m.cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "Percentage of processor usage.",
	})
	prometheus.MustRegister(m.cpuUsage)
	m.cpuSystemUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_system_usage",
		Help: "x",
	})
	prometheus.MustRegister(m.cpuSystemUsage)
	m.cpuUserUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_user_usage",
		Help: "x",
	})
	prometheus.MustRegister(m.cpuUserUsage)
	m.cpuLoad1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_load_1",
		Help: "Processor load average in one minute.",
	})
	prometheus.MustRegister(m.cpuLoad1)
	m.cpuLoad5 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_load_5",
		Help: "Processor load average in five minutes.",
	})
	prometheus.MustRegister(m.cpuLoad5)
	m.cpuLoad15 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_load_15",
		Help: "Processor load average in fifthteen minutes.",
	})
	prometheus.MustRegister(m.cpuLoad15)

	// Memory
	m.memUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage",
		Help: "Percentage of memory usage.",
	})
	prometheus.MustRegister(m.memUsage)
	m.memTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_total",
		Help: "The size of the total memory (MB).",
	})
	prometheus.MustRegister(m.memTotal)
	m.memBuffers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_buffers",
		Help: "The size of the buffered memory (MB).",
	})
	prometheus.MustRegister(m.memBuffers)
	m.memCached = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_cached",
		Help: "The size of the cached memory (MB).",
	})
	prometheus.MustRegister(m.memCached)
	m.memUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_used",
		Help: "The size of the used memory (MB).",
	})
	prometheus.MustRegister(m.memUsed)
	m.memFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_free",
		Help: "The size of the free memory (MB).",
	})
	prometheus.MustRegister(m.memFree)
	//
	m.swapUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "swap_memory_usage",
		Help: "Percentage of swap memory usage.",
	})
	prometheus.MustRegister(m.swapUsage)
	m.swapTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "swap_memory_total",
		Help: "The size of the total swap memory (MB).",
	})
	prometheus.MustRegister(m.swapTotal)
	m.swapUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "swap_memory_used",
		Help: "The size of the used swap memory (MB).",
	})
	prometheus.MustRegister(m.swapUsed)
	m.swapFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "swap_memory_free",
		Help: "The size of the free swap memory (MB).",
	})
	prometheus.MustRegister(m.swapFree)

	// Disk
	m.diskRead = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_read_speed",
		Help: "Disk read speed (Byte/s).",
	})
	prometheus.MustRegister(m.diskRead)
	m.diskWrite = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_write_speed",
		Help: "Disk write speed (Byte/s).",
	})
	prometheus.MustRegister(m.diskWrite)
	m.diskUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_usage",
		Help: "Percentage of disk usage.",
	})
	prometheus.MustRegister(m.diskUsage)
	m.diskUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_used",
		Help: "The size of the used disk (MB).",
	})
	prometheus.MustRegister(m.diskUsed)
	m.diskFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_free",
		Help: "The size of the free disk (MB).",
	})
	prometheus.MustRegister(m.diskFree)
	m.diskTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_total",
		Help: "The size of the total disk (MB).",
	})
	prometheus.MustRegister(m.diskTotal)

	// Network
	m.networkIn = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_inbound",
		Help: "Network inbound speed (Byte/s).",
	})
	prometheus.MustRegister(m.networkIn)
	m.networkOut = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_outbound",
		Help: "Network outbound speed (Byte/s).",
	})
	prometheus.MustRegister(m.networkOut)
	m.networkInTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_inbound_total",
		Help: "The size of the total inbound bytes.",
	})
	prometheus.MustRegister(m.networkInTotal)
	m.networkOutTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_outbound_total",
		Help: "The size of the total outbound bytes.",
	})
	prometheus.MustRegister(m.networkOutTotal)
	m.networkInPkt = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_inbound_packets",
		Help: "The count of the total inbound packets.",
	})
	prometheus.MustRegister(m.networkInPkt)
	m.networkOutPkt = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_outbound_packets",
		Help: "The count of the total outbound packets.",
	})
	prometheus.MustRegister(m.networkOutPkt)

	// Requests & Responses
	m.reqTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_total",
			Help: "Total number of HTTP requests made.",
		},
		[]string{"code", "method", "handler"},
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

	// Events
	m.evntRecvTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "event_received_total",
			Help: "Total number of the receivied events.",
		},
		[]string{"handler"},
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
			Help: "Total count of the error occurred while sending the events.",
		},
		[]string{"event"},
	)
	prometheus.MustRegister(m.evntErrTotal)
	m.evntUnsent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "event_unsent",
		Help: "The total of the unsent events.",
	})
	prometheus.MustRegister(m.evntUnsent)
	m.evntRecv = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "event_received",
		Help: "The amount of the received events in a second.",
	})
	prometheus.MustRegister(m.evntRecv)
	m.evntSent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "event_sent",
		Help: "The amount of the sent events in a second.",
	})
	prometheus.MustRegister(m.evntSent)
	m.evntOnline = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "event_online",
		Help: "Returns 1 if connected to the Event Store, otherwise 0.",
	})
	prometheus.MustRegister(m.evntOnline)

	// Messages
	m.msgRecvTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "message_received_total",
			Help: "Total number of the receivied messages.",
		},
		[]string{"handler"},
	)
	prometheus.MustRegister(m.msgRecvTotal)
	m.msgSentTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "message_sent_total",
			Help: "Total number of the sent messages.",
		},
		[]string{"topic"},
	)
	prometheus.MustRegister(m.msgSentTotal)
	m.msgErrTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "message_error_total",
			Help: "Total count of the error occurred while sending the messages.",
		},
		[]string{"topic"},
	)
	prometheus.MustRegister(m.msgErrTotal)
	m.msgUnsent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "message_unsent",
		Help: "The total of the unsent messages.",
	})
	prometheus.MustRegister(m.msgUnsent)
	m.msgRecv = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "message_received",
		Help: "The amount of the received messages in a second.",
	})
	prometheus.MustRegister(m.msgRecv)
	m.msgSent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "message_sent",
		Help: "The amount of the sent messages in a second.",
	})
	prometheus.MustRegister(m.msgSent)
	m.msgOnline = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "message_online",
		Help: "Returns 1 if connected to the NSQ Producer, otherwise 0.",
	})
	prometheus.MustRegister(m.msgOnline)

	return m
}

// computeReqSize computes the byte size of the request,
// from https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go
func computeReqSize(r *http.Request, out chan int) {
	s := 0
	if r.URL != nil {
		s = len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.
	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	out <- s
}
