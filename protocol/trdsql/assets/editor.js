((g) => {
  "use strict"

  var d = g.document;
  var localStorage = g.localStorage;

  var queryForm = d.querySelector("#query");
  var submitButton = d.querySelector("#submit");
  var sourceForm = d.querySelector("#source_url");
  var composedForm = d.querySelector("#composed_url");
  var outputForm = d.querySelector("#output");
  var optIhCheckbox = d.querySelector("#opt_ih");
  var optOhCheckbox = d.querySelector("#opt_oh");
  var optIfmtSelect = d.querySelector("#opt_ifmt");
  var optOfmtSelect = d.querySelector("#opt_ofmt");
  var optInullCheckbox = d.querySelector("#opt_inull");
  var optInullForm = d.querySelector("#opt_inull_text");
  var optOnullCheckbox = d.querySelector("#opt_onull");
  var optOnullForm = d.querySelector("#opt_onull_text");
  var resetOptsButton = d.querySelector("#resetOptions");

  keepElementValue(queryForm, "query");
  queryForm.addEventListener("keydown", (ev) => {
    // Ctrl+Enter: do query
    if (ev.ctrlKey && ev.keyCode == 13) {
      ev.preventDefault();
      doQuery();
      return false;
    }
  });

  keepElementValue(sourceForm, "source_url");

  submitButton.addEventListener("click", () => doQuery());

  keepElementValue(optIhCheckbox, "optIh");
  keepElementValue(optOhCheckbox, "optOh");
  keepElementValue(optIfmtSelect, "optIfmt");
  keepElementValue(optOfmtSelect, "optOfmt");
  keepElementValue(optInullCheckbox, "optInull");
  keepElementValue(optInullForm, "optInullText");
  keepElementValue(optOnullCheckbox, "optOnull");
  keepElementValue(optOnullForm, "optOnullText");

  bindCheckboxToReadonly(optInullCheckbox, optInullForm);
  bindCheckboxToReadonly(optOnullCheckbox, optOnullForm);

  resetOptsButton.addEventListener("click", () => {
    optIhCheckbox.checked = false;
    optOhCheckbox.checked = false;
    optIfmtSelect.value = "CSV";
    optOfmtSelect.value = "CSV";
    optInullCheckbox.checked = false;
    optInullForm.text = null;
    optOnullCheckbox.checked = false;
    optOnullForm.text = null;
    // TODO: save to localStorage
    // TODO: apply UI changes
  });

  function keepElementValue(el, id) {
    el.addEventListener("input", () => {
      localStorage.setItem(id, getElementValue(el));
    });
    setElementValue(el, localStorage.getItem(id));
  }

  function getElementValue(el) {
    switch (el.type) {
      case "checkbox":
      case "radio":
        return el.checked;
      default:
        return el.value;
    }
  }

  function setElementValue(el, v) {
    if (!v) {
      return;
    }
    switch (el.type) {
      case "checkbox":
      case "radio":
        el.checked = v == "true";
        break;
      default:
        el.value = v;
        break;
    }
  }

  function bindCheckboxToReadonly(cb, target) {
    cb.addEventListener("change", () => target.readOnly = !cb.checked);
    target.readOnly = !cb.checked;
  }

  function doQuery() {
    var url = buildURL();
    composedForm.value = url;
    // do query, and show result
    fetch(url, { mode: "cors" })
      .then((r) => r.text())
      .then((v) => outputForm.value = v);
  }

  function buildURL() {
    var s = sourceForm.value;
    s += s.includes("?") ? "&" : "?";
    s += "trdsql=";
    s += "q:"+encodeQuery(queryForm.value);
    if (optIhCheckbox.checked) {
      s += "%3Bih:true"
    }
    if (optOhCheckbox.checked) {
      s += "%3Boh:true"
    }
    if (optIfmtSelect.value != "CSV") {
      s += "%3Bifmt:"+optIfmtSelect.value;
    }
    if (optOfmtSelect.value != "CSV") {
      s += "%3Bofmt:"+optOfmtSelect.value;
    }
    if (optInullCheckbox.checked) {
      s += "%3Binull:"+encodeQuery(optInullForm.value);
    }
    if (optOnullCheckbox.checked) {
      s += "%3Bonull:"+encodeQuery(optOnullForm.value);
    }
    // FIXME: support more options
    return s;
  }

  function encodeQuery(s) {
    s = g.encodeURIComponent(s);
    return s;
  }

})(this);
