package trdsql

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewCSVWriter(t *testing.T) {
	outStream := new(bytes.Buffer)
	type args struct {
		writeOpts *WriteOpts
	}
	tests := []struct {
		name string
		args args
		want rune
	}{
		{
			name: "testDefault",
			args: args{
				writeOpts: &WriteOpts{
					OutDelimiter: ",",
					OutStream:    outStream,
				},
			},
			want: ',',
		},
		{
			name: "invalidDelimiter",
			args: args{
				writeOpts: &WriteOpts{
					OutDelimiter: "--",
					OutStream:    outStream,
				},
			},
			want: ',',
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCSVWriter(tt.args.writeOpts)
			if !reflect.DeepEqual(got.outDelimiter, tt.want) {
				t.Errorf("NewCSVWriter() = %v, want %v", got.outDelimiter, tt.want)
			}
		})
	}
}

func TestCSVWriter_PreWrite(t *testing.T) {
	type args struct {
		columns []string
		types   []string
	}
	tests := []struct {
		name      string
		writeOpts WriteOpts
		args      args
		want      []string
		wantErr   bool
	}{
		{
			name: "empty",
			writeOpts: WriteOpts{
				OutDelimiter: ",",
			},
			args: args{
				columns: []string{},
				types:   []string{},
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "emptyHeader",
			writeOpts: WriteOpts{
				OutDelimiter: ",",
				OutHeader:    true,
			},
			args: args{
				columns: []string{"h1", "h2"},
				types:   []string{"text", "text"},
			},
			want:    make([]string, 2),
			wantErr: false,
		},
		{
			name: "noHeader",
			writeOpts: WriteOpts{
				OutDelimiter: ",",
				OutHeader:    false,
			},
			args: args{
				columns: []string{"v1", "v2"},
				types:   []string{"text", "text"},
			},
			want:    make([]string, 2),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCSVWriter(&tt.writeOpts)
			if err := w.PreWrite(tt.args.columns, tt.args.types); (err != nil) != tt.wantErr {
				t.Errorf("CSVWriter.PreWrite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
