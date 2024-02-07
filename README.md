
## Ping

Performs 4 scenarios:
 - 1 process spawns 'pong'-process locally and sends 3M messages
 - N processes spawn 'pong'-process locally and send 1M messages (N = number of CPU)
 - 1 process spawn 'pong'-process on a remote node and sends 3M messages
 - N processes spawn 'pong'-process on a remote node and send 1M messages (N = number of CPU)

![image](https://github.com/ergo-services/benchmarks/assets/118860/f33285b7-5cf9-4195-aa56-c0f4b867d420)

## memory usage

Performs the following scenario:
 - Takes node information that includes memory usage value.
 - Starts 1M processes
 - Takes node information 3 times with 1s interval to make sure the GC has freed unused memory

![image](https://github.com/ergo-services/benchmarks/assets/118860/2003f3e0-c217-4a8c-aa11-63d2d8c50702)
