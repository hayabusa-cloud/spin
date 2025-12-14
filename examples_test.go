// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"fmt"
	"sync/atomic"
	"time"

	"code.hybscloud.com/spin"
)

func ExampleLock() {
	var lk spin.Lock
	lk.Lock()
	// critical section
	lk.Unlock()
	fmt.Println("locked")
	// Output: locked
}

func ExampleWait() {
	var sw spin.Wait
	var ready atomic.Bool
	go func() {
		time.Sleep(50 * time.Microsecond)
		ready.Store(true)
	}()
	for !ready.Load() {
		sw.Once()
	}
	fmt.Println("ready")
	// Output: ready
}

func ExampleYield() {
	spin.Yield(10 * time.Millisecond)
	fmt.Println("yielded")
}
