[![Gitbook Documentation](https://img.shields.io/badge/GitBook-Documentation-f37f40?style=plastic&logo=gitbook&logoColor=white&style=flat)](https://docs.ergo.services)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Telegram Community](https://img.shields.io/badge/Telegram-ergo__services-229ed9?style=flat&logo=telegram&logoColor=white)](https://t.me/ergo_services)
[![Twitter](https://img.shields.io/badge/Twitter-ergo__services-00acee?style=flat&logo=twitter&logoColor=white)](https://twitter.com/ergo_services)
[![Reddit](https://img.shields.io/badge/Reddit-r/ergo__services-ff4500?style=plastic&logo=reddit&logoColor=white&style=flat)](https://reddit.com/r/ergo_services)

# Benchmarks of the Ergo Framework 3.0 (and above)

The tests below are performed on the laptop Macbook Air M3 (2024)

## Ping

Performs 4 scenarios:
 - 1 process spawns 'pong'-process locally and sends 3M messages
 - N processes spawn 'pong'-process locally and send 1M messages (N = number of CPU)
 - 1 process spawns 'pong'-process on a remote node and sends 3M messages
 - N processes spawn 'pong'-process on a remote node and send 1M messages (N = number of CPU)

![image](https://github.com/ergo-services/benchmarks/assets/118860/31e17b33-ce92-4ef1-8dec-d6bcac0ab66f)

## Memory usage (per process)

Performs the following scenario:
 - Takes node information that includes memory usage value.
 - Starts 1M processes
 - Takes node information 3 times with 1s intervals to make sure the GC has freed unused memory

![image](https://github.com/ergo-services/benchmarks/assets/118860/ead567bb-beae-40bf-b881-519e89ce1190)

## Serialization benchmarks: EDF vs Protobuf vs Gob

These benchmarks compare EDF, EDF (+cache), Protobuf, and Gob serialization performance across common data types.
EDF and Gob rely on runtime reflection, which dynamically inspects and serializes data structures at runtime. Protobuf uses code generation, producing static type-safe marshalling and unmarshalling logic.
| Data Type | EDF | EDF (+cache) | Protobuf | Gob | Winner | EDF Advantage |
|-----------|-----|--------------|----------|-----|---------|---------------|
| String Encode | 29.96ns, 53B, 0a | **23.45ns, 0B, 0a** | 44.37ns, 32B, 1a | 76.42ns, 16B, 1a | **EDF+Cache** | 47% faster than Protobuf, 69% faster than Gob |
| String Decode | 76.63ns, 72B, 4a | 76.62ns, 72B, 4a | **67.45ns, 96B, 2a** | 429.5ns, 1000B, 19a | **Protobuf** | EDF 14% slower, but 6x faster than Gob |
| Map Encode | 310.5ns, 325B, 10a | **224.7ns, 204B, 5a** | 339.9ns, 112B, 5a | 234.0ns, 32B, 2a | **EDF+Cache** | 34% faster than Protobuf, competitive with Gob |
| Map Decode | 557.4ns, 955B, 25a | 468.8ns, 846B, 22a | **353.3ns, 528B, 13a** | 6733ns, 8256B, 185a | **Protobuf** | EDF+Cache 33% slower, but 14x faster than Gob |
| Complex Struct Encode | 277.8ns, 307B, 6a | **269.8ns, 306B, 6a** | 474.3ns, 224B, 9a | 357.1ns, 88B, 5a | **EDF+Cache** | 43% faster than Protobuf, 24% faster than Gob |
| Complex Struct Decode | 740.1ns, 1364B, 41a | 700.8ns, 1368B, 41a | **597.2ns, 824B, 25a** | 9335ns, 10872B, 255a | **Protobuf** | EDF+Cache 17% slower, but 13x faster than Gob |
| Nested Struct Encode | **739.9ns, 732B, 17a** | 796.2ns, 846B, 17a | 1557ns, 640B, 27a | 901.4ns, 256B, 15a | **EDF** | 52% faster than Protobuf, 18% faster than Gob |
| Nested Struct Decode | 2137ns, 4594B, 107a | 2291ns, 5438B, 120a | **1684ns, 2544B, 71a** | 12729ns, 14712B, 342a | **Protobuf** | EDF 27% slower, but 6x faster than Gob |

*Format: `time ns/op, memory B/op, allocations/op`*

*Hardware: `Apple M4 Max`*



