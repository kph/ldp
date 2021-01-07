// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package ldp

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Sink holds the context for received data. It contains state information,
// residual data, and counters.
type Sink struct {
	GoodFrames  int64  // Count of good frames received
	FramingErr  int64  // Count of framing errors
	ProtocolErr int64  // Count of protocol errors
	ChecksumErr int64  // Count of checksum errors
	PatternErr  int64  // Count of pattern errors
	sync        bool   // Do we expect that we are synchronized
	b           []byte // Unprocessed input bytes
}

var patternRegex, sha256Regex, lengthRegex *regexp.Regexp

// NewSink() creates a new sink. If we expect that we are not going
// to be initially synchronized (in the case that the agent is started
// independently on two machines) set preSync to false to indicate that
// it is expected that we will be initially unsynchronized. For loopback
// tests, set preSync to true. This ensures we count spurious bytes
// correctly.
func NewSink(preSync bool) (s *Sink) {
	return &Sink{sync: preSync}
}

// Write() - write into the sink. Write bytes received into the
// sink, frame and parse.
func (s *Sink) Write(p []byte) (n int, err error) {
	s.b = append(s.b, p...)

	for {
		// Scan the input looking for the start of the header in the
		// input stream
		for len(s.b) >= 2+len(Pattern) && (s.b[0] != '\n' ||
			s.b[1] != '\n' ||
			string(s.b[2:2+len(Pattern)]) != Pattern) {
			if s.sync {
				s.sync = false // We are re-syncing
				s.FramingErr++ // Count this as a framing error
			}
			s.b = s.b[1:]
		}
		// Get out if we don't have enough data to parse
		if len(s.b) <= 2+len(Pattern) {
			break
		}

		// Split into header and data
		sp := strings.SplitAfterN(string(s.b[2:]), "\n\n", 2)
		if len(sp) < 1 || !strings.HasSuffix(sp[0], "\n\n") {
			// TODO: possibly bail on an impossibly long header.
			// i.e. we are this far because we have seen the
			// introduction of the header, but not an end. Should
			// eventually resync. Think about a test for this case.
			break
		}

		// TODO: Do this with one regex
		patternMatch := patternRegex.FindAllStringSubmatch(sp[0], 2)
		lengthMatch := lengthRegex.FindAllStringSubmatch(sp[0], 2)
		sha256Match := sha256Regex.FindAllStringSubmatch(sp[0], 2)

		dataLen := 0
		patternStr := ""
		lengthStr := ""
		sha256Str := ""

		if len(patternMatch) > 0 && len(patternMatch[0]) > 2 &&
			len(lengthMatch) > 0 && len(patternMatch[0]) > 2 &&
			len(sha256Match) > 0 && len(sha256Match[0]) > 2 {
			patternStr = patternMatch[0][2]
			lengthStr = lengthMatch[0][2]
			sha256Str = sha256Match[0][2]

			dataLen, err = strconv.Atoi(lengthStr)
		} else {
			err = errors.New("Protocol error")
		}

		if err != nil {
			s.sync = false      // Protocol error - resync
			s.ProtocolErr++     // Count as a protocol error
			s.b = []byte(sp[1]) // Scan after the header
			continue            // Continue parsing
		}

		// We have a complete header. Do we have complete data?
		patData := []byte{}
		if dataLen > 0 {
			if len(sp) < 2 || len(sp[1]) < dataLen {
				break
			}
			s.b = []byte(sp[1][dataLen:])
			patData = []byte(sp[1][:dataLen])
		} else {
			s.b = s.b[len(sp[0])+2:]
		}

		// Complete data, we are synced
		s.sync = true

		// Check the data sequence
		sha256Calc := fmt.Sprintf("%x", sha256.Sum256(patData))
		shaGood := sha256Calc == sha256Str
		pat, patFound := PatternMap[patternStr]
		patErr := error(nil)
		if patFound {
			patErr = pat.Sinker(pat.SequenceData, patData)
		} else {
			patErr = ErrPatternUnknown
		}
		if shaGood && patErr == nil {
			s.GoodFrames++
			continue
		}
		if !shaGood {
			fmt.Printf("Expected SHA256 %s got %s\n",
				sha256Str, sha256Calc)
			s.ChecksumErr++
		}
		if patErr != nil {
			fmt.Printf("Pattern %s error: %s\n", patternStr, patErr)
			s.PatternErr++
		}
	}
	return len(p), nil
}

// CheckPattern checks that the pattern received matches the expected
// data sequence. A pattern is defined as a potentially truncated,
// potentially repeated sequence.
// truncated
func CheckPattern(s []byte, p []byte) (err error) {
	residual := p
	for len(residual) > 0 {
		m := len(residual)
		if m > len(s) {
			m = len(s)
		}
		if !reflect.DeepEqual(s[:m], residual[:m]) {
			fmt.Printf("mismatch: m:%d s:%v r:%v\n",
				m, s[:m], residual[:m])
			return ErrPatternMismatch
		}
		residual = residual[m:]
	}
	return nil
}

func init() {
	patternRegex = regexp.MustCompile(`(?m)(^` + Pattern + `)([\p{L}\d_]+)$`)
	lengthRegex = regexp.MustCompile(`(?m)(^` + Length + `)(\d+)$`)
	sha256Regex = regexp.MustCompile(`(?m)(^` + SHA256 + `)([0-9a-fA-F]+)$`)
}
