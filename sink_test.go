// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package ldp

import (
	"reflect"
	"testing"
)

var sinkTestAlpha = []byte(`

Pattern: Alpha
Length: 43
Sha256: d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592

The quick brown fox jumps over the lazy dog`)

func (s *Sink) testIsExpected(t *testing.T,
	sync bool, goodFrames int64, framingErr int64,
	protocolErr int64, checksumErr int64, patternErr int64, b []byte) {
	if sync != s.sync {
		t.Errorf("Sync expected %v got %v", sync, s.sync)
	}
	if goodFrames != s.GoodFrames {
		t.Errorf("GoodFrames expected %d got %d", goodFrames,
			s.GoodFrames)
	}
	if framingErr != s.FramingErr {
		t.Errorf("FramingErr expected %d got %d", framingErr,
			s.FramingErr)
	}
	if protocolErr != s.ProtocolErr {
		t.Errorf("ProtocolErr expected %d got %d", protocolErr,
			s.ProtocolErr)
	}
	if checksumErr != s.ChecksumErr {
		t.Errorf("ChecksumErr expected %d got %d", checksumErr,
			s.ChecksumErr)
	}
	if patternErr != s.PatternErr {
		t.Errorf("PatternErr expected %d got %d", patternErr,
			s.PatternErr)
	}
	if !reflect.DeepEqual(b, s.b) {
		t.Errorf("Residual data mismatch, got %v expected %v",
			b, s.b)
	}
}

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
		s.testIsExpected(t, true, i+1, 0, 0, 0, 0, []byte{})
	}
}

func TestSinkSyncUnsynced(t *testing.T) {
	s := NewSink(false) // Testing non pre-sync case
	for i := int64(0); i < 10; i++ {
		written, err := s.Write(sinkTestAlpha)
		if err != nil {
			t.Errorf("Unexpected error %s\n", err)
		}
		if written != len(sinkTestAlpha) {
			t.Errorf("Unexpcted length %d expected %d\n",
				written, len(sinkTestAlpha))
		}
		s.testIsExpected(t, true, i+1, 0, 0, 0, 0, []byte{})
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
				s.testIsExpected(t, true, i, 0, 0, 0, 0, sinkTestAlpha[:j+1])
			} else {
				s.testIsExpected(t, true, i+1, 0, 0, 0, 0, []byte{})
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
				s.testIsExpected(t, i != 0, i, 0, 0, 0, 0, sinkTestAlpha[:j+1])
			} else {
				s.testIsExpected(t, true, i+1, 0, 0, 0, 0, []byte{})
			}
		}
	}
}
