package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/ergo/net/edf"
	"ergo.services/logger/colored"
	. "github.com/klauspost/cpuid/v2"
)

type startPublish struct{}

type eventMessage struct {
	Payload string
}

func init() {
	// Register types for network transmission
	if err := edf.RegisterTypeOf(eventMessage{}); err != nil {
		panic(err)
	}
}

var (
	WGready    sync.WaitGroup // Tracks readiness (producer registered + all consumers subscribed)
	WGpublish  sync.WaitGroup // Tracks publish completion
	WGreceive  sync.WaitGroup // Tracks message reception by all consumers
	NCPU       int            = runtime.NumCPU()
	EVENT_NAME gen.Atom       = "benchmark.event"
)

func runFullBenchmark() {
	const (
		numConsumerNodes   = 10
		subscribersPerNode = 100_000
		totalSubscribers   = numConsumerNodes * subscribersPerNode
	)

	fmt.Printf("=================================================================\n")
	fmt.Printf("Distributed Pub/Sub Benchmark: 1 Event -> %d Subscribers\n", totalSubscribers)
	fmt.Printf("=================================================================\n")
	fmt.Printf("Go Version : %s\n", runtime.Version())
	fmt.Printf("CPU: %s (Physical Cores: %d)\n", CPU.BrandName, CPU.PhysicalCores)
	fmt.Printf("Runtime CPUs: %d\n", NCPU)
	fmt.Printf("\n")

	// Create logger
	loggercolored, err := colored.CreateLogger(colored.Options{
		TimeFormat: time.DateTime,
	})
	if err != nil {
		panic(err)
	}

	// Prepare node options
	options := gen.NodeOptions{}
	options.Network.Cookie = "benchmark_cookie"
	options.Log.DefaultLogger.Disable = true
	options.Log.Loggers = append(
		options.Log.Loggers,
		gen.Logger{Name: "colored", Logger: loggercolored},
	)

	fmt.Printf("Step 1: Starting producer node...\n")
	producerNode, err := ergo.StartNode("producer@localhost", options)
	if err != nil {
		panic(err)
	}

	// Start consumer nodes
	fmt.Printf("Step 2: Starting %d consumer nodes...\n", numConsumerNodes)
	consumerNodes := make([]gen.Node, numConsumerNodes)
	for i := 0; i < numConsumerNodes; i++ {
		nodeName := fmt.Sprintf("consumer%d@localhost", i+1)
		node, err := ergo.StartNode(gen.Atom(nodeName), options)
		if err != nil {
			panic(err)
		}
		consumerNodes[i] = node
	}

	// Connect all consumer nodes to producer node
	fmt.Printf("Step 3: Connecting nodes...\n")
	for i := 0; i < numConsumerNodes; i++ {
		if _, err := consumerNodes[i].Network().GetNode(producerNode.Name()); err != nil {
			panic(err)
		}
		producerNode.Log().Info("Connected to %s", consumerNodes[i].Name())
	}

	// Spawn producer process
	fmt.Printf("Step 4: Starting producer process...\n")
	WGready.Add(1)
	producerPID, err := producerNode.Spawn(factory_producer, gen.ProcessOptions{}, EVENT_NAME)
	if err != nil {
		panic(err)
	}
	WGready.Wait() // Wait for producer to register event
	producerNode.Log().Info("Producer process started: %s", producerPID)

	event := gen.Event{
		Node: producerNode.Name(),
		Name: EVENT_NAME,
	}

	// Spawn consumers on each node
	fmt.Printf("Step 5: Spawning %d consumers (%d per node)...\n", totalSubscribers, subscribersPerNode)
	startSpawn := time.Now()
	for i := 0; i < numConsumerNodes; i++ {
		WGready.Add(subscribersPerNode)
		for j := 0; j < subscribersPerNode; j++ {
			_, err := consumerNodes[i].Spawn(factory_consumer, gen.ProcessOptions{}, event)
			if err != nil {
				panic(err)
			}
		}
		consumerNodes[i].Log().Info("Spawned %d consumers on node %d", subscribersPerNode, i+1)
	}

	fmt.Printf("Step 6: Waiting for all consumers to subscribe...\n")
	WGready.Wait() // Wait for all consumers to subscribe
	spawnDuration := time.Since(startSpawn)
	producerNode.Log().Info("All %d consumers subscribed in %s", totalSubscribers, spawnDuration)

	// Prepare for benchmark
	fmt.Printf("\n")
	fmt.Printf("=================================================================\n")
	fmt.Printf("BENCHMARK START: Publishing 1 event to %d subscribers\n", totalSubscribers)
	fmt.Printf("=================================================================\n")

	WGpublish.Add(1)
	WGreceive.Add(totalSubscribers)

	// Trigger publish
	benchmarkStart := time.Now()
	if err := producerNode.Send(producerPID, startPublish{}); err != nil {
		panic(err)
	}
	WGpublish.Wait() // Wait for producer to finish publishing
	publishDuration := time.Since(benchmarkStart)

	// Wait for all consumers to receive
	WGreceive.Wait()
	totalDuration := time.Since(benchmarkStart)

	// Results
	fmt.Printf("\n")
	fmt.Printf("=================================================================\n")
	fmt.Printf("BENCHMARK RESULTS\n")
	fmt.Printf("=================================================================\n")
	fmt.Printf("Total subscribers:       %d\n", totalSubscribers)
	fmt.Printf("Consumer nodes:          %d\n", numConsumerNodes)
	fmt.Printf("Subscribers per node:    %d\n", subscribersPerNode)
	fmt.Printf("\n")
	fmt.Printf("Time to publish:         %s\n", publishDuration)
	fmt.Printf("Time to deliver all:     %s\n", totalDuration)
	fmt.Printf("Network messages sent:   %d (1 per consumer node)\n", numConsumerNodes)
	fmt.Printf("Delivery rate:           %.0f msg/sec\n", float64(totalSubscribers)/totalDuration.Seconds())
	fmt.Printf("=================================================================\n")

	// Give some time for cleanup
	time.Sleep(2 * time.Second)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		fmt.Println("Running small test version...")
		testSmall()
	} else {
		fmt.Println("Running full 1M benchmark...")
		fmt.Println("(Use 'go run . test' for small test version)")
		fmt.Println()
		runFullBenchmark()
	}
}
