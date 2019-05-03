// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package ldp

import (
	"fmt"
	"regexp"
	"strings"
)

type Sink struct {
	sync     bool // Do we expect that we are synchronized
	Spurious byte // Spurious - unexpected bytes
	b        []byte
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
	// Scan the input looking for the start of the header in the
	// input stream
	for len(s.b) >= 2+len(Pattern) && s.b[0] != '\n' && s.b[1] != '\n' &&
		string(s.b[2:2+len(Pattern)]) != Pattern {
		if !s.sync {
			s.Spurious++
		}
		s.b = s.b[1:]
	}
	// Get out if we don't have enough data to parse
	if len(s.b) <= 2+len(Pattern) {
		return len(p), nil
	}
	// Split into header and data
	sp := strings.SplitAfterN(string(s.b[2:]), "\n\n", 2)
	if len(sp) < 2 {
		return
	}
	fmt.Printf("Parsed header\n%s", sp[0])
	fmt.Printf("%s %s\n", Pattern, patternRegex.FindAllStringSubmatch(sp[0], 2)[0][2])
	fmt.Printf("%s %s\n", Length, lengthRegex.FindAllStringSubmatch(sp[0], 2)[0][2])
	fmt.Printf("%s %s\n", SHA256, sha256Regex.FindAllStringSubmatch(sp[0], 2)[0][2])

	return len(p), nil
}

func init() {
	patternRegex = regexp.MustCompile(`(?m)(^` + Pattern + `)([\p{L}\d_]+)$`)
	lengthRegex = regexp.MustCompile(`(?m)(^` + Length + `)(\d+)$`)
	sha256Regex = regexp.MustCompile(`(?m)(^` + SHA256 + `)([0-9a-fA-F]+)$`)
}
