// Copyright 2025 Redpanda Data, Inc.

package cli

import (
	"errors"

	"github.com/urfave/cli/v2"

	"github.com/redpanda-data/benthos/v4/internal/bloblang/parser"
	"github.com/redpanda-data/benthos/v4/internal/bloblang/plugins"
	"github.com/redpanda-data/benthos/v4/internal/bundle"
	"github.com/redpanda-data/benthos/v4/internal/cli/blobl"
	"github.com/redpanda-data/benthos/v4/internal/cli/common"
	"github.com/redpanda-data/benthos/v4/internal/cli/studio"
	clitemplate "github.com/redpanda-data/benthos/v4/internal/cli/template"
	"github.com/redpanda-data/benthos/v4/internal/cli/test"
	"github.com/redpanda-data/benthos/v4/internal/config"
	"github.com/redpanda-data/benthos/v4/internal/docs"
	"github.com/redpanda-data/benthos/v4/internal/filepath"
	"github.com/redpanda-data/benthos/v4/internal/filepath/ifs"
	"github.com/redpanda-data/benthos/v4/internal/template"
)

func runCliCommand(opts *common.CLIOpts) *cli.Command {
	flags := common.RunFlags(opts, false)
	flags = append(flags, common.EnvFileAndTemplateFlags(opts, false)...)

	return &cli.Command{
		Name:  "run",
		Usage: opts.ExecTemplate("Run {{.ProductName}} in normal mode against a specified config file"),
		Flags: flags,
		Before: func(c *cli.Context) error {
			return common.PreApplyEnvFilesAndTemplates(c, opts)
		},
		Description: opts.ExecTemplate(`
Run a {{.ProductName}} config.

  {{.BinaryName}} run ./foo.yaml`)[1:],
		Action: func(c *cli.Context) error {
			if c.Args().Len() > 0 {
				if c.Args().Len() > 1 || opts.RootFlags.Config != "" {
					return errors.New("a maximum of one config must be specified with the run command")
				}
				opts.RootFlags.Config = c.Args().First()
			}
			return common.RunService(c, opts, false)
		},
	}

	app.OnUsageError = func(context *cli.Context, err error, isSubcommand bool) error {
		fmt.Printf("Usage error: %v\n", err)
		_ = cli.ShowAppHelp(context)
		return err
	}
	return app
}

// Run the Benthos service, if the pipeline is started successfully then this
// call blocks until either the pipeline shuts down or a termination signal is
// received.
func Run(ctx context.Context) {
	if err := plugins.Register(); err != nil {
		panic(err)
	}

	_ = App().RunContext(ctx, os.Args)
}
