package cli

import (
	"fmt"
	"os"
	"path/filepath"
	clp "github.com/urfave/cli"
	"github.com/opwire/opwire-testa/lib/bootstrap"
	"github.com/opwire/opwire-testa/lib/utils"
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
	app.Name = "opwire-testa"
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
			Name: "req",
			Aliases: []string{"curl"},
			Usage: "Make an HTTP request",
			Flags: []clp.Flag{
				clp.StringFlag{
					Name: "request, X",
					Usage: "Specify request command to use",
				},
				clp.StringFlag{
					Name: "url",
					Usage: "URL to work with",
				},
				clp.StringSliceFlag{
					Name: "header, H",
					Usage: "Pass custom header(s) to server",
				},
				clp.StringFlag{
					Name: "data, d",
					Usage: "HTTP POST data",
				},
			},
			Action: func(c *clp.Context) error {
				f := new(CmdReqFlags)
				f.Method = c.String("request")
				f.Url = c.String("url")
				f.Header = c.StringSlice("header")
				f.Body = c.String("data")
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

type CmdReqFlags struct {
	Method string
	Url string
	Header []string
	Body string
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
	if a.SpecDirs == nil || len(a.SpecDirs) == 0 {
		testDir := filepath.Join(utils.FindWorkingDir(), "tests")
		if utils.IsDir(testDir) {
			a.SpecDirs = []string{testDir}
		}
	}
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
