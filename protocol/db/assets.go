package db

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assetscd9af3cf17d5a0546ede38b9af7947b522cbd711 = "// restore.js dummy\n(function(g) {\n  'use strict'\n})(this);\n"
var _Assetse63c31a38a562037297091840b1433d3b16cacb4 = "/* restore.css dummy */\n"
var _Assets008ca05abb21a6c52ec7712e1318522991975c13 = "/* update.css dummy */\n"
var _Assetsd62688e30fe02ff0df109515268604d16d69f82e = "(function(g) {\n  'use strict'\n\n  var d = g.document;\n  var query = d.querySelector('#query');\n  var html = d.querySelector('#submit');\n  var raw = d.querySelector('#submit_raw');\n  html.addEventListener('click', function(ev) {\n    ev.preventDefault();\n    var s = g.encodeURIComponent(query.value);\n    g.location.href += s + \"?htmltable\";\n  });\n  raw.addEventListener('click', function(ev) {\n    ev.preventDefault();\n    var s = g.encodeURIComponent(query.value);\n    g.location.href += s;\n  });\n})(this);\n"
var _Assetsb0590a31f1a55ffe1e02dc8e22bf87051f1402c4 = "<!DOCTYPE! html>\n<meta charset=\"UTF-8\">\n<link rel=\"stylesheet\" type=\"text/css\" href=\"./restore.css\">\n<h1>Restore {{.name}}</h1>\n<form class=\"container\" method=\"post\" enctype=\"multipart/form-data\">\n  <input type=\"file\" name=\"file00\"></input>\n  <input type=\"submit\" value=\"Restore\"></input>\n  <input type=\"reset\" value=\"Reset\"></input>\n</form>\n<script src=\"./restore.js\"></script>\n"
var _Assets8f58f95489d6e7a2ba3c4d4fb40efd576bb92995 = "<!DOCTYPE! html>\n<meta charset=\"UTF-8\">\n<link rel=\"stylesheet\" type=\"text/css\" href=\"./update.css\">\n<h1>Update {{.name}}</h1>\n<form class=\"container\" method=\"post\" enctype=\"multipart/form-data\">\n  <input type=\"file\" name=\"file00\"></input>\n  <input type=\"submit\" value=\"Restore\"></input>\n  <input type=\"reset\" value=\"Reset\"></input>\n</form>\n<script src=\"./update.js\"></script>\n"
var _Assets3e678865eb64e5ce333408c3cf15af47d15d18fe = "// update.js dummy\n(function(g) {\n  'use strict'\n})(this);\n"
var _Assets7b28d06da7d2ff49b0679c2f4a7e2f4b81ab28d9 = "/* dump.css dummy */\n"
var _Assets19787aaa6d9704d88a58fc595d8fe37c352949ac = "(function(g) {\n  'use strict'\n})(this);\n"
var _Assets7c8c0da38f3dc61f5f9279ef483848b78a4633e7 = "<!DOCTYPE! html>\n<meta charset=\"UTF-8\">\n<link rel=\"stylesheet\" type=\"text/css\" href=\"./assets/dump.css\">\n<h1>Dump</h1>\n<script src=\"./assets/dump.js\"></script>\n"
var _Assetsf0a1e684333a9c2a21538b70ce7e4d8779474dac = "<!DOCTYPE! html>\n<meta charset=\"UTF-8\">\n<link rel=\"stylesheet\" type=\"text/css\" href=\"./assets/query.css\">\n<h1>Query</h1>\n<form class=\"container\">\n  <textarea id=\"query\" rows=\"12\"></textarea>\n  <div class=\"vpanel\">\n    <button id=\"submit\">Query</button>\n    <button id=\"submit_raw\">Query (LTSV)</button>\n    <input type=\"reset\" value=\"Reset\"></input>\n  </div>\n</form>\n<script src=\"./assets/query.js\"></script>\n"
var _Assets6c0f5b9d340720d38a128145a10e4db1bbdb23f9 = ".container {\n  display: flex;\n}\n\n.container #query {\n  flex: 1 1 auto;\n}\n\n.container > *:not(:first-child) {\n  margin-left: 0.5em;\n}\n\n.vpanel {\n  display: flex;\n  flex-direction: column;\n}\n\n.vpanel > * {\n  padding: 1ex 1em;\n}\n.vpanel > *:not(:first-child) {\n  margin-top: 1ex;\n}\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{}, map[string]*assets.File{
	"update.css": {
		Path:     "update.css",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1498013769, 1498013769163692500),
		Data:     []byte(_Assets008ca05abb21a6c52ec7712e1318522991975c13),
	}, "restore.css": {
		Path:     "restore.css",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1498013067, 1498013067896368300),
		Data:     []byte(_Assetse63c31a38a562037297091840b1433d3b16cacb4),
	}, "dump.js": {
		Path:     "dump.js",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1498011449, 1498011449203356200),
		Data:     []byte(_Assets19787aaa6d9704d88a58fc595d8fe37c352949ac),
	}, "query.js": {
		Path:     "query.js",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1490755363, 1490755363383749100),
		Data:     []byte(_Assetsd62688e30fe02ff0df109515268604d16d69f82e),
	}, "restore.html": {
		Path:     "restore.html",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1498014342, 1498014342775390800),
		Data:     []byte(_Assetsb0590a31f1a55ffe1e02dc8e22bf87051f1402c4),
	}, "update.html": {
		Path:     "update.html",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1498014342, 1498014342791930600),
		Data:     []byte(_Assets8f58f95489d6e7a2ba3c4d4fb40efd576bb92995),
	}, "update.js": {
		Path:     "update.js",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1498013773, 1498013773201457500),
		Data:     []byte(_Assets3e678865eb64e5ce333408c3cf15af47d15d18fe),
	}, "dump.css": {
		Path:     "dump.css",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1498011466, 1498011466284450200),
		Data:     []byte(_Assets7b28d06da7d2ff49b0679c2f4a7e2f4b81ab28d9),
	}, "dump.html": {
		Path:     "dump.html",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1498011357, 1498011357381579400),
		Data:     []byte(_Assets7c8c0da38f3dc61f5f9279ef483848b78a4633e7),
	}, "index.html": {
		Path:     "index.html",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1498015180, 1498015180303297300),
		Data:     []byte(_Assetsf0a1e684333a9c2a21538b70ce7e4d8779474dac),
	}, "query.css": {
		Path:     "query.css",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1498015180, 1498015180308802200),
		Data:     []byte(_Assets6c0f5b9d340720d38a128145a10e4db1bbdb23f9),
	}, "restore.js": {
		Path:     "restore.js",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1498013089, 1498013089195290900),
		Data:     []byte(_Assetscd9af3cf17d5a0546ede38b9af7947b522cbd711),
	}}, "")
