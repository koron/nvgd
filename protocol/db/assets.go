package db

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assetsf0a1e684333a9c2a21538b70ce7e4d8779474dac = "<!DOCTYPE! html>\n<meta charset=\"UTF-8\">\n<link rel=\"stylesheet\" type=\"text/css\" href=\"./assets/query.css\">\n<h1>Query</h1>\n<form class=\"container\">\n  <textarea id=\"query\" rows=\"12\"></textarea>\n  <div class=\"vpanel\">\n    <button id=\"submit\">Query</button>\n    <button id=\"submit_raw\">Query (LTSV)</button>\n    <input type=\"reset\" value=\"Reset\"></input>\n  </div>\n</form>\n<script src=\"./assets/query.js\"></script>\n"
var _Assets6c0f5b9d340720d38a128145a10e4db1bbdb23f9 = ".container {\n  display: flex;\n}\n\n.container #query {\n  flex: 1 1 auto;\n}\n\n.container > *:not(:first-child) {\n  margin-left: 0.5em;\n}\n\n.vpanel {\n  display: flex;\n  flex-direction: column;\n}\n\n.vpanel > * {\n  padding: 1ex 1em;\n}\n.vpanel > *:not(:first-child) {\n  margin-top: 1ex;\n}\n"
var _Assetsd62688e30fe02ff0df109515268604d16d69f82e = "(function(g) {\n  'use strict'\n\n  var d = g.document;\n  var query = d.querySelector('#query');\n  var html = d.querySelector('#submit');\n  var raw = d.querySelector('#submit_raw');\n  html.addEventListener('click', function(ev) {\n    ev.preventDefault();\n    var s = g.encodeURIComponent(query.value);\n    g.location.href += s + \"?htmltable\";\n  });\n  raw.addEventListener('click', function(ev) {\n    ev.preventDefault();\n    var s = g.encodeURIComponent(query.value);\n    g.location.href += s;\n  });\n})(this);\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{}, map[string]*assets.File{
	"query.css": &assets.File{
		Path:     "query.css",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1491359751, 1491359751909445200),
		Data:     []byte(_Assets6c0f5b9d340720d38a128145a10e4db1bbdb23f9),
	}, "query.js": &assets.File{
		Path:     "query.js",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1490755363, 1490755363383749100),
		Data:     []byte(_Assetsd62688e30fe02ff0df109515268604d16d69f82e),
	}, "index.html": &assets.File{
		Path:     "index.html",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1491359996, 1491359996372379300),
		Data:     []byte(_Assetsf0a1e684333a9c2a21538b70ce7e4d8779474dac),
	}}, "")
