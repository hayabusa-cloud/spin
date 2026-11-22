// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"runtime"
	"sync"
	"testing"

	"code.hybscloud.com/spin"
)

// ----------------------------------------------------------------------------
// Lock Benchmarks
// ----------------------------------------------------------------------------

// BenchmarkLock measures lock performance across contention levels.
func BenchmarkLock(b *testing.B) {
	b.Run("Uncontended", func(b *testing.B) {
		var lk spin.Lock
		for range b.N {
			lk.Lock()
			lk.Unlock()
		}
	})

	b.Run("Try/Uncontended", func(b *testing.B) {
		var lk spin.Lock
		for range b.N {
			if lk.Try() {
				lk.Unlock()
			}
		}
	})

	// Contention scaling benchmarks
	for _, procs := range []int{2, 4, 8, 16} {
		b.Run("Contended/GOMAXPROCS="+itoa(procs), func(b *testing.B) {
			setBenchmarkProcs(b, procs)

			var lk spin.Lock
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					lk.Lock()
					lk.Unlock()
				}
			})
		})
	}
}

// BenchmarkLockVsMutex compares spin.Lock against sync.Mutex.
func BenchmarkLockVsMutex(b *testing.B) {
	b.Run("spin.Lock/Uncontended", func(b *testing.B) {
		var lk spin.Lock
		for range b.N {
			lk.Lock()
			lk.Unlock()
		}
	})

	b.Run("sync.Mutex/Uncontended", func(b *testing.B) {
		var mu sync.Mutex
		for range b.N {
			mu.Lock()
			mu.Unlock()
		}
	})

	b.Run("spin.Lock/Contended", func(b *testing.B) {
		var lk spin.Lock
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				lk.Lock()
				lk.Unlock()
			}
		})
	})

	b.Run("sync.Mutex/Contended", func(b *testing.B) {
		var mu sync.Mutex
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				mu.Lock()
				mu.Unlock()
			}
		})
	})
}

// BenchmarkLockWithWork records 64-byte packets with two Go processors.
func BenchmarkLockWithWork(b *testing.B) {
	const procs = 2
	type counters struct{ packets, bytes uint64 }
	for _, workers := range []int{2, 8} {
		run := func(name string, record func(*testing.PB, *counters)) {
			b.Run(name+"/GOMAXPROCS=2/workers="+itoa(workers), func(b *testing.B) {
				setBenchmarkProcs(b, procs)
				b.SetParallelism(workers / procs)
				var total counters
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) { record(pb, &total) })
				b.StopTimer()
				if total.packets != uint64(b.N) || total.bytes != 64*uint64(b.N) {
					b.Fatalf("recorded %d packets and %d bytes for %d updates", total.packets, total.bytes, b.N)
				}
			})
		}

		var lk spin.Lock
		run("spin.Lock", func(pb *testing.PB, total *counters) {
			for pb.Next() {
				lk.Lock()
				total.packets++
				total.bytes += 64
				lk.Unlock()
			}
		})

		var mu sync.Mutex
		run("sync.Mutex", func(pb *testing.PB, total *counters) {
			for pb.Next() {
				mu.Lock()
				total.packets++
				total.bytes += 64
				mu.Unlock()
			}
		})
	}
}

// ----------------------------------------------------------------------------
// Wait Benchmarks
// ----------------------------------------------------------------------------

// BenchmarkWait measures Wait performance.
func BenchmarkWait(b *testing.B) {
	b.Run("Once/Fresh", func(b *testing.B) {
		for range b.N {
			var sw spin.Wait
			sw.Once()
		}
	})

	b.Run("Once/Repeated", func(b *testing.B) {
		var sw spin.Wait
		for range b.N {
			sw.Once()
		}
	})

	b.Run("WillYield", func(b *testing.B) {
		// Record states from every yield interval before timing the prediction.
		var states [128]spin.Wait
		var sw spin.Wait
		for i := range states {
			states[i] = sw
			sw.Once()
		}

		predicted := 0
		b.ResetTimer()
		for i := range b.N {
			if states[i%len(states)].WillYield() {
				predicted++
			}
		}
		b.StopTimer()
		b.ReportMetric(float64(predicted)/float64(b.N), "predicted-yields/op")
	})

	b.Run("Reset", func(b *testing.B) {
		var sw spin.Wait
		for i := range b.N {
			if i%16 == 0 {
				sw.Reset()
			}
			sw.Once()
		}
	})
}

// ----------------------------------------------------------------------------
// Pause Benchmarks
// ----------------------------------------------------------------------------

// BenchmarkPauseCycles measures Pause across different cycle counts.
func BenchmarkPauseCycles(b *testing.B) {
	for _, cycles := range []int{1, 10, 30, 50, 100} {
		b.Run("cycles="+itoa(cycles), func(b *testing.B) {
			for range b.N {
				spin.Pause(cycles)
			}
		})
	}
}

// BenchmarkPauseLoop measures pause in a tight loop (common pattern).
func BenchmarkPauseLoop(b *testing.B) {
	b.Run("10_iterations", func(b *testing.B) {
		for range b.N {
			for range 10 {
				spin.Pause()
			}
		}
	})

	b.Run("100_iterations", func(b *testing.B) {
		for range b.N {
			for range 100 {
				spin.Pause()
			}
		}
	})
}

// ----------------------------------------------------------------------------
// Spin-Wait Pattern Benchmarks
// ----------------------------------------------------------------------------

// BenchmarkSpinWaitPattern measures realistic spin-wait scenarios.
func BenchmarkSpinWaitPattern(b *testing.B) {
	// Pattern 1: Spin until condition (Wait helper)
	b.Run("Wait/UntilCondition", func(b *testing.B) {
		for range b.N {
			var sw spin.Wait
			for range 16 {
				sw.Once()
			}
		}
	})

	// Pattern 2: Raw pause loop
	b.Run("Pause/RawLoop", func(b *testing.B) {
		for range b.N {
			for range 16 {
				spin.Pause()
			}
		}
	})
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func setBenchmarkProcs(b *testing.B, procs int) {
	b.Helper()
	previous := runtime.GOMAXPROCS(procs)
	b.Cleanup(func() { runtime.GOMAXPROCS(previous) })
	if got := runtime.GOMAXPROCS(0); got != procs {
		b.Skipf("requested GOMAXPROCS=%d, runtime permits %d", procs, got)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
