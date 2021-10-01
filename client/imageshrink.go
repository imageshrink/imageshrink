package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/imageshrink/imageshrink/common"
)

func buildWorkers(nproc int, remotes string) []common.Worker {
	workers := make([]common.Worker, 0)
	_, err := exec.LookPath("convert")
	if nil == err {
		for i := 0; i < nproc; i++ {
			workers = append(workers, common.MakeLocalWorker())
			fmt.Println("[Worker] Added a local worker")
		}
	}
	splits := strings.Split(remotes, ",")
	for _, host := range splits {
		if len(host) == 0 {
			continue
		}
		workers = append(workers, common.MakeRemoteWorker(host))
		fmt.Printf("[Worker] Added a remote worker: %v\n", host)
	}
	return workers
}

func main() {
	remotes := flag.String(
		"remotes", "",
		"remotes workers: host1:port1,host2:port2",
	)
	nproc := flag.Int(
		"nproc", 1,
		"number of local processes",
	)
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Printf("Usage: imageshrink [path to scan]\n")
		return
	}
	common.DoImageShrink(args[0], buildWorkers(*nproc, *remotes))
}
