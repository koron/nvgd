package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"time"

	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
	"github.com/tealeg/xlsx"
)

type DumpHandler struct {
}

func init() {
	protocol.MustRegister("db-dump", &DumpHandler{})
}

func (dh *DumpHandler) Open(u *url.URL) (*resource.Resource, error) {
	c, err := openDB(u)
	if err != nil {
		return nil, err
	}
	xf := xlsx.NewFile()
	table := path(u)
	rows, err := c.db.Query("SELECT * FROM " + table)
	if err != nil {
		return nil, err
	}
	_, err = addSheet(xf, table, rows)
	if err != nil {
		return nil, err
	}
	buf, err := saveToBuffer(xf)
	if err != nil {
		return nil, err
	}
	rs := resource.New(ioutil.NopCloser(buf))
	n, _ := extractNames(u)
	t := time.Now().Format("20060102T150405MST")
	rs.PutFilename(fmt.Sprintf("%s-%s.xlsx", n, t))
	rs.Put(protocol.Small, true)
	return rs, nil
}

func saveToBuffer(xf *xlsx.File) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := xf.Write(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func addSheet(xf *xlsx.File, name string, rows *sql.Rows) (*xlsx.Sheet, error) {
	xs, err := xf.AddSheet(name)
	if err != nil {
		return nil, err
	}
	ct, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	// convert column types and add as the header to xlsx.  and setup receiver.
	h1 := xs.AddRow()
	vals := make([]interface{}, len(ct))
	for i, t := range ct {
		h1.AddCell().SetString(t.Name())
		vals[i] = reflect.New(t.ScanType()).Interface()
	}
	// convert values to xlsx'x cells
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			return nil, err
		}
		drows := xs.AddRow()
		for _, v := range vals {
			c := drows.AddCell()
			c.SetValue(*(v.(*interface{})))
		}
	}
	return xs, nil
}
