// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"testing"

	"code.hybscloud.com/spin"
)

// Ensure that after sufficient progress the adaptive backoff enters
// a yield-heavy phase (WillYield is consistently true for a while),
// and Reset() returns it to the initial phase.
func TestWait_EventuallyAlwaysYieldThenReset(t *testing.T) {
	var sw spin.Wait

	// Drive the internal state forward with many steps.
	for i := 0; i < 2000; i++ {
		sw.Once()
	}

	// In the late phase, WillYield should be consistently true for a series
	// of checks (state doesn't change when only calling WillYield).
	for i := 0; i < 16; i++ {
		if !sw.WillYield() {
			t.Fatalf("expected WillYield=true in late phase at i=%d", i)
		}
	}

	// Reset should roll back to the initial phase where WillYield is false.
	sw.Reset()
	if sw.WillYield() {
		t.Fatalf("expected WillYield=false immediately after Reset")
	}
}

func BenchmarkWaitOnce_Baseline(b *testing.B) {
	var sw spin.Wait
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.Once()
	}
}

func BenchmarkWaitOnce_YieldHeavy(b *testing.B) {
	var sw spin.Wait
	// Move into yield-heavy phase before timing.
	for i := 0; i < 4000; i++ {
		sw.Once()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.Once()
	}
}
