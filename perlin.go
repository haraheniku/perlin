package perlin

import (
	"math"
	"math/rand"
)

const (
	defaultOctaves     = 8
	defaultPersistence = 0.5
)

var defaultNoise = New()

func Noise1D(x float64) float64       { return defaultNoise.Noise1D(x) }
func Noise2D(x, y float64) float64    { return defaultNoise.Noise2D(x, y) }
func Noise3D(x, y, z float64) float64 { return defaultNoise.Noise3D(x, y, z) }
func Seed(seed int64)                 { defaultNoise.Seed(seed) }

type Perlin struct {
	p           []int
	Octaves     int
	Persistence float64
}

func New() *Perlin {
	perlin := &Perlin{
		Octaves:     defaultOctaves,
		Persistence: defaultPersistence,
	}
	perlin.Seed(0)
	return perlin
}

func (perlin *Perlin) Seed(seed int64) {
	src := rand.NewSource(seed)
	rng := rand.New(src)

	p := make([]int, 512)
	copy(p, rng.Perm(256))
	for i := 0; i < 256; i++ {
		p[256+i] = p[i]
	}
	perlin.p = p
}

func (perlin *Perlin) Noise1D(x float64) float64 {
	return perlin.Noise3D(x, 0, 0)
}

func (perlin *Perlin) Noise2D(x, y float64) float64 {
	return perlin.Noise3D(x, y, 0)
}

func (perlin *Perlin) Noise3D(x, y, z float64) float64 {
	var (
		total     float64
		maxValue  float64
		frequency float64 = 1
		amplitude float64 = 1
	)
	for i := 0; i < perlin.Octaves; i++ {
		total += perlin.noise(x*frequency, y*frequency, z*frequency) * amplitude
		maxValue += amplitude
		amplitude *= perlin.Persistence
		frequency *= 2
	}
	return total / maxValue
}

func (perlin *Perlin) noise(x, y, z float64) float64 {
	xi := int(x) & 255
	yi := int(y) & 255
	zi := int(z) & 255

	xf := x - math.Floor(x)
	yf := y - math.Floor(y)
	zf := z - math.Floor(z)

	u := fade(xf)
	v := fade(yf)
	w := fade(zf)

	p := perlin.p
	a := p[xi] + yi
	aa := p[a] + zi
	ab := p[a+1] + zi
	b := p[xi+1] + yi
	ba := p[b] + zi
	bb := p[b+1] + zi

	x1 := lerp(u, grad(p[aa], xf, yf, zf), grad(p[ba], xf-1, yf, zf))
	x2 := lerp(u, grad(p[ab], xf, yf-1, zf), grad(p[bb], xf-1, yf-1, zf))
	y1 := lerp(v, x1, x2)
	x3 := lerp(u, grad(p[aa+1], xf, yf, zf-1), grad(p[ba+1], xf-1, yf, zf-1))
	x4 := lerp(u, grad(p[ab+1], xf, yf-1, zf-1), grad(p[bb+1], xf-1, yf-1, zf-1))
	y2 := lerp(v, x3, x4)
	return (lerp(w, y1, y2) + 1) / 2
}

func fade(t float64) float64 {
	// 6t^5 - 15t^4 + 10t^3x
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func grad(hash int, x, y, z float64) float64 {
	var u, v float64
	h := hash & 15
	if h < 8 {
		u = x
	} else {
		u = y
	}
	if h < 4 {
		v = y
	} else if h == 12 || h == 14 {
		v = x
	} else {
		v = z
	}

	var g float64
	if (h & 1) == 0 {
		g += u
	} else {
		g -= u
	}
	if (h & 2) == 0 {
		g += v
	} else {
		g -= v
	}
	return g
}
