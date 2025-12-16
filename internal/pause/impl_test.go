// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pause

import "testing"

func TestPause_Branches(t *testing.T) {
	// Exercise both branches in Pause.
	Pause(-1)
	Pause(0)
	Pause(1)
	Pause(2)
	Pause(64)
}

func TestPause_NoAllocs(t *testing.T) {
	allocs := testing.AllocsPerRun(1000, func() {
		Pause(1)
		Pause(2)
	})
	if allocs != 0 {
		t.Fatalf("Pause allocated: got %v allocs/op, want 0", allocs)
	}
}
