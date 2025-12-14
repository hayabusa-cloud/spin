[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/spin.svg)](https://pkg.go.dev/code.hybscloud.com/spin)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/spin)](https://goreportcard.com/report/github.com/hayabusa-cloud/spin)
[![Coverage Status](https://codecov.io/gh/hayabusa-cloud/spin/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/spin)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

# spin

高性能・低レイテンシ向けに、スピンロック（`Lock`）とスピン待機（`Wait`）のプリミティブを提供します。

- `Lock` — 極短いクリティカルセクション向けのスピンロック
- `Wait` — 適応型スピン待機
- `Pause` — タイトループ向けの CPU ヒント
- `Yield` — 非ホットパス向けの協調的譲歩/スリープ

言語: [English](./README.md) | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | 日本語 | [Français](./README.fr.md)

## インストール

```shell
go get code.hybscloud.com/spin
```

## クイックスタート

```go
func workReady() bool { return true }

func main() {
    var sl spin.Lock
    sl.Lock()
    // クリティカルセクション
    sl.Unlock()
    fmt.Println("ok")

    var sw spin.Wait
    for !workReady() {
        sw.Once()
    }

    spin.Pause() // ホットループでの CPU ヒント
    spin.Yield() // 協調的譲歩（非ホットパス）
}
```

## API 概要

- `type Lock`
  - `Lock()` 適応的バックオフで取得できるまでスピンします。
  - `Unlock()` 解放します。
  - `Try()` 待たずに取得を試みます。成功なら `true` を返します。

- `type Wait`
  - `Once()` 1 回の適応ステップ（CPU `Pause` または協調的譲歩）を実行します。タイトループに適します。
  - `WillYield()` 次の `Once()` が pause ではなく譲歩するかどうかを示します。
  - `Reset()` 内部カウンタをクリアします。

- `func Pause(cycles ...int)`
  - アーキテクチャ依存の CPU ヒントを発行し、スケジューラをブロック/譲歩しません。

- `func Yield(duration ...time.Duration)`
  - 協調的に譲歩します。既定では短時間 sleep し、非正値なら `runtime.Gosched()` にフォールバックします。

- `func SetYieldDuration(d time.Duration)`
  - `Yield()` に明示引数がない場合の基準 sleep 時間を設定します。

注意:
- `Lock` は非フェアであり、汎用ミューテックスとして使うべきではありません。
- スピンループではアドホックな `for {}` + `runtime.Gosched()` より `Wait` を優先してください。

対応アーキテクチャ: amd64, arm64, 386, arm, riscv64, ppc64le, s390x, loong64, wasm.

## 使いどころ

- `Lock` は極短いクリティカルセクション専用。非フェアで、特殊用途を想定しています。
- `Wait` は近い将来に進捗が見込めるときの適応型スピン待機に。`Pause` から協調的譲歩へ段階的にエスカレートします。
- `Pause` はホットなポーリングループで使用し、消費電力とハイパースレッド間の競合を抑えます。
- `Yield` は非ホットパスで使用し、他の goroutine に実行機会を協調的に与えます（短時間の sleep を伴う場合があります）。

汎用の相互排他には `sync.Mutex` を推奨します。

## ライセンス

MIT — [LICENSE](./LICENSE) を参照してください。

©2025 Hayabusa Cloud Co., Ltd.
