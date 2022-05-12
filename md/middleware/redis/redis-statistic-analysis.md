# redis 数据统计分析
常情况下，我们面临的用户数量以及访问量都是巨大的，比如百万、千万级别的用户数量，或者千万级别、甚至亿级别的访问信息。  
所以，我们必须要选择能够非常高效地统计大量数据（例如亿级）的集合类型。  
如何选择合适的数据集合，我们首先要了解常用的统计模式，并运用合理的数据来解决实际问题。   

四种统计类型：  

1. 二值状态统计；
2. 聚合统计；
3. 排序统计；
4. 基数统计。  

本文将用到 String、Set、Zset、List、hash 以外的拓展数据类型 Bitmap、HyperLogLog来实现。  

今天我们来看下剩下的三种统计类型。

文章涉及到的指令可以通过在线 Redis 客户端运行调试，地址：https://try.redis.io/，超方便的说。  

## 基数统计
`基数统计：统计一个集合中不重复元素的个数，常见于计算独立用户数（UV）`  

- PV(访问量)：即Page View, 即页面浏览量或点击量，用户每次刷新即被计算一次。
- UV(独立访客)：即Unique Visitor,访问您网站的一台电脑客户端为一个访客。00:00-24:00内相同的客户端只被计算一次。  
- IP(独立IP)：即Internet Protocol,指独立IP数。00:00-24:00内相同IP地址之被计算一次。  

实现基数统计最直接的方法，就是采用集合（Set）这种数据结构，当一个元素从未出现过时，便在集合中增加一个元素；如果出现过，那么集合仍保持不变。  

当页面访问量巨大，就需要一个超大的 Set 集合来统计，将会浪费大量空间。另外，这样的数据也**不需要很精确**，到底有没有更好的方案呢？  

这个问题问得好，`Redis` 提供了 `HyperLogLog` 数据结构就是用来解决种种场景的统计问题。  

`HyperLogLog` 是一种不精确的去重基数方案，它的统计规则是基于概率实现的，标准误差 `0.81%`，这样的精度足以满足 UV 统计需求了。  

### Set方案  

```shell
> sadd uv-set xiao
(integer) 1
> sadd uv-set ming
(integer) 1
> sadd uv-set hong
(integer) 1
> sadd uv-set lan
(integer) 1
> SMEMBERS uv-set
1) "xiao"
2) "lan"
3) "hong"
4) "ming"
```

### Hash 方案

`利用 Hash 类型实现，将用户 ID 作为 Hash 集合的 key，访问页面则执行 HSET 命令将 value 设置成 1。`  

即使用户重复访问，重复执行命令，也只会把这个 userId 的值设置成 “1"。

最后，利用 HLEN 命令统计 Hash 集合中的元素个数就是 UV。  


```
> hset uv-hset xiao:id5 1
1
> hset uv-hset ming:id5 1
1
> hset uv-hset hong:id7 1
1
> hlen uv-hset
3
> hkeys uv-hset
1) "xiao:id5"
2) "ming:id5"
3) "hong:id7"
```  

### HyperLogLog 方案  

利用  Redis 提供的 HyperLogLog 高级数据结构（不要只知道 Redis 的五种基础数据类型了）。这是一种用于基数统计的数据集合类型，即使数据量很大，计算基数需要的空间也是固定的。

每个 HyperLogLog 最多只需要花费 12KB 内存就可以计算 2 的 64 次方个元素的基数。

Redis 对 HyperLogLog 的存储进行了优化，在计数比较小的时候，存储空间采用系数矩阵，占用空间很小。

只有在计数很大，稀疏矩阵占用的空间超过了阈值才会转变成稠密矩阵，占用 12KB 空间。  

> 什么是基数?
> 比如数据集 {1, 3, 5, 7, 5, 7, 8}， 那么这个数据集的基数集为 {1, 3, 5 ,7, 8}, 基数(不重复元素)为5。 基数估计就是在误差可接受的范围内，快速计算基数。  

```shell
> PFADD hll1 foo bar zap a
(integer) 1
> PFADD hll2 a b c foo
(integer) 1
> PFMERGE hll3 hll1 hll2
"OK"
> PFCOUNT hll3
(integer) 6
```

<br>
<div align=center>
    <img src="../../../res/hyperlog对比.png" width="70%" height="70%" title="const与指针"></img>  
</div>
<br>  

> 三种方式的内存消耗对比。  











