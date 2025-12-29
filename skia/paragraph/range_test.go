package paragraph

import (
	"testing"
)

func TestNewRange(t *testing.T) {
	r := NewRange(10, 20)
	if r.Start != 10 || r.End != 20 {
		t.Errorf("NewRange(10, 20) = {%d, %d}, want {10, 20}", r.Start, r.End)
	}
}

func TestRangeWidth(t *testing.T) {
	tests := []struct {
		name  string
		start int
		end   int
		want  int
	}{
		{"positive width", 10, 20, 10},
		{"zero width", 5, 5, 0},
		{"single unit", 0, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRange(tt.start, tt.end)
			if got := r.Width(); got != tt.want {
				t.Errorf("Range{%d, %d}.Width() = %d, want %d", tt.start, tt.end, got, tt.want)
			}
		})
	}
}

func TestRangeShift(t *testing.T) {
	tests := []struct {
		name      string
		start     int
		end       int
		delta     int
		wantStart int
		wantEnd   int
	}{
		{"shift right", 10, 20, 5, 15, 25},
		{"shift left", 10, 20, -3, 7, 17},
		{"no shift", 10, 20, 0, 10, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRange(tt.start, tt.end)
			r.Shift(tt.delta)
			if r.Start != tt.wantStart || r.End != tt.wantEnd {
				t.Errorf("after Shift(%d): got {%d, %d}, want {%d, %d}",
					tt.delta, r.Start, r.End, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestRangeContains(t *testing.T) {
	tests := []struct {
		name   string
		r1     Range[int]
		r2     Range[int]
		expect bool
	}{
		{"contains fully", NewRange(0, 100), NewRange(10, 50), true},
		{"contains equal", NewRange(10, 50), NewRange(10, 50), true},
		{"contains start edge", NewRange(0, 50), NewRange(0, 25), true},
		{"contains end edge", NewRange(0, 50), NewRange(25, 50), true},
		{"not contains - before", NewRange(10, 20), NewRange(5, 15), false},
		{"not contains - after", NewRange(10, 20), NewRange(15, 25), false},
		{"not contains - outer", NewRange(10, 20), NewRange(5, 25), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r1.Contains(tt.r2); got != tt.expect {
				t.Errorf("Range{%d, %d}.Contains({%d, %d}) = %v, want %v",
					tt.r1.Start, tt.r1.End, tt.r2.Start, tt.r2.End, got, tt.expect)
			}
		})
	}
}

func TestRangeIntersects(t *testing.T) {
	tests := []struct {
		name   string
		r1     Range[int]
		r2     Range[int]
		expect bool
	}{
		{"overlapping", NewRange(0, 20), NewRange(10, 30), true},
		{"contained", NewRange(0, 100), NewRange(10, 50), true},
		{"touching", NewRange(0, 10), NewRange(10, 20), true},
		{"same", NewRange(10, 20), NewRange(10, 20), true},
		{"separated", NewRange(0, 10), NewRange(20, 30), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r1.Intersects(tt.r2); got != tt.expect {
				t.Errorf("Range{%d, %d}.Intersects({%d, %d}) = %v, want %v",
					tt.r1.Start, tt.r1.End, tt.r2.Start, tt.r2.End, got, tt.expect)
			}
		})
	}
}

func TestRangeIntersection(t *testing.T) {
	tests := []struct {
		name      string
		r1        Range[int]
		r2        Range[int]
		wantStart int
		wantEnd   int
	}{
		{"overlapping", NewRange(0, 20), NewRange(10, 30), 10, 20},
		{"contained", NewRange(0, 100), NewRange(10, 50), 10, 50},
		{"touching", NewRange(0, 10), NewRange(10, 20), 10, 10},
		{"same", NewRange(10, 20), NewRange(10, 20), 10, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.r1.Intersection(tt.r2)
			if result.Start != tt.wantStart || result.End != tt.wantEnd {
				t.Errorf("Range{%d, %d}.Intersection({%d, %d}) = {%d, %d}, want {%d, %d}",
					tt.r1.Start, tt.r1.End, tt.r2.Start, tt.r2.End,
					result.Start, result.End, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestRangeEmpty(t *testing.T) {
	// Empty range
	empty := EmptyRange
	if !empty.Empty() {
		t.Errorf("EmptyRange.Empty() = false, want true")
	}

	// Non-empty range
	normal := NewRange(0, 10)
	if normal.Empty() {
		t.Errorf("Range{0, 10}.Empty() = true, want false")
	}
}

func TestRangeIsValid(t *testing.T) {
	tests := []struct {
		name   string
		r      Range[int]
		expect bool
	}{
		{"valid positive width", NewRange(0, 10), true},
		{"valid zero width", NewRange(5, 5), true},
		{"invalid negative width", NewRange(10, 5), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.IsValid(); got != tt.expect {
				t.Errorf("Range{%d, %d}.IsValid() = %v, want %v",
					tt.r.Start, tt.r.End, got, tt.expect)
			}
		})
	}
}

func TestRangeEquals(t *testing.T) {
	r1 := NewRange(10, 20)
	r2 := NewRange(10, 20)
	r3 := NewRange(10, 25)

	if !r1.Equals(r2) {
		t.Errorf("Range{10, 20}.Equals({10, 20}) = false, want true")
	}
	if r1.Equals(r3) {
		t.Errorf("Range{10, 20}.Equals({10, 25}) = true, want false")
	}
}

func TestTextRange(t *testing.T) {
	tr := NewTextRange(0, 100)
	if tr.Start != 0 || tr.End != 100 {
		t.Errorf("NewTextRange(0, 100) = {%d, %d}, want {0, 100}", tr.Start, tr.End)
	}
	if tr.Width() != 100 {
		t.Errorf("TextRange{0, 100}.Width() = %d, want 100", tr.Width())
	}
}

func TestBlockRange(t *testing.T) {
	br := NewBlockRange(5, 15)
	if br.Start != 5 || br.End != 15 {
		t.Errorf("NewBlockRange(5, 15) = {%d, %d}, want {5, 15}", br.Start, br.End)
	}
	if br.Width() != 10 {
		t.Errorf("BlockRange{5, 15}.Width() = %d, want 10", br.Width())
	}
}

func TestEmptyIndex(t *testing.T) {
	// EmptyIndex should be a very large number
	if EmptyIndex <= 0 {
		t.Errorf("EmptyIndex = %d, should be positive", EmptyIndex)
	}
}
