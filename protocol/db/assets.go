package db

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assetsf0a1e684333a9c2a21538b70ce7e4d8779474dac = "<!DOCTYPE! html>\n<meta charset=\"UTF-8\">\n<link rel=\"stylesheet\" type=\"text/css\" href=\"./assets/query.css\">\n<h1>Query</h1>\n<form class=\"container\">\n  <input id=\"query\" type=\"text\" name=\"query\"></input>\n  <button id=\"submit\">Query</button>\n  <button id=\"submit_raw\">Query (LTSV)</button>\n  <input type=\"reset\" value=\"Reset\"></input>\n</form>\n<script src=\"./assets/query.js\"></script>\n"
var _Assets6c0f5b9d340720d38a128145a10e4db1bbdb23f9 = ".container {\n  display: flex;\n}\n\n.container #query {\n  flex: 1 1 auto;\n}\n\n.container > *:not(:first-child) {\n  margin-left: 0.5em;\n}\n"
var _Assetsd62688e30fe02ff0df109515268604d16d69f82e = "(function(g) {\n  'use strict'\n\n  var d = g.document;\n  var query = d.querySelector('#query');\n  var html = d.querySelector('#submit');\n  var raw = d.querySelector('#submit_raw');\n  html.addEventListener('click', function(ev) {\n    ev.preventDefault();\n    var s = g.encodeURIComponent(query.value);\n    g.location.href += s + \"?htmltable\";\n  });\n  raw.addEventListener('click', function(ev) {\n    ev.preventDefault();\n    var s = g.encodeURIComponent(query.value);\n    g.location.href += s;\n  });\n})(this);\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{}, map[string]*assets.File{
	"index.html": &assets.File{
		Path:     "index.html",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1490162563, 1490162563096477300),
		Data:     []byte(_Assetsf0a1e684333a9c2a21538b70ce7e4d8779474dac),
	}, "query.css": &assets.File{
		Path:     "query.css",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1490162910, 1490162910846002900),
		Data:     []byte(_Assets6c0f5b9d340720d38a128145a10e4db1bbdb23f9),
	}, "query.js": &assets.File{
		Path:     "query.js",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1490162664, 1490162664129530000),
		Data:     []byte(_Assetsd62688e30fe02ff0df109515268604d16d69f82e),
	}}, "")
