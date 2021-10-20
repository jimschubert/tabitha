package tabitha

import (
	"bytes"
	"os"
	"testing"
)

func TestWriter_WriteTo(t *testing.T) {
	type fields struct {
		header []string
		rows   [][]string
	}
	tests := []struct {
		name    string
		fields  fields
		inits   func(w *Writer)
		result  string
		wantN   int64
		wantErr bool
	}{
		{
			name: "Basic output with headers",
			fields: fields{
				header: []string{"Status", "Name", "Details"},
				rows: [][]string{
					{`new`, "tabitha", "This project"},
					{`released`, "docked", "Another project"},
				},
			},
			result:  "  Status\t   Name\t        Details\n     new\ttabitha\t   This project\nreleased\t docked\tAnother project\n",
			wantN:   int64(99),
			wantErr: false,
		},
		{
			name: "Basic output without headers",
			fields: fields{
				rows: [][]string{
					{`new`, "tabitha", "This project"},
					{`released`, "docked", "Another project"},
				},
			},
			result:  "     new\ttabitha\t   This project\nreleased\t docked\tAnother project\n",
			wantN:   int64(66),
			wantErr: false,
		},
		{
			name: "Basic output with headers custom padding character",
			fields: fields{
				header: []string{"Status", "Name", "Details"},
				rows: [][]string{
					{`new`, "tabitha", "This project"},
					{`released`, "docked", "Another project"},
				},
			},
			inits: func(w *Writer) {
				w.padding = true
				w.padChar = '.'
				w.separator = '\t'
			},
			result:  "..Status\t...Name\t........Details\n.....new\ttabitha\t...This project\nreleased\t.docked\tAnother project\n",
			wantN:   int64(99),
			wantErr: false,
		},
		{
			name: "Basic output with headers no padding",
			fields: fields{
				header: []string{"Status", "Name", "Details"},
				rows: [][]string{
					{`new`, "tabitha", "This project"},
					{`released`, "docked", "Another project"},
				},
			},
			inits: func(w *Writer) {
				w.padding = false
				w.separator = '\t'
			},
			result:  "Status\tName\tDetails\nnew\ttabitha\tThis project\nreleased\tdocked\tAnother project\n",
			wantN:   int64(77),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWriter()
			if tt.inits != nil {
				tt.inits(w)
			}
			err := w.Header(tt.fields.header...)
			if err != nil {
				t.Errorf("WriteTo() errored calling Header() = %v", err)
				return
			}

			for _, row := range tt.fields.rows {
				if err := w.AddLine(row...); err != nil {
					t.Errorf("WriteTo() errored calling AddLine() = %v", err)
					return
				}
			}

			writer := &bytes.Buffer{}
			gotN, err := w.WriteTo(writer)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.result {
				t.Errorf("WriteTo() gotWriter =\n %v\nwant =\n %v\n", gotWriter, tt.result)
			}
			if gotN != tt.wantN {
				t.Errorf("WriteTo() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func ExampleWriter_WriteTo() {
	tt := NewWriter()
	_ = tt.Header("a", "bb", "cc", "dd")
	_ = tt.SpacerLine()
	_ = tt.AddLine("Line 1", "Under bb", "Third", "4th")
	_ = tt.AddLine("cat", "dog", "bird", "frog")
	_ = tt.SpacerLine()
	_, _ = tt.WriteTo(os.Stdout)

	// Output:
	//     a	      bb	   cc	  dd
	// ------	--------	-----	----
	// Line 1	Under bb	Third	 4th
	//    cat	     dog	 bird	frog
	// ------	--------	-----	----
}

func ExampleWriter_WriteTo_withInitialization() {
	tt := NewWriter()
	tt.WithPadding(true)
	tt.LinesEndWith('|')
	tt.LinesStartWith('|')
	tt.PaddingCharacter(' ')
	tt.CellSeparator(' ')
	tt.IgnoreAnsiWidths(true)

	// Table example from https://www.markdownguide.org/extended-syntax/#alignment
	_ = tt.Header("Syntax", "Description", "Test Text")
	_ = tt.AddLine(":---", ":----:", "---:")
	_ = tt.AddLine("Header", "Title", "Here's this")
	_ = tt.AddLine("Paragraph", "Text", "And more")
	_, _ = tt.WriteTo(os.Stdout)

	// Output:
	// |   Syntax |Description |  Test Text|
	// |     :--- |     :----: |       ---:|
	// |   Header |      Title |Here's this|
	// |Paragraph |       Text |   And more|
}

//
// func TestNew(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want *Writer
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NewWriter(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewWriter() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestWriter_AddLine(t *testing.T) {
// 	type fields struct {
// 		header []string
// 		rows   []cells
// 		widths []int
// 	}
// 	type args struct {
// 		input []string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := Writer{
// 				header: tt.fields.header,
// 				rows:   tt.fields.rows,
// 				widths: tt.fields.widths,
// 			}
// 			if err := w.AddLine(tt.args.input...); (err != nil) != tt.wantErr {
// 				t.Errorf("AddLine() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
//
// func TestWriter_Header(t *testing.T) {
// 	type fields struct {
// 		header []string
// 		rows   []cells
// 		widths []int
// 	}
// 	type args struct {
// 		input []string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := Writer{
// 				header: tt.fields.header,
// 				rows:   tt.fields.rows,
// 				widths: tt.fields.widths,
// 			}
// 			if err := w.Header(tt.args.input...); (err != nil) != tt.wantErr {
// 				t.Errorf("Header() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
//

//
// func TestWriter_handlePanic(t *testing.T) {
// 	type fields struct {
// 		header []string
// 		rows   []cells
// 		widths []int
// 	}
// 	type args struct {
// 		err   *error
// 		where string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := Writer{
// 				header: tt.fields.header,
// 				rows:   tt.fields.rows,
// 				widths: tt.fields.widths,
// 			}
// 		})
// 	}
// }
//
// func TestWriter_initWidths(t *testing.T) {
// 	type fields struct {
// 		header []string
// 		rows   []cells
// 		widths []int
// 	}
// 	type args struct {
// 		input []string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := Writer{
// 				header: tt.fields.header,
// 				rows:   tt.fields.rows,
// 				widths: tt.fields.widths,
// 			}
// 		})
// 	}
// }
//
// func TestWriter_reset(t *testing.T) {
// 	type fields struct {
// 		header []string
// 		rows   []cells
// 		widths []int
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := Writer{
// 				header: tt.fields.header,
// 				rows:   tt.fields.rows,
// 				widths: tt.fields.widths,
// 			}
// 		})
// 	}
// }
//
// func TestWriter_updateWidths(t *testing.T) {
// 	type fields struct {
// 		header []string
// 		rows   []cells
// 		widths []int
// 	}
// 	type args struct {
// 		input []string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := Writer{
// 				header: tt.fields.header,
// 				rows:   tt.fields.rows,
// 				widths: tt.fields.widths,
// 			}
// 		})
// 	}
// }
