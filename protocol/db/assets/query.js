(function(g) {
  'use strict'

  var d = g.document;
  var query = d.querySelector('#query');
  var html = d.querySelector('#submit');
  var raw = d.querySelector('#submit_raw');
  html.addEventListener('click', function(ev) {
    ev.preventDefault();
    var s = g.encodeURIComponent(query.value);
    g.location.href += s + "?htmltable";
  });
  raw.addEventListener('click', function(ev) {
    ev.preventDefault();
    var s = g.encodeURIComponent(query.value);
    g.location.href += s;
  });

  var q = g.sessionStorage.getItem('query');
  if (q) {
    query.value = q;
    g.sessionStorage.removeItem('query');
  }
})(this);
