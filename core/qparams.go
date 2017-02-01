package core

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

type qparams []*qparamItem

func (qp qparams) split(n string) (matched qparams, others qparams) {
	for _, item := range qp {
		if item.name == n {
			matched = append(matched, item)
			continue
		}
		others = append(others, item)
	}
	return matched, others
}

func (qp qparams) String() string {
	b := &bytes.Buffer{}
	b.WriteString("[")
	for i, p := range qp {
		if i != 0 {
			b.WriteString(" ")
		}
		fmt.Fprintf(b, "%+v", *p)
	}
	b.WriteString("]")
	return b.String()
}

func qparamsParse(qs string) (qparams, error) {
	var qp qparams
	for qs != "" {
		k := qs
		if i := strings.Index(k, "&"); i >= 0 {
			k, qs = k[:i], k[i+1:]
		} else {
			qs = ""
		}
		if k == "" {
			continue
		}
		v := ""
		if i := strings.Index(k, "="); i >= 0 {
			k, v = k[:i], k[i+1:]
		}
		k, err := url.QueryUnescape(k)
		if err != nil {
			return nil, err
		}
		v, err = url.QueryUnescape(v)
		if err != nil {
			return nil, err
		}
		qp = append(qp, &qparamItem{name: k, value: v})
	}
	return qp, nil
}

type qparamItem struct {
	name  string
	value string
}
