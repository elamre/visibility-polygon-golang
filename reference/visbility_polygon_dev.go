package reference

/*
	This is a direct translation from the javascript version. See:
	https://github.com/byronknoll/visibility-polygon-js/blob/master/visibility_polygon_dev.js
	Made by Made by Byron Knoll.
*/

import (
	"math"
	"sort"
)

type VisibilityPolygon struct{}

func Compute(position [2]float32, segments [][2][2]float32) [][2]float32 {
	var bounded = make([][2][2]float32, 0)
	var minX = position[0]
	var minY = position[1]
	var maxX = position[0]
	var maxY = position[1]
	for i := 0; i < len(segments); i++ {
		for j := 0; j < 2; j++ {
			minX = float32(math.Min(float64(minX), float64(segments[i][j][0])))
			minY = float32(math.Min(float64(minY), float64(segments[i][j][1])))
			maxX = float32(math.Max(float64(maxX), float64(segments[i][j][0])))
			maxY = float32(math.Max(float64(maxY), float64(segments[i][j][1])))
		}
		bounded = append(bounded, [2][2]float32{{segments[i][0][0], segments[i][0][1]}, {segments[i][1][0], segments[i][1][1]}})
	}
	minX--
	minY--
	maxX++
	maxY++

	bounded = append(bounded, [2][2]float32{{minX, minY}, {maxX, minY}})
	bounded = append(bounded, [2][2]float32{{maxX, minY}, {maxX, maxY}})
	bounded = append(bounded, [2][2]float32{{maxX, maxY}, {minX, maxY}})
	bounded = append(bounded, [2][2]float32{{minX, maxY}, {minX, minY}})
	polygon := make([][2]float32, 0)
	sorted := sortPoints(position, bounded)

	mmap := make([]int, len(bounded))
	for i := 0; i < len(mmap); i++ {
		mmap[i] = -1
	}

	heap := new([]int)
	start := [2]float32{position[0] + 1, position[1]}
	for i := 0; i < len(bounded); i++ {
		a1 := angle(bounded[i][0], position)
		a2 := angle(bounded[i][1], position)
		active := false
		if a1 > -180 && a1 <= 0 && a2 <= 180 && a2 >= 0 && a2-a1 > 180 {
			active = true
		}
		if a2 > -180 && a2 <= 0 && a1 <= 180 && a1 >= 0 && a1-a2 > 180 {
			active = true
		}
		if active {
			insert(i, heap, position, bounded, start, mmap)
		}
	}
	for i := 0; i < len(sorted); {
		extend := false
		shorten := false
		orig := i
		vertex := bounded[int(sorted[i][0])][int(sorted[i][1])]
		oldSegment := (*heap)[0]
		for ; ; {
			if mmap[int(sorted[i][0])] != -1 {
				if int(sorted[i][0]) == oldSegment {
					extend = true
					vertex = bounded[int(sorted[i][0])][int(sorted[i][1])]
				}

				remove(mmap[int(sorted[i][0])], (*[]int)(heap), position, bounded, vertex, mmap)
			} else {
				insert(int(sorted[i][0]), heap, position, bounded, vertex, mmap)
				if (*heap)[0] != oldSegment {
					shorten = true
				}
			}
			i++
			if i == len(sorted) {
				break
			}
			if sorted[i][2] >= (sorted[orig][2] + epsilon()) {
				break
			}
		}

		if extend {
			polygon = append(polygon, vertex)
			_, cur := intersectLines(bounded[(*heap)[0]][0], bounded[(*heap)[0]][1], position, vertex)
			if !equal(cur, vertex) {
				polygon = append(polygon, cur)

			}
		} else if shorten {
			_, p := intersectLines(bounded[oldSegment][0], bounded[oldSegment][1], position, vertex)
			polygon = append(polygon, p)
			_, p = intersectLines(bounded[(*heap)[0]][0], bounded[(*heap)[0]][1], position, vertex)
			polygon = append(polygon, p)
		}
	}
	return polygon
}

