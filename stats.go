package stats

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/quipo/statsd"
)

type Stats struct {
	Host              string        //
	Port              int16         //
	Prefix            string        //
	FlushInterval     time.Duration //
	ExportRuntimeData bool          //
	Statsd            *statsd.StatsdBuffer
}

var logger *log.Logger

func (s *Stats) Start() {
	statsdclient := statsd.NewStatsdClient(fmt.Sprintf("%s:%d", s.Host, s.Port), s.Prefix)
	err := statsdclient.CreateSocket()
	if nil != err {
		logger.Fatalln("Unable to create socket")
	}
	stats := statsd.NewStatsdBuffer(s.FlushInterval*time.Second, statsdclient)
	s.Statsd = stats
	if s.ExportRuntimeData {
		go s.runtimeStats(statsdclient, s.FlushInterval)
	}
}

func (s *Stats) track(start time.Time, name string) {
	s.Statsd.PrecisionTiming(name, time.Since(start))
}
func (s *Stats) runtimeStats(statsd *statsd.StatsdClient, interval time.Duration) {
	for {
		<-time.After(interval * time.Second)
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		statsd.Gauge("cpu.goroutines", int64(runtime.NumGoroutine()))
		statsd.Gauge("cpu.cgocalls", int64(runtime.NumCgoCall()))
		statsd.Gauge("cpu.cpunum", int64(runtime.NumCPU()))
		// Memory
		statsd.Gauge("memory.alloc", int64(mem.Alloc))
		statsd.Gauge("memory.total", int64(mem.TotalAlloc))
		statsd.Gauge("memory.othersys", int64(mem.OtherSys))
		statsd.Gauge("memory.sys", int64(mem.Sys))
		statsd.Gauge("memory.lookups", int64(mem.Lookups))
		statsd.Gauge("memory.mallocs", int64(mem.Mallocs))
		statsd.Gauge("memory.frees", int64(mem.Frees))
		// Stack
		statsd.Gauge("stack.inuse", int64(mem.StackInuse))
		statsd.Gauge("stack.sys", int64(mem.StackSys))
		statsd.Gauge("stack.mspan_inuse", int64(mem.MSpanInuse))
		statsd.Gauge("stack.mspan_sys", int64(mem.MSpanSys))
		statsd.Gauge("stack.mcache_inuse", int64(mem.MCacheInuse))
		statsd.Gauge("stack.mcache_sys", int64(mem.MCacheSys))
		// Heap
		statsd.Gauge("heap.alloc", int64(mem.HeapAlloc))
		statsd.Gauge("heap.sys", int64(mem.HeapSys))
		statsd.Gauge("heap.idle", int64(mem.HeapIdle))
		statsd.Gauge("heap.inuse", int64(mem.HeapInuse))
		statsd.Gauge("heap.released", int64(mem.HeapReleased))
		statsd.Gauge("heap.objects", int64(mem.HeapObjects))
		// GC
		statsd.Gauge("gc.next", int64(mem.NextGC))
		statsd.Gauge("gc.last", int64(mem.LastGC))
		statsd.Gauge("gc.count", int64(mem.NumGC))
		statsd.Gauge("gc.sys", int64(mem.GCSys))
		statsd.Gauge("gc.pause_total", int64(mem.PauseTotalNs))
		statsd.Gauge("gc.pause", int64(mem.PauseNs[(mem.NumGC+255)%255]))
	}

}
