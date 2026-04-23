// Package studio contains domain utilities for the Studio content pipeline.
// Business orchestration lives in business/v1; this package holds pure helpers
// (fractional indexing, default phase templates, etc.) that are easy to unit
// test in isolation.
package studio

import (
	"fmt"
	"strings"
)

// Base62 alphabet used for fractional-index position strings. Chosen so that
// byte-wise comparisons on ASCII produce the desired ordering: '0' < '9' < 'A'
// < 'Z' < 'a' < 'z'.
const (
	posDigits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	posBase   = 62
)

func digitIndex(c byte) int { return strings.IndexByte(posDigits, c) }

// GeneratePositionBetween returns a position string strictly between prev and
// next in lexicographic order. An empty prev means "before everything"; an
// empty next means "after everything". If both are empty this returns a seed
// value that lies in the middle of the keyspace.
//
// The returned key is guaranteed to satisfy prev < key < next (for non-empty
// bounds) using byte-wise ASCII comparison.
func GeneratePositionBetween(prev, next string) (string, error) {
	if prev != "" && next != "" && prev >= next {
		return "", fmt.Errorf("studio: prev %q must be strictly less than next %q", prev, next)
	}
	return midpoint(prev, next), nil
}

// midpoint implements the core fractional-index algorithm. It walks past any
// shared prefix of a and b and then finds a digit that lies strictly between
// the next differing digit of each side, descending when adjacent digits force
// a longer key.
func midpoint(a, b string) string {
	n := 0
	for n < len(a) && n < len(b) && a[n] == b[n] {
		n++
	}
	if n > 0 {
		return a[:n] + midpoint(a[n:], b[n:])
	}

	ac := 0
	if len(a) > 0 {
		ac = digitIndex(a[0])
	}
	bc := posBase
	if len(b) > 0 {
		bc = digitIndex(b[0])
	}

	if bc-ac > 1 {
		mid := (ac + bc) / 2
		return string(posDigits[mid])
	}

	// Adjacent digits: need to descend to build a longer key.
	if len(b) > 1 {
		// Pin the first digit of b and recurse to find something below its tail.
		return string(b[0]) + midpoint("", b[1:])
	}
	var tail string
	if len(a) > 0 {
		tail = a[1:]
	}
	return string(posDigits[ac]) + midpoint(tail, "")
}

// PositionBefore returns a key that sorts strictly before k.
func PositionBefore(k string) (string, error) {
	return GeneratePositionBetween("", k)
}

// PositionAfter returns a key that sorts strictly after k.
func PositionAfter(k string) (string, error) {
	return GeneratePositionBetween(k, "")
}

// FirstPosition returns a reasonable seed key used for the first item in a phase.
func FirstPosition() string {
	return midpoint("", "")
}