func ComputeViewport(position [2]float32, segments [][2][2]float32, viewportMinCorner [2]float32, viewportMaxCorner [2]float32) [][2]float32 {
	var brokenSegments = make([][2][2]float32, 0)
	viewport := [4][2]float32{{viewportMinCorner[0], viewportMinCorner[1]}, {viewportMaxCorner[0], viewportMinCorner[1]}, {viewportMaxCorner[0], viewportMaxCorner[1]}, {viewportMinCorner[0], viewportMaxCorner[1]}}
	for i := 0; i < len(segments); i++ {
		if segments[i][0][0] < viewportMinCorner[0] && segments[i][1][0] < viewportMinCorner[0] {
			continue
		}
		if segments[i][0][1] < viewportMinCorner[1] && segments[i][1][1] < viewportMinCorner[1] {
			continue
		}
		if segments[i][0][0] > viewportMaxCorner[0] && segments[i][1][0] > viewportMaxCorner[0] {
			continue
		}
		if segments[i][0][1] > viewportMaxCorner[1] && segments[i][1][1] > viewportMaxCorner[1] {
			continue
		}
		intersections := make([][2]float32, 0)
		for j := 0; j < len(viewport); j++ {
			k := j + 1
			if k == len(viewport) {
				k = 0
			}
			if doLineSegmentsIntersect(segments[i][0][0], segments[i][0][1], segments[i][1][0], segments[i][1][1], viewport[j][0], viewport[j][1], viewport[k][0], viewport[k][1]) {
				doIntersect, intersect := intersectLines(segments[i][0], segments[i][1], viewport[j], viewport[k])
				if !doIntersect {
					continue
				}
				if equal(intersect, segments[i][0]) || equal(intersect, segments[i][1]) {
					continue
				}
				intersections = append(intersections, intersect)
			}
		}
		var start = [2]float32{segments[i][0][0], segments[i][0][1]}
		for ; len(intersections) > 0; {
			var endIndex = 0
			var endDis = distance(start, intersections[0])
			for j := 1; j < len(intersections); j++ {
				var dis = distance(start, intersections[j])
				if dis < endDis {
					endDis = dis
					endIndex = j
				}
			}
			brokenSegments = append(brokenSegments, [2][2]float32{{start[0], start[1]}, {intersections[endIndex][0], intersections[endIndex][1]}})
			start[0] = intersections[endIndex][0]
			start[1] = intersections[endIndex][1]
			intersections = append(intersections[:endIndex], intersections[endIndex+1:]...)
		}
		brokenSegments = append(brokenSegments, [2][2]float32{start, {segments[i][1][0], segments[i][1][1]}})
	}

	viewportSegments := make([][2][2]float32, 0)
	for i := 0; i < len(brokenSegments); i++ {
		if InViewport(brokenSegments[i][0], viewportMinCorner, viewportMaxCorner) && InViewport(brokenSegments[i][1], viewportMinCorner, viewportMaxCorner) {
			viewportSegments = append(viewportSegments, [2][2]float32{{brokenSegments[i][0][0], brokenSegments[i][0][1]}, {brokenSegments[i][1][0], brokenSegments[i][1][1]}})
		}
	}
	var eps = epsilon() * 10
	viewportSegments = append(viewportSegments, [2][2]float32{{viewportMinCorner[0] - eps, viewportMinCorner[1] - eps}, {viewportMaxCorner[0] + eps, viewportMinCorner[1] - eps}})
	viewportSegments = append(viewportSegments, [2][2]float32{{viewportMaxCorner[0] + eps, viewportMinCorner[1] - eps}, {viewportMaxCorner[0] + eps, viewportMaxCorner[1] + eps}})
	viewportSegments = append(viewportSegments, [2][2]float32{{viewportMaxCorner[0] + eps, viewportMaxCorner[1] + eps}, {viewportMinCorner[0] - eps, viewportMaxCorner[1] + eps}})
	viewportSegments = append(viewportSegments, [2][2]float32{{viewportMinCorner[0] - eps, viewportMaxCorner[1] + eps}, {viewportMinCorner[0] - eps, viewportMinCorner[1] - eps}})
	return Compute(position, viewportSegments)
}

func InViewport(position [2]float32, viewportMinCorner [2]float32, viewportMaxCorner [2]float32) bool {
	if (position[0]) < viewportMinCorner[0]-epsilon() {
		return false
	}
	if (position[1]) < viewportMinCorner[1]-epsilon() {
		return false
	}
	if (position[0]) > viewportMaxCorner[0]+epsilon() {
		return false
	}
	if (position[1]) > viewportMaxCorner[1]+epsilon() {
		return false
	}
	return true
}

func angle(a, b [2]float32) float32 {
	return float32(math.Atan2(float64(b[1]-a[1]), float64(b[0]-a[0])) * 180 / math.Pi)
}

func angle2(a, b, c [2]float32) float32 {
	a1 := angle(a, b)
	a2 := angle(b, c)
	a3 := a1 - a2
	if a3 < 0 {
		a3 += 360
	}
	if a3 > 360 {
		a3 -= 360
	}
	return a3
}

func sortPoints(position [2]float32, segments [][2][2]float32) [][3]float32 {
	points := make([][3]float32, len(segments)*2)
	for i := 0; i < len(segments); i++ {
		for j := 0; j < 2; j++ {
			a := angle(segments[i][j], position)
			points[2*i+j] = [3]float32{float32(i), float32(j), a}
		}
	}
	sort.SliceStable(points, func(i, j int) bool {
		return points[i][2]-points[j][2] < 0
	})
	return points
}

