package cli

import (
	"fmt"
	"os"
	clp "github.com/urfave/cli"
	"github.com/opwire/opwire-qakit/lib/bootstrap"
)

type Commander struct {
	app *clp.App
}

func NewCommander(manifest Manifest) (*Commander, error) {
	if manifest == nil {
		return nil, fmt.Errorf("Manifest must not be nil")
	}

	c := new(Commander)

	clp.HelpFlag = clp.BoolFlag{
		Name: "help",
	}
	if info, ok := manifest.String(); ok {
		clp.AppHelpTemplate = fmt.Sprintf("%s\nNOTES:\n   %s\n\n", clp.AppHelpTemplate, info)
	}
	clp.VersionFlag = clp.BoolFlag{
		Name: "version",
	}
	clp.VersionPrinter = func(c *clp.Context) {
		fmt.Fprintf(c.App.Writer, "%s\n", c.App.Version)
	}

	app := clp.NewApp()
	app.Name = "opwire-qakit"
	app.Usage = "Testing toolkit for opwire-agent"
	app.Version = manifest.GetVersion()

	app.Commands = []clp.Command {
		{
			Name: "run",
			Aliases: []string{"start"},
			Usage: "Run the testcases",
			Flags: []clp.Flag{
				clp.StringFlag{
					Name: "config-path, c",
					Usage: "Explicit configuration file",
				},
				clp.StringSliceFlag{
					Name: "spec-dirs, test-dirs, d",
					Usage: "The testcases directories",
				},
			},
			Action: func(c *clp.Context) error {
				f := new(CmdFlags)
				f.ConfigPath = c.String("config-path")
				f.SpecDirs = c.StringSlice("spec-dirs")
				f.manifest = manifest
				tester, err := bootstrap.NewTestRunner(f)
				if err != nil {
					return err
				}
				tester.RunTests()
				return nil
			},
		},
		{
			Name: "help",
			Usage: "Shows a list of commands or help for one command",
		},
	}
	c.app = app
	return c, nil
}

func (c *Commander) Run() error {
	if c.app == nil {
		return fmt.Errorf("Commander has not initialized properly")
	}
	return c.app.Run(os.Args)
}

type Manifest interface {
	GetRevision() string
	GetVersion() string
	String() (string, bool)
}

type CmdFlags struct {
	ConfigPath string
	SpecDirs []string
	manifest Manifest
}

func (a *CmdFlags) GetConfigPath() string {
	return a.ConfigPath
}

func (a *CmdFlags) GetSpecDirs() []string {
	return a.SpecDirs
}

func (a *CmdFlags) GetRevision() string {
	if a.manifest == nil {
		return ""
	}
	return a.manifest.GetRevision()
}

func (a *CmdFlags) GetVersion() string {
	if a.manifest == nil {
		return ""
	}
	return a.manifest.GetVersion()
}
