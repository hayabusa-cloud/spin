# spin

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/spin.svg)](https://pkg.go.dev/code.hybscloud.com/spin)
[![Benchmarks](https://github.com/hayabusa-cloud/spin/actions/workflows/benchmarks.yml/badge.svg)](https://github.com/hayabusa-cloud/spin/actions/workflows/benchmarks.yml)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/spin/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/spin)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Go で短時間の競合を扱うためのスピンロックと適応型スピン待機を提供します。

[English](./README.md) | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | **日本語** | [Français](./README.fr.md)

## 概要

`Lock` はごく短いクリティカルセクションで共有状態を保護します。`Wait` は、間もなく変化すると見込まれる条件を goroutine が繰り返し確認するときに使います。競合が起きると、どちらも最初は短いポーズを挟み、待機が長引くにつれて Go のスケジューラに実行権を譲る頻度が上がります。

## インストール

Go 1.27 以降が必要です。

```shell
go get code.hybscloud.com/spin
```

## 使い方

書き込み側の goroutine はすべて `mu` と 2 つのカウンターを共有します。64 バイトのパケットを 1 個数える処理は、短いクリティカルセクションに収まります。

```go
var mu spin.Lock
var packets, bytes int

mu.Lock()
packets++
bytes += 64
mu.Unlock()
```

処理の準備がすぐに整う見込みなら、再試行のたびに同じ `Wait` を使います。準備ができたかどうかを調べる `workReady` は呼び出し元が用意します。

```go
var wait spin.Wait
for !workReady() {
	wait.Once()
}
```

回数を決めて短くポーリングする場合は、次の確認の前に `Pause` を呼びます。

```go
spin.Pause()
```

次の再試行まで設定された時間だけスリープできる場合は、`Yield` を使います。

```go
spin.Yield()
```

## Lock と Wait

- `Lock` は取得順を保証せず、ゼロ値は未ロック状態です。`Try()` はロックの取得を一度だけ試み、成功すると `true` を返します。汎用の排他制御には `sync.Mutex` を使います。
- ポーリングする goroutine ごとに、ゼロ値の `Wait` を用意します。`WillYield()` は、次の `Once()` が `runtime.Gosched()` を直接呼ぶかどうかを返します。`Reset()` を呼ぶと、ポーズの間隔が初期値に戻ります。

## Pause と Yield

- アセンブリ実装を持つアーキテクチャでは、ポーリングループの中で `Pause(n)` を呼ぶと、プロセッサへのヒント命令を `n` 回繰り返します。引数を省略するか 0 以下の値を指定すると、繰り返し回数は 30 回になります。所要時間はプロセッサによって異なります。
- `d` が正なら `Yield(d)` は `d` の間スリープし、それ以外では `runtime.Gosched()` を呼びます。`Yield()` が使うデフォルトの待機時間は初期状態で 250 マイクロ秒で、`SetYieldDuration(d)` で変更できます。

## 対応プラットフォーム

`386`、`amd64`、`arm64`、`loong64`、`ppc64le`、`riscv64`、`s390x` では、`Pause` はアセンブリで実装されています。`arm` や `wasm` を含む他のターゲットでは、`Pause` を呼ぶたびに `runtime.Gosched()` で一度だけ実行権を譲ります。

## ライセンス

[MIT ライセンス](./LICENSE)

©2025 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com)
