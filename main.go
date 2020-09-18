package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	//"math"
	//"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	//oscillationPeriod = flag.Duration("oscillation-period", 10*time.Minute, "The duration of the rate oscillation period.")
)

var (
	labels = []string{"link", "responseCode"}
)

var (
	// Create a summary to track fictional interservice RPC latencies for three
	// distinct services with different latency distributions. These services are
	// differentiated via a "service" label.
	rpcDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "url_health_probe",
			Help: "RPC latency distributions.",
			//Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		labels,
	)
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(rpcDurations)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}

func testurl() (filepaths string) {
	filepaths, _ = filepath.Abs("./urls")
	fmt.Println(filepaths)
	return filepaths
}

func readfile() {
	var urls []string
	r, err := os.Open(testurl())
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer r.Close()
	var s = bufio.NewScanner(r)
	for s.Scan() { // 循环直到文件结束
		line := s.Text() // 这个 line 就是每一行的文本了，string 类型
		fmt.Println(line)
		urls = append(urls, line)
	}
	fmt.Println(urls)
	for k, v := range urls {
		fmt.Printf("url[%d]:%s\r\n", k, v)
	}
}

func main() {
	flag.Parse()

	go func() {
		for {
			u, _ := url.Parse("http://www.baidu.com")
			//q := u.Query()
			//u.RawQuery = q.Encode()
			res, err := http.Get(u.String())
			if err != nil {
				fmt.Println("0")
				time.Sleep(60 * time.Second)
				continue
			}
			resCode := res.StatusCode
			err = res.Body.Close()
			if err != nil {
				fmt.Println("0")
				return
			}
			if resCode == 200 {
				fmt.Printf("%s success , http_status is %d \r\n", u.String(), resCode)
				fmt.Printf("%s is success", u.String())
			} else {
				fmt.Printf("%s failed , http_status is %d \r\n", u.String(), resCode)
			}
			v := rand.ExpFloat64()
			link := u.String()
			responseCode := strconv.Itoa(resCode)
			rpcDurations.WithLabelValues(link, responseCode).Observe(v)
			time.Sleep(10 * time.Second)
		}
	}()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
