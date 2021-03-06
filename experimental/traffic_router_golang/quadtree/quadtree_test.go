package quadtree

import (
	"testing"
)

func expectNearest(lat float64, lon float64, expected DataT, qt *Quadtree, t *testing.T) {
	d, ok := qt.Nearest(lat, lon)
	if !ok {
		t.Errorf("no nearest")
	}
	if d.Obj != expected.Obj {
		t.Errorf("Nearest %f,%f expected '%v', actual '%v'", lat, lon, expected.Obj, d.Obj)
	}
	if d.Lat != expected.Lat {
		t.Errorf("Nearest %f,%s expected latitude '%f', actual '%f'", lat, lon, expected.Lat, d.Lat)
	}
	if d.Lon != expected.Lon {
		t.Errorf("Nearest %f,%s expected longitude '%f', actual '%f'", lat, lon, expected.Lon, d.Lon)
	}
}

func TestNearest(t *testing.T) {
	pts := []DataT{
		{100, 100, "a"},
		{90, 60, "b"},
		{50, 50, "c"},
		{49, 40, "d"},
		{40, 49, "e"},
	}
	qt := New()
	for _, pt := range pts {
		qt.Insert(pt)
	}

	expectNearest(100, 100, pts[0], qt, t)
	expectNearest(100, 99, pts[0], qt, t)
	expectNearest(99, 99, pts[0], qt, t)

	expectNearest(90, 50, pts[1], qt, t)
	expectNearest(80, 59, pts[1], qt, t)

	expectNearest(50, 50, pts[2], qt, t)
	expectNearest(51, 51, pts[2], qt, t)

	expectNearest(49, 40, pts[3], qt, t)
	expectNearest(48, 41, pts[3], qt, t)

	expectNearest(40, 49, pts[4], qt, t)
	expectNearest(41, 48, pts[4], qt, t)
}

func expectNearestIn(lat float64, lon float64, expecteds []DataT, qt *Quadtree, t *testing.T) {
	d, ok := qt.Nearest(lat, lon)
	if !ok {
		t.Errorf("no nearest")
	}

	found := false
	for _, expected := range expecteds {
		if d.Obj == expected.Obj && d.Lat == expected.Lat && d.Lon == expected.Lon {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Nearest %f,%f expected in %+v actual %+v", lat, lon, expecteds, d)
	}
}

func TestDuplicatePos(t *testing.T) {
	pts := []DataT{
		{100.0, 100.0, "a"},
		{100.0, 100.0, "b"},
		{0.0, 0.0, "c"},
		{0, 0, "d"},
		{39.578968, -104.934333, "e"},
		{39.578968, -104.934333, "f"},
		{39.578967, -104.934332, "g"},
	}
	qt := New()
	for _, pt := range pts {
		qt.Insert(pt)
	}

	expectNearestIn(100, 100, []DataT{pts[0], pts[1]}, qt, t)
	expectNearestIn(0, 0, []DataT{pts[2], pts[3]}, qt, t)
	expectNearestIn(39.578968, -104.934333, []DataT{pts[4], pts[5]}, qt, t)
	expectNearest(39.578967, -104.934332, pts[6], qt, t)
	expectNearest(39.578967, -104.934330, pts[6], qt, t)
}
