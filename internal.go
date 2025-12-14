// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin

// noCopy may be embedded to prevent copying by static analyzers.
// It follows the convention used across the Go standard library.
//
//go:nocopy
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
