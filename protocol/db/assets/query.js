(function(g) {
  'use strict'

  var d = g.document;
  var query = d.querySelector('#query');
  var html = d.querySelector('#submit');
  var raw = d.querySelector('#submit_raw');

  function encodeQuery(s) {
    if (s.indexOf('%') !== -1) {
      s = s.replace(/%/g, '%25');
    }
    s = g.encodeURIComponent(s);
    return s
  }

  function queryHTML() {
    var s = encodeQuery(query.value);
    g.location.href += s + "?htmltable";
  }
  function queryLTSV() {
    var s = encodeQuery(query.value);
    g.location.href += s;
  }

  html.addEventListener('click', function(ev) {
    ev.preventDefault();
    queryHTML();
  });
  raw.addEventListener('click', function(ev) {
    ev.preventDefault();
    queryLTSV();
  });

  query.addEventListener('keydown', function(ev) {
    if (ev.ctrlKey && ev.keyCode == 13) {
      ev.preventDefault();
      queryHTML();
      return false;
    }
    if (ev.altKey && ev.keyCode == 13) {
      ev.preventDefault();
      queryLTSV();
      return false;
    }
  });

  var q = g.sessionStorage.getItem('query');
  if (q) {
    query.value = q;
    g.sessionStorage.removeItem('query');
  }
})(this);
