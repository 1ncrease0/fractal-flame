package application

import (
	"errors"
	"fractal-flame/internal/domain"
	"math/rand"
	"time"
)

type Params struct {
	affine     domain.AffineList
	variations domain.VariationList
	iterations int64
	symmetry   int
	threads    int
	gamma      float64
	seed       int64
	width      int
	height     int
}

var (
	ErrInvalidIters    = errors.New("invalid iterations count")
	ErrInvalidSymmetry = errors.New("invalid symmetry")
	ErrEmptyAffineList = errors.New("empty affine list")
	ErrInvalidThreads  = errors.New("invalid number of threads")
	ErrInvalidSize     = errors.New("invalid number of size")
)

func NewParams(iters, seed int64, symmetry, threads, width, height int, gamma float64, affine domain.AffineList, variations domain.VariationList) (*Params, error) {
	if iters <= 0 {
		return nil, ErrInvalidIters
	}
	if symmetry < 1 {
		return nil, ErrInvalidSymmetry
	}
	if threads <= 0 {
		return nil, ErrInvalidThreads
	}
	if affine.Len() == 0 {
		return nil, ErrEmptyAffineList
	}
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidSize
	}

	return &Params{
		affine:     affine,
		variations: variations,
		iterations: iters,
		symmetry:   symmetry,
		threads:    threads,
		gamma:      gamma,
		seed:       seed,
		width:      width,
		height:     height,
	}, nil
}

func RandomParams() *Params {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	affineList := domain.NewAffineList()
	nAffine := r.Intn(6) + 3

	for i := 0; i < nAffine; i++ {
		affine := domain.Affine{
			A:      r.Float64()*2 - 1,
			B:      r.Float64()*2 - 1,
			C:      r.Float64()*2 - 1,
			D:      r.Float64()*2 - 1,
			E:      r.Float64()*2 - 1,
			F:      r.Float64()*2 - 1,
			Weight: 0.2 + r.Float64()*0.8,
			Color: domain.Color{
				R: r.Float64(),
				G: r.Float64(),
				B: r.Float64(),
			},
		}
		_ = affineList.Add(affine)
	}

	variationList := domain.NewVariationList()
	allNames := domain.GetVariationNames()
	r.Shuffle(len(allNames), func(i, j int) { allNames[i], allNames[j] = allNames[j], allNames[i] })

	nVar := r.Intn(3) + 2
	for i := 0; i < nVar; i++ {
		v, _ := domain.GetVariation(allNames[i])
		_ = variationList.Add(v, 0.2+r.Float64()*0.8)
	}

	symmetry := r.Intn(5) + 1
	threads := 4

	iters := int64(5000000 + r.Intn(10000000))
	gamma := 2.2

	width := 1920
	height := 1080

	seedNew := time.Now().UnixNano()
	params, err := NewParams(iters, seedNew, symmetry, threads, width, height, gamma, *affineList, *variationList)
	if err != nil {
		return nil
	}
	return params
}

func (p *Params) Affine() domain.AffineList {
	return p.affine
}

func (p *Params) Variations() domain.VariationList {
	return p.variations
}

func (p *Params) Iterations() int64 {
	return p.iterations
}

func (p *Params) Symmetry() int {
	return p.symmetry
}

func (p *Params) Threads() int {
	return p.threads
}

func (p *Params) Gamma() float64 {
	return p.gamma
}

func (p *Params) Seed() int64 {
	return p.seed
}

func (p *Params) Width() int {
	return p.width
}

func (p *Params) Height() int {
	return p.height
}
