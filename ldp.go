// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

// Package ldp implements the Platina Link Diagnostic Protocol
package ldp

const (
	// Pattern header
	Pattern   = "Pattern: "
	Length    = "Length: "
	SHA256    = "Sha256: "
	Alpha     = "Alpha"
	AlphaTest = "The quick brown fox jumps over the lazy dog"
)

type PatternEntry struct {
	Name        string                       // Pattern name
	Sourcer     func(s []byte, l int) []byte // Function to create pattern
	SourcerData []byte                       // Data to pass to the sourcer
}

var PatternList = []PatternEntry{
	{"00", SequenceToPattern, []byte{0x00}},
	{"FF", SequenceToPattern, []byte{0xff}},
	{"AA", SequenceToPattern, []byte{0xaa}},
	{"55", SequenceToPattern, []byte{0x55}},
	{"AA55", SequenceToPattern, []byte{0xaa, 0x55}},
	{"55AA", SequenceToPattern, []byte{0x55, 0xaa}},
	{"00FF", SequenceToPattern, []byte{0x00, 0xff}},
	{"FF00", SequenceToPattern, []byte{0xff, 0x00}},
	{Alpha, SequenceToPattern, []byte(AlphaTest)},
}

//Pattern: Random
//This is random data. This is the only pattern that the receiver can not check.
