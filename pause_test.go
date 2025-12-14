// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"code.hybscloud.com/spin"
)

func TestPause(t *testing.T) {
	// Default (20 cycles)
	spin.Pause()

	// Single cycle
	spin.Pause(1)

	// Multiple iterations should not hang or panic
	for i := 0; i < 100; i++ {
		spin.Pause()
	}
}

func TestPauseWithCycles(t *testing.T) {
	cases := []int{-1, 1, 10, 20, 50, 100}
	for _, c := range cases {
		t.Run(fmt.Sprintf("cycles=%d", c), func(t *testing.T) {
			// Should not crash or hang
			spin.Pause(c)
		})
	}
}

func TestPauseInSpinLoop(t *testing.T) {
	var counter atomic.Int32
	done := make(chan struct{})

	go func() {
		time.Sleep(5 * time.Millisecond)
		counter.Store(1)
		close(done)
	}()

	for counter.Load() == 0 {
		spin.Pause()
	}

	<-done
	if counter.Load() != 1 {
		t.Fatalf("expected 1, got %d", counter.Load())
	}
}

func BenchmarkPause(b *testing.B) {
	b.Run("default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			spin.Pause()
		}
	})
	b.Run("1 cycle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			spin.Pause(1)
		}
	})
}

func BenchmarkSpinLoopWithPauseDefault(b *testing.B) {
	var counter atomic.Int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Store(0)
			go func() {
				runtime.Gosched()
				counter.Store(1)
			}()
			for counter.Load() == 0 {
				spin.Pause()
			}
		}
	})
}
