package ring

import (
	"math"

	"github.com/tidwall/boxtree/d2"
)

type treeRing struct {
	points []Point
	rect   Rect
	convex bool
	tree   d2.BoxTree
}

func newTreeRing(points []Point) *treeRing {
	var ring treeRing
	ring.points = make([]Point, len(points))
	copy(ring.points, points)
	ring.init()
	return &ring
}

func (ring *treeRing) Points() []Point {
	return ring.points
}

func (ring *treeRing) Search(
	rect Rect, iter func(seg Segment, index int) bool,
) {
	ring.tree.Search(
		[]float64{rect.Min.X, rect.Min.Y},
		[]float64{rect.Max.X, rect.Max.Y},
		func(_, _ []float64, value interface{}) bool {
			index := value.(int)
			var seg Segment
			seg.A = ring.points[index]
			if index == len(ring.points)-1 {
				seg.B = ring.points[0]
			} else {
				seg.B = ring.points[index+1]
			}
			if !iter(seg, index) {
				return false
			}
			return true
		},
	)
}

func (ring *treeRing) Scan(iter func(seg Segment) bool) {
	for i := 0; i < len(ring.points); i++ {
		var seg Segment
		seg.A = ring.points[i]
		if i == len(ring.points)-1 {
			if seg.A == ring.points[0] {
				break
			}
			seg.B = ring.points[0]
		} else {
			seg.B = ring.points[i+1]
		}
		if !iter(seg) {
			return
		}
	}
}

func (ring *treeRing) Rect() Rect {
	return ring.rect
}
func (ring *treeRing) Convex() bool {
	return ring.convex
}

func (ring *treeRing) init() {
	ring.convex, ring.rect = pointsConvexRect(ring.points)
	if false {
		for i := 0; i < len(ring.points); i++ {
			var seg Segment
			seg.A = ring.points[i]
			if i == len(ring.points)-1 {
				if seg.A == ring.points[0] {
					break
				}
				seg.B = ring.points[0]
			} else {
				seg.B = ring.points[i+1]
			}
			rect := seg.Rect()
			ring.tree.Insert(
				[]float64{rect.Min.X, rect.Min.Y},
				[]float64{rect.Max.X, rect.Max.Y},
				i,
			)
		}
	} else {
		var items []d2.BulkItem
		for i := 0; i < len(ring.points); i++ {
			var seg Segment
			seg.A = ring.points[i]
			if i == len(ring.points)-1 {
				if seg.A == ring.points[0] {
					break
				}
				seg.B = ring.points[0]
			} else {
				seg.B = ring.points[i+1]
			}

			rect := seg.Rect()
			items = append(items, d2.BulkItem{
				Min:   []float64{rect.Min.X, rect.Min.Y},
				Max:   []float64{rect.Max.X, rect.Max.Y},
				Value: i,
			})
		}
		ring.tree.BulkLoad(items)
	}
}

func (ring *treeRing) IntersectsSegment(seg Segment, allowOnEdge bool) bool {
	var intersects bool
	rect := seg.Rect()
	ring.tree.Search(
		[]float64{rect.Min.X, rect.Min.Y},
		[]float64{rect.Max.X, rect.Max.Y},
		func(_, _ []float64, value interface{}) bool {
			index := value.(int)
			var other Segment
			other.A = ring.points[index]
			if index == len(ring.points)-1 {
				other.B = ring.points[0]
			} else {
				other.B = ring.points[index+1]
			}
			if segmentsIntersect(seg.A, seg.B, other.A, other.B) {
				if !allowOnEdge {
					if raycast(seg.A, other.A, other.B).on ||
						raycast(seg.B, other.A, other.B).on ||
						raycast(other.A, seg.A, seg.B).on ||
						raycast(other.B, seg.A, seg.B).on {
						intersects = false
						return false
					}
				}
				intersects = true
				return false
			}
			return true
		})
	return intersects
}

func (ring *treeRing) ContainsPoint(point Point, allowOnEdge bool) bool {
	rect := Rect{
		Min: Point{math.Inf(-1), point.Y},
		Max: Point{math.Inf(+1), point.Y},
	}
	in := false
	ring.Search(rect, func(seg Segment, index int) bool {
		res := raycast(point, seg.A, seg.B)
		if res.on {
			in = allowOnEdge
			return false
		}
		if res.in {
			in = !in
		}
		return true
	})
	return in
}

func (ring *treeRing) ContainsRing(other Ring, allowOnEdge bool) bool {
	return ringContainsRing(ring, other, allowOnEdge)
}

func (ring *treeRing) IntersectsRing(other Ring, allowOnEdge bool) bool {
	return ringIntersectsRing(ring, other, allowOnEdge)
}
