[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/spin.svg)](https://pkg.go.dev/code.hybscloud.com/spin)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/spin)](https://goreportcard.com/report/github.com/hayabusa-cloud/spin)
[![Coverage Status](https://codecov.io/gh/hayabusa-cloud/spin/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/spin)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

# spin

提供用于性能敏感场景的自旋锁（`Lock`）与自旋等待（`Wait`）原语。

- `Lock` — 面向极短临界区的自旋锁
- `Wait` — 自适应自旋等待
- `Pause` — 面向紧循环的 CPU 提示
- `Yield` — 非热点路径上的协作式让出/睡眠

语言: [English](./README.md) | 简体中文 | [Español](./README.es.md) | [日本語](./README.ja.md) | [Français](./README.fr.md)

## 安装

```shell
go get code.hybscloud.com/spin
```

## 快速开始

```go
func workReady() bool { return true }

func main() {
    var sl spin.Lock
    sl.Lock()
    // 临界区
    sl.Unlock()
    fmt.Println("ok")

    var sw spin.Wait
    for !workReady() {
        sw.Once()
    }

    spin.Pause() // 在热点循环里给出 CPU 提示
    spin.Yield() // 协作式让出（非热点路径）
}
```

## API 概览

- `type Lock`
  - `Lock()` 采用自适应回退自旋直到获取。
  - `Unlock()` 释放。
  - `Try()` 非阻塞尝试获取；成功返回 `true`。

- `type Wait`
  - `Once()` 执行一次自适应步骤（CPU `Pause` 或协作让出），适合紧循环。
  - `WillYield()` 指示下一次 `Once()` 是否会让出而非暂停。
  - `Reset()` 重置内部计数。

- `func Pause(cycles ...int)`
  - 发出与体系结构相关的 CPU 提示，不会阻塞或让出调度器。

- `func Yield(duration ...time.Duration)`
  - 协作式让出。默认睡眠一个很短的时间；若参数为非正值，则回退到 `runtime.Gosched()`。

- `func SetYieldDuration(d time.Duration)`
  - 设置当 `Yield()` 未显式提供参数时所用的基础睡眠时长。

注意：
- `Lock` 为非公平实现，不应作为通用互斥锁使用。
- 在自旋循环中优先使用 `Wait`，避免临时性的 `for {}` + `runtime.Gosched()` 方案。

支持架构：amd64、arm64、386、arm、riscv64、ppc64le、s390x、loong64、wasm。

## 适用场景

- `Lock`：仅用于极短临界区；该锁非公平，面向特定场景。
- `Wait`：当预期很快就能取得进展时用于自适应自旋等待；其会从 `Pause` 升级到协作式让出。
- `Pause`：用于热点轮询循环，降低功耗并减少超线程之间的竞争。
- `Yield`：用于非热点路径，协作式地让其他 goroutine 运行（可选地睡眠短时）。

通用互斥推荐使用 `sync.Mutex`。

## 许可证

MIT — 参见 [LICENSE](./LICENSE)。

©2025 Hayabusa Cloud Co., Ltd.
