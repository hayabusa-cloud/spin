# spin

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/spin.svg)](https://pkg.go.dev/code.hybscloud.com/spin)
[![Benchmarks](https://github.com/hayabusa-cloud/spin/actions/workflows/benchmarks.yml/badge.svg)](https://github.com/hayabusa-cloud/spin/actions/workflows/benchmarks.yml)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/spin/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/spin)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Verrous tournants et attente adaptative en cas de contention de courte durée en Go.

[English](./README.md) | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | [日本語](./README.ja.md) | **Français**

## Vue d’ensemble

`Lock` protège l’état partagé dans des sections critiques très brèves. Une goroutine peut utiliser `Wait` pour vérifier une condition qui devrait bientôt changer. En cas de contention, les deux commencent par de brèves pauses et rendent la main à l’ordonnanceur Go de plus en plus souvent à mesure que l’attente se prolonge.

## Installation

Nécessite Go 1.27 ou une version ultérieure.

```shell
go get code.hybscloud.com/spin
```

## Utilisation

Les goroutines qui écrivent partagent `mu` et les deux compteurs. Une courte section critique suffit pour comptabiliser un paquet de 64 octets.

```go
var mu spin.Lock
var packets, bytes int

mu.Lock()
packets++
bytes += 64
mu.Unlock()
```

Si du travail doit être disponible sous peu, réutilisez le même `Wait` à chaque tentative. L’appelant fournit la fonction `workReady` pour vérifier si le travail est prêt.

```go
var wait spin.Wait
for !workReady() {
	wait.Once()
}
```

Dans une courte boucle de scrutation dont le nombre d’essais est limité, appelez `Pause` entre deux vérifications.

```go
spin.Pause()
```

Utilisez `Yield` si la prochaine tentative peut attendre le délai configuré.

```go
spin.Yield()
```

## Lock et Wait

- `Lock` ne garantit aucun ordre d’acquisition et sa valeur zéro correspond à un verrou libre. `Try()` tente une seule fois de prendre le verrou et renvoie `true` en cas de succès. Pour l’exclusion mutuelle d’usage général, utilisez `sync.Mutex`.
- Chaque goroutine qui attend utilise son propre `Wait`, initialisé à sa valeur zéro. `WillYield()` indique si le prochain appel à `Once()` appellera directement `runtime.Gosched()`. `Reset()` rétablit l’intervalle de pause initial.

## Pause et Yield

- Sur les architectures disposant d’une implémentation en assembleur, `Pause(n)` permet de répéter `n` fois une instruction de pause du processeur dans une boucle de scrutation. Si l’argument est omis, nul ou négatif, le nombre de répétitions vaut 30. La durée dépend du processeur.
- Si `d` est positif, `Yield(d)` suspend la goroutine pendant `d`. Sinon, la fonction appelle `runtime.Gosched()`. `Yield()` utilise la durée par défaut, initialement de 250 µs et modifiable avec `SetYieldDuration(d)`.

## Plateformes prises en charge

`Pause` est implémentée en assembleur sur `386`, `amd64`, `arm64`, `loong64`, `ppc64le`, `riscv64` et `s390x`. Sur les autres architectures, dont `arm` et `wasm`, chaque appel à `Pause` rend la main une fois en appelant `runtime.Gosched()`.

## Licence

[Licence MIT](./LICENSE)

©2025 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com)
