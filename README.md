[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/spin.svg)](https://pkg.go.dev/code.hybscloud.com/spin)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/spin)](https://goreportcard.com/report/github.com/hayabusa-cloud/spin)
[![Coverage Status](https://codecov.io/gh/hayabusa-cloud/spin/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/spin)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

# spin

Provides spinlock(`Lock`) and spin-wait(`Wait`) primitives for performance-sensitive spinning.

- `Lock` — spinlock for extremely short critical sections
- `Wait` — spin wait for adaptive spinning
- `Pause` — CPU hint for tight loops
- `Yield` — cooperative yield/sleep for non-hot paths

Languages: English | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | [日本語](./README.ja.md) | [Français](./README.fr.md)

## Installation

```shell
go get code.hybscloud.com/spin
```

## Quick start

```go
func workReady() bool { return true }

func main() {
    var sl spin.Lock
    sl.Lock()
    // critical section
    sl.Unlock()
    fmt.Println("ok")

    var sw spin.Wait
    for !workReady() {
        sw.Once()
    }

    spin.Pause() // CPU hint in a hot loop
    spin.Yield() // cooperative yield (non-hot path)
}
```

## API overview

- `type Lock`
  - `Lock()` spins until acquired with adaptive backoff.
  - `Unlock()` releases.
  - `Try()` attempts to acquire without waiting; returns true on success.

- `type Wait`
  - `Once()` performs one adaptive step (CPU `Pause` or cooperative yield), suitable for tight loops.
  - `WillYield()` reports whether the next `Once()` will yield instead of pausing.
  - `Reset()` clears internal counters.

- `func Pause(cycles ...int)`
  - Issues an architecture-specific CPU hint and must not block or yield the scheduler.

- `func Yield(duration ...time.Duration)`
  - Cooperatively yields. By default sleeps for a small duration; if non-positive, falls back to `runtime.Gosched()`.

- `func SetYieldDuration(d time.Duration)`
  - Sets the base sleep duration used by `Yield()` when no explicit argument is provided.

Notes:
- `Lock` is non-fair and should not be used as a general-purpose mutex.
- Prefer `Wait` in spin loops instead of ad-hoc `for {}` + `runtime.Gosched()`.

Architectures: amd64, arm64, 386, arm, riscv64, ppc64le, s390x, loong64, wasm.

## When to use

- Use `Lock` only for extremely short critical sections; it is non-fair and intended for specialized scenarios.
- Use `Wait` for adaptive spin-waiting when progress is expected very soon; it escalates from `Pause` to cooperative yield.
- Use `Pause` in tight polling loops on hot paths to reduce power and contention between hyper-threads.
- Use `Yield` in non-hot paths to cooperatively let other goroutines run (optionally sleeping a short duration).

For general-purpose mutual exclusion, prefer `sync.Mutex`.

## License

MIT — see [LICENSE](./LICENSE).

©2025 Hayabusa Cloud Co., Ltd.
