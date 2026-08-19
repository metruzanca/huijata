package seed

import (
	"math"
)

type RandomPos struct {
	X int32
	Y int32
}

type Rng interface {
	Next() float64
	Random(min, max int32) int32
	SetRandomSeed(x, y float64)
	SetWorldSeed(ws uint32)
}

type NollaPrng struct {
	seed      float64
	worldSeed uint32
}

func NewNollaPrng(seed float64) NollaPrng {
	r := NollaPrng{seed: seed}
	r.Next()
	return r
}

func (r *NollaPrng) SetWorldSeed(ws uint32) {
	r.worldSeed = ws
}

func (r *NollaPrng) Next() float64 {
	s := int64(r.seed)
	v4 := 16807*s - 2147483647*(s/127773)
	if v4 <= 0 {
		v4 += 2147483647
	}
	r.seed = float64(v4)
	return r.seed / 2147483647.0
}

func (r *NollaPrng) NextU() uint32 {
	r.Next()
	return uint32(r.seed * 4.656612875e-10 * 2147483645.0)
}

func jsToInt32(f float64) uint32 {
	if math.IsNaN(f) || f == 0 {
		return 0
	}
	u := uint64(f) % (1 << 32)
	return uint32(int32(u))
}

func (r *NollaPrng) SetRandomFromWorldSeed(s float64) {
	r.seed = s
	if r.seed >= 2147483647.0 {
		r.seed = s * 0.5
	}
}

func (r *NollaPrng) SetRandomSeed(x, y float64) {
	ws := r.worldSeed
	a := ws ^ 0x93262e6f
	b := a & 0xfff
	c := (a >> 12) & 0xfff

	xAdjusted := x + float64(b)
	yAdjusted := y + float64(c)

	rVal := xAdjusted * 134217727.0
	e := jsToInt32(rVal)

	_y := math.Float64bits(yAdjusted) & 0x7fffffffffffffff
	_x := math.Float64bits(xAdjusted) & 0x7fffffffffffffff

	if math.Float64frombits(_y) >= 102400.0 || math.Float64frombits(_x) <= 1.0 {
		rVal = yAdjusted * 134217727.0
	} else {
		yWork := yAdjusted*3483.328 + float64(e)
		yAdjusted *= yWork
		rVal = yAdjusted
	}

	f := jsToInt32(rVal)

	g := helper2(uint32(e), uint32(f), ws)

	t := g
	if g < 0x80000000 {
		t += 1
	}
	if g == 0 {
		t += 1
	}
	t -= g / 252645135
	diddleTable := []uint32{0, 4, 6, 25, 12, 39, 52, 9, 21, 64, 78, 92, 104, 118, 18, 32, 44}
	idx := g / 252645135
	if idx >= 0 && int(idx) < len(diddleTable) {
		if g%252645135 < diddleTable[idx] && (g < 0xc3c3c3c3+4 || g >= 0xc3c3c3c3+62) {
			t += 1
		}
	}
	if g > 0x80000000 {
		t += 1
	}
	t = t >> 1
	if g == 0xffffffff {
		t += 1
	}

	r.seed = float64(t)
	r.Next()

	h := ws & 3
	for h > 0 {
		r.Next()
		h--
	}
}

func helper2(a, b, ws uint32) uint32 {
	v2 := ((a - b) - ws) ^ (ws >> 13)
	v1 := ((b - v2) - ws) ^ (v2 << 8)
	v3 := ((ws - v2) - v1) ^ (v1 >> 13)
	v2 = ((v2 - v1) - v3) ^ (v3 >> 12)
	v1 = ((v1 - v2) - v3) ^ (v2 << 16)
	v3 = ((v3 - v2) - v1) ^ (v1 >> 5)
	v2 = ((v2 - v1) - v3) ^ (v3 >> 3)
	v1 = ((v1 - v2) - v3) ^ (v2 << 10)
	return (((v3 - v2) - v1) ^ (v1 >> 15))
}

func (r *NollaPrng) Random(a, b int32) int32 {
	return a + int32(float64(b+1-a)*r.Next())
}

func randomNextF(worldSeed uint32, pos *RandomPos, min, max float64) float64 {
	var rng NollaPrng
	rng.worldSeed = worldSeed
	rng.SetRandomSeed(float64(pos.X), float64(pos.Y))
	pos.Y++
	return min + (max-min)*rng.Next()
}

func randomNextI(worldSeed uint32, pos *RandomPos, min, max int32) int32 {
	var rng NollaPrng
	rng.worldSeed = worldSeed
	rng.SetRandomSeed(float64(pos.X), float64(pos.Y))
	pos.Y++
	return rng.Random(min, max)
}
