# spin

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/spin.svg)](https://pkg.go.dev/code.hybscloud.com/spin)
[![Benchmarks](https://github.com/hayabusa-cloud/spin/actions/workflows/benchmarks.yml/badge.svg)](https://github.com/hayabusa-cloud/spin/actions/workflows/benchmarks.yml)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/spin/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/spin)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Bloqueos de espera activa y esperas adaptativas para situaciones de contención breve en Go.

[English](./README.md) | [简体中文](./README.zh-CN.md) | **Español** | [日本語](./README.ja.md) | [Français](./README.fr.md)

## Descripción general

`Lock` protege el estado compartido en secciones críticas muy breves. Una goroutine puede usar `Wait` para consultar una condición que espera que cambie pronto. Cuando hay contención, ambos empiezan con pausas breves y ceden la ejecución al planificador de Go con mayor frecuencia a medida que se prolonga la espera.

## Instalación

Requiere Go 1.27 o posterior.

```shell
go get code.hybscloud.com/spin
```

## Uso

Las goroutines que escriben comparten `mu` y ambos contadores. Basta una sección crítica breve para registrar un paquete de 64 bytes.

```go
var mu spin.Lock
var packets, bytes int

mu.Lock()
packets++
bytes += 64
mu.Unlock()
```

Si se espera que haya trabajo disponible en breve, use el mismo `Wait` en todos los reintentos. La aplicación proporciona `workReady` para comprobar si ya se puede continuar.

```go
var wait spin.Wait
for !workReady() {
	wait.Once()
}
```

En un bucle de sondeo breve con un número limitado de intentos, use `Pause` entre comprobaciones.

```go
spin.Pause()
```

Use `Yield` si puede esperar el intervalo configurado antes de volver a intentarlo.

```go
spin.Yield()
```

## Lock y Wait

- `Lock` no garantiza un orden de adquisición y su valor cero está desbloqueado. `Try()` intenta adquirir el bloqueo una sola vez y devuelve `true` si lo consigue. Para la exclusión mutua de uso general, utilice `sync.Mutex`.
- Cada goroutine que espera usa su propio `Wait`, inicializado con el valor cero. `WillYield()` indica si la siguiente llamada a `Once()` invocará directamente `runtime.Gosched()`. `Reset()` restablece el intervalo de pausa inicial.

## Pause y Yield

- En las arquitecturas con implementación en ensamblador, `Pause(n)` permite repetir `n` veces una instrucción de pausa del procesador dentro de un bucle de sondeo. Si se omite el argumento o se pasa un valor menor o igual que cero, se realizan 30 repeticiones. La duración depende del procesador.
- Si `d` es positivo, `Yield(d)` suspende la goroutine durante `d`. En caso contrario, llama a `runtime.Gosched()`. `Yield()` usa la duración predeterminada, inicialmente de 250 µs. `SetYieldDuration(d)` permite cambiarla.

## Plataformas compatibles

`Pause` está implementado en ensamblador para `386`, `amd64`, `arm64`, `loong64`, `ppc64le`, `riscv64` y `s390x`. En las demás arquitecturas, incluidas `arm` y `wasm`, cada llamada a `Pause` cede la ejecución una vez mediante `runtime.Gosched()`.

## Licencia

[Licencia MIT](./LICENSE)

©2025 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com)
