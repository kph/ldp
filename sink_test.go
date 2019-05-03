// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package ldp

import (
	"testing"
)

var sinkTestAlpha = []byte(`

Pattern: Alpha
Length: 43
Sha256: d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592

The quick brown fox jumps over the lazy dog`)

func TestSink(t *testing.T) {
	s := NewSink(true) // Testing pre-sync case
	s.Write(sinkTestAlpha)
}
