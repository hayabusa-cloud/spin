[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/spin.svg)](https://pkg.go.dev/code.hybscloud.com/spin)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/spin)](https://goreportcard.com/report/github.com/hayabusa-cloud/spin)
[![Coverage Status](https://codecov.io/gh/hayabusa-cloud/spin/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/spin)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

# spin

Fournit des primitives de spinlock (`Lock`) et d’attente active (`Wait`) pour des chemins sensibles aux performances.

- `Lock` — spinlock pour des sections critiques extrêmement courtes
- `Wait` — attente active adaptative
- `Pause` — indice CPU pour les boucles serrées
- `Yield` — cession coopérative/sommeil pour les chemins non chauds

Langues : [English](./README.md) | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | [日本語](./README.ja.md) | Français

## Installation

```shell
go get code.hybscloud.com/spin
```

## Démarrage rapide

```go
func workReady() bool { return true }

func main() {
    var sl spin.Lock
    sl.Lock()
    // section critique
    sl.Unlock()
    fmt.Println("ok")

    var sw spin.Wait
    for !workReady() {
        sw.Once()
    }

    spin.Pause() // indice CPU dans une boucle chaude
    spin.Yield() // cession coopérative (chemin non chaud)
}
```

## Vue d’ensemble de l’API

- `type Lock`
  - `Lock()` tourne jusqu’à acquisition avec repli adaptatif.
  - `Unlock()` libère.
  - `Try()` tente d’acquérir sans attendre ; renvoie `true` en cas de succès.

- `type Wait`
  - `Once()` effectue une étape adaptative (CPU `Pause` ou cession coopérative), adaptée aux boucles serrées.
  - `WillYield()` indique si le prochain `Once()` cèdera au lieu de pauser.
  - `Reset()` remet à zéro les compteurs internes.

- `func Pause(cycles ...int)`
  - Émet un indice CPU spécifique à l’architecture et ne doit ni bloquer ni céder au planificateur.

- `func Yield(duration ...time.Duration)`
  - Cède de manière coopérative. Par défaut dort une courte durée ; si la valeur est non positive, se rabat sur `runtime.Gosched()`.

- `func SetYieldDuration(d time.Duration)`
  - Définit la durée de sommeil de base utilisée par `Yield()` quand aucun argument explicite n’est fourni.

Notes :
- `Lock` est non équitable et ne doit pas être utilisé comme mutex générique.
- Préférez `Wait` dans les boucles d’attente active plutôt que `for {}` + `runtime.Gosched()` ad hoc.

Architectures : amd64, arm64, 386, arm, riscv64, ppc64le, s390x, loong64, wasm.

## Quand l’utiliser

- Utilisez `Lock` uniquement pour des sections critiques extrêmement courtes ; il est non équitable et prévu pour des scénarios spécialisés.
- Utilisez `Wait` pour l’attente active adaptative lorsque l’on s’attend à un progrès imminent ; escalade de `Pause` vers la cession coopérative.
- Utilisez `Pause` dans des boucles de sondage chaudes pour réduire la consommation et la contention entre hyper‑threads.
- Utilisez `Yield` sur des chemins non chauds pour laisser coopérativement s’exécuter d’autres goroutines (en dormant éventuellement un court instant).

Pour l’exclusion mutuelle générique, préférez `sync.Mutex`.

## Licence

MIT — voir [LICENSE](./LICENSE).

©2025 Hayabusa Cloud Co., Ltd.
