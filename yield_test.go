// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"testing"
	"time"

	"code.hybscloud.com/spin"
)

func BenchmarkYield_Default(b *testing.B) {
	// Use default configured duration
	spin.SetYieldDuration(250 * time.Microsecond)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spin.Yield()
	}
}

func BenchmarkYield_Gosched(b *testing.B) {
	// Force gosched path (non-positive duration)
	spin.SetYieldDuration(0)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spin.Yield()
	}
}

func BenchmarkYield_ShortSleep(b *testing.B) {
	spin.SetYieldDuration(0) // ensure override path is used
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spin.Yield(50 * time.Microsecond)
	}
}