func doLineSegmentsIntersect(x1, y1, x2, y2, x3, y3, x4, y4 float32) bool {
	d1 := computeDirection(x3, y3, x4, y4, x1, y1)
	d2 := computeDirection(x3, y3, x4, y4, x2, y2)
	d3 := computeDirection(x1, y1, x2, y2, x3, y3)
	d4 := computeDirection(x1, y1, x2, y2, x4, y4)
	return (((d1 > 0 && d2 < 0) || (d1 < 0 && d2 > 0)) &&
		((d3 > 0 && d4 < 0) || (d3 < 0 && d4 > 0))) ||
		(d1 == 0 && isOnSegment(x3, y3, x4, y4, x1, y1)) ||
		(d2 == 0 && isOnSegment(x3, y3, x4, y4, x2, y2)) ||
		(d3 == 0 && isOnSegment(x1, y1, x2, y2, x3, y3)) ||
		(d4 == 0 && isOnSegment(x1, y1, x2, y2, x4, y4))
}

func intersectLines(a1, a2, b1, b2 [2]float32) (bool, [2]float32) {
	dbx := b2[0] - b1[0]
	dby := b2[1] - b1[1]
	dax := a2[0] - a1[0]
	day := a2[1] - a1[1]

	u_b := dby*dax - dbx*day
	if u_b != 0 {
		ua := (dbx*(a1[1]-b1[1]) - dby*(a1[0]-b1[0])) / u_b
		return true, [2]float32{a1[0] - ua*-dax, a1[1] - ua*-day}
	}
	return false, [2]float32{}
}

func distance(a, b [2]float32) float32 {
	var dx = a[0] - b[0]
	var dy = a[1] - b[1]
	return dx*dx + dy*dy
}

func isOnSegment(xi, yi, xj, yj, xk, yk float32) bool {
	return (xi <= xk || xj <= xk) && (xk <= xi || xk <= xj) && (yi <= yk || yj <= yk) && (yk <= yi || yk <= yj)
}

