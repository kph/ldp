// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package ldp

import (
	"math/rand"
	"reflect"
	"testing"
)

// TestSequenceTruncate makes sure that we properly truncate sequences
// We use random (but seeded) numbers for reproducability, and use
// a geometric progression to test many combinations within a reasonable
// amount of time.
func TestSequenceTruncate(t *testing.T) {
	for i := 0; i < 100000; i++ {
		s := make([]byte, i)
		_, _ = rand.Read(s)
		for j := 0; j < i; j = (j * 2) + 1 {
			p := SequenceToPattern(s, j)
			if len(p) != j {
				t.Errorf("Expected length %d got %d",
					j, len(p))
			}
			if !reflect.DeepEqual(s[:j], p) {
				t.Errorf("Truncated data did not match")
			}
		}
	}
}

// TestSequenceExpandWithoutRemainder tests the functionality of
// expanding a sequence to an exact multiple of the test data
func TestSequenceExpandWithoutRemainder(t *testing.T) {
	for i := 1; i < 500; i++ {
		s := make([]byte, i)
		_, _ = rand.Read(s)
		for j := 1; j < 100; j++ {
			k := i * j
			p := SequenceToPattern(s, k)
			if len(p) != k {
				t.Errorf("Expected length %d got %d",
					k, len(p))
			}
			for l := 0; l < k; l += i {
				if !reflect.DeepEqual(s, p[l:l+i]) {
					t.Errorf("Expanded data did not match")
				}
			}
		}
	}
}

// TestSequenceExpandWithRemainder tests the functionality of
// expanding a sequence to a non-exact multiple, thus duplicating some
// of the test data
func TestSequenceExpandWithRemainder(t *testing.T) {
	for i := 1; i < 100; i++ {
		s := make([]byte, i)
		_, _ = rand.Read(s)
		for j := 1; j < 100; j++ {
			for k := 1; k < i; k++ {
				l := i*j + k
				p := SequenceToPattern(s, l)
				if len(p) != l {
					t.Errorf("Expected length %d got %d",
						l, len(p))
				}
				for l := 0; l < (i * j); l += i {
					if !reflect.DeepEqual(s, p[l:l+i]) {
						t.Errorf("Expanded data did not match")
					}
				}
				if !reflect.DeepEqual(s[:k], p[i*j:]) {
					t.Errorf("Truncated data did not match")
				}
			}
		}
	}
}

// TestPatternAlpha tests the message generation code using a predefined
// string and precalculated SHA256.
func TestPatternToMsg(t *testing.T) {
	x := []byte(`

Pattern: Alpha
Length: 43
Sha256: d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592

The quick brown fox jumps over the lazy dog`)
	m := PatternToMsg(Alpha, []byte(AlphaTest))
	if !reflect.DeepEqual(m, x) {
		t.Errorf("Unexpected header response, expected %s\n got %s\n",
			x, m)
	}
}
