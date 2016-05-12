package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/benclarkwood/trez"
	"github.com/pkg/profile"
	"github.com/rcrowley/go-metrics"
)

var images = metrics.NewTimer()

func printStats() {
	fmt.Printf("mean: % 12s, min: % 12s, max: % 12s, %%99: % 12s, stdDev: % 12s, rate: % 8.2f, count: % 8d\n",
		time.Duration(images.Mean()),
		time.Duration(images.Min()),
		time.Duration(images.Max()),
		time.Duration(images.Percentile(0.99)),
		time.Duration(images.StdDev()),
		images.RateMean(),
		images.Count(),
	)
}

func main() {
	p := profile.Start(profile.MemProfile, profile.ProfilePath("."))
	defer p.Stop()
	times := 10000
	workers := runtime.GOMAXPROCS(0)
	filename := ""
	size := "200x200,400x180,800x600"
	algo := "fit,fill"

	flag.IntVar(&times, "times", times, "number of resizes")
	flag.IntVar(&workers, "workers", workers, "number of workers")
	flag.StringVar(&filename, "file", filename, "file to test")
	flag.StringVar(&size, "size", size, "comma separated list of sizes")
	flag.StringVar(&algo, "algo", algo, "comma separated list of algos")
	flag.Parse()

	log.Printf("GOMAXPROCS: %d, WORKERS: %d", runtime.GOMAXPROCS(0), workers)

	if filename == "" {
		fmt.Println("usage: bench filename [options]")
		flag.PrintDefaults()
		log.Fatal("please specify a filename")
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	src, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	opts := []trez.Options{}

	for _, s := range strings.Split(size, ",") {
		size := strings.Split(s, "x")
		w, err := strconv.Atoi(size[0])
		if err != nil {
			log.Fatal(err)
		}
		h, err := strconv.Atoi(size[1])
		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, trez.Options{Width: w, Height: h})
	}

	if len(opts) == 0 {
		log.Fatal("you must provide at least one size")
	}

	algos := []trez.Algo{}
	for _, s := range strings.Split(algo, ",") {
		switch s {
		case "fit":
			algos = append(algos, trez.FIT)
		case "fill":
			algos = append(algos, trez.FILL)
		}
	}

	go func() {
		for _ = range time.Tick(1e9) {
			printStats()
		}
	}()

	wg := new(sync.WaitGroup)
	ch := make(chan trez.Options, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for opt := range ch {
				images.Time(func() {
					_, err := trez.Resize(src, opt)
					if err != nil {
						log.Fatal(err)
					}
				})
			}
		}()
	}

	for i := 0; i < times; i++ {
		for _, opt := range opts {
			for _, algo := range algos {
				opt.Algo = algo
				ch <- opt
			}
		}
	}

	close(ch)
	wg.Wait()
	printStats()
}
