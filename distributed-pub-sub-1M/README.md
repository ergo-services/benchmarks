# Distributed Pub/Sub Benchmark: 1M Subscribers

This benchmark demonstrates Ergo Framework's pub/sub performance by measuring event delivery time from 1 producer to 1,000,000 subscribers distributed across 10 nodes.

## Scenario

- **1 producer node** with 1 producer process
- **10 consumer nodes** with 100,000 consumers each (1M total)
- **1 event** published from producer
- **Measurement**: Time from publish to all 1M consumers receiving the message

## What This Tests

This benchmark validates the subscription sharing optimization described in the documentation:

- Without optimization: Would require **1,000,000 network messages** (one per subscriber)
- With optimization: Requires only **10 network messages** (one per consumer node)

The optimization transforms O(N) network cost into O(M) cost, where:
- N = total consumers (1,000,000)
- M = number of nodes (10)

## Running the Benchmark

```bash
go run .
```

## Expected Results

The benchmark measures:
- **Time to publish**: How long it takes the producer to publish the event
- **Time to deliver all**: Total time for all 1M consumers to receive the event
- **Delivery rate**: Messages delivered per second

The network optimization means:
- Only 10 network messages are sent (one per consumer node)
- Each node locally distributes to its 100K subscribers
- Network cost is constant regardless of subscriber count per node

### Results on Apple M4 Max 

```
...
Step 5: Spawning 1000000 consumers (100000 per node)...
2025-12-01 16:26:05 [info] 336B493D: new connection with 'consumer10@localhost' (9E2174F9)
2025-12-01 16:26:05 [info] CF9DAE57: Spawned 100000 consumers on node 1
2025-12-01 16:26:06 [info] DCB8560A: Spawned 100000 consumers on node 2
2025-12-01 16:26:06 [info] 4B5802C0: Spawned 100000 consumers on node 3
2025-12-01 16:26:06 [info] FAF3A6B0: Spawned 100000 consumers on node 4
2025-12-01 16:26:06 [info] 6D13F27A: Spawned 100000 consumers on node 5
2025-12-01 16:26:06 [info] 7E360A27: Spawned 100000 consumers on node 6
2025-12-01 16:26:07 [info] E9D65EED: Spawned 100000 consumers on node 7
2025-12-01 16:26:07 [info] B66447C4: Spawned 100000 consumers on node 8
2025-12-01 16:26:07 [info] 2184130E: Spawned 100000 consumers on node 9
2025-12-01 16:26:07 [info] 9E2174F9: Spawned 100000 consumers on node 10
Step 6: Waiting for all consumers to subscribe...
2025-12-01 16:26:07 [info] 336B493D: All 1000000 consumers subscribed in 1.842276s

=================================================================
BENCHMARK START: Publishing 1 event to 1000000 subscribers
=================================================================
2025-12-01 16:26:07 [info] <336B493D.0.1004>: Producer publishing event...

=================================================================
BENCHMARK RESULTS
=================================================================
Total subscribers:       1000000
Consumer nodes:          10
Subscribers per node:    100000

Time to publish:         64.125µs
Time to deliver all:     342.414375ms
Network messages sent:   10 (1 per consumer node)
Delivery rate:           2920438 msg/sec
=================================================================
``````

## Architecture

```
Producer Node
  └─ Producer Process (registers event, publishes message)
       ↓ (10 network messages)
       ├─ Consumer Node 1 → 100K local subscribers
       ├─ Consumer Node 2 → 100K local subscribers
       ├─ Consumer Node 3 → 100K local subscribers
       ├─ Consumer Node 4 → 100K local subscribers
       ├─ Consumer Node 5 → 100K local subscribers
       ├─ Consumer Node 6 → 100K local subscribers
       ├─ Consumer Node 7 → 100K local subscribers
       ├─ Consumer Node 8 → 100K local subscribers
       ├─ Consumer Node 9 → 100K local subscribers
       └─ Consumer Node 10 → 100K local subscribers
```

## Memory Considerations

This benchmark spawns 1 million goroutines (one per consumer process). Ensure your system has sufficient memory:
- Approximate memory per process: ~2-4KB
- Total memory for 1M processes: ~2-4GB
- Plus overhead for nodes, network connections, etc.


