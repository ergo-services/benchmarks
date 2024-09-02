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
