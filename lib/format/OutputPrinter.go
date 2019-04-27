package format

import (
	"fmt"
	"strings"
	"github.com/gookit/color"
	"github.com/opwire/opwire-testa/lib/utils"
)

type OutputPrinterOptions interface {
	GetNoColor() bool
}

type OutputPrinter struct {
	options OutputPrinterOptions
}

func NewOutputPrinter(opts OutputPrinterOptions) (ref *OutputPrinter, err error) {
	ref = new(OutputPrinter)
	ref.options = opts
	return ref, err
}

func (w *OutputPrinter) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Printf(format, args...)
}

func (w *OutputPrinter) Println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

func (w *OutputPrinter) Heading(title string) string {
	pen := w.GetPen(HeadingPen)
	return pen(title)
}

func (w *OutputPrinter) ContextInfo(name string, values ...string) string {
	pen := w.GetPen(ContextNamePen)
	text := "[+] " + pen(name)
	number := len(values)
	if number == 1 && len(values[0]) > 0 {
		text = fmt.Sprintf("%s: %s", text, values[0])
	}
	if number > 1 {
		lines := []string{text + ":"}
		for _, val := range values {
			lines = append(lines, "    - " + val)
		}
		text = strings.Join(lines, "\n")
	}
	return text
}

func (w *OutputPrinter) TestSuiteTitle(filepath string) string {
	pen := w.GetPen(TestSuiteTitlePen)
	return fmt.Sprintf("[#] %s", pen(filepath))
}

func (w *OutputPrinter) TestCase(title string) string {
	pen := w.GetPen(TestCaseTitlePen)
	return fmt.Sprintf("[=] %s", pen(title))
}

func (w *OutputPrinter) Skipped(title string) string {
	pen := w.GetPen(SkippedPen)
	return fmt.Sprintf("[%s] %s", pen("-"), title)
}

func (w *OutputPrinter) Success(title string) string {
	pen := w.GetPen(SuccessPen)
	return fmt.Sprintf("[%s] %s", pen("v"), title)
}

func (w *OutputPrinter) Failure(title string) string {
	pen := w.GetPen(FailurePen)
	return fmt.Sprintf("[%s] %s", pen("x"), title)
}

func (w *OutputPrinter) Cracked(title string) string {
	pen := w.GetPen(CrackedPen)
	return fmt.Sprintf("[%s] %s", pen("~"), title)
}

func (w *OutputPrinter) SectionTitle(title string) string {
	return fmt.Sprintf("--- %s", title)
}

func (w *OutputPrinter) Section(block string) string {
	lines := strings.Split(block, "\n")
	lines = utils.Map(lines, func(line string, number int) string {
		return "    " + line
	})
	return strings.Join(lines, "\n")
}

func (w *OutputPrinter) IsColorized() bool {
	if w.options != nil && w.options.GetNoColor() {
		return false
	}
	return true
}

func (w *OutputPrinter) GetPen(name PenType) Renderer {
	pen := ColorlessPen
	if w.IsColorized() {
		if Pens == nil {
			Pens = make(map[PenType]Renderer, 0)
		}
		if val, ok := Pens[name]; ok {
			pen = val
		} else {
			switch(name) {
			case HeadingPen:
				pen = color.Style{color.FgCyan, color.OpUnderscore, color.OpBold}.Render
			case ContextNamePen:
				pen = color.Style{color.FgCyan, color.OpBold}.Render
			case TestSuiteTitlePen:
				pen = color.Style{color.FgYellow}.Render
			case TestCaseTitlePen:
				pen = color.Style{color.FgLightYellow}.Render
			case SkippedPen:
				pen = color.Style{color.FgYellow}.Render
			case SuccessPen:
				pen = color.Style{color.FgGreen}.Render
			case FailurePen:
				pen = color.Style{color.FgRed}.Render
			case CrackedPen:
				pen = color.Style{color.FgRed}.Render
			}
			Pens[name] = pen
		}
	}
	return pen
}

type Renderer func(a ...interface{}) string

var ColorlessPen = func(a ...interface{}) string {
	text := ""
	if a == nil || len(a) == 0 {
		return text
	}
	for _, s := range a {
		text = text + fmt.Sprintf("%v", s)
	}
	return text
}

type PenType int

const (
	_ PenType = iota
	ContextNamePen
	SectionTitlePen
	SectionBodyPen
	HeadingPen
	TestSuiteTitlePen
	TestCaseTitlePen
	SuccessPen
	FailurePen
	CrackedPen
	SkippedPen
)

var Pens map[PenType]Renderer
