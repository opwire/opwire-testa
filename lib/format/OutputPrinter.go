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

func (w *OutputPrinter) ContextInfo(name string, message string, values ...string) string {
	pen := w.GetPen(ContextNamePen)
	text := "[+] " + pen(name)
	if len(message) > 0 {
		text = fmt.Sprintf("%s: %s", text, message)
	}
	number := len(values)
	if number > 0 {
		if len(message) == 0 {
			text = text + ":"
		}
		lines := []string{text}
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

func (w *OutputPrinter) PositiveTag(tag string) string {
	pen := w.GetPen(PositiveTagPen)
	return pen(tag)
}

func (w *OutputPrinter) NegativeTag(tag string) string {
	pen := w.GetPen(NegativeTagPen)
	return pen(tag)
}

func (w *OutputPrinter) RegularTag(tag string) string {
	pen := w.GetPen(RegularTagPen)
	return pen(tag)
}

func (w *OutputPrinter) Pending(title string) string {
	pen := w.GetPen(PendingPen)
	return fmt.Sprintf("[%s] %s", pen("|"), title)
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
		if number == 0 {
			return " - " + line
		}
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
			case PositiveTagPen:
				pen = color.Style{color.FgGreen}.Render
			case NegativeTagPen:
				pen = color.Style{color.FgRed}.Render
			case RegularTagPen:
				pen = color.Style{color.FgGray}.Render
			case PendingPen:
				pen = color.Style{color.FgYellow, color.OpBold}.Render
			case SkippedPen:
				pen = color.Style{color.FgYellow, color.OpBold}.Render
			case SuccessPen:
				pen = color.Style{color.FgGreen, color.OpBold}.Render
			case FailurePen:
				pen = color.Style{color.FgRed, color.OpBold}.Render
			case CrackedPen:
				pen = color.Style{color.FgRed, color.OpBold}.Render
			}
			Pens[name] = pen
		}
	}
	return pen
}

type Renderer func(a ...interface{}) string

var ColorlessPen = func(a ...interface{}) string {
	if a == nil || len(a) == 0 {
		return ""
	}
	items := make([]string, 0, len(a))
	for _, v := range a {
		s := fmt.Sprintf("%v", v)
		if len(s) > 0 {
			items = append(items, s)
		}
	}
	return strings.Join(items, " ")
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
	PositiveTagPen
	NegativeTagPen
	RegularTagPen
	PendingPen
	SkippedPen
	SuccessPen
	FailurePen
	CrackedPen
)

var Pens map[PenType]Renderer
