package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/goreman"
)

var (
	servers = env.Get("LSIF_NUM_SERVERS", "1", "the number of server instances to run (defaults to one)")
	workers = env.Get("LSIF_NUM_WORKERS", "1", "the number of worker instances to run (defaults to one)")

	// Set in docker image
	prometheusStorageDir       = os.Getenv("PROMETHEUS_STORAGE_DIR")
	prometheusConfigurationDir = os.Getenv("PROMETHEUS_CONFIGURATION_DIR")
)

const (
	serverPort = 3186
	workerPort = 3187
)

func main() {
	numServers, err := strconv.ParseInt(servers, 10, 64)
	if err != nil || numServers < 0 || numServers > 1 {
		log.Fatalf("Invalid int %q for LSIF_NUM_SERVERS: %s", servers, err)
	}

	numWorkers, err := strconv.ParseInt(workers, 10, 64)
	if err != nil || numWorkers < 0 {
		log.Fatalf("Invalid int %q for LSIF_NUM_WORKERS: %s", workers, err)
	}

	if err := ioutil.WriteFile(
		filepath.Join(prometheusConfigurationDir, "targets.yml"),
		[]byte(makePrometheusTargets(numServers, numWorkers)),
		0644,
	); err != nil {
		log.Fatalf("Writing prometheus config: %v", err)
	}

	// This mirrors the behavior from cmd/start
	if err := goreman.Start([]byte(makeProcfile(numServers, numWorkers)), goreman.Options{
		RPCAddr:        "127.0.0.1:5005",
		ProcDiedAction: goreman.Shutdown,
	}); err != nil {
		log.Fatalf("Starting goreman: %v", err)
	}
}

func makeProcfile(numServers, numWorkers int64) string {
	procfile := []string{}
	addProcess := func(name, command string) {
		procfile = append(procfile, fmt.Sprintf("%s: %s", name, command))
	}

	if numServers > 0 {
		addProcess("lsif-server", "node /lsif/out/server/server.js")
	}

	for i := 0; i < int(numWorkers); i++ {
		addProcess(
			fmt.Sprintf("lsif-worker-%d", i),
			fmt.Sprintf("env WORKER_METRICS_PORT=%d node /lsif/out/worker/worker.js", workerPort+i),
		)
	}

	addProcess("prometheus", fmt.Sprintf("prometheus '--storage.tsdb.path=%s' '--config.file=%s/prometheus.yml'",
		prometheusStorageDir,
		prometheusConfigurationDir,
	))

	return strings.Join(procfile, "\n") + "\n"
}

func makePrometheusTargets(numServers, numWorkers int64) string {
	content := []string{"---"}
	addTarget := func(job string, port int) {
		content = append(content,
			"- labels:",
			fmt.Sprintf("    job: %s", job),
			"  targets:",
			fmt.Sprintf("    - 127.0.0.1:%d", port),
		)
	}

	if numServers > 0 {
		addTarget("lsif-server", serverPort)
	}

	for i := 0; i < int(numWorkers); i++ {
		addTarget("lsif-worker", workerPort+i)
	}

	return strings.Join(content, "\n") + "\n"
}
