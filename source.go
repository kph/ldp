// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package ldp

import (
	"crypto/sha256"
	"fmt"
)

// PatternToMsg takes a raw pattern named "n" with the contents "p"
// and formats a protocol message ready for the wire.
func PatternToMsg(n string, p []byte) (m []byte) {
	m = append([]byte(fmt.Sprintf("\n\n%s%s\n%s%d\n%s%x\n\n",
		Pattern, n,
		Length, len(p),
		SHA256, sha256.Sum256(p))), p...)
	return
}

// SequenceToPattern converts a sequence of bytes into a pattern of
// the specified length. The pattern will be repeated up to the maximum
// length, and truncated. So if the input sequence is longer than the
// length it will simply be truncated. If the length is a multiple of
// the pattern length, the output pattern won't be truncated. If its
// not a multiple, then the bytes that fit will be repeated.
func SequenceToPattern(s []byte, l int) (p []byte) {
	if l == 0 {
		return []byte{}
	}
	p = s
	for len(p) < l {
		p = append(p, p...) // Double sequence geometrically
	}
	p = p[:l] // Truncate
	return
}
