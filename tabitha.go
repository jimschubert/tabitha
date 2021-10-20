package tabitha

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"
)

// writerError mirrors behavior of osError in text/tabwriter, essentially allowing a "thrown" error using panics
type writerError struct {
	err error
}

type cell struct {
	data   string
	format string
}

type cells []cell

// A Writer provides a tabbed output with padding. This differs from text/tabwriter in the standard library in that
// Writer has fewer formatting options and honors rune-width of user inputs, whereas text/tabwriter assumes that
// all characters have the same width.
type Writer struct {
	lineStart        *rune
	lineEnd          *rune
	separator        rune
	padding          bool
	padChar          rune
	ignoreAnsiWidths bool
	header           []string
	rows             []cells
	widths           []int
}

// LinesStartWith allows the user to define a starting rune for the line.
func (w *Writer) LinesStartWith(start rune) *Writer {
	w.lineStart = &start
	return w
}

// LinesEndWith allows the user to modify the line ending character. This is the character which closes a cell, and comes before the newline.
func (w *Writer) LinesEndWith(end rune) *Writer {
	w.lineEnd = &end
	return w
}

// CellSeparator allows the user to modify the cell separator character. Default is newline ('\t').
func (w *Writer) CellSeparator(separator rune) *Writer {
	w.separator = separator
	return w
}

// PaddingCharacter allows the user to modify a character (specified as rune) to pad a cell when padding is enabled.
// Default is space (' ')
func (w *Writer) PaddingCharacter(padChar rune) *Writer {
	w.padChar = padChar
	return w
}

// WithPadding allows the user to specify whether padding to the width of displayable text is enabled.
// When IgnoreAnsiWidths is passed true, ANSI codes are extracted from cell text to return "displayable text".
// When IgnoreAnsiWidths is passed false, all text including ANSI codes are considered in cell width calculations.
func (w *Writer) WithPadding(padding bool) *Writer {
	w.padding = padding
	return w
}

// IgnoreAnsiWidths allows the user configure how displayable text is calculated. Pass true to drop ANSI escape codes from
// width calculations, or false to include all ANSI escape codes in those calculations.
func (w *Writer) IgnoreAnsiWidths(ignoreAnsiWidths bool) *Writer {
	w.ignoreAnsiWidths = ignoreAnsiWidths
	return w
}

func (w *Writer) reset() {
	w.header = make([]string, 0)
	w.rows = make([]cells, 0)
	w.widths = make([]int, 0)
}

// handlePanic logic is taken from standard library tabwriter
func (w *Writer) handlePanic(err *error, where string) {
	if e := recover(); e != nil {
		if where == "WriteTo" {
			w.reset()
		}
		// allows for known error to be raised
		if nerr, ok := e.(writerError); ok {
			*err = nerr.err
			return
		}
		panic("tabitha: panic during " + where)
	}
}

func (w *Writer) calculateWidth(input string) int {
	var length int
	if w.ignoreAnsiWidths {
		pattern := `(\x1b\[[0-9;]+[a-zA-Z~])`
		re := regexp.MustCompilePOSIX(pattern)
		length = utf8.RuneCountInString(re.ReplaceAllString(input, ""))
	} else {
		length = utf8.RuneCountInString(input)
	}

	if length >= 0 {
		return length
	}

	return 0
}

func (w *Writer) initWidths(input ...string) {
	if len(w.widths) > 0 {
		panic("attempted to re-initialize widths at an unexpected location")
	}
	for _, i := range input {
		w.widths = append(w.widths, w.calculateWidth(i))
	}
}

func (w *Writer) updateWidths(input ...string) {
	if len(w.widths) != len(input) {
		panic(writerError{err: errors.New("invalid column count")})
	}

	for i, current := range input {
		knownWidth := w.widths[i]
		currentWidth := w.calculateWidth(current)
		if currentWidth > knownWidth {
			w.widths[i] = currentWidth
		}
	}
}

// Header defines the text for the table header. This is a semantic helper, differing from AddLine in that calling Header
// more than once will result in a panic.
func (w *Writer) Header(input ...string) (err error) {
	defer w.handlePanic(&err, "Header")
	w.header = append(w.header, input...)
	w.initWidths(input...)
	return
}

// SpacerLine registers a line of '-' characters by default. The width of the line is calculated when WriteTo is called.
// The width is calculated against the collected header and lines.
func (w *Writer) SpacerLine() (err error) {
	defer w.handlePanic(&err, "SpacerLine")
	// spacer doesn't have a width… it's formatted later
	// however, if there are no known widths (i.e. no header or lines), that's a programmatic error
	if len(w.widths) == 0 {
		panic(writerError{err: errors.New("spacer cannot be called before Header or AddLine")})
	}

	row := make([]cell, 0)
	for _, _ = range w.widths {
		c := cell{data: "-", format: "%*s"}
		row = append(row, c)
	}
	w.rows = append(w.rows, row)
	return
}

// AddLine collects the input strings for deferred evaluation of overall table width.
func (w *Writer) AddLine(input ...string) (err error) {
	defer w.handlePanic(&err, "AddLine")
	if len(w.widths) == 0 {
		w.initWidths(input...)
	} else {
		w.updateWidths(input...)
	}

	row := make([]cell, 0)
	for _, i := range input {
		c := cell{data: i}
		row = append(row, c)
	}
	w.rows = append(w.rows, row)
	return
}

func (w *Writer) write(writer io.Writer, cell cell, index int) (int64, error) {
	width := w.widths[index]
	buf := bytes.Buffer{}
	isLast := index == len(w.widths)-1

	value := cell.data
	if cell.format != "" {
		if !strings.Contains(cell.format, "*") {
			panic(writerError{err: errors.New("invalid spacer format (requires * width pattern)")})
		}
		formatted := fmt.Sprintf(cell.format, width, value)
		value = strings.Replace(formatted, " ", cell.data, -1)
	}

	if w.lineStart != nil {
		buf.WriteRune(*w.lineStart)
	}

	if w.padding {
		v := w.calculateWidth(value)
		if v > 0 {
			for i := 0; i < (width - v); i++ {
				buf.WriteRune(w.padChar)
			}
		}
		buf.WriteString(value)
	} else {
		buf.WriteString(value)
	}

	if !isLast {
		buf.WriteRune(w.separator)
	} else {
		if w.lineEnd != nil {
			buf.WriteRune(*w.lineEnd)
		}
		buf.WriteRune('\n')
	}
	n, err := fmt.Fprint(writer, buf.String())
	return int64(n), err
}

// WriteTo an io.Writer for all collected contents in tabitha.Writer. Returns runes written and an error (if populated).
func (w *Writer) WriteTo(writer io.Writer) (n int64, err error) {
	defer w.handlePanic(&err, "Write")
	var nn int64
	for i, h := range w.header {
		nn, err = w.write(writer, cell{data: h}, i)
		n += nn
		if err != nil {
			panic(writerError{err: err})
		}
	}

	for _, row := range w.rows {
		for i, c := range row {
			nn, err = w.write(writer, c, i)
			n += nn
			if err != nil {
				panic(writerError{err: err})
			}
		}
	}
	return
}

// NewWriter returns a new Writer with default configuration options
func NewWriter() *Writer {
	writer := Writer{
		separator:        '\t',
		padding:          true,
		padChar:          ' ',
		ignoreAnsiWidths: false,
	}
	writer.reset()
	return &writer
}
