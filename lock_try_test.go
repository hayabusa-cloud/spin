// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"sync"
	"testing"

	"code.hybscloud.com/spin"
)

func TestLockTry(t *testing.T) {
	var lk spin.Lock
	// First try should succeed
	if !lk.Try() {
		t.Fatalf("expected Try to succeed on unlocked lock")
	}
	// While locked, Try should fail
	if lk.Try() {
		t.Fatalf("expected Try to fail when already locked")
	}
	lk.Unlock()
	// After unlock, Try should succeed again
	if !lk.Try() {
		t.Fatalf("expected Try to succeed after Unlock")
	}
	lk.Unlock()
}

func BenchmarkLockContention(b *testing.B) {
	var lk spin.Lock
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lk.Lock()
			lk.Unlock()
		}
	})
}

func BenchmarkLockCriticalSection(b *testing.B) {
	var lk spin.Lock
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			lk.Lock()
			lk.Unlock()
		}
	}()
	wg.Wait()
}
