- # 内存及性能测试

## 内存
### 存储空间
在 Redis 中，存储空间的大小并不是直接设置的，但你可以设置一些参数来控制 Redis 使用的最大内存和相应的行为。你可以使用 `maxmemory` 和 `maxmemory-policy` 配置选项来实现这一点。

1. **设置最大内存使用量**:

你可以使用 `maxmemory` 选项来设置 Redis 实例可以使用的最大内存量。

例如，要设置最大内存为 100MB:
```
maxmemory 100mb
```

2. **设置内存超出最大值时的行为**:

当内存使用超过 `maxmemory` 设置的值时，你可以使用 `maxmemory-policy` 选项来定义 Redis 如何行动。

可用的策略包括:
- `volatile-lru`: 从设置了过期时间的键中选择最近最少使用的键进行移除。
- `allkeys-lru`: 从所有键中选择最近最少使用的键进行移除。
- `volatile-random`: 从设置了过期时间的键中随机移除。
- `allkeys-random`: 从所有键中随机移除。
- `volatile-ttl`: 从设置了过期时间的键中选择将要过期的键进行移除。
- `noeviction`: 不移除任何键，只返回一个错误当添加新的数据并超过 `maxmemory` 设置的值。

例如，使用 LRU 算法从所有键中移除键:
```
maxmemory-policy allkeys-lru
```

3. **如何在运行时更改配置**:

你可以使用 Redis 的 `CONFIG` 命令在运行时更改这些配置。

例如:
```sh
CONFIG SET maxmemory 100mb
CONFIG SET maxmemory-policy allkeys-lru
```

请注意，如果你设置了一个较小的 `maxmemory` 值，而 Redis 实例的数据大小超过这个值，Redis 会开始根据你选择的策略删除键，直到内存使用量满足条件。因此，确保选择的 `maxmemory` 值和策略是符合你的需求的。

### 内存信息查看
```sh
> INFO memory
# Memory
used_memory:19279624
used_memory_human:18.39M
used_memory_rss:25473024
used_memory_rss_human:24.29M
used_memory_peak:19279624
used_memory_peak_human:18.39M
used_memory_peak_perc:100.26%
used_memory_overhead:1682394
used_memory_startup:796232
used_memory_dataset:17597230
used_memory_dataset_perc:95.21%
allocator_allocated:19399856
allocator_active:19779584
allocator_resident:22761472
total_system_memory:16742916096
total_system_memory_human:15.59G
used_memory_lua:41984
used_memory_lua_human:41.00K
used_memory_scripts:0
used_memory_scripts_human:0B
number_of_cached_scripts:0
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
allocator_frag_ratio:1.02
allocator_frag_bytes:379728
allocator_rss_ratio:1.15
allocator_rss_bytes:2981888
rss_overhead_ratio:1.12
rss_overhead_bytes:2711552
mem_fragmentation_ratio:1.32
mem_fragmentation_bytes:6243840
mem_not_counted_for_evict:0
mem_replication_backlog:0
mem_clients_slaves:0
mem_clients_normal:217842
mem_aof_buffer:0
mem_allocator:jemalloc-5.2.1
active_defrag_running:0
lazyfree_pending_objects:0
```

其中 `used_memory` 表示 Redis 分配的内存总量，单位是字节。`used_memory_human` 是一个更容易读懂的格式，它表示相同的信息，但是用 KB、MB 等单位。

此外，`maxmemory` 表示 Redis 的最大内存配置（如果有的话）。如果 Redis 达到这个限制，它将根据 `maxmemory_policy` 指定的策略来处理新的插入。如果你未设置最大内存，那么 `maxmemory` 的值将为0。

> 如果 Redis 的 `maxmemory` 没有被设置（默认为0），则 Redis 将无限制地使用内存，直到宿主机的物理内存耗尽。

