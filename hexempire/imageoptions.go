package hexempire

// WaterImageOptions and LandImageOptions describe which sprite variant to
// draw for a field and how to orient it. They are chosen during map
// generation (so they're deterministic per seed) and consumed later by
// MapRenderer.
type WaterImageOptions struct {
	WaterBgImgIndex int
	FlipH           int
	FlipV           int
	RotateDegrees   int
}

type LandImageOptions struct {
	BgDirtImgIndex  int
	BgGrassImgIndex int
	FlipH           int
	FlipV           int
	RotateDegrees   int
}
