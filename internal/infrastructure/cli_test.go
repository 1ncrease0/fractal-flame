package infrastructure

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIParams_ParseAffineParams(t *testing.T) {
	t.Parallel()
	t.Run("valid single affine", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		affineStr := "1.0,0.0,0.0,0.0,1.0,0.0"
		cli.AffineParams = &affineStr
		seed := int64(42)
		cli.Seed = &seed

		affines, err := cli.ParseAffineParams()
		require.NoError(t, err)
		require.Len(t, affines, 1)
		assert.Equal(t, 1.0, affines[0].A)
		assert.Equal(t, 0.0, affines[0].B)
		assert.Equal(t, 0.0, affines[0].C)
		assert.Equal(t, 0.0, affines[0].D)
		assert.Equal(t, 1.0, affines[0].E)
		assert.Equal(t, 0.0, affines[0].F)
	})

	t.Run("valid multiple affines", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		affineStr := "1.0,0.0,0.0,0.0,1.0,0.0/2.0,0.0,0.0,0.0,2.0,0.0"
		cli.AffineParams = &affineStr
		seed := int64(42)
		cli.Seed = &seed

		affines, err := cli.ParseAffineParams()
		require.NoError(t, err)
		require.Len(t, affines, 2)
		assert.Equal(t, 1.0, affines[0].A)
		assert.Equal(t, 2.0, affines[1].A)
	})

	t.Run("invalid format - wrong number of values", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		affineStr := "1.0,0.0,0.0"
		cli.AffineParams = &affineStr
		seed := int64(42)
		cli.Seed = &seed

		affines, err := cli.ParseAffineParams()
		assert.Error(t, err)
		assert.Nil(t, affines)
	})

	t.Run("invalid format - non-numeric", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		affineStr := "a,b,c,d,e,f"
		cli.AffineParams = &affineStr
		seed := int64(42)
		cli.Seed = &seed

		affines, err := cli.ParseAffineParams()
		assert.Error(t, err)
		assert.Nil(t, affines)
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		affineStr := ""
		cli.AffineParams = &affineStr
		seed := int64(42)
		cli.Seed = &seed

		affines, err := cli.ParseAffineParams()
		require.NoError(t, err)
		assert.Len(t, affines, 0)
	})

	t.Run("with whitespace", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		affineStr := " 1.0 , 0.0 , 0.0 , 0.0 , 1.0 , 0.0 "
		cli.AffineParams = &affineStr
		seed := int64(42)
		cli.Seed = &seed

		affines, err := cli.ParseAffineParams()
		require.NoError(t, err)
		require.Len(t, affines, 1)
		assert.Equal(t, 1.0, affines[0].A)
	})
}

func TestCLIParams_ParseFunctions(t *testing.T) {
	t.Parallel()
	t.Run("valid single function", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		funcStr := "linear:1.0"
		cli.Functions = &funcStr

		funcs, err := cli.ParseFunctions()
		require.NoError(t, err)
		require.Len(t, funcs, 1)
		assert.Equal(t, "linear", funcs[0].Name)
		assert.Equal(t, 1.0, funcs[0].Weight)
	})

	t.Run("valid multiple functions", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		funcStr := "linear:1.0,sinusoidal:2.0"
		cli.Functions = &funcStr

		funcs, err := cli.ParseFunctions()
		require.NoError(t, err)
		require.Len(t, funcs, 2)
		assert.Equal(t, "linear", funcs[0].Name)
		assert.Equal(t, 1.0, funcs[0].Weight)
		assert.Equal(t, "sinusoidal", funcs[1].Name)
		assert.Equal(t, 2.0, funcs[1].Weight)
	})

	t.Run("invalid format - missing weight", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		funcStr := "linear"
		cli.Functions = &funcStr

		funcs, err := cli.ParseFunctions()
		assert.Error(t, err)
		assert.Nil(t, funcs)
	})

	t.Run("invalid format - non-numeric weight", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		funcStr := "linear:abc"
		cli.Functions = &funcStr

		funcs, err := cli.ParseFunctions()
		assert.Error(t, err)
		assert.Nil(t, funcs)
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		funcStr := ""
		cli.Functions = &funcStr

		funcs, err := cli.ParseFunctions()
		require.NoError(t, err)
		assert.Len(t, funcs, 0)
	})

	t.Run("with whitespace", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		funcStr := " linear : 1.0 "
		cli.Functions = &funcStr

		funcs, err := cli.ParseFunctions()
		require.NoError(t, err)
		require.Len(t, funcs, 1)
		assert.Equal(t, "linear", funcs[0].Name)
		assert.Equal(t, 1.0, funcs[0].Weight)
	})
}