func computeDirection(xi, yi, xj, yj, xk, yk float32) int {
	a := (xk - xi) * (yj - yi)
	b := (xj - xi) * (yk - yi)
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

func InPolygon(position [2]float32, polygon [][2]float32) bool {
	var val = polygon[0][0]
	for i := 0; i < len(polygon); i++ {
		val = float32(math.Min(float64(polygon[i][0]), float64(val)))
		val = float32(math.Min(float64(polygon[i][1]), float64(val)))
	}
	edge := [2]float32{val - 1, val - 1}
	var parity = 0
	for i := 0; i < len(polygon); i++ {
		var j = i + 1
		if j == len(polygon) {
			j = 0
		}
		if doLineSegmentsIntersect(edge[0], edge[1], position[0], position[1], polygon[i][0], polygon[i][1], polygon[j][0], polygon[j][1]) {
			_, intersect := intersectLines(edge, position, polygon[i], polygon[j])
			if equal(position, intersect) {
				return true
			}
			if equal(intersect, polygon[i]) {
				if angle2(position, edge, polygon[j]) < 180 {
					parity++
				}
			} else if equal(intersect, polygon[j]) {
				if angle2(position, edge, polygon[i]) < 180 {
					parity++
				}
			} else {
				parity++
			}
		}
	}
	return (parity % 2) != 0
}

func ConvertToSegments(polygons [][][2]float32) [][2][2]float32 {
	segments := make([][2][2]float32, 0)
	for i := 0; i < len(polygons); i++ {
		for j := 0; j < len(polygons[i]); j++ {
			k := j + 1
			if k == len(polygons[i]) {
				k = 0
			}
			segments = append(segments, [2][2]float32{{polygons[i][j][0], polygons[i][j][1]}, {polygons[i][k][0], polygons[i][k][1]}})
		}
	}
	return segments
}

func BreakIntersections(segments [][2][2]float32) [][2][2]float32 {
	output := make([][2][2]float32, 0)
	for i := 0; i < len(segments); i++ {
		intersections := make([][2]float32, 0)
		for j := 0; j < len(segments); j++ {
			if i == j {
				continue
			}
			if doLineSegmentsIntersect(segments[i][0][0], segments[i][0][1], segments[i][1][0], segments[i][1][1], segments[j][0][0], segments[j][0][1], segments[j][1][0], segments[j][1][1]) {
				doIntersect, intersect := intersectLines(segments[i][0], segments[i][1], segments[j][0], segments[j][1])
				if !doIntersect {
					continue
				}
				if equal(intersect, segments[i][0]) || equal(intersect, segments[i][1]) {
					continue
				}
				intersections = append(intersections, intersect)
			}
		}
		start := [2]float32{segments[i][0][0], segments[i][0][1]}
		for ; len(intersections) > 0; {
			endIndex := 0
			endDis := distance(start, intersections[0])
			for j := 1; j < len(intersections); j++ {
				dis := distance(start, intersections[j])
				if dis < endDis {
					endDis = dis
					endIndex = j
				}
			}
			output = append(output, [2][2]float32{{start[0], start[1]}, {intersections[endIndex][0], intersections[endIndex][1]}})
			start[0] = intersections[endIndex][0]
			start[1] = intersections[endIndex][1]
			intersections = append(intersections[:endIndex], intersections[endIndex+1:]...)
		}
		output = append(output, [2][2]float32{start, {segments[i][1][0], segments[i][1][1]}})
	}
	return output
}

func remove(index int, heap *[]int, position [2]float32, segments [][2][2]float32, destination [2]float32, mmap []int) {
	mmap[(*heap)[index]] = -1
	if index == len(*heap)-1 {
		temp := (*heap)[:len(*heap)-1]
		*heap = temp
		return
	}
	elem := (*heap)[len(*heap)-1]
	temp := (*heap)[:len(*heap)-1]
	*heap = temp
	(*heap)[index] = elem
	mmap[(*heap)[index]] = index
	var cur = index
	var curParent = parent(cur)
	if cur != 0 && lessThan((*heap)[cur], (*heap)[curParent], position, segments, destination) {
		for ; ; {
			var parent = parent(cur)
			if !lessThan((*heap)[cur], (*heap)[parent], position, segments, destination) {
				break
			}
			mmap[(*heap)[parent]] = cur
			mmap[(*heap)[cur]] = parent
			var temp = (*heap)[cur]
			(*heap)[cur] = (*heap)[parent]
			(*heap)[parent] = temp
			cur = parent
		}
	} else {
		for ; ; {
			var left = child(cur)
			var right = left + 1
			if left < len((*heap)) && lessThan((*heap)[left], (*heap)[cur], position, segments, destination) &&
				(right == len(*heap) || lessThan((*heap)[left], (*heap)[right], position, segments, destination)) {
				mmap[(*heap)[left]] = cur
				mmap[(*heap)[cur]] = left
				var temp = (*heap)[left]
				(*heap)[left] = (*heap)[cur]
				(*heap)[cur] = temp
				cur = left
			} else if right < len(*heap) && lessThan((*heap)[right], (*heap)[cur], position, segments, destination) {
				mmap[(*heap)[right]] = cur
				mmap[(*heap)[cur]] = right
				var temp = (*heap)[right]
				(*heap)[right] = (*heap)[cur]
				(*heap)[cur] = temp
				cur = right
			} else {
				break
			}
		}
	}
}

func insert(index int, heap *[]int, position [2]float32, segments [][2][2]float32, destination [2]float32, mmap []int) {
	found, _ := intersectLines(segments[index][0], segments[index][1], position, destination)
	if !found {
		return
	}
	cur := len(*heap)
	*heap = append(*heap, index)
	mmap[index] = cur
	for ; cur > 0; {
		parent := parent(cur)
		if !lessThan((*heap)[cur], (*heap)[parent], position, segments, destination) {
			break
		}
		mmap[(*heap)[parent]] = cur
		mmap[(*heap)[cur]] = parent
		temp := (*heap)[cur]
		(*heap)[cur] = (*heap)[parent]
		(*heap)[parent] = temp
		cur = parent
	}
}

func epsilon() float32 {
	return 0.00001
}

func equal(a, b [2]float32) bool {
	if float32(math.Abs(float64(a[0]-b[0]))) < epsilon() && float32(math.Abs(float64(a[1]-b[1]))) < epsilon() {
		return true
	}
	return false
}

func lessThan(index1, index2 int, position [2]float32, segments [][2][2]float32, destination [2]float32) bool {
	_, inter1 := intersectLines(segments[index1][0], segments[index1][1], position, destination)
	_, inter2 := intersectLines(segments[index2][0], segments[index2][1], position, destination)
	if !equal(inter1, inter2) {
		d1 := distance(inter1, position)
		d2 := distance(inter2, position)
		return d1 < d2
	}
	end1 := 0
	if equal(inter1, segments[index1][0]) {
		end1 = 1
	}
	end2 := 0
	if equal(inter2, segments[index2][0]) {
		end2 = 1
	}
	var a1 = angle2(segments[index1][end1], inter1, position)
	var a2 = angle2(segments[index2][end2], inter2, position)
	if a1 < 180 {
		if a2 > 180 {
			return true
		}
		return a2 < a1
	}
	return a1 < a2
}

func parent(index int) int {
	if index == 0 {
		return -1
	}
	return int(math.Floor(float64((index - 1) / 2)))
}

func child(index int) int {
	return 2*index + 1
}
