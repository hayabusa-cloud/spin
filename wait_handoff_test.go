// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"code.hybscloud.com/spin"
)

// BenchmarkWaitHandoff measures one request and acknowledgment between two
// goroutines. On Linux, cpu-ns/op reports process user and system CPU time per
// round trip when the run spans at least one clock tick and the counter advances.
func BenchmarkWaitHandoff(b *testing.B) {
	var ticksPerSecond uint64
	if runtime.GOOS == "linux" {
		ticksPerSecond = linuxClockTicksPerSecond(b)
	}
	for _, procs := range []int{1, 2} {
		run := func(name string, wait func(*atomic.Bool, bool)) {
			b.Run(name+"/GOMAXPROCS="+itoa(procs), func(b *testing.B) {
				setBenchmarkProcs(b, procs)
				benchmarkHandoff(b, wait, ticksPerSecond)
			})
		}
		run("Wait", func(turn *atomic.Bool, want bool) {
			var wait spin.Wait
			for turn.Load() != want {
				wait.Once()
			}
		})
		run("Gosched", func(turn *atomic.Bool, want bool) {
			for turn.Load() != want {
				runtime.Gosched()
			}
		})
	}
}

func benchmarkHandoff(b *testing.B, wait func(*atomic.Bool, bool), ticksPerSecond uint64) {
	var turn atomic.Bool // true requests a handoff, false acknowledges it.
	ready := make(chan struct{})
	rounds := make(chan int)
	done := make(chan struct{})
	go func() {
		close(ready)
		for range <-rounds {
			wait(&turn, true)
			turn.Store(false)
		}
		close(done)
	}()
	<-ready
	defer func() {
		close(rounds)
		<-done
	}()

	var startTicks uint64
	if ticksPerSecond != 0 {
		startTicks = linuxProcessTicks(b)
	}
	rounds <- b.N
	b.ResetTimer()
	for range b.N {
		turn.Store(true)
		wait(&turn, false)
	}
	b.StopTimer()
	if ticksPerSecond == 0 || b.Elapsed() < time.Second/time.Duration(ticksPerSecond) {
		return
	}
	if ticks := linuxProcessTicks(b) - startTicks; ticks > 0 {
		nanos := float64(ticks) * float64(time.Second) / float64(ticksPerSecond)
		b.ReportMetric(nanos/float64(b.N), "cpu-ns/op")
	}
}
