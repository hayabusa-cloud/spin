// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin

import (
	"runtime"
	"time"
)

var yieldDuration = 250 * time.Microsecond

// Yield cooperatively yields execution to reduce CPU burn in tight loops.
//
// By default, Yield sleeps for a small duration (configured by SetYieldDuration)
// to give other goroutines and the OS scheduler a chance to run.
//
// If an explicit duration is provided, it overrides the default for this call.
// A non-positive duration disables sleeping and falls back to runtime.Gosched()
// (a purely cooperative yield without a timer sleep).
//
// For automatic adaptive backoff in tight loops, use Wait instead.
func Yield(duration ...time.Duration) {
	d := yieldDuration
	if len(duration) > 0 {
		d = max(0, duration[0])
	}
	if d > 0 {
		time.Sleep(d)
		return
	}
	runtime.Gosched()
}

// SetYieldDuration sets the base duration unit for Yield().
// Recommended: 50-250 microseconds for real-time systems, 1-4 ms for general workloads.
func SetYieldDuration(d time.Duration) {
	if d < 0 {
		d = 0
	}
	yieldDuration = d
}
