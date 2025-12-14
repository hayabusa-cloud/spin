[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/spin.svg)](https://pkg.go.dev/code.hybscloud.com/spin)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/spin)](https://goreportcard.com/report/github.com/hayabusa-cloud/spin)
[![Coverage Status](https://codecov.io/gh/hayabusa-cloud/spin/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/spin)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

# spin

Proporciona primitivas de spinlock (`Lock`) y espera activa (`Wait`) para escenarios sensibles al rendimiento.

- `Lock` — spinlock para secciones críticas extremadamente cortas
- `Wait` — espera activa adaptativa
- `Pause` — pista de CPU para bucles ajustados
- `Yield` — cesión cooperativa/sueño para rutas no calientes

Idiomas: [English](./README.md) | [简体中文](./README.zh-CN.md) | Español | [日本語](./README.ja.md) | [Français](./README.fr.md)

## Instalación

```shell
go get code.hybscloud.com/spin
```

## Inicio rápido

```go
func workReady() bool { return true }

func main() {
    var sl spin.Lock
    sl.Lock()
    // sección crítica
    sl.Unlock()
    fmt.Println("ok")

    var sw spin.Wait
    for !workReady() {
        sw.Once()
    }

    spin.Pause() // pista de CPU en un bucle caliente
    spin.Yield() // cesión cooperativa (ruta no caliente)
}
```

## Descripción de la API

- `type Lock`
  - `Lock()` gira hasta adquirir con retroceso adaptativo.
  - `Unlock()` libera.
  - `Try()` intenta adquirir sin esperar; devuelve `true` si tiene éxito.

- `type Wait`
  - `Once()` realiza un paso adaptativo (CPU `Pause` o cesión cooperativa), adecuado para bucles ajustados.
  - `WillYield()` indica si el próximo `Once()` cederá en lugar de pausar.
  - `Reset()` limpia contadores internos.

- `func Pause(cycles ...int)`
  - Emite una pista específica de la arquitectura y no debe bloquear ni ceder al planificador.

- `func Yield(duration ...time.Duration)`
  - Cede cooperativamente. Por defecto duerme un intervalo corto; si es no positivo, recurre a `runtime.Gosched()`.

- `func SetYieldDuration(d time.Duration)`
  - Establece la duración base usada por `Yield()` cuando no se proporciona argumento explícito.

Notas:
- `Lock` es no justo y no debe usarse como mutex de propósito general.
- Prefiere `Wait` en bucles de espera activa en lugar de `for {}` + `runtime.Gosched()` ad hoc.

Arquitecturas: amd64, arm64, 386, arm, riscv64, ppc64le, s390x, loong64, wasm.

## Cuándo usar

- Use `Lock` solo para secciones críticas extremadamente cortas; es no justo y pensado para escenarios especializados.
- Use `Wait` para espera activa adaptativa cuando se espera progreso muy pronto; escala de `Pause` a cesión cooperativa.
- Use `Pause` en bucles de sondeo calientes para reducir potencia y la contención entre hiperhilos.
- Use `Yield` en rutas no calientes para permitir cooperativamente que se ejecuten otras goroutines (opcionalmente durmiendo un corto intervalo).

Para exclusión mutua de propósito general, prefiera `sync.Mutex`.

## Licencia

MIT — ver [LICENSE](./LICENSE).

©2025 Hayabusa Cloud Co., Ltd.
