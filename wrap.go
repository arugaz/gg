// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gg

import (
	"strings"
	"unicode"
)

// measureStringer is an interface for objects that can measure the width and height of a string
// when rendered with a specific font and style. Implementing this interface allows objects to
// provide text measurement capabilities for layout and rendering.
type measureStringer interface {
	MeasureString(s string) (w, h float64)
}

// splitOnSpace splits the input string by whitespace characters and returns a slice
// of substrings. It breaks the input string at each occurrence of one or more consecutive
// whitespace characters, treating multiple spaces as a single separator.
func splitOnSpace(x string) []string {
	var result []string
	pi := 0
	ps := false

	for i, c := range x {
		s := unicode.IsSpace(c)
		if s != ps && i > 0 {
			result = append(result, x[pi:i])
			pi = i
		}
		ps = s
	}

	result = append(result, x[pi:])

	return result
}

// wordWrap performs word wrapping on a given input string, breaking it into lines
// based on a specified maximum width. The function uses a measureStringer to calculate
// the width of the text. It takes into account spaces and line breaks and ensures that
// words are not split in the middle.
func wordWrap(m measureStringer, s string, width float64) []string {
	var result []string

	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) == "" {
			result = append(result, line)
			continue
		}

		fields := splitOnSpace(line)

		if len(fields)%2 == 1 {
			fields = append(fields, "")
		}

		x := ""

		for i := 0; i < len(fields); i += 2 {
			w, _ := m.MeasureString(x + fields[i])

			if w > width {
				if x == "" {
					result = append(result, fields[i])
					x = ""
					continue
				} else {
					result = append(result, x)
					x = ""
				}
			}

			x += fields[i] + fields[i+1]
		}

		if x != "" {
			result = append(result, x)
		}
	}

	for i, line := range result {
		result[i] = strings.TrimSpace(line)
	}

	return result
}