func TestCLIParams_toConfig(t *testing.T) {
	t.Parallel()
	t.Run("full CLI params", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		width := 1920
		height := 1080
		seed := int64(42)
		iters := 1000000
		outputPath := "test.png"
		threads := 4
		symmetry := 2
		gamma := 2.2
		gammaCorrection := true

		cli.Width = &width
		cli.Height = &height
		cli.Seed = &seed
		cli.IterationCount = &iters
		cli.OutputPath = &outputPath
		cli.Threads = &threads
		cli.SymmetryLevel = &symmetry
		cli.Gamma = &gamma
		cli.GammaCorrection = &gammaCorrection

		affineStr := "1.0,0.0,0.0,0.0,1.0,0.0"
		cli.AffineParams = &affineStr
		funcStr := "linear:1.0"
		cli.Functions = &funcStr

		cfg, err := cli.toConfig()
		require.NoError(t, err)
		require.NotNil(t, cfg)
		assert.Equal(t, 1920, cfg.Size.Width)
		assert.Equal(t, 1080, cfg.Size.Height)
		assert.NotNil(t, cfg.Seed)
		assert.Equal(t, int64(42), *cfg.Seed)
		assert.NotNil(t, cfg.IterationCount)
		assert.Equal(t, int64(1000000), *cfg.IterationCount)
		assert.Equal(t, "test.png", *cfg.OutputPath)
		assert.Equal(t, 4, *cfg.Threads)
		assert.Equal(t, 2, *cfg.SymmetryLevel)
		assert.Equal(t, 2.2, *cfg.Gamma)
		assert.NotNil(t, cfg.GammaCorrection)
		assert.True(t, *cfg.GammaCorrection)
	})

	t.Run("minimal CLI params", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}

		cfg, err := cli.toConfig()
		require.NoError(t, err)
		require.NotNil(t, cfg)
		assert.Equal(t, 0, cfg.Size.Width)
		assert.Equal(t, 0, cfg.Size.Height)
		assert.Nil(t, cfg.Seed)
	})
}

func TestCLIParams_LoadConfig(t *testing.T) {
	t.Parallel()
	t.Run("from file", func(t *testing.T) {
		t.Parallel()
		configData := `{
			"size": {
				"width": 1920,
				"height": 1080
			},
			"seed": 42,
			"functions": [
				{"name": "linear", "weight": 1.0}
			],
			"affine_params": [
				{"a": 1.0, "b": 0.0, "c": 0.0, "d": 0.0, "e": 1.0, "f": 0.0}
			]
		}`

		tmpFile, err := os.CreateTemp("", "test_config_*.json")
		require.NoError(t, err)
		defer func(name string) {
			_ = os.Remove(name)
		}(tmpFile.Name())

		_, err = tmpFile.WriteString(configData)
		require.NoError(t, err)
		require.NoError(t, tmpFile.Close())

		cli := &CLIParams{}
		configPath := tmpFile.Name()
		cli.ConfigPath = &configPath

		cfg, err := cli.LoadConfig()
		require.NoError(t, err)
		require.NotNil(t, cfg)
		assert.Equal(t, 1920, cfg.Size.Width)
		assert.Equal(t, 1080, cfg.Size.Height)
	})

	t.Run("from CLI params", func(t *testing.T) {
		t.Parallel()
		cli := &CLIParams{}
		width := 1920
		height := 1080
		cli.Width = &width
		cli.Height = &height

		affineStr := "1.0,0.0,0.0,0.0,1.0,0.0"
		cli.AffineParams = &affineStr
		funcStr := "linear:1.0"
		cli.Functions = &funcStr

		cfg, err := cli.LoadConfig()
		require.NoError(t, err)
		require.NotNil(t, cfg)
		assert.Equal(t, 1920, cfg.Size.Width)
		assert.Equal(t, 1080, cfg.Size.Height)
	})
}

func TestParseCLI(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	t.Run("basic parsing", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = []string{"test", "-w", "1920", "-h", "1080", "-seed", "42"}

		params := ParseCLI()
		require.NotNil(t, params)
		if params.Width != nil {
			assert.Equal(t, 1920, *params.Width)
		}
		if params.Height != nil {
			assert.Equal(t, 1080, *params.Height)
		}
		if params.Seed != nil {
			assert.Equal(t, int64(42), *params.Seed)
		}
	})
}
