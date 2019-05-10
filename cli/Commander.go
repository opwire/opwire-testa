package cli

import (
	"fmt"
	"os"
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

	testSourceFlags := []clp.Flag{
		clp.StringFlag{
			Name: "config-path, c",
			Usage: "Explicit configuration file",
		},
		clp.StringSliceFlag{
			Name: "test-dirs, spec-dirs, d",
			Usage: "The testcases directories",
		},
		clp.StringSliceFlag{
			Name: "incl-files, included-files, i",
			Usage: "Matching sub-string/pattern to include files",
		},
		clp.StringSliceFlag{
			Name: "excl-files, excluded-files, e",
			Usage: "Matching sub-string/pattern to exclude files",
		},
		clp.StringFlag{
			Name: "test-name, n",
			Usage: "Test title/name matching pattern",
		},
		clp.StringSliceFlag{
			Name: "tags, g",
			Usage: "Conditional tags for selecting tests",
		},
		clp.BoolFlag{
			Name: "no-color",
			Usage: "Display output in plain text, without color",
		},
	}

	app := clp.NewApp()
	app.Name = "opwire-testa"
	app.Usage = "Testing toolkit for opwire-agent"
	app.Version = manifest.GetVersion()

	app.Commands = []clp.Command {
		{
			Name: "run",
			Aliases: []string{"start"},
			Usage: "Run tests",
			Flags: append([]clp.Flag{}, testSourceFlags...),
			Action: func(c *clp.Context) error {
				o := readScriptSourceFlags(manifest, c)
				ctl, err := bootstrap.NewRunController(o)
				if err != nil {
					return err
				}
				ctl.Execute(&CmdRunFlags{})
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
						clp.StringFlag{
							Name: "export",
							Usage: "Output format (testcase syntax)",
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
						f.Format = c.String("export")
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
					Flags: append([]clp.Flag{}, testSourceFlags...),
					Action: func(c *clp.Context) error {
						o := readScriptSourceFlags(manifest, c)
						ctl, err := bootstrap.NewGenController(o)
						if err != nil {
							return err
						}
						ctl.Execute(&CmdGenFlags{})
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

func readScriptSourceFlags(manifest Manifest, c *clp.Context) *ControllerOptions {
	o := &ControllerOptions{ manifest: manifest }
	o.ConfigPath = c.String("config-path")
	o.TestDirs = c.StringSlice("test-dirs")
	o.InclFiles = c.StringSlice("incl-files")
	o.ExclFiles = c.StringSlice("excl-files")
	o.TestName = c.String("test-name")
	o.Tags = c.StringSlice("tags")
	o.NoColor = c.Bool("no-color")
	return o
}

type Manifest interface {
	GetRevision() string
	GetVersion() string
	String() (string, bool)
}

type ControllerOptions struct {
	ConfigPath string
	TestDirs []string
	InclFiles []string
	ExclFiles []string
	TestName string
	Tags []string
	NoColor bool
	manifest Manifest
}

func (a *ControllerOptions) GetConfigPath() string {
	return a.ConfigPath
}

func (a *ControllerOptions) GetTestDirs() []string {
	return a.TestDirs
}

func (a *ControllerOptions) GetInclFiles() []string {
	return a.InclFiles
}

func (a *ControllerOptions) GetExclFiles() []string {
	return a.ExclFiles
}

func (a *ControllerOptions) GetTestName() string {
	return a.TestName
}

func (a *ControllerOptions) GetConditionalTags() []string {
	return a.Tags
}

func (a *ControllerOptions) GetNoColor() bool {
	return a.NoColor
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
	Format string
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

func (f *CmdReqFlags) GetFormat() string {
	if f.Snapshot {
		return "testcase"
	}
	return f.Format
}

type CmdRunFlags struct {
}

type CmdGenFlags struct {
}
