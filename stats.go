package stats

import (
	"log"
	"runtime"
	"time"

	"github.com/quipo/statsd"
)

type Stats struct {
	Host     string //
	Port     int16  //
	FilePath string

	Prefix string //

	FlushInterval time.Duration //

	ExportRuntimeData bool //

	Statsd *statsd.StatsdBuffer
}

var logger *log.Logger

func (s *Stats) Start() {
	statsdclient := statsd.NewStdoutClient(s.FilePath, s.Prefix)
	err := statsdclient.CreateSocket()
	if nil != err {
		logger.Fatalln("Unable to create socket")
	}
	stats := statsd.NewStatsdBuffer(s.FlushInterval*time.Second, statsdclient)
	s.Statsd = stats
	if s.ExportRuntimeData {
		go s.runtime_stats(statsdclient, s.FlushInterval)
	}
}

func (s *Stats) track(start time.Time, name string) {
	s.Statsd.PrecisionTiming(name, time.Since(start))
}
func (s *Stats) runtime_stats(statsd *statsd.StdoutClient, interval time.Duration) {
	for {
		<-time.After(interval * time.Second)
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		statsd.Gauge("goroutines", int64(runtime.NumGoroutine()))
		statsd.Gauge("CpuNum", int64(runtime.NumCPU()))
		statsd.Gauge("Gomaxprocs", int64(runtime.GOMAXPROCS(0)))
		statsd.Gauge("memory.alloc", int64(mem.Alloc))
		statsd.Gauge("memory.total.alloc", int64(mem.TotalAlloc))
		statsd.Gauge("memory.sys", int64(mem.Sys))
		statsd.Gauge("memory.lookups", int64(mem.Lookups))
		statsd.Gauge("memory.mallocs", int64(mem.Mallocs))
		statsd.Gauge("memory.frees", int64(mem.Frees))
		statsd.Gauge("stackInUse", int64(mem.StackInuse))
		statsd.Gauge("heap.alloc", int64(mem.HeapAlloc))
		statsd.Gauge("heap.sys", int64(mem.HeapSys))
		statsd.Gauge("heap.idle", int64(mem.HeapIdle))
		statsd.Gauge("heap.inuse", int64(mem.HeapInuse))
		statsd.Gauge("heap.released", int64(mem.HeapReleased))
		statsd.Gauge("heap.objects", int64(mem.HeapObjects))
		statsd.Gauge("gc.next", int64(mem.NextGC))
		statsd.Gauge("gc.last", int64(mem.LastGC))
		statsd.Gauge("gc.num", int64(mem.NumGC))
	}

}
