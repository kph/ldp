// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package ldp

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Sink struct {
	Sync        bool   // Do we expect that we are synchronized
	GoodFrames  int64  // Count of good frames received
	FramingErr  int64  // Count of framing errors
	ProtocolErr int64  // Count of protocol errors
	ChecksumErr int64  // Count of checksum errors
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
	return &Sink{Sync: preSync}
}

// Write() - write into the sink. Write bytes received into the
// sink, frame and parse.
func (s *Sink) Write(p []byte) (n int, err error) {
	s.b = append(s.b, p...)

	for {
		// Scan the input looking for the start of the header in the
		// input stream
		for len(s.b) >= 2+len(Pattern) && s.b[0] != '\n' && s.b[1] != '\n' &&
			string(s.b[2:2+len(Pattern)]) != Pattern {
			if !s.Sync {
				s.Sync = true  // We are re-syncing
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
		if len(sp) < 2 {
			// TODO: possibly bail on an impossibly long header.
			// i.e. we are this far because we have seen the
			// introduction of the header, but not an end. Should
			// eventually resync. Think about a test for this case.
			break
		}

		fmt.Printf("Parsed header\n%s", sp[0])
		//patternStr := patternRegex.FindAllStringSubmatch(sp[0], 2)[0][2]
		lengthStr := lengthRegex.FindAllStringSubmatch(sp[0], 2)[0][2]
		sha256Str := sha256Regex.FindAllStringSubmatch(sp[0], 2)[0][2]

		dataLen, err := strconv.Atoi(lengthStr)
		if err != nil {
			s.Sync = true       // Protocol error - resync
			s.ProtocolErr++     // Count as a protocol error
			s.b = []byte(sp[1]) // Scan after the header
			continue            // Continue parsing
		}

		// We have a complete header. Do we have complete data?
		if len(sp[1]) < dataLen {
			break
		}

		// Complete data, save the residual for next time
		s.b = []byte(sp[1][dataLen:])

		sha256Calc := fmt.Sprintf("%x", sha256.Sum256([]byte(sp[1][:dataLen])))
		if sha256Calc != sha256Str {
			fmt.Printf("Expected SHA256 %s got %s\n",
				sha256Str, sha256Calc)
			s.ChecksumErr++
			continue
		}
		s.GoodFrames++
	}
	return len(p), nil
}

func init() {
	patternRegex = regexp.MustCompile(`(?m)(^` + Pattern + `)([\p{L}\d_]+)$`)
	lengthRegex = regexp.MustCompile(`(?m)(^` + Length + `)(\d+)$`)
	sha256Regex = regexp.MustCompile(`(?m)(^` + SHA256 + `)([0-9a-fA-F]+)$`)
}
