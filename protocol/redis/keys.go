package redis

import (
	"bytes"
	"html/template"
	"io"
	"strings"

	"github.com/go-redis/redis/v7"
	"github.com/koron/nvgd/resource"
)

func hasKeysMeta(s string) bool {
	return strings.ContainsAny(s, "?*[")
}

func keys(c *redis.Client, args []string) (*resource.Resource, error) {
	q := strings.Join(args, "/")
	if !hasKeysMeta(q) {
		q += "*"
	}
	r, err := c.Keys(q).Result()
	if err != nil {
		return nil, err
	}
	return resource.NewString(strings.Join(r, "\n")), nil
}

var keysTmpl = template.Must(template.New("htmltable").Parse(`<!DOCTYPE html>
<meta charset="UTF-8">
<meta name="referrer" content="no-referrer">
<style>
#query {
  display: block;
  width: 100%;
  margin-bottom: 1em;
}
#result {
  display: block;
  width: 100%;
  margin-bottom: 1em;
}
</style>
Keys query:<br>
<input type="text" id="query">
Results:<br>
<textarea id="result" rows="10">
</textarea>
<script>
(function(g) {
  'use strict'
  var d = g.document;
  var query = d.querySelector('#query');
  var result = d.querySelector('#result');
  var last;
  query.addEventListener('input', function(ev) {
	ev.preventDefault();
	if (query.value == last) {
	  return;
	}
	last = query.value;
	fetch('./keys/' + encodeURIComponent(query.value))
	  .then(function(resp) {
	    resp.text().then(function(text) {
		  result.value = text;
		});
	  });
  });
})(this);
</script>
`))

func keysForm(c *redis.Client, args []string) (*resource.Resource, error) {
	buf := new(bytes.Buffer)
	err := keysTmpl.Execute(buf, nil)
	if err != nil {
		return nil, err
	}
	return resource.New(io.NopCloser(buf)), nil
}
