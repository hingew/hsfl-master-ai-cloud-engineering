package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type TestConfig struct {
	NumSteps                int    `yaml:"numSteps"`
	StepSize                int    `yaml:"stepSize"`
	NoCoalescingUrl         string `yaml:"noCoalescingUrl"`
	ControllerCoalescingUrl string `yaml:"controllerCoalescingUrl"`
	GatewayCoalescingUrl    string `yaml:"gatewayCoalescingUrl"`
}

func readConfig(configPath string) TestConfig {
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config TestConfig
	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func doRequest(url string, wg *sync.WaitGroup, durations chan<- time.Duration) {
	defer wg.Done()

	startTime := time.Now()

	response, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("Couldn't send request: %s", err))
	}
	defer response.Body.Close()

	duration := time.Since(startTime)

	if response.StatusCode != http.StatusOK {
		log.Default().Printf("Response status not ok: %d", response.StatusCode)
		durations <- -1
	} else {
		durations <- duration
	}
}

func doMultipleRequests(url string, numRequests int) <-chan time.Duration {
	durations := make(chan time.Duration, numRequests)
	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go doRequest(url, &wg, durations)
	}

	go func() {
		wg.Wait()
		close(durations)
	}()

	return durations
}

func calculateAverageDuration(durations <-chan time.Duration) *time.Duration {
	var totalDuration time.Duration
	var validDurations []time.Duration

	for d := range durations {
		if d <= -1 {
			return nil
		}
		validDurations = append(validDurations, d)
	}

	sort.Slice(validDurations, func(i, j int) bool {
		return validDurations[i] < validDurations[j]
	})

	lowerBound := int(0.1 * float64(len(validDurations)))
	upperBound := int(0.9 * float64(len(validDurations)))

	for i := lowerBound; i < upperBound; i++ {
		totalDuration += validDurations[i]
	}

	averageDuration := totalDuration / time.Duration(upperBound-lowerBound)

	return &averageDuration
}

func measureTimeToDoRequests(url string, numRequests int) *time.Duration {
	durations := doMultipleRequests(url, numRequests)
	return calculateAverageDuration(durations)
}

func createCSVFile() *csv.Writer {
	file, err := os.Create("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer
}

func writeCSVHeader(writer *csv.Writer) {
	writer.Write([]string{"Num Request", "No Coalescing (ms)", "Controller Coalescing (ms)", "Gateway Coalescing (ms)"})
	fmt.Println("Num Request | No Coalescing (ms) | Controller Coalescing (ms) | Gateway Coalescing (ms)")
}

func writeCSVRow(writer *csv.Writer, numRequests int, durationNoCoalescing *time.Duration, durationControllerCoalescing *time.Duration, durationGatewayCoalescing *time.Duration) {
	requestStr := fmt.Sprintf("%-12d", numRequests)
	var noCoalescingStr string
	var controllerCoalescingStr string
	var gatewayCoalescingStr string

	if durationNoCoalescing != nil {
		noCoalescingStr = fmt.Sprintf("%-20d", durationNoCoalescing.Milliseconds())
	} else {
		noCoalescingStr = fmt.Sprintf("%-20s", "N/A")
	}
	if durationNoCoalescing != nil {
		controllerCoalescingStr = fmt.Sprintf("%-20d", durationControllerCoalescing.Milliseconds())
	} else {
		controllerCoalescingStr = fmt.Sprintf("%-20s", "N/A")
	}
	if durationNoCoalescing != nil {
		gatewayCoalescingStr = fmt.Sprintf("%-20d", durationGatewayCoalescing.Milliseconds())
	} else {
		gatewayCoalescingStr = fmt.Sprintf("%-20s", "N/A")
	}

	fmt.Printf("%s | %s | %s | %s\n", requestStr, noCoalescingStr, controllerCoalescingStr, gatewayCoalescingStr)
	err := writer.Write([]string{
		fmt.Sprintf("%s", requestStr),
		fmt.Sprintf("%s", noCoalescingStr),
		fmt.Sprintf("%s", controllerCoalescingStr),
		fmt.Sprintf("%s", gatewayCoalescingStr),
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config := readConfig(*configPath)

	writer := createCSVFile()

	writeCSVHeader(writer)

	for steps := 1; steps <= config.NumSteps; steps++ {
		numRequests := steps * config.StepSize

		durationNoCoalescing := measureTimeToDoRequests(config.NoCoalescingUrl, numRequests)
		durationControllerCoalescing := measureTimeToDoRequests(config.ControllerCoalescingUrl, numRequests)
		durationGatewayCoalescing := measureTimeToDoRequests(config.GatewayCoalescingUrl, numRequests)

		writeCSVRow(writer, numRequests, durationNoCoalescing, durationControllerCoalescing, durationGatewayCoalescing)
	}
}
