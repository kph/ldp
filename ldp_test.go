// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package ldp

import (
	"reflect"
	"testing"
)

func (s *Sink) testIsExpected(t *testing.T,
	sync bool, goodFrames int64, framingErr int64,
	protocolErr int64, checksumErr int64, patternErr int64) {
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
}

func (s *Sink) testIsExpectedResidual(t *testing.T,
	sync bool, goodFrames int64, framingErr int64,
	protocolErr int64, checksumErr int64, patternErr int64, b []byte) {
	s.testIsExpected(t, sync, goodFrames, framingErr,
		protocolErr, checksumErr, patternErr)

	if !reflect.DeepEqual(b, s.b) {
		t.Errorf("Residual data mismatch, got %v expected %v",
			b, s.b)
	}
}

func TestLDP(t *testing.T) {
	s := NewSink(true)
	m := NewRandomMessage()
	n, err := s.Write(m)
	if n != len(m) {
		t.Errorf("Expected written length %d got %d",
			n, len(m))
	}
	if err != nil {
		t.Errorf("Unexpcted error %s", err)
	}
	s.testIsExpectedResidual(t, true, 1, 0, 0, 0, 0, []byte{})
}
