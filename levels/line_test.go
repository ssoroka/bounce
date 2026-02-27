package levels

import (
	"testing"
)

func TestIntersect(t *testing.T) {
	tests := []struct {
		name      string
		l1        Line
		l2        Line
		wantPoint Point
		wantHit   bool
	}{
		{
			name:      "simple cross at origin",
			l1:        Line{From: Point{-1, 0}, To: Point{1, 0}},
			l2:        Line{From: Point{0, -1}, To: Point{0, 1}},
			wantPoint: Point{0, 0},
			wantHit:   true,
		},
		{
			name:      "cross at (5,5)",
			l1:        Line{From: Point{0, 5}, To: Point{10, 5}},
			l2:        Line{From: Point{5, 0}, To: Point{5, 10}},
			wantPoint: Point{5, 5},
			wantHit:   true,
		},
		{
			name:      "diagonal cross",
			l1:        Line{From: Point{0, 0}, To: Point{10, 10}},
			l2:        Line{From: Point{0, 10}, To: Point{10, 0}},
			wantPoint: Point{5, 5},
			wantHit:   true,
		},
		{
			name:    "parallel horizontal lines",
			l1:      Line{From: Point{0, 0}, To: Point{10, 0}},
			l2:      Line{From: Point{0, 5}, To: Point{10, 5}},
			wantHit: false,
		},
		{
			name:    "parallel vertical lines",
			l1:      Line{From: Point{0, 0}, To: Point{0, 10}},
			l2:      Line{From: Point{5, 0}, To: Point{5, 10}},
			wantHit: false,
		},
		{
			name:    "segments don't reach",
			l1:      Line{From: Point{0, 0}, To: Point{5, 0}},
			l2:      Line{From: Point{10, 0}, To: Point{10, 10}},
			wantHit: false,
		},
		{
			name:    "would intersect if extended",
			l1:      Line{From: Point{0, 0}, To: Point{2, 2}},
			l2:      Line{From: Point{5, 0}, To: Point{5, 3}},
			wantHit: false,
		},
		{
			name:      "T intersection at endpoint",
			l1:        Line{From: Point{0, 5}, To: Point{10, 5}},
			l2:        Line{From: Point{5, 5}, To: Point{5, 10}},
			wantPoint: Point{5, 5},
			wantHit:   true,
		},
		{
			name:      "boundary edge vs line crossing",
			l1:        Line{From: Point{0, 0}, To: Point{100, 0}},
			l2:        Line{From: Point{50, -10}, To: Point{50, 10}},
			wantPoint: Point{50, 0},
			wantHit:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPoint, gotHit := tt.l1.Intersect(tt.l2)

			if gotHit != tt.wantHit {
				t.Errorf("Intersect() hit = %v, want %v", gotHit, tt.wantHit)
				return
			}

			if tt.wantHit {
				const eps = 0.01
				if abs(gotPoint.X-tt.wantPoint.X) > eps || abs(gotPoint.Y-tt.wantPoint.Y) > eps {
					t.Errorf("Intersect() point = %v, want %v", gotPoint, tt.wantPoint)
				}
			}
		})
	}
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func TestContainsPoint(t *testing.T) {
	tests := []struct {
		name string
		line Line
		p    Point
		want bool
	}{
		{
			name: "point at start",
			line: Line{From: Point{0, 0}, To: Point{10, 0}},
			p:    Point{0, 0},
			want: true,
		},
		{
			name: "point at end",
			line: Line{From: Point{0, 0}, To: Point{10, 0}},
			p:    Point{10, 0},
			want: true,
		},
		{
			name: "point in middle horizontal",
			line: Line{From: Point{0, 0}, To: Point{10, 0}},
			p:    Point{5, 0},
			want: true,
		},
		{
			name: "point in middle vertical",
			line: Line{From: Point{0, 0}, To: Point{0, 10}},
			p:    Point{0, 5},
			want: true,
		},
		{
			name: "point in middle diagonal",
			line: Line{From: Point{0, 0}, To: Point{10, 10}},
			p:    Point{5, 5},
			want: true,
		},
		{
			name: "point off line horizontally",
			line: Line{From: Point{0, 0}, To: Point{10, 0}},
			p:    Point{5, 1},
			want: false,
		},
		{
			name: "point before start",
			line: Line{From: Point{0, 0}, To: Point{10, 0}},
			p:    Point{-1, 0},
			want: false,
		},
		{
			name: "point after end",
			line: Line{From: Point{0, 0}, To: Point{10, 0}},
			p:    Point{11, 0},
			want: false,
		},
		{
			name: "point collinear but outside diagonal",
			line: Line{From: Point{0, 0}, To: Point{10, 10}},
			p:    Point{15, 15},
			want: false,
		},
		{
			name: "point slightly off diagonal",
			line: Line{From: Point{0, 0}, To: Point{10, 10}},
			p:    Point{5, 5.1},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.line.ContainsPoint(tt.p)
			if got != tt.want {
				t.Errorf("ContainsPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}