### 内存测试
测试用例
```sh
package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"time"
)

var ctx = context.Background()

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root",
		DB:       1,
	})

	rdb.FlushDB(ctx)

	numKeys := 100000
	keySize := 10
	valueSize := 1024

	fmt.Println("start", time.Now().Format("2006-01-02 15:04:05"))
	for i := 0; i < numKeys; i++ {
		key := strconv.Itoa(i) + "-" + randomString(keySize)
		value := randomString(valueSize)
		err := rdb.Set(ctx, key, value, 0).Err()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("end", time.Now().Format("2006-01-02 15:04:05"))
}

```
没有设置`maxmemory`时，存储的值占用你的内存大小为:
- `10万`数据需要`8s` 占用内存`120M`  每个键占用 `1.06KB`  
- `100万`数据需要`1分15s` 占用内存`1.28G`  
- `1000万`数据需要`13分钟` 占用内存`12.79G`  


设置`maxmemory`为`1G`  
```sh
CONFIG SET maxmemory 100mb
CONFIG SET maxmemory-policy allkeys-lru
```

这时内存占用12G，设置了最大内存占用1G，redis日志输出警告
```sh
93944:M 27 Oct 2023 18:35:49.147 * Background saving started by pid 194533
194533:C 27 Oct 2023 18:37:18.304 * DB saved on disk
194533:C 27 Oct 2023 18:37:18.379 * RDB: 0 MB of memory used by copy-on-write
193944:M 27 Oct 2023 18:37:18.511 * Background saving terminated with success
193944:M 28 Oct 2023 10:30:17.094 # WARNING: the new maxmemory value set via CONFIG SET is smaller than the current memory usage. This will result in key eviction and/or the inability to accept new write commands depending on the maxmemory-policy.
```
> 警告：通过 CONFIG SET 设置的新 maxmemory 值小于当前内存使用量。 这将导致密钥驱逐和/或无法接受新的写入命令，具体取决于最大内存策略。  

现在已经无法通过redis-cli连接redis并下发指令了。如果放在`redis.conf`中，可以启动，但是无法创建新值（超过阈值的部分也不会删除）  
```sh
127.0.0.1:6379> set dd bb
(error) OOM command not allowed when used memory > 'maxmemory'.
```

默认配置:
```sh
# The default is:
#
# maxmemory-policy noeviction
```

那就清空数据，重新开始

设置内存为100M,存储10万数据。使用三种内存策略`volatile-lru`、`allkeys-lru`与`noeviction`  
> 所有键都没有设置过期时间，`volatile-lru` 什么现象？  

redis启动后查看状态
```sh
127.0.0.1:6379[1]> CONFIG GET maxmemory
1) "maxmemory"
2) "104857600"
127.0.0.1:6379[1]> CONFIG GET maxmemory-policy
1) "maxmemory-policy"
2) "volatile-lru"
```

`volatile-lru`模式时，内存到达上限后，直接报错:`panic: OOM command not allowed when used memory > 'maxmemory'.`  


```sh
127.0.0.1:6379[1]> CONFIG GET maxmemory
1) "maxmemory"
2) "104857600"
127.0.0.1:6379[1]> CONFIG SET maxmemory-policy allkeys-lru
OK

127.0.0.1:6379[1]> CONFIG GET maxmemory-policy 
1) "maxmemory-policy"
2) "allkeys-lru"
```

内存存储达到上限之后，会删除之前的数据，所以总数是固定的  
```sh
127.0.0.1:6379[1]> dbsize
(integer) 75656
127.0.0.1:6379[1]> dbsize
(integer) 75656
```

## 性能测试

Redis 的性能测试通常使用 `redis-benchmark` 工具，这是 Redis 官方提供的一个性能测试工具。使用 `redis-benchmark`，您可以测试 Redis 服务器的各种命令和数据结构的性能。

以下是如何使用 `redis-benchmark` 进行基本的性能测试，以及如何针对不同的数据类型进行测试的说明：

1. **安装 Redis**:
   如果您已经安装了 Redis，那么 `redis-benchmark` 工具应该已经包含在内。如果没有，您需要安装 Redis。

2. **基本测试**:
   在命令行中运行以下命令以执行默认的性能测试：
   ```bash
   redis-benchmark
   ```

