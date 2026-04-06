package infrastructure

import (
	"encoding/json"
	"fractal-flame/internal/application"
	"fractal-flame/internal/domain"
	"os"
	"reflect"
)

func ExportParamsToConfig(params *application.Params) (*Config, error) {
	cfg := &Config{}

	cfg.Size.Width = params.Width()
	cfg.Size.Height = params.Height()

	iterations := params.Iterations()
	cfg.IterationCount = &iterations

	outputPath := "random.png"
	cfg.OutputPath = &outputPath

	seed := params.Seed()
	cfg.Seed = &seed

	threads := params.Threads()
	cfg.Threads = &threads

	symmetry := params.Symmetry()
	cfg.SymmetryLevel = &symmetry

	gamma := params.Gamma()
	cfg.Gamma = &gamma

	if gamma == 1.0 {
		gammaCorrection := false
		cfg.GammaCorrection = &gammaCorrection
	}

	affineList := params.Affine()
	affine := affineList.Affine()
	cfg.AffineParams = make([]struct {
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
	}, len(affine))

	for i, a := range affine {
		cfg.AffineParams[i].A = a.A
		cfg.AffineParams[i].B = a.B
		cfg.AffineParams[i].C = a.C
		cfg.AffineParams[i].D = a.D
		cfg.AffineParams[i].E = a.E
		cfg.AffineParams[i].F = a.F
		weight := a.Weight
		cfg.AffineParams[i].Weight = &weight
		color := struct {
			R float64 `json:"r"`
			G float64 `json:"g"`
			B float64 `json:"b"`
		}{
			R: a.Color.R,
			G: a.Color.G,
			B: a.Color.B,
		}
		cfg.AffineParams[i].Color = &color
	}

	variationList := params.Variations()
	variations, weights := variationList.Variations()

	allVariationNames := domain.GetVariationNames()
	variationMap := make(map[uintptr]string)
	for _, name := range allVariationNames {
		v, _ := domain.GetVariation(name)
		variationMap[reflect.ValueOf(v).Pointer()] = name
	}

	cfg.Functions = make([]struct {
		Name   string  `json:"name"`
		Weight float64 `json:"weight"`
	}, 0, len(variations))

	for i, v := range variations {
		if i >= len(weights) {
			break
		}
		vPtr := reflect.ValueOf(v).Pointer()
		name, ok := variationMap[vPtr]
		if !ok {
			continue
		}
		cfg.Functions = append(cfg.Functions, struct {
			Name   string  `json:"name"`
			Weight float64 `json:"weight"`
		}{
			Name:   name,
			Weight: weights[i],
		})
	}

	return cfg, nil
}

func SaveConfig(cfg *Config, path string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
