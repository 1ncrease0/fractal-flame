package main

import (
	"fractal-flame/internal/application"
	"fractal-flame/internal/infrastructure"
	"log/slog"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cli := infrastructure.ParseCLI()

	if cli.Random != nil && *cli.Random {
		params := application.RandomParams()
		if params == nil {
			logger.Error("Error generating random params")
			os.Exit(1)
		}
		const randomImagePath = "random.png"
		const randomConfigPath = "config.json"

		logger.Info("Generated random params", "seed", params.Seed())

		start := time.Now()

		renderer := application.NewRenderer(*params, logger)

		histogram := renderer.Render()

		logger.Info("Applying image correction")

		logger.Info("Rendering complete", "duration", time.Since(start))

		img := infrastructure.ToImage(histogram.Correction())
		if err := infrastructure.SaveImage(img, randomImagePath); err != nil {
			logger.Error("Error saving image", "error", err)
			os.Exit(1)
		}

		cfg, err := infrastructure.ExportParamsToConfig(params)
		if err != nil {
			logger.Error("Error exporting config", "error", err)
			os.Exit(1)
		}

		if err := infrastructure.SaveConfig(cfg, randomConfigPath); err != nil {
			logger.Error("Error saving config", "error", err)
			os.Exit(1)
		}

		logger.Info("Image successfully saved", "path", randomImagePath)
		logger.Info("Config successfully saved", "path", randomConfigPath)
		return
	}

	cfg, err := cli.LoadConfig()
	if err != nil {
		logger.Error("Error loading config", "error", err)
		os.Exit(2)
	}

	params, outputPath, err := cfg.BuildParams()
	if err != nil {
		logger.Error("Error building params", "error", err)
		os.Exit(2)
	}

	logger.Info("Starting rendering", "seed", params.Seed(), "iterations", params.Iterations(),
		"threads", params.Threads(), "size", slog.Group("size", "width", params.Width(), "height", params.Height()))

	start := time.Now()

	renderer := application.NewRenderer(*params, logger)

	histogram := renderer.Render()

	logger.Info("Applying gamma correction")
	histogram.Correction()

	logger.Info("Rendering complete", "duration", time.Since(start))

	img := infrastructure.ToImage(histogram.Correction())
	if err := infrastructure.SaveImage(img, outputPath); err != nil {
		logger.Error("Error saving image", "error", err, "path", outputPath)
		os.Exit(1)
	}

	logger.Info("Image successfully saved", "path", outputPath)
}
