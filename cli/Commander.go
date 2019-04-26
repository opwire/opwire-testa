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
				o := &ControllerOptions{ manifest: manifest }
				o.ConfigPath = c.String("config-path")
				ctl, err := bootstrap.NewRunController(o)
				if err != nil {
					return err
				}
				f := &CmdRunFlags{
					TestDirs: c.StringSlice("test-dirs"),
				}
				ctl.Execute(f)
				return nil
			},
		},
		{
			Name: "req",
			Usage: "Make an HTTP request",
			Subcommands: []clp.Command{
				{
					Name: "curl",
					Usage: "Make an HTTP request using curl syntax",
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
						clp.BoolFlag{
							Name: "snapshot",
							Usage: "Create a snapshot of testcase",
						},
					},
					Action: func(c *clp.Context) error {
						o := &ControllerOptions{ manifest: manifest }
						broker, err := bootstrap.NewReqController(o)
						if err != nil {
							return err
						}
						f := new(CmdReqFlags)
						f.Method = c.String("request")
						f.Url = c.String("url")
						f.Header = c.StringSlice("header")
						f.Body = c.String("data")
						f.Snapshot = c.Bool("snapshot")
						broker.Execute(f)
						return nil
					},
				},
			},
		},
		{
			Name: "gen",
			Usage: "Generation commands",
			Subcommands: []clp.Command{
				{
					Name: "curl",
					Usage: "Generate curl style of a testcase",
					Flags: []clp.Flag{
						clp.StringFlag{
							Name: "config-path, c",
							Usage: "Explicit configuration file",
						},
						clp.StringSliceFlag{
							Name: "spec-dirs, test-dirs, d",
							Usage: "The testcases directories",
						},
						clp.StringFlag{
							Name: "test-file, f",
							Usage: "Suffix of path to testing script file (.i.e ... your-test.yml)",
						},
						clp.StringFlag{
							Name: "test-case, t",
							Usage: "Prefix of testcase title/name",
						},
					},
					Action: func(c *clp.Context) error {
						o := &ControllerOptions{ manifest: manifest }
						o.ConfigPath = c.String("config-path")
						ctl, err := bootstrap.NewGenController(o)
						if err != nil {
							return err
						}
						f := &CmdGenFlags{
							TestDirs: c.StringSlice("test-dirs"),
							TestFile: c.String("test-file"),
							TestCase: c.String("test-case"),
						}
						ctl.Execute(f)
						return nil
					},
				},
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

type ControllerOptions struct {
	ConfigPath string
	manifest Manifest
}

func (a *ControllerOptions) GetConfigPath() string {
	return a.ConfigPath
}

func (a *ControllerOptions) GetVersion() string {
	if a.manifest == nil {
		return ""
	}
	return utils.StandardizeVersion(a.manifest.GetVersion())
}

func (a *ControllerOptions) GetRevision() string {
	if a.manifest == nil {
		return ""
	}
	return a.manifest.GetRevision()
}

type CmdReqFlags struct {
	Method string
	Url string
	Header []string
	Body string
	Snapshot bool
}

func (f *CmdReqFlags) GetMethod() string {
	return f.Method
}

func (f *CmdReqFlags) GetUrl() string {
	return f.Url
}

func (f *CmdReqFlags) GetHeader() []string {
	return f.Header
}

func (f *CmdReqFlags) GetBody() string {
	return f.Body
}

func (f *CmdReqFlags) GetSnapshot() bool {
	return f.Snapshot
}

type CmdRunFlags struct {
	TestDirs []string
}

func (a *CmdRunFlags) GetTestDirs() []string {
	a.TestDirs = initDefaultDirs(a.TestDirs)
	return a.TestDirs
}

type CmdGenFlags struct {
	TestDirs []string
	TestFile string
	TestCase string
}

func (a *CmdGenFlags) GetTestDirs() []string {
	a.TestDirs = initDefaultDirs(a.TestDirs)
	return a.TestDirs
}

func (a *CmdGenFlags) GetTestFile() string {
	return a.TestFile
}

func (a *CmdGenFlags) GetTestCase() string {
	return a.TestCase
}

func initDefaultDirs(testDirs []string) []string {
	if testDirs == nil || len(testDirs) == 0 {
		testDir := filepath.Join(utils.FindWorkingDir(), "tests")
		if utils.IsDir(testDir) {
			testDirs = []string{testDir}
		}
	}
	return testDirs
}
