package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/rancher/security-scan/pkg/kb-summarizer/summarizer"
	cli "github.com/urfave/cli/v3"
)

const (
	K8SVersionFlag                = "k8s-version"
	BenchmarkVersionFlag          = "benchmark-version"
	ControlsDirFlag               = "controls-dir"
	InputDirFlag                  = "input-dir"
	OutputDirFlag                 = "output-dir"
	OutputFileNameFlag            = "output-filename"
	FailuresOnlyFlag              = "failures-only"
	UserSkipConfigFileFlag        = "user-skip-config-file"
	UserSkipConfigFileEnvVar      = "USER_SKIP_CONFIG_FILE"
	DefaultSkipConfigFileFlag     = "default-skip-config-file"
	DefaultSkipConfigFileEnvVar   = "DEFAULT_SKIP_CONFIG_FILE"
	NotApplicableConfigFileFlag   = "not-applicable-config-file"
	NotApplicableConfigFileEnvVar = "NOT_APPLICABLE_CONFIG_FILE"
)

var (
	VERSION = "v0.0.0-dev"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	app := &cli.Command{
		Name:    "kb-summarizer",
		Version: VERSION,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  K8SVersionFlag,
				Value: "",
			},
			&cli.StringFlag{
				Name:  BenchmarkVersionFlag,
				Value: "",
			},
			&cli.StringFlag{
				Name:  ControlsDirFlag,
				Value: summarizer.DefaultControlsDirectory,
			},
			&cli.StringFlag{
				Name:  InputDirFlag,
				Value: "",
			},
			&cli.StringFlag{
				Name:  OutputDirFlag,
				Value: "",
			},
			&cli.StringFlag{
				Name:  OutputFileNameFlag,
				Value: summarizer.DefaultOutputFileName,
			},
			&cli.StringFlag{
				Name:    UserSkipConfigFileFlag,
				Sources: cli.EnvVars(UserSkipConfigFileEnvVar),
				Value:   "",
			},
			&cli.StringFlag{
				Name:    DefaultSkipConfigFileFlag,
				Sources: cli.EnvVars(DefaultSkipConfigFileEnvVar),
				Value:   "",
			},
			&cli.StringFlag{
				Name:    NotApplicableConfigFileFlag,
				Sources: cli.EnvVars(NotApplicableConfigFileEnvVar),
				Value:   "",
			},
			&cli.BoolFlag{
				Name: FailuresOnlyFlag,
			},
		},
		Action: run,
	}

	if err := app.Run(context.TODO(), os.Args); err != nil {
		slog.Error("fatal error running application",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}

func run(ctx context.Context, c *cli.Command) error {
	slog.Info("Running Summarizer")
	k8sversion := c.String(K8SVersionFlag)
	benchmarkVersion := c.String(BenchmarkVersionFlag)
	controlsDir := c.String(ControlsDirFlag)
	inputDir := c.String(InputDirFlag)
	outputDir := c.String(OutputDirFlag)
	outputFilename := c.String(OutputFileNameFlag)
	failuresOnly := c.Bool(FailuresOnlyFlag)
	userSkipConfigFile := c.String(UserSkipConfigFileFlag)
	defaultSkipConfigFile := c.String(DefaultSkipConfigFileFlag)
	notApplicableConfigFile := c.String(NotApplicableConfigFileFlag)
	if k8sversion == "" && benchmarkVersion == "" {
		return fmt.Errorf("error: either of the flags %v, %v not specified", K8SVersionFlag, BenchmarkVersionFlag)
	}
	if k8sversion != "" && benchmarkVersion != "" {
		return fmt.Errorf("error: both flags %v, %v can not be specified at the same time", K8SVersionFlag, BenchmarkVersionFlag)
	}
	if controlsDir == "" {
		return fmt.Errorf("error: %v not specified", ControlsDirFlag)
	}
	if inputDir == "" {
		return fmt.Errorf("error: %v not specified", InputDirFlag)
	}
	if outputDir == "" {
		return fmt.Errorf("error: %v not specified", OutputDirFlag)
	}
	s, err := summarizer.NewSummarizer(
		k8sversion,
		benchmarkVersion,
		controlsDir,
		inputDir,
		outputDir,
		outputFilename,
		userSkipConfigFile,
		defaultSkipConfigFile,
		notApplicableConfigFile,
		failuresOnly,
	)
	if err != nil {
		return fmt.Errorf("error creating summarizer: %w", err)
	}
	if err := s.Summarize(); err != nil {
		return fmt.Errorf("error summarizing: %w", err)
	}
	return nil
}
