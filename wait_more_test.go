// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"testing"

	"code.hybscloud.com/spin"
)

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
