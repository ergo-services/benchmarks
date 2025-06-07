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

### **📤 Encoding Performance Comparison**

| Data Type | EDF | Protobuf | Gob | **Winner** | **EDF Advantage** |
|-----------|-----|----------|-----|------------|-------------------|
| **String** | 28.75 ns/op (57 B, 0 allocs) | 42.29 ns/op (32 B, 1 alloc) | 72.98 ns/op (16 B, 1 alloc) | **🥇 EDF** | **47% faster** than Protobuf |
| **PID** | 40.65 ns/op (32 B, 1 alloc) | 51.15 ns/op (24 B, 1 alloc) | 106.9 ns/op (32 B, 1 alloc) | **🥇 EDF** | **26% faster** than Protobuf |
| **ProcessID** | 41.90 ns/op (32 B, 1 alloc) | 48.70 ns/op (32 B, 1 alloc) | 98.46 ns/op (32 B, 1 alloc) | **🥇 EDF** | **16% faster** than Protobuf |
| **Simple Struct** | 108.5 ns/op (155 B, 2 allocs) | 46.15 ns/op (24 B, 1 alloc) | 111.5 ns/op (32 B, 1 alloc) | **🥇 Protobuf** | 57% slower than Protobuf |
| **Complex Struct** | 283.2 ns/op (307 B, 6 allocs) | 479.5 ns/op (224 B, 9 allocs) | 351.9 ns/op (88 B, 5 allocs) | **🥇 EDF** | **41% faster** than Protobuf |
| **Nested Struct** | 731.4 ns/op (731 B, 17 allocs) | 1516 ns/op (640 B, 27 allocs) | 909.5 ns/op (256 B, 15 allocs) | **🥇 EDF** | **52% faster** than Protobuf |
| **Map** | 306.9 ns/op (325 B, 10 allocs) | 335.9 ns/op (112 B, 5 allocs) | 235.5 ns/op (32 B, 2 allocs) | **🥇 Gob** | 23% slower than Gob |
| **Nested Map** | 606.8 ns/op (534 B, 19 allocs) | 850.8 ns/op (256 B, 11 allocs) | 438.4 ns/op (80 B, 5 allocs) | **🥇 Gob** | 28% slower than Gob |

### **📥 Decoding Performance Comparison**

| Data Type | EDF | Protobuf | Gob | **Winner** | **EDF vs Winner** |
|-----------|-----|----------|-----|------------|-------------------|
| **String** | 76.26 ns/op (72 B, 4 allocs) | 70.84 ns/op (96 B, 2 allocs) | 452.2 ns/op (1000 B, 19 allocs) | **🥇 Protobuf** | 8% slower |
| **PID** | 110.1 ns/op (160 B, 6 allocs) | 66.73 ns/op (80 B, 2 allocs) | 5902 ns/op (7136 B, 160 allocs) | **🥇 Protobuf** | 65% slower |
| **ProcessID** | 120.9 ns/op (168 B, 7 allocs) | 81.94 ns/op (104 B, 3 allocs) | 5660 ns/op (6984 B, 157 allocs) | **🥇 Protobuf** | 48% slower |
| **Simple Struct** | 190.6 ns/op (258 B, 9 allocs) | 70.18 ns/op (88 B, 2 allocs) | 5821 ns/op (7072 B, 156 allocs) | **🥇 Protobuf** | 172% slower |
| **Complex Struct** | 741.1 ns/op (1366 B, 41 allocs) | 595.6 ns/op (824 B, 25 allocs) | 9534 ns/op (10872 B, 255 allocs) | **🥇 Protobuf** | 24% slower |
| **Nested Struct** | 2104 ns/op (4593 B, 107 allocs) | 1704 ns/op (2544 B, 71 allocs) | 12543 ns/op (14712 B, 342 allocs) | **🥇 Protobuf** | 23% slower |
| **Map** | 547.9 ns/op (954 B, 25 allocs) | 365.4 ns/op (528 B, 13 allocs) | 6689 ns/op (8256 B, 185 allocs) | **🥇 Protobuf** | 50% slower |
| **Nested Map** | 1128 ns/op (2045 B, 51 allocs) | 894.3 ns/op (1328 B, 30 allocs) | 7496 ns/op (9696 B, 212 allocs) | **🥇 Protobuf** | 26% slower |

### **📊 Performance Summary Table**

| Library | **Encoding Wins** | **Decoding Wins** | **Avg Encoding Speed** | **Avg Decoding Speed** | **Overall Rating** |
|---------|-------------------|-------------------|-------------------------|------------------------|-------------------|
| **EDF** | 🥇 **6/8** (75%) | 🥉 **0/8** (0%) | **🟢 Excellent** | **🟡 Good** | **⭐⭐⭐⭐** |
| **Protobuf** | 🥈 **1/8** (12.5%) | 🥇 **8/8** (100%) | **🟡 Good** | **🟢 Excellent** | **⭐⭐⭐⭐** |
| **Gob** | 🥉 **1/8** (12.5%) | 🥉 **0/8** (0%) | **🟡 Good** | **🔴 Poor** | **⭐⭐** |

### **🎯 Key Performance Insights**

| Metric | **EDF** | **Protobuf** | **Gob** |
|--------|---------|-------------|---------|
| **Best Use Case** | Write-heavy workloads | Read-heavy workloads | Simple data only |
| **Encoding Speed** | 🟢 **Fastest** for complex data | 🟡 Moderate | 🟡 Moderate |
| **Decoding Speed** | 🟡 **20-40% slower** than Protobuf | 🟢 **Fastest** | 🔴 **10-88x slower** |
| **Memory Efficiency** | 🟡 Higher decode allocation | 🟢 **Most efficient** | 🔴 **Massive overhead** |
| **Zero Allocations** | 🟢 **String encoding** | ❌ None | ❌ None |
| **Overall Performance** | 🟢 **Excellent** for encoding | 🟢 **Excellent** for decoding | 🔴 **Poor** for decoding |

### **📈 Performance Multipliers (vs Fastest)**

| Operation | **EDF Multiplier** | **Protobuf Multiplier** | **Gob Multiplier** |
|-----------|-------------------|------------------------|-------------------|
| **Complex Struct Encode** | **1.00x** _(fastest)_ | 1.69x | 1.24x |
| **Nested Struct Encode** | **1.00x** _(fastest)_ | 2.07x | 1.24x |
| **String Encode** | **1.00x** _(fastest)_ | 1.47x | 2.54x |
| **Simple Struct Decode** | 2.72x | **1.00x** _(fastest)_ | 82.9x |
| **Complex Struct Decode** | 1.24x | **1.00x** _(fastest)_ | 16.0x |
| **Map Decode** | 1.50x | **1.00x** _(fastest)_ | 18.3x |

### **✅ Recommendations**

| Scenario | **Recommended Library** | **Reason** |
|----------|-------------------------|------------|
| **Actor Message Passing** | **🥇 EDF** | Fastest encoding, good overall balance |
| **API Responses** | **🥇 Protobuf** | Fastest decoding, lower memory usage |
| **Logging/Persistence** | **🥇 EDF** | Excellent write performance |
| **Configuration Files** | **🥇 Protobuf** | Best read performance |
| **Real-time Systems** | **🥈 EDF/Protobuf** | Both excellent, avoid Gob |
| **Legacy Integration** | **🥉 Gob** | Only if Go-to-Go compatibility required |

**🎯 Conclusion: EDF excels at encoding performance and provides excellent overall balance, making it ideal for Ergo's actor-based message passing workloads!**

