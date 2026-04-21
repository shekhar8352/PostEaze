package studio

import (
	"testing"
)

func TestGeneratePositionBetween_Empty(t *testing.T) {
	got, err := GeneratePositionBetween("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Fatalf("expected non-empty position")
	}
}

func TestGeneratePositionBetween_Ordering(t *testing.T) {
	// Simulate inserting 20 items at the end and then shuffling in the middle.
	var positions []string
	var last string
	for i := 0; i < 20; i++ {
		p, err := GeneratePositionBetween(last, "")
		if err != nil {
			t.Fatalf("append %d: %v", i, err)
		}
		if last != "" && p <= last {
			t.Fatalf("append not strictly increasing: last=%q next=%q", last, p)
		}
		positions = append(positions, p)
		last = p
	}

	// Insert between positions[5] and positions[6] 10 times — should never
	// collide and must remain strictly between its neighbors.
	left, right := positions[5], positions[6]
	for i := 0; i < 10; i++ {
		mid, err := GeneratePositionBetween(left, right)
		if err != nil {
			t.Fatalf("insert %d: %v", i, err)
		}
		if !(mid > left && mid < right) {
			t.Fatalf("insert out of range: left=%q mid=%q right=%q", left, mid, right)
		}
		right = mid
	}
}

func TestGeneratePositionBetween_RejectsInvalid(t *testing.T) {
	if _, err := GeneratePositionBetween("Z", "A"); err == nil {
		t.Fatalf("expected error when prev >= next")
	}
}
