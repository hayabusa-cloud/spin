// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin

import "runtime"

// Wait is a lightweight adaptive spin-wait helper used in tight polling loops.
//
// It escalates from a short CPU pause to a cooperative scheduler yield
// based on the recent history of calls. The zero value is ready to use.
type Wait struct {
	counter, n uint32
}

// Once performs a single adaptive step in the spin-wait strategy:
// it either issues a light `Pause()` or yields the scheduler via `runtime.Gosched()`.
func (s *Wait) Once() {
	y := s.WillYield()
	s.counter++
	if y {
		s.n++
		runtime.Gosched()
		return
	}
	Pause(defaultPauseCycles)
}

// WillYield reports whether the next call to `Once` will yield the processor
// (`runtime.Gosched`) instead of performing a light `Pause()`.
//
// It is a pure "peek": it does not modify the internal state.
func (s *Wait) WillYield() bool {
	k := s.n >> 2
	if k > 4 {
		k = 4
	}
	mask := (uint32(1) << (4 - k)) - 1
	if (s.counter+1)&mask == 0 {
		return true
	}
	return false
}

// Reset resets the counters in Wait.
func (s *Wait) Reset() {
	s.counter = 0
	s.n = 0
}