3. **指定命令进行测试**:
   使用 `-t` 标志 followed by command names，可以测试特定的 Redis 命令。例如，要测试 `SET` 和 `GET` 命令，您可以使用：
   ```bash
   redis-benchmark -t set,get
   ```

   ```sh
	redis-benchmark -t set,get
	====== redis-benchmark -t set,get ======
	100000 requests completed in 0.49 seconds
	50 parallel clients
	3 bytes payload
	keep alive: 1

	99.99% <= 1 milliseconds
	100.00% <= 1 milliseconds
	204081.62 requests per second
   ```

4. **测试不同的数据类型**:
   为了测试不同的数据类型，您需要指定相应的 Redis 命令：

   - **Strings**: 使用 `set` 和 `get`
     ```bash
     redis-benchmark -t set,get
     ```

	 ```sh
     ====== SET ======
       100000 requests completed in 0.51 seconds
       50 parallel clients
       3 bytes payload
       keep alive: 1
     
     99.95% <= 1 milliseconds
     100.00% <= 1 milliseconds
     194552.53 requests per second
     
     ====== GET ======
       100000 requests completed in 0.55 seconds
       50 parallel clients
       3 bytes payload
       keep alive: 1
     
     100.00% <= 0 milliseconds
	 ```

   - **Hashes**: 使用 `hset` 和 `hget`
     ```bash
     redis-benchmark -t hset,hget
     ```

   - **Lists**: 使用 `lpush` 和 `lpop`
     ```bash
     redis-benchmark -t lpush,lpop
     ```

   - **Sets**: 使用 `sadd` 和 `spop`
     ```bash
     redis-benchmark -t sadd,spop
     ```

   - **Sorted Sets (Zsets)**: 使用 `zadd` 和 `zrange`
     ```bash
     redis-benchmark -t zadd,zrange
     ```

5. **自定义测试参数**:
   `redis-benchmark` 提供了许多选项，允许您自定义测试参数。例如，您可以指定并发客户端的数量、每个客户端要发送的请求数等。

   ```bash
   redis-benchmark -t set,get -c 100 -n 100000
   ```
   上述命令使用 100 个并发客户端，总共发送 100,000 个请求进行测试。

   ```sh
    ====== SET ======
      100000 requests completed in 0.55 seconds
      100 parallel clients
      3 bytes payload
      keep alive: 1
    
    99.87% <= 1 milliseconds
    99.90% <= 14 milliseconds
    99.94% <= 15 milliseconds
    100.00% <= 15 milliseconds
    182815.36 requests per second
    
    ====== GET ======
      100000 requests completed in 0.61 seconds
      100 parallel clients
      3 bytes payload
      keep alive: 1
    
    99.96% <= 1 milliseconds
    100.00% <= 1 milliseconds
    164744.64 requests per second
   ```

6. **查看所有选项**:
   要查看所有的 `redis-benchmark` 选项，您可以使用 `-h` 标志：
   ```bash
   redis-benchmark -h
   ```


### 平均性能
Redis 的读写性能受到许多因素的影响，包括硬件、网络、配置、数据大小和结构以及具体的操作。在一台典型的现代硬件上，Redis 可以达到十万到几十万的 QPS（每秒查询数）。

以下是一些常见的性能指标：

1. **在本地环境**（localhost）上进行测试时：
   - `PING` 测试的性能：约 200,000 QPS。
   - 简单的 `SET` 或 `GET` 操作：约 100,000 to 150,000 QPS。

2. **在网络环境**中进行测试时：
   - 根据网络延迟和带宽，性能可能会有所下降。
   - 在一般的 1 Gbps 网络环境中，性能可能在 50,000 - 80,000 QPS 之间。

3. **与数据大小和类型相关的性能**:
   - 小的字符串键/值对（例如，< 100 bytes）的性能通常最高。
   - 大型数据结构（例如，包含数千个元素的列表或散列）的性能会下降。
   - 执行复杂操作（如 `SORT`）也会影响性能。

4. **持久性选项**:
   - 如果您启用了 RDB 快照或 AOF 日志记录，这会对性能产生影响。具体的性能下降取决于快照或日志记录的频率。
   - 如果 AOF 以每秒 fsync 的方式配置，性能会受到较大影响。

5. **多线程**:
   - 在 Redis 6 及更高版本中，某些操作可以利用 I/O 线程并行化，这可以提高性能，尤其是在多核机器上。
