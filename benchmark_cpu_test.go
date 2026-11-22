// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin_test

import (
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"testing"
)

func linuxClockTicksPerSecond(b *testing.B) uint64 {
	b.Helper()
	auxv, err := os.ReadFile("/proc/self/auxv")
	if err != nil {
		b.Fatal(err)
	}
	// Each auxiliary vector entry is a pair of native words. AT_CLKTCK
	// gives the same clock frequency used by /proc/self/stat.
	const atClktck = 17
	const wordSize = strconv.IntSize / 8
	word := func(p []byte) uint64 {
		if wordSize == 4 {
			return uint64(binary.NativeEndian.Uint32(p))
		}
		return binary.NativeEndian.Uint64(p)
	}
	for len(auxv) >= 2*wordSize {
		if word(auxv) == atClktck {
			return word(auxv[wordSize:])
		}
		auxv = auxv[2*wordSize:]
	}
	b.Fatal("missing AT_CLKTCK in /proc/self/auxv")
	return 0
}

func linuxProcessTicks(b *testing.B) uint64 {
	b.Helper()
	stat, err := os.ReadFile("/proc/self/stat")
	if err != nil {
		b.Fatal(err)
	}
	// The process name in field 2 can contain spaces and parentheses.
	// Fields 14 and 15 are process-wide user and system CPU ticks.
	end := bytes.LastIndexByte(stat, ')')
	fields := bytes.Fields(stat[end+1:])
	if end < 0 || len(fields) < 13 {
		b.Fatal("missing CPU times in /proc/self/stat")
	}
	var ticks uint64
	for _, field := range fields[11:13] {
		value, err := strconv.ParseUint(string(field), 10, 64)
		if err != nil {
			b.Fatal(err)
		}
		ticks += value
	}
	return ticks
}
