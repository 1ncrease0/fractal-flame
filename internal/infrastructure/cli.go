package infrastructure

import (
	"flag"
	"fmt"
	"fractal-flame/internal/domain"
	"strconv"
	"strings"
)

type CLIParams struct {
	Width           *int
	Height          *int
	Seed            *int64
	IterationCount  *int
	OutputPath      *string
	Threads         *int
	AffineParams    *string
	Functions       *string
	ConfigPath      *string
	SymmetryLevel   *int
	GammaCorrection *bool
	Gamma           *float64
	Random          *bool
}

func ParseCLI() *CLIParams {
	params := &CLIParams{}

	params.Width = flag.Int("w", 0, "width of the image")
	widthLong := flag.Int("width", 0, "width of the image")
	params.Height = flag.Int("h", 0, "height of the image")
	heightLong := flag.Int("height", 0, "height of the image")
	params.Seed = flag.Int64("seed", 0, "seed for random generator")
	params.IterationCount = flag.Int("i", 0, "iteration count")
	iterationCountLong := flag.Int("iteration-count", 0, "iteration count")
	params.OutputPath = flag.String("o", "", "output path")
	outputPathLong := flag.String("output-path", "", "output path")
	params.Threads = flag.Int("t", 0, "number of threads")
	threadsLong := flag.Int("threads", 0, "number of threads")
	params.AffineParams = flag.String("ap", "", "affine params")
	affineParamsLong := flag.String("affine-params", "", "affine params")
	params.Functions = flag.String("f", "", "functions")
	functionsLong := flag.String("functions", "", "functions")
	params.ConfigPath = flag.String("config", "", "config file path")
	params.SymmetryLevel = flag.Int("s", 0, "symmetry level")
	symmetryLevelLong := flag.Int("symmetry-level", 0, "symmetry level")
	params.GammaCorrection = flag.Bool("g", false, "gamma correction")
	gammaCorrectionLong := flag.Bool("gamma-correction", false, "gamma correction")
	params.Gamma = flag.Float64("gamma", 0, "gamma value")
	params.Random = flag.Bool("r", false, "generate random params")
	randomLong := flag.Bool("random", false, "generate random params")

	flag.Parse()

	if *widthLong > 0 {
		params.Width = widthLong
	}
	if *heightLong > 0 {
		params.Height = heightLong
	}
	if *iterationCountLong > 0 {
		params.IterationCount = iterationCountLong
	}
	if *outputPathLong != "" {
		params.OutputPath = outputPathLong
	}
	if *threadsLong > 0 {
		params.Threads = threadsLong
	}
	if *affineParamsLong != "" {
		params.AffineParams = affineParamsLong
	}
	if *functionsLong != "" {
		params.Functions = functionsLong
	}
	if *symmetryLevelLong > 0 {
		params.SymmetryLevel = symmetryLevelLong
	}
	if *gammaCorrectionLong {
		params.GammaCorrection = gammaCorrectionLong
	}
	if *randomLong {
		params.Random = randomLong
	}

	return params
}

func (cli *CLIParams) LoadConfig() (*Config, error) {
	if cli.ConfigPath != nil && *cli.ConfigPath != "" {
		return LoadConfig(*cli.ConfigPath)
	}

	return cli.toConfig()
}

func (cli *CLIParams) toConfig() (*Config, error) {
	cfg := &Config{}

	if cli.Width != nil && *cli.Width > 0 {
		cfg.Size.Width = *cli.Width
	}
	if cli.Height != nil && *cli.Height > 0 {
		cfg.Size.Height = *cli.Height
	}
	if cli.Seed != nil && *cli.Seed != 0 {
		seed := *cli.Seed
		cfg.Seed = &seed
	}
	if cli.IterationCount != nil && *cli.IterationCount > 0 {
		iters := int64(*cli.IterationCount)
		cfg.IterationCount = &iters
	}
	if cli.OutputPath != nil && *cli.OutputPath != "" {
		cfg.OutputPath = cli.OutputPath
	}
	if cli.Threads != nil && *cli.Threads > 0 {
		cfg.Threads = cli.Threads
	}
	if cli.SymmetryLevel != nil && *cli.SymmetryLevel > 0 {
		cfg.SymmetryLevel = cli.SymmetryLevel
	}
	if cli.GammaCorrection != nil && *cli.GammaCorrection {
		cfg.GammaCorrection = cli.GammaCorrection
	}
	if cli.Gamma != nil && *cli.Gamma > 0 {
		cfg.Gamma = cli.Gamma
	}

	if cli.AffineParams != nil && *cli.AffineParams != "" {
		affine, err := cli.ParseAffineParams()
		if err != nil {
			return nil, err
		}
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
		}
	}

	if cli.Functions != nil && *cli.Functions != "" {
		funcs, err := cli.ParseFunctions()
		if err != nil {
			return nil, err
		}
		cfg.Functions = make([]struct {
			Name   string  `json:"name"`
			Weight float64 `json:"weight"`
		}, len(funcs))
		for i, f := range funcs {
			cfg.Functions[i].Name = f.Name
			cfg.Functions[i].Weight = f.Weight
		}
	}

	return cfg, nil
}

func (cli *CLIParams) ParseAffineParams() ([]domain.Affine, error) {
	parts := strings.Split(*cli.AffineParams, "/")
	affine := make([]domain.Affine, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		values := strings.Split(part, ",")
		if len(values) != 6 {
			return nil, fmt.Errorf("invalid affine params format: expected 6 values, got %d", len(values))
		}

		var a, b, c, d, e, f float64
		var err error

		if a, err = strconv.ParseFloat(strings.TrimSpace(values[0]), 64); err != nil {
			return nil, fmt.Errorf("invalid value for a: %v", err)
		}
		if b, err = strconv.ParseFloat(strings.TrimSpace(values[1]), 64); err != nil {
			return nil, fmt.Errorf("invalid value for b: %v", err)
		}
		if c, err = strconv.ParseFloat(strings.TrimSpace(values[2]), 64); err != nil {
			return nil, fmt.Errorf("invalid value for c: %v", err)
		}
		if d, err = strconv.ParseFloat(strings.TrimSpace(values[3]), 64); err != nil {
			return nil, fmt.Errorf("invalid value for d: %v", err)
		}
		if e, err = strconv.ParseFloat(strings.TrimSpace(values[4]), 64); err != nil {
			return nil, fmt.Errorf("invalid value for e: %v", err)
		}
		if f, err = strconv.ParseFloat(strings.TrimSpace(values[5]), 64); err != nil {
			return nil, fmt.Errorf("invalid value for f: %v", err)
		}

		affine = append(affine, domain.Affine{
			A:      a,
			B:      b,
			C:      c,
			D:      d,
			E:      e,
			F:      f,
			Weight: DefaultAffineWeight,
		})
	}

	return affine, nil
}

func (cli *CLIParams) ParseFunctions() ([]struct {
	Name   string
	Weight float64
}, error) {
	parts := strings.Split(*cli.Functions, ",")
	funcs := make([]struct {
		Name   string
		Weight float64
	}, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		parts2 := strings.Split(part, ":")
		if len(parts2) != 2 {
			return nil, fmt.Errorf("invalid function format: expected name:weight, got %s", part)
		}

		name := strings.TrimSpace(parts2[0])
		weight, err := strconv.ParseFloat(strings.TrimSpace(parts2[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid weight for function %s: %v", name, err)
		}

		funcs = append(funcs, struct {
			Name   string
			Weight float64
		}{
			Name:   name,
			Weight: weight,
		})
	}

	return funcs, nil
}
