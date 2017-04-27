package main

import (
	"bytes"
	"errors"
	"log"
	"os"

	"github.com/tealeg/xlsx"
)

type column struct {
	name string
}

func parseColumns(xs *xlsx.Sheet) ([]column, error) {
	if len(xs.Rows) < 2 {
		return nil, errors.New("less rows")
	}
	var cols []column
	for _, xc := range xs.Rows[0].Cells {
		cols = append(cols, column{name: xc.Value})
	}
	return cols, nil
}

func run(filename string) error {
	xf, err := xlsx.OpenFile(filename)
	if err != nil {
		return err
	}
	for _, xs := range xf.Sheets {
		cols, err := parseColumns(xs)
		if err != nil {
			return err
		}
		bb := new(bytes.Buffer)
		for _, col := range cols {
			if bb.Len() > 0 {
				bb.WriteString("\n  ")
			}
			bb.WriteString(col.name)
		}
		//log.Printf("PREPARED: INSERT INTO %s (%s)", xs.Name, bb.String())
		for i, xr := range xs.Rows[1:] {
			for j, xc := range xr.Cells {
				if xc.Value == "" {
					log.Printf("%2d,%2d: (SKIPPED) %s", i, j, cols[j].name)
					continue
				}
				log.Printf("%2d,%2d: %s = %s", i, j, cols[j].name, xc.Value)
			}
		}
	}
	return nil
}

func main() {
	err := run(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}
