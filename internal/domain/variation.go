package domain

import (
	"errors"
	"fmt"
	"math"
)

const eps = 1e-12

var (
	ErrNegativeVariationWeight = errors.New("variation weight cannot be negative")
	ErrUnsupportedVariation    = errors.New("unsupported variation")
)

type Variation func(x, y float64) (float64, float64)

type VariationList struct {
	variations  []Variation
	weights     []float64
	totalWeight float64
}

func NewVariationList() *VariationList {
	return &VariationList{
		variations: make([]Variation, 0),
		weights:    make([]float64, 0),
	}
}

func (l *VariationList) Len() int {
	return len(l.variations)
}

func (l *VariationList) Add(variation Variation, weight float64) error {
	if weight < 0 {
		return ErrNegativeVariationWeight
	}
	l.weights = append(l.weights, weight)
	l.variations = append(l.variations, variation)
	l.totalWeight += weight
	return nil
}

func (l *VariationList) Apply(x, y float64) (xs float64, ys float64) {
	if l.totalWeight == 0 || len(l.variations) == 0 {
		vx, vy := x, y
		return vx, vy
	}

	for i, variation := range l.variations {
		vx, vy := variation(x, y)
		xs += vx * l.weights[i] / l.totalWeight
		ys += vy * l.weights[i] / l.totalWeight
	}
	return xs, ys
}

func (l *VariationList) TotalWeight() float64 {
	return l.totalWeight
}

func (l *VariationList) Variations() ([]Variation, []float64) {
	return l.variations, l.weights
}

var variations = map[string]Variation{
	"linear":       Linear,
	"sinusoidal":   Sinusoidal,
	"spherical":    Spherical,
	"swirl":        Swirl,
	"horseshoe":    Horseshoe,
	"polar":        Polar,
	"handkerchief": Handkerchief,
	"heart":        Heart,
	"disc":         Disc,
	"spiral":       Spiral,
	"hyperbolic":   Hyperbolic,
	"diamond":      Diamond,
	"ex":           Ex,
	"bent":         Bent,
	"fisheye":      Fisheye,
	"eyefish":      Eyefish,
	"bubble":       Bubble,
	"cylinder":     Cylinder,
	"tangent":      Tangent,
	"cross":        Cross,
	"power":        Power,
}

func GetVariation(name string) (Variation, error) {
	if v, ok := variations[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrUnsupportedVariation, name)
}

func GetVariationNames() []string {
	names := make([]string, 0)
	for k := range variations {
		names = append(names, k)
	}
	return names
}

func Linear(x, y float64) (float64, float64) {
	return x, y
}

func Sinusoidal(x, y float64) (float64, float64) {
	return math.Sin(x), math.Sin(y)
}

func Spherical(x, y float64) (float64, float64) {
	r2 := x*x + y*y
	if r2 < eps {
		return 0, 0
	}
	return x / r2, y / r2
}

func Swirl(x, y float64) (float64, float64) {
	r2 := x*x + y*y
	return x*math.Sin(r2) - y*math.Cos(r2), x*math.Cos(r2) + y*math.Sin(r2)
}

func Horseshoe(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	if r < eps {
		return 0, 0
	}
	return 1 / r * (x - y) * (x + y), 1 / r * 2 * x * y
}

func Disc(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	t := math.Atan2(y, x)
	return t / math.Pi * math.Sin(math.Pi*r), t / math.Pi * math.Cos(math.Pi*r)
}

func Diamond(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	t := math.Atan2(y, x)
	return math.Sin(t) * math.Cos(r), math.Sin(r) * math.Cos(t)
}

func Ex(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	t := math.Atan2(y, x)
	p0 := math.Sin(t + r)
	p1 := math.Cos(t - r)

	return r * (p0*p0*p0 + p1*p1*p1), r * (p0*p0*p0 - p1*p1*p1)
}

func Polar(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	t := math.Atan2(y, x)
	return t / math.Pi, r - 1
}

func Handkerchief(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	t := math.Atan2(y, x)
	return r * math.Sin(t+r), r * math.Cos(t-r)
}

func Heart(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	t := math.Atan2(y, x)
	return r * math.Sin(t*r), -r * math.Cos(t*r)
}

func Spiral(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	if r < eps {
		return 0, 0
	}
	t := math.Atan2(y, x)
	invR := 1.0 / r
	return invR * (math.Cos(t) + math.Sin(r)), invR * (math.Sin(t) - math.Cos(r))
}

func Hyperbolic(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	if r < eps {
		return 0, 0
	}
	t := math.Atan2(y, x)
	return math.Sin(t) / r, r * math.Cos(t)
}

func Bent(x, y float64) (float64, float64) {
	if x >= 0 && y >= 0 {
		return x, y
	} else if x < 0 && y >= 0 {
		return 2 * x, y
	} else if x >= 0 && y < 0 {
		return x, y / 2
	} else {
		return 2 * x, y / 2
	}
}

func Fisheye(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	if r < eps {
		return 0, 0
	}
	re := 2.0 / (r + 1)
	return re * y, re * x
}

func Eyefish(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	if r < eps {
		return 0, 0
	}
	re := 2.0 / (r + 1)
	return re * x, re * y
}

func Bubble(x, y float64) (float64, float64) {
	r2 := x*x + y*y
	if r2 < eps {
		return 0, 0
	}
	re := 4.0 / (r2 + 4)
	return re * x, re * y
}

func Cylinder(x, y float64) (float64, float64) {
	return math.Sin(x), y
}

func Tangent(x, y float64) (float64, float64) {
	cosY := math.Cos(y)
	if math.Abs(cosY) < eps {
		return 0, 0
	}
	return math.Sin(x) / cosY, math.Tan(y)
}

func Cross(x, y float64) (float64, float64) {
	diff2 := x*x - y*y
	if math.Abs(diff2) < eps {
		return 0, 0
	}
	s := math.Sqrt(1.0 / (diff2 * diff2))
	return s * x, s * y
}

func Power(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	if r < eps {
		return 0, 0
	}
	t := math.Atan2(y, x)
	rsth := math.Pow(r, math.Sin(t))
	return rsth * math.Cos(t), rsth * math.Sin(t)
}
