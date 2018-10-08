package geojson

import "testing"

func TestFeature(t *testing.T) {
	p := expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]}}`, nil)
	if p.Center() != P(1, 2) {
		t.Fatalf("expected '%v', got '%v'", P(1, 2), p.Center())
	}
	expectJSON(t, `{"type":"Feature"}`, errGeometryMissing)
	expectJSON(t, `{"type":"Feature","geometry":null}`, errDataInvalid)
	expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"bbox":null}`, errBBoxInvalid)
}

func TestFeaturePoly(t *testing.T) {
	p := expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]}}`, nil)
	expect(t, p.Intersects(P(1, 2)))
	expect(t, p.Contains(P(1, 2)))
	expect(t, p.Within(P(1, 2)))

}