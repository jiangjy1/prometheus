package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.bb.local/jiangjunyu/dlog"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
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

//noinspection ALL
func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(rpcDurations)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
	//logFile, err := os.OpenFile(`./urls.log`, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	panic(err)
	//}
	// 设置存储位置
	//log.SetOutput(logFile)
	var logSaveDays = 3
	var logFile, _ = filepath.Abs("./urls.log")

	// 日志初始化
	logger.Config(logFile, 4)
	logger.SetSaveMode(1)
	logger.SetSaveDays(logSaveDays)
	logger.Infof("loginit success ! \n")
}

func testurl() (filepaths string) {
	filepaths, _ = filepath.Abs("./urls")
	fmt.Println(filepaths)
	return filepaths
}

func readfile() ([]string, error) {
	var urls []string
	r, err := os.Open(testurl())
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return urls, err
	}
	defer r.Close()
	var s = bufio.NewScanner(r)
	for s.Scan() { // 循环直到文件结束
		line := s.Text() // 这个 line 就是每一行的文本了，string 类型
		//fmt.Println(line)
		urls = append(urls, line)
	}
	return urls, err
	//fmt.Println(urls)
	//for k, v := range urls {
	//	fmt.Printf("url[%d]:%s\r\n", k, v)
	//}
}

func listenURL(url1 string) {
	for {
		u, _ := url.Parse(url1)
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
			//fmt.Printf("%s success , http_status is %d \r\n", u.String(), resCode)
			log.Printf("%s success , http_status is %d \r\n", u.String(), resCode)

		} else {
			log.Printf("%s failed , http_status is %d \r\n", u.String(), resCode)
		}
		v := rand.ExpFloat64()
		link := u.String()
		responseCode := strconv.Itoa(resCode)
		rpcDurations.WithLabelValues(link, responseCode).Observe(v)
		time.Sleep(10 * time.Second)
	}
}

func main() {
	flag.Parse()

	urls, err := readfile()
	if err != nil {
		log.Print("read urls err :", err)
	}
	for _, v := range urls {
		go listenURL(v)
	}

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
