# spin

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/spin.svg)](https://pkg.go.dev/code.hybscloud.com/spin)
[![Benchmarks](https://github.com/hayabusa-cloud/spin/actions/workflows/benchmarks.yml/badge.svg)](https://github.com/hayabusa-cloud/spin/actions/workflows/benchmarks.yml)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/spin/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/spin)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

用于在 Go 中处理短暂竞争的自旋锁和自适应等待原语。

[English](./README.md) | **简体中文** | [Español](./README.es.md) | [日本語](./README.ja.md) | [Français](./README.fr.md)

## 概述

`Lock` 在极短的临界区内保护共享状态。goroutine 可以用 `Wait` 轮询预计很快就会改变的条件。发生竞争时，两者起初只做短暂停顿，随着等待时间延长，会更频繁地将执行权交还给 Go 调度器。

## 安装

需要 Go 1.27 或更高版本。

```shell
go get code.hybscloud.com/spin
```

## 用法

所有写入方共享 `mu` 和两个计数器。统计一个 64 字节的数据包只需要一个很短的临界区。

```go
var mu spin.Lock
var packets, bytes int

mu.Lock()
packets++
bytes += 64
mu.Unlock()
```

如果本地工作预计很快就绪，在重试时复用同一个 `Wait`。调用方提供 `workReady` 来检查是否就绪。

```go
var wait spin.Wait
for !workReady() {
	wait.Once()
}
```

在重试次数有限的短轮询循环中，可在两次检查之间调用 `Pause`。

```go
spin.Pause()
```

如果下次重试前可以休眠一段时间，可调用 `Yield`，使用配置的休眠时长。

```go
spin.Yield()
```

## Lock 与 Wait

- `Lock` 是非公平锁，零值表示未加锁。`Try()` 尝试获取锁一次，成功时返回 `true`。通用互斥场景请使用 `sync.Mutex`。
- 每个轮询的 goroutine 使用自己的零值 `Wait`。`WillYield()` 表示下一次 `Once()` 是否会直接调用 `runtime.Gosched()`。`Reset()` 恢复初始的暂停间隔。

## Pause 与 Yield

- 在采用汇编实现的平台上，可以在紧密的轮询循环中调用 `Pause(n)`，将处理器提示指令重复执行 `n` 次。省略参数或传入非正数时，默认重复 30 次。实际耗时取决于处理器。
- `d` 为正时，`Yield(d)` 休眠 `d` 指定的时长，否则调用 `runtime.Gosched()`。`Yield()` 使用默认时长，初始值为 250 微秒，可通过 `SetYieldDuration(d)` 修改。

## 平台支持

`Pause` 在 `386`、`amd64`、`arm64`、`loong64`、`ppc64le`、`riscv64` 和 `s390x` 上采用汇编实现。在 `arm`、`wasm` 等其他目标平台上，每次调用 `Pause` 都会通过 `runtime.Gosched()` 让出一次执行权。

## 许可证

[MIT 许可证](./LICENSE)

©2025 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com)
