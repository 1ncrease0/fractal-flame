package infrastructure

import (
	"encoding/json"
	"errors"
	"fractal-flame/internal/application"
	"fractal-flame/internal/domain"
	"math/rand"
	"os"
)

type Config struct {
	Size struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"size"`

	IterationCount *int64  `json:"iteration_count,omitempty"`
	OutputPath     *string `json:"output_path,omitempty"`

	Seed            *int64   `json:"seed,omitempty"`
	Threads         *int     `json:"threads,omitempty"`
	SymmetryLevel   *int     `json:"symmetry_level,omitempty"`
	GammaCorrection *bool    `json:"gamma_correction,omitempty"`
	Gamma           *float64 `json:"gamma,omitempty"`

	Functions []struct {
		Name   string  `json:"name"`
		Weight float64 `json:"weight"`
	} `json:"functions"`

	AffineParams []struct {
		A      float64  `json:"a"`
		B      float64  `json:"b"`
		C      float64  `json:"c"`
		D      float64  `json:"d"`
		E      float64  `json:"e"`
		F      float64  `json:"f"`
		Weight *float64 `json:"weight,omitempty"`
		Color  *struct {
			R float64 `json:"r"`
			G float64 `json:"g"`
			B float64 `json:"b"`
		} `json:"color,omitempty"`
	} `json:"affine_params"`
}

var (
	ErrEmptyAffineParams    = errors.New("empty affine params")
	ErrEmptyVariationParams = errors.New("empty variation params")
)

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) BuildParams() (*application.Params, string, error) {
	width := DefaultWidth
	if cfg.Size.Width > 0 {
		width = cfg.Size.Width
	}

	height := DefaultHeight
	if cfg.Size.Height > 0 {
		height = cfg.Size.Height
	}

	seed := DefaultSeed
	if cfg.Seed != nil {
		seed = *cfg.Seed
	}

	iterations := DefaultIterationCount
	if cfg.IterationCount != nil && *cfg.IterationCount > 0 {
		iterations = *cfg.IterationCount
	}

	outputPath := DefaultOutputPath
	if cfg.OutputPath != nil && *cfg.OutputPath != "" {
		outputPath = *cfg.OutputPath
	}

	threads := DefaultThreads
	if cfg.Threads != nil && *cfg.Threads > 0 {
		threads = *cfg.Threads
	}

	symmetry := DefaultSymmetryLevel
	if cfg.SymmetryLevel != nil && *cfg.SymmetryLevel >= 1 {
		symmetry = *cfg.SymmetryLevel
	}

	gamma := DefaultGamma
	if cfg.GammaCorrection != nil && !*cfg.GammaCorrection {
		gamma = 1.0
	} else if cfg.Gamma != nil && *cfg.Gamma > 0 {
		gamma = *cfg.Gamma
	}

	affineList, err := cfg.BuildAffineList()
	if err != nil {
		return nil, "", err
	}

	variationList, err := cfg.BuildVariationList()
	if err != nil {
		return nil, "", err
	}

	if affineList.Len() == 0 {
		return nil, "", ErrEmptyAffineParams
	}

	if variationList.Len() == 0 {
		return nil, "", ErrEmptyVariationParams
	}

	params, err := application.NewParams(
		iterations,
		seed,
		symmetry,
		threads,
		width,
		height,
		gamma,
		*affineList,
		*variationList,
	)
	if err != nil {
		return nil, "", err
	}

	return params, outputPath, nil
}

func (cfg *Config) seed() int64 {
	if cfg.Seed != nil {
		return *cfg.Seed
	}
	return DefaultSeed
}

func (cfg *Config) BuildVariationList() (*domain.VariationList, error) {
	variationList := domain.NewVariationList()

	for _, f := range cfg.Functions {
		v, err := domain.GetVariation(f.Name)
		if err != nil {
			return nil, err
		}
		err = variationList.Add(v, f.Weight)
		if err != nil {
			return nil, err
		}
	}

	return variationList, nil
}

func (cfg *Config) BuildAffineList() (*domain.AffineList, error) {
	affineList := domain.NewAffineList()
	colorRnd := rand.New(rand.NewSource(cfg.seed()))

	for _, p := range cfg.AffineParams {
		weight := DefaultAffineWeight
		if p.Weight != nil {
			weight = *p.Weight
		}

		var color domain.Color
		if p.Color != nil {
			color = domain.Color{
				R: p.Color.R,
				G: p.Color.G,
				B: p.Color.B,
			}
		} else {
			color = domain.Color{
				R: colorRnd.Float64(),
				G: colorRnd.Float64(),
				B: colorRnd.Float64(),
			}
		}

		affine := domain.Affine{
			A:      p.A,
			B:      p.B,
			C:      p.C,
			D:      p.D,
			E:      p.E,
			F:      p.F,
			Weight: weight,
			Color:  color,
		}
		err := affineList.Add(affine)
		if err != nil {
			return nil, err
		}
	}

	return affineList, nil
}
