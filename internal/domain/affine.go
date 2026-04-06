package domain

import (
	"errors"
	"math/rand"
)

type Affine struct {
	A      float64
	B      float64
	C      float64
	D      float64
	E      float64
	F      float64
	Color  Color
	Weight float64
}

func (a Affine) Apply(x, y float64) (float64, float64) {
	return a.A*x + a.B*y + a.C, a.D*x + a.E*y + a.F
}

type AffineList struct {
	affine      []Affine
	totalWeight float64
}

func NewAffineList() *AffineList {
	return &AffineList{
		affine: make([]Affine, 0),
	}
}

var ErrNegativeAffineWeight = errors.New("affine weight cannot be negative")

func (l *AffineList) Add(affine Affine) error {
	if affine.Weight < 0 {
		return ErrNegativeAffineWeight
	}
	l.affine = append(l.affine, affine)
	l.totalWeight += affine.Weight
	return nil
}

func (l *AffineList) Len() int {
	return len(l.affine)
}

func (l *AffineList) Affine() []Affine {
	return l.affine
}

func (l *AffineList) Random(rnd *rand.Rand) *Affine {
	if len(l.affine) == 0 {
		return nil
	}
	r := rnd.Float64() * l.totalWeight
	sum := 0.0
	for i, w := range l.affine {
		sum += w.Weight
		if r <= sum {
			return &l.affine[i]
		}
	}
	return &l.affine[len(l.affine)-1]

}
