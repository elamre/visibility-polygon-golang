package reference

import "fmt"

func Example() {
	polygons := make([][][2]float32, 0)
	polygons = append(polygons, [][2]float32{{-1, -1}, {501, -1}, {501, 501}, {-1, 501}})
	polygons = append(polygons, [][2]float32{{250, 100}, {260, 140}, {240, 140}})
	segments := ConvertToSegments(polygons)
	segments = BreakIntersections(segments)
	position := [2]float32{60, 60}
	if InPolygon(position, polygons[0]) {
		visibility := Compute(position, segments)
		fmt.Printf("visibility: %+v\n", visibility)
	}
	viewPortVisbility := ComputeViewport(position, segments, [2]float32{50, 50}, [2]float32{450, 450})
	fmt.Printf("viewport: %+v\n", viewPortVisbility)
	// Output:
	// visibility: [[501 152.8421] [250 100] [240 140] [501 256] [501 501] [-1 501] [-1 -1] [501 -1]]
	// viewport: [[450.0001 142.10529] [250 100] [240 140] [450.0001 233.3334] [450.0001 450.0001] [49.9999 450.0001] [49.9999 49.9999] [450.0001 49.9999]]
}
