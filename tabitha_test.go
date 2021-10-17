package tabitha

import (
	"bytes"
	"testing"
)

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
// 			if got := New(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("New() = %v, want %v", got, tt.want)
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

func TestWriter_WriteTo(t *testing.T) {
	type fields struct {
		header []string
		rows   [][]string
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		wantN      int64
		wantErr    bool
	}{
		{
			name: "Basic output with headers",
			fields: fields{
				header: []string{"Status", "Name", "Details"},
				rows: [][]string {
					{ `new`, "tabitha", "This project" },
					{ `released`, "docked", "Another project"},
				},
			},
			wantWriter: "  Status\t   Name\t        Details\t\n     new\ttabitha\t   This project\t\nreleased\t docked\tAnother project\t\n",
			wantN: int64(102),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
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
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("WriteTo() gotWriter =\n %v\nwant =\n %v\n", gotWriter, tt.wantWriter)
			}
			if gotN != tt.wantN {
				t.Errorf("WriteTo() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
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
