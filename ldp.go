// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

// Package ldp implements the Platina Link Diagnostic Protocol
package ldp

import (
	"errors"
	"math/rand"
)

const (
	// Pattern header
	Pattern = "Pattern: "
	Length  = "Length: "
	SHA256  = "Sha256: "

	// Special patterns
	Alpha     = "Alpha"
	AlphaTest = "The quick brown fox jumps over the lazy dog"
)

var ErrPatternMismatch = errors.New("Pattern data mismatch")
var ErrPatternUnknown = errors.New("Pattern unknown")

type PatternEntry struct {
	Sourcer      func(s []byte, l int) []byte   // Function to create pattern
	Sinker       func(s []byte, p []byte) error // Function to check pattern
	SequenceData []byte                         // Sequence bytes
}

var PatternMap = map[string]PatternEntry{
	"00":   {SequenceToPattern, CheckPattern, []byte{0x00}},
	"FF":   {SequenceToPattern, CheckPattern, []byte{0xff}},
	"AA":   {SequenceToPattern, CheckPattern, []byte{0xaa}},
	"55":   {SequenceToPattern, CheckPattern, []byte{0x55}},
	"AA55": {SequenceToPattern, CheckPattern, []byte{0xaa, 0x55}},
	"55AA": {SequenceToPattern, CheckPattern, []byte{0x55, 0xaa}},
	"00FF": {SequenceToPattern, CheckPattern, []byte{0x00, 0xff}},
	"FF00": {SequenceToPattern, CheckPattern, []byte{0xff, 0x00}},
	Alpha:  {SequenceToPattern, CheckPattern, []byte(AlphaTest)},
}

//Pattern: Random to do
//This is random data. This is the only pattern that the receiver can not check.

// (pm *PatternEntry)NewMessage returns a protocol message using the
// specified PatternEntry and length
func (pat *PatternEntry) NewMessage(pn string, l int) []byte {
	return PatternToMsg(pn, pat.Sourcer(pat.SequenceData, l))
}

// NewMessage returns a protocol message from the named pattern of the
// specified length
func NewMessage(pn string, l int) (p []byte, err error) {
	pat, patFound := PatternMap[pn]
	if !patFound {
		return nil, ErrPatternUnknown
	}
	return pat.NewMessage(pn, l), nil
}

// NewRandomMessage returns a protocol message of a random length of
// a random type.
func NewRandomMessage() []byte {
	mtype := rand.Int31n(int32(len(PatternMap)))
	mlen := int(rand.Int31n(0xffff))
	for k, pat := range PatternMap {
		if mtype == 0 {
			return pat.NewMessage(k, mlen)
		}
		mtype--
	}
	panic("weird randomness")
}
