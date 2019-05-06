// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

// Package ldp implements the Platina Link Diagnostic Protocol
package ldp

import (
	"errors"
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

//Pattern: Random
//This is random data. This is the only pattern that the receiver can not check.
