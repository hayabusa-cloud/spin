// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"code.hybscloud.com/spin"
)

func TestSpinLockBasic(t *testing.T) {
	var lk spin.Lock
	var v int64
	const G = 16
	const N = 1000

	var wg sync.WaitGroup
	wg.Add(G)
	for g := 0; g < G; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < N; i++ {
				lk.Lock()
				v++
				lk.Unlock()
			}
		}()
	}
	wg.Wait()

	if want, got := int64(G*N), v; want != got {
		t.Fatalf("counter mismatch: want=%d got=%d", want, got)
	}
}

// Attempt to exercise contention paths (pause and scheduler yield).
func TestSpinLockContention(t *testing.T) {
	var lk spin.Lock
	var ok atomic.Int32
	const G = 32
	var wg sync.WaitGroup
	wg.Add(G)
	for i := 0; i < G; i++ {
		go func() {
			defer wg.Done()
			// Create rapid contention
			for j := 0; j < 200; j++ {
				lk.Lock()
				// force some delay to keep others spinning
				runtime.Gosched()
				ok.Add(1)
				lk.Unlock()
			}
		}()
	}
	wg.Wait()
	if ok.Load() != G*200 {
		t.Fatalf("unexpected total: %d", ok.Load())
	}
}

// Ensure Unlock allows immediate re-lock and that Try fails while held.
func TestLock_UnlockImmediateRelockAndTry(t *testing.T) {
	var lk spin.Lock
	lk.Lock()
	if lk.Try() {
		t.Fatalf("Try should fail while lock is held")
	}
	lk.Unlock()
	if !lk.Try() {
		t.Fatalf("Try should succeed immediately after Unlock")
	}
	lk.Unlock()
}

// Try under contention: only one goroutine should succeed at a time.
func TestLock_TryUnderContention(t *testing.T) {
	var lk spin.Lock
	const G = 32
	const N = 200
	var successes atomic.Int64

	var wg sync.WaitGroup
	wg.Add(G)
	for g := 0; g < G; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < N; i++ {
				if lk.Try() {
					// simulate small work to let others attempt concurrently
					runtime.Gosched()
					successes.Add(1)
					lk.Unlock()
				} else {
					runtime.Gosched()
				}
			}
		}()
	}
	wg.Wait()

	// At most one success per attempt window; ensure the counter is sensible (>0)
	if successes.Load() == 0 {
		t.Fatalf("expected some successful Try() operations under contention")
	}
}
