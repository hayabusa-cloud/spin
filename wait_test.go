// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"testing"
	"time"

	"code.hybscloud.com/spin"
)

func TestWaitYieldSchedule(t *testing.T) {
	// Four yields at each interval lead to a sustained yield on every call.
	phases := []struct{ interval, calls int }{
		{16, 64}, {8, 32}, {4, 16}, {2, 8}, {1, 2000},
	}
	var wait spin.Wait
	for pass := range 2 {
		for _, phase := range phases {
			for step := 1; step <= phase.calls; step++ {
				want := step%phase.interval == 0
				if got := wait.WillYield(); got != want {
					t.Fatalf("pass %d, interval %d, step %d: WillYield() = %t, want %t", pass+1, phase.interval, step, got, want)
				}
				wait.Once()
			}
		}
		wait.Reset()
	}
}

func TestWaitResetBeforeFirstYield(t *testing.T) {
	var wait spin.Wait
	for range 15 {
		wait.Once()
	}
	wait.Reset()
	if wait.WillYield() {
		t.Fatal("first step after Reset should pause")
	}
}

func TestYield(t *testing.T) {
	spin.Yield()
	spin.Yield(-1)
	spin.Yield(3)
}

func TestSetYieldDuration(t *testing.T) {
	spin.SetYieldDuration(100 * time.Microsecond)
	spin.Yield()
	spin.SetYieldDuration(0)
	spin.Yield(5)
	spin.SetYieldDuration(-time.Millisecond)
	spin.Yield()
}
