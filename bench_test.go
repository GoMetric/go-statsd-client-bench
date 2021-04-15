package bench

import (
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	gometric "github.com/GoMetric/go-statsd-client"
	cactus "github.com/cactus/go-statsd-client/statsd"
	g2s "github.com/peterbourgon/g2s"
	quipo "github.com/quipo/statsd"
	ac "gopkg.in/alexcesaro/statsd.v2"
)

const (
	host        = "localhost"
	port        = 8125
	prefix      = "prefix."
	prefixNoDot = "prefix"
	counterKey  = "foo.bar.counter"
	gaugeKey    = "foo.bar.gauge"
	gaugeValue  = 42
	timingKey   = "foo.bar.timing"
	timingValue = 153 * time.Millisecond
	flushPeriod = 100 * time.Millisecond
)

func BenchmarkQuipo(b *testing.B) {
	client := quipo.NewStatsdClient(
		host+":"+strconv.Itoa(port),
		prefix,
	)

	err := client.CreateSocket()
	if nil != err {
		log.Println(err)
		os.Exit(1)
	}

	clientBuffer := quipo.NewStatsdBuffer(
		flushPeriod,
		client,
	)

	clientBuffer.Verbose = false

	defer clientBuffer.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Incr(counterKey, 1)
		client.Gauge(gaugeKey, gaugeValue)
		client.Timing(timingKey, int64(timingValue))
	}

	client.Close()
}

func BenchmarkG2s(b *testing.B) {
	client, err := g2s.Dial("udp", host+":"+strconv.Itoa(port))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Counter(1, counterKey, 1)
		client.Gauge(1, gaugeKey, strconv.Itoa(gaugeValue))
		client.Timing(1, timingKey, timingValue)
	}
}

func BenchmarkGoMetric(b *testing.B) {
	client := gometric.NewClient(host, port)
	client.SetPrefix(prefix)
	client.Open()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Count(counterKey, 1, 1)
		client.Gauge(gaugeKey, gaugeValue)
		client.Timing(timingKey, int64(timingValue), 1)
	}

	client.Close()
}

func BenchmarkAlexcesaro(b *testing.B) {
	client, err := ac.New(
		ac.Address(host+":"+strconv.Itoa(port)),
		ac.Prefix(prefixNoDot),
		ac.FlushPeriod(flushPeriod),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Increment(counterKey)
		client.Gauge(gaugeKey, gaugeValue)
		client.Timing(timingKey, timingValue)
	}
	client.Close()
}

func BenchmarkCactus(b *testing.B) {
	client, err := cactus.NewBufferedClient(
		host+":"+strconv.Itoa(port),
		prefix,
		flushPeriod,
		1432,
	)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Inc(counterKey, 1, 1)
		client.Gauge(gaugeKey, gaugeValue, 1)
		client.Timing(timingKey, int64(timingValue), 1)
	}
	client.Close()
}
