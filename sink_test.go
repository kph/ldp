// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package ldp

import (
	"io"
	"math/rand"
	"testing"
)

var sinkTestAlpha = []byte(`

Pattern: Alpha
Length: 43
Sha256: d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592

The quick brown fox jumps over the lazy dog`)

func TestSinkSimpleSynced(t *testing.T) {
	s := NewSink(true) // Testing pre-sync case
	for i := int64(0); i < 10; i++ {
		written, err := s.Write(sinkTestAlpha)
		if err != nil {
			t.Errorf("Unexpected error %s\n", err)
		}
		if written != len(sinkTestAlpha) {
			t.Errorf("Unexpcted length %d expected %d\n",
				written, len(sinkTestAlpha))
		}
		s.testIsExpectedResidual(t, true, i+1, 0, 0, 0, 0, []byte{})
	}
}

func (s *Sink) testWriteRandom(t *testing.T) {
	rgen := rand.New(rand.NewSource(0)) // Seed with zero for determinism
	rcnt := int64(rgen.Int31n(0xffff))
	written, err := io.CopyN(s, rgen, rcnt)
	if err != nil {
		t.Errorf("Error writing random bytes %s", err)
	}
	if written != rcnt {
		t.Errorf("Random bytes copied %d expected %d", written, rcnt)
	}
}

func TestSinkSyncUnsynced(t *testing.T) {
	s := NewSink(false) // Testing non pre-sync case
	s.testWriteRandom(t)
	s.testIsExpected(t, false, 0, 0, 0, 0, 0)
	for i := int64(0); i < 10; i++ {
		written, err := s.Write(sinkTestAlpha)
		if err != nil {
			t.Errorf("Unexpected error %s\n", err)
		}
		if written != len(sinkTestAlpha) {
			t.Errorf("Unexpcted length %d expected %d\n",
				written, len(sinkTestAlpha))
		}
		s.testIsExpectedResidual(t, true, i+1, 0, 0, 0, 0, []byte{})
	}
}

func TestSinkBytewiseSynced(t *testing.T) {
	s := NewSink(true)
	for i := int64(0); i < 10; i++ {
		for j := 0; j < len(sinkTestAlpha); j++ {
			written, err := s.Write([]byte{sinkTestAlpha[j]})
			if err != nil {
				t.Errorf("Unexpected error %s\n", err)
			}
			if written != 1 {
				t.Errorf("Unexpcted length %d expected 1\n",
					written)
			}
			if j < len(sinkTestAlpha)-1 {
				s.testIsExpectedResidual(t, true, i, 0, 0, 0, 0, sinkTestAlpha[:j+1])
			} else {
				s.testIsExpectedResidual(t, true, i+1, 0, 0, 0, 0, []byte{})
			}
		}
	}
}

func TestSinkBytewiseUnsynced(t *testing.T) {
	s := NewSink(false)
	for i := int64(0); i < 10; i++ {
		for j := 0; j < len(sinkTestAlpha); j++ {
			written, err := s.Write([]byte{sinkTestAlpha[j]})
			if err != nil {
				t.Errorf("Unexpected error %s\n", err)
			}
			if written != 1 {
				t.Errorf("Unexpcted length %d expected 1\n",
					written)
			}
			if j < len(sinkTestAlpha)-1 {
				s.testIsExpectedResidual(t, i != 0, i, 0, 0, 0, 0, sinkTestAlpha[:j+1])
			} else {
				s.testIsExpectedResidual(t, true, i+1, 0, 0, 0, 0, []byte{})
			}
		}
	}
}
