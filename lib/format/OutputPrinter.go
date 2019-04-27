package format

import (
	"fmt"
	"strings"
	"github.com/gookit/color"
	"github.com/opwire/opwire-testa/lib/utils"
)

type OutputPrinterOptions struct {}

type OutputPrinter struct {
}

func NewOutputPrinter(opts *OutputPrinterOptions) (ref *OutputPrinter, err error) {
	ref = new(OutputPrinter)
	return ref, err
}

func (w *OutputPrinter) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Printf(format, args...)
}

func (w *OutputPrinter) Println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

func (w *OutputPrinter) Heading(title string) string {
	if HeadingPen == nil {
		HeadingPen = color.Style{color.FgCyan, color.OpUnderscore, color.OpBold}.Render
	}
	pen := HeadingPen
	return pen(title)
}

func (w *OutputPrinter) ContextInfo(name string, values ...string) string {
	if ContextNamePen == nil {
		ContextNamePen = color.Style{color.FgCyan, color.OpBold}.Render
	}
	pen := ContextNamePen
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
	if TestSuiteTitlePen == nil {
		TestSuiteTitlePen = color.Style{color.FgYellow}.Render
	}
	pen := TestSuiteTitlePen
	return fmt.Sprintf("[#] %s", pen(filepath))
}

func (w *OutputPrinter) TestCase(title string) string {
	if TestCaseTitlePen == nil {
		TestCaseTitlePen = color.Style{color.FgLightYellow}.Render
	}
	pen := TestCaseTitlePen
	return fmt.Sprintf("[=] %s", pen(title))
}

func (w *OutputPrinter) Skipped(title string) string {
	if SkippedPen == nil {
		SkippedPen = color.Style{color.FgGreen}.Render
	}
	pen := SkippedPen
	return fmt.Sprintf("[%s] %s", pen("-"), title)
}

func (w *OutputPrinter) Success(title string) string {
	if SuccessPen == nil {
		SuccessPen = color.Style{color.FgGreen}.Render
	}
	pen := SuccessPen
	return fmt.Sprintf("[%s] %s", pen("v"), title)
}

func (w *OutputPrinter) Failure(title string) string {
	if FailurePen == nil {
		FailurePen = color.Style{color.FgRed}.Render
	}
	pen := FailurePen
	return fmt.Sprintf("[%s] %s", pen("x"), title)
}

func (w *OutputPrinter) Cracked(title string) string {
	if CrackedPen == nil {
		CrackedPen = color.Style{color.FgGreen}.Render
	}
	pen := CrackedPen
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

type Renderer func(a ...interface{}) string

var ContextNamePen Renderer
var SectionTitlePen, SectionBodyPen Renderer
var HeadingPen, TestSuiteTitlePen, TestCaseTitlePen Renderer
var SuccessPen, FailurePen, CrackedPen, SkippedPen Renderer
