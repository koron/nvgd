((g) => {
  "use strict"

  var d = g.document;
  var localStorage = g.localStorage;

  var sourceForm = d.querySelector("#source_url");
  var seriesDirSelector = d.querySelector("#series_direction");
  var formatSelector = d.querySelector("#format");
  var chartTypeSelector = d.querySelector("#chart_type");
  var submitButton = d.querySelector("#submit");
  var composedForm = d.querySelector("#composed_url");
  var outputFrame = d.querySelector("#output");
  var copyUrlButton = d.querySelector("#copy_url");

  var titleTextForm = d.querySelector("#title_text");
  var titleSubtextForm = d.querySelector("#title_subtext");
  var legendShowCheckbox = d.querySelector("#legend_show");

  keepElementValue(sourceForm, "source_url");
  keepElementValue(seriesDirSelector, "series_dir");
  keepElementValue(formatSelector, "format");
  keepElementValue(chartTypeSelector, "chart_type");

  submitButton.addEventListener("click", () => doQuery());

  d.addEventListener("keydown", (ev) => {
    // Ctrl+Enter: do query
    if (ev.ctrlKey && ev.keyCode == 13) {
      ev.preventDefault();
      doQuery();
      return false;
    }
  });

  copyUrlButton.addEventListener("click", () => {
    g.navigator.clipboard.writeText(composedForm.value)
      .then(() => {}, () => g.alert("Copy failed"))
  });

  function doQuery() {
    var url = buildURL();
    composedForm.value = url;
    outputFrame.src = url;
  }

  var storagePrefix = "echarts_";

  function keepElementValue(el, rawID) {
    var id = storagePrefix + rawID;
    el.addEventListener("input", () => {
      localStorage.setItem(id, getElementValue(el));
    });
    setElementValue(el, localStorage.getItem(id));
  }

  function saveItem(id, value) {
    localStorage.setItem(storagePrefix + id, value);
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

  function getTitleOpts() {
    var opts = {};
    if (titleTextForm.value != "") {
      opts["text"] = titleTextForm.value;
    }
    if (titleSubtextForm.value != "") {
      opts["subtext"] = titleSubtextForm.value;
    }
    return opts;
  }

  function getLegendOpts() {
    var opts = {};
    if (legendShowCheckbox.checked) {
      opts["show"] = true;
    }
    return opts;
  }

  function buildURL() {
    var s = sourceForm.value;
    s += s.includes("?") ? "&" : "?";
    s += "echarts=";
    s += "t:" + chartTypeSelector.value;
    if (seriesDirSelector.value != "column") {
      s += "%3Bd:" + seriesDirSelector.value;
    }
    if (formatSelector.value != "CSV") {
      s += "%3Bf:" + formatSelector.value;
    }
    var titleOpts = getTitleOpts();
    if (Object.keys(titleOpts).length > 0) {
      s += "%3BtitleOpts:" + JSON.stringify(titleOpts);
    }
    var legendOpts = getLegendOpts();
    if (Object.keys(legendOpts).length > 0) {
      s += "%3BlegendOpts:" + JSON.stringify(legendOpts);
    }
    return s;
  }

  function encodeQuery(s) {
    s = g.encodeURIComponent(s);
    return s;
  }

})(this);
