// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin

import "runtime"

// Wait adapts to short local contention. The zero value is ready to use.
type Wait struct {
	counter, n uint32
}

// Once takes one wait step. It starts with [Pause] and calls [runtime.Gosched]
// more often as the wait continues.
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

// WillYield reports whether the next [Wait.Once] calls [runtime.Gosched]
// directly. It leaves the wait state unchanged. On fallback targets, [Pause]
// can also yield to the scheduler.
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

// Reset starts the wait again at its initial pause interval.
func (s *Wait) Reset() {
	s.counter = 0
	s.n = 0
}
