package tabitha

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

// writerError mirrors behavior of osError in text/tabwriter, essentially allowing a "thrown" error using panics
type writerError struct {
	err error
}

type cell struct {
	data string
}

type cells []cell

// A Writer provides a tabbed output with padding. This differs from text/tabwriter in the standard library in that
// Writer has fewer formatting options and honors rune-width of user inputs, whereas text/tabwriter assumes that
// all characters have the same width.
type Writer struct {
	header []string
	rows   []cells
	widths []int
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

func (w *Writer) initWidths(input ...string) {
	if len(w.widths) > 0 {
		panic("attempted to re-initialize widths at an unexpected location")
	}
	for _, i := range input {
		w.widths = append(w.widths, utf8.RuneCountInString(i))
	}
}

func (w *Writer) updateWidths(input ...string) {
	if len(w.widths) != len(input) {
		panic(writerError{err: errors.New("invalid column count")})
	}

	for i, current := range input {
		knownWidth := w.widths[i]
		currentWidth := utf8.RuneCountInString(current)
		if currentWidth > knownWidth {
			w.widths[i] = currentWidth
		}
	}
}

func (w *Writer) Header(input... string) (err error)  {
	defer w.handlePanic(&err, "Header")
	w.header = append(w.header, input...)
	w.initWidths(input...)
	return
}

func (w *Writer) AddLine(input ...string) (err error) {
	defer w.handlePanic(&err, "AddLine")
	if len(w.widths) == 0 {
		w.initWidths(input...)
	} else {
		w.updateWidths(input...)
	}

	row := make([]cell, 0)
	for _, i := range input {
		c := cell{ data: i }
		row = append(row, c)
	}
	w.rows = append(w.rows, row)
	return
}

func (w *Writer) WriteTo(writer io.Writer) (n int64, err error) {
	defer w.handlePanic(&err, "Write")
	widthFormat := make([]string, 0)
	for _, width := range w.widths {
		format := "%" + strconv.Itoa(width) + "s\t"
		widthFormat = append(widthFormat, format)
	}

	for i, h := range w.header {
		nn, nerr := fmt.Fprintf(writer, widthFormat[i], h)
		n += int64(nn)
		if nerr != nil {
			panic(writerError{err: nerr})
		}
	}
	if len(w.header) > 0 {
		nn, nerr := fmt.Fprint(writer, "\n")
		n += int64(nn)
		if nerr != nil {
			panic(writerError{err: nerr})
		}
	}

	for _, row := range w.rows {
		for i, c := range row {
			nn, nerr := fmt.Fprintf(writer, widthFormat[i], c.data)
			n += int64(nn)
			if nerr != nil {
				panic(writerError{err: nerr})
			}
		}
		nn, nerr := fmt.Fprint(writer, "\n")
		n += int64(nn)
		if nerr != nil {
			panic(writerError{err: nerr})
		}
	}
	return
}

func New() *Writer {
	writer := Writer{}
	writer.reset()
	return &writer
}
