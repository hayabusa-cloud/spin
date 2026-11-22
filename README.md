# spin

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/spin.svg)](https://pkg.go.dev/code.hybscloud.com/spin)
[![Benchmarks](https://github.com/hayabusa-cloud/spin/actions/workflows/benchmarks.yml/badge.svg)](https://github.com/hayabusa-cloud/spin/actions/workflows/benchmarks.yml)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/spin/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/spin)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Spin locks and adaptive waits for short-lived contention in Go.

**English** | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | [日本語](./README.ja.md) | [Français](./README.fr.md)

## Overview

`Lock` protects shared state in very short critical sections. `Wait` lets a goroutine poll a condition it expects to change soon. Under contention, both pause briefly at first and yield to the Go scheduler more often as the wait continues.

## Installation

Requires Go 1.27 or later.

```shell
go get code.hybscloud.com/spin
```

## Usage

All writers share `mu` and both counters. Recording a 64-byte packet takes one short critical section.

```go
var mu spin.Lock
var packets, bytes int

mu.Lock()
packets++
bytes += 64
mu.Unlock()
```

When local work should be ready soon, keep one `Wait` across retries. The caller supplies the `workReady` check.

```go
var wait spin.Wait
for !workReady() {
	wait.Once()
}
```

Use `Pause` between checks in a short, bounded polling loop.

```go
spin.Pause()
```

Use `Yield` when the next retry can wait for the configured sleep interval.

```go
spin.Yield()
```

## Lock and Wait

- `Lock` is non-fair and its zero value is unlocked. `Try()` attempts to acquire it once and returns `true` on success. For general-purpose locking, use `sync.Mutex`.
- Each polling goroutine uses its own zero-value `Wait`. `WillYield()` reports whether the next `Once()` will call `runtime.Gosched()` directly. `Reset()` restores the initial pause interval.

## Pause and Yield

- On targets with an assembly implementation, `Pause(n)` repeats a processor hint `n` times in a tight polling loop. An omitted or non-positive count selects 30 repetitions. The elapsed time depends on the processor.
- `Yield(d)` sleeps for `d` if the duration is positive and calls `runtime.Gosched()` otherwise. `Yield()` uses the default duration, initially 250 µs. `SetYieldDuration(d)` changes that default.

## Platform support

`Pause` uses assembly on `386`, `amd64`, `arm64`, `loong64`, `ppc64le`, `riscv64`, and `s390x`. On other targets, including `arm` and `wasm`, each call to `Pause` yields once through `runtime.Gosched()`.

## License

[MIT License](./LICENSE)

©2025 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com)
