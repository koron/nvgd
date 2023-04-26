package trdsql

import (
	"strings"

	"github.com/noborus/tbln"
)

// TBLNWriter provides methods of the Writer interface.
type TBLNWriter struct {
	writer   *tbln.Writer
	outNULL  string
	results  []string
	needNULL bool
}

// NewTBLNWriter returns TBLNWriter.
func NewTBLNWriter(writeOpts *WriteOpts) *TBLNWriter {
	w := &TBLNWriter{}
	w.writer = tbln.NewWriter(writeOpts.OutStream)
	w.needNULL = writeOpts.OutNeedNULL
	w.outNULL = writeOpts.OutNULL
	return w
}

// PreWrite is prepare tbln definition body.
func (w *TBLNWriter) PreWrite(columns []string, types []string) error {
	d := tbln.NewDefinition()

	if err := d.SetNames(columns); err != nil {
		return err
	}
	if err := d.SetTypes(ConvertTypes(types)); err != nil {
		return err
	}
	if err := w.writer.WriteDefinition(d); err != nil {
		return err
	}
	w.results = make([]string, len(columns))
	return nil
}

// WriteRow is row write.
func (w *TBLNWriter) WriteRow(values []interface{}, columns []string) error {
	for i, col := range values {
		str := ""
		if col == nil && w.needNULL {
			str = w.outNULL
		} else {
			str = ValString(col)
		}
		w.results[i] = strings.ReplaceAll(str, "\n", "\\n")
	}
	return w.writer.WriteRow(w.results)
}

// PostWrite is nil.
func (w *TBLNWriter) PostWrite() error {
	return nil
}

// ConvertTypes is converts database types to common types.
func ConvertTypes(dbTypes []string) []string {
	ret := make([]string, len(dbTypes))
	for i, t := range dbTypes {
		ret[i] = convertType(t)
	}
	return ret
}

func convertType(dbType string) string {
	switch strings.ToLower(dbType) {
	case "smallint", "integer", "int", "int2", "int4", "smallserial", "serial":
		return "int"
	case "bigint", "int8", "bigserial":
		return "bigint"
	case "float", "decimal", "numeric", "real", "double precision":
		return "numeric"
	case "bool":
		return "bool"
	case "timestamp", "timestamptz", "date", "time":
		return "timestamp"
	case "string", "text", "char", "varchar":
		return "text"
	default:
		return "text"
	}
}
