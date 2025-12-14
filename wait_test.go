// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"testing"
	"time"

	"code.hybscloud.com/spin"
)

func TestSpinWait(t *testing.T) {
	var sw spin.Wait
	for i := 0; i < 15; i++ {
		if sw.WillYield() {
			t.Fatalf("expected WillYield=false before threshold, got true at i=%d", i)
		}
		sw.Once()
	}
	if !sw.WillYield() {
		t.Fatalf("expected WillYield=true after threshold, got false")
	}
	sw.Reset()
	for i := 0; i < 15; i++ {
		if sw.WillYield() {
			t.Fatalf("expected WillYield=false after Reset before threshold")
		}
		sw.Once()
	}
	sw.Reset()
	if sw.WillYield() {
		t.Fatalf("expected WillYield=false immediately after Reset")
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

// Ensure we also execute the branch inside Once() where WillYield() is true.
func TestSpinWait_YieldBranch(t *testing.T) {
	var sw spin.Wait
	// Advance until a yield is expected
	for !sw.WillYield() {
		sw.Once()
	}
	// Now execute the yielding path
	sw.Once()
	// After yielding, it's still valid to continue spinning
	for i := 0; i < 3; i++ {
		sw.Once()
	}
}
