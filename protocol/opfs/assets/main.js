// OPFS directory explorer

const opfs = {
  // Current directory level ([FileSystemDirectoryHandle])
  dirs: [],

  // Delegate  event to base eventTarget.
  eventTarget: new EventTarget(),

  addEventListener() {
    this.eventTarget.addEventListener.apply(this.eventTarget, arguments);
  },

  removeEventListener() {
    this.eventTarget.removeEventListener.apply(this.eventTarget, arguments);
  },

  dispatchEvent() {
    this.eventTarget.dispatchEvent.apply(this.eventTarget, arguments);
  },

  // Change current dir with adding the last of hierarchy.
  pushDir(dir, name) {
    this.dirs.push({name: name ? name : dir.name, dir: dir});
  },

  popDir(n) {
    if (this.dirs.length <= n) {
      alert('No parent directory');
      return;
    }
    this.dirs.splice(this.dirs.length - n, n);
  },

  // Get current direcotry (FileSystemDirectoryHandle)
  get currDir() {
    return this.dirs[this.dirs.length - 1].dir;
  },

  get currPath() {
    return this.dirs.map((d) => d.name).join('/');
  },

  async setCurrPath(path) {
    this.dirs.splice(0);
    const root = await navigator.storage.getDirectory();
    this.pushDir(root, '(Root)');
    if (path && path != '/') {
      const entries = path.replace(/^\/|\/$/g, '').split('/');
      for (const entry of entries) {
        const dir = await this.currDir.getDirectoryHandle(entry);
        this.pushDir(dir);
      }
    }
  },

  absPath(path='.') {
    if (path.startsWith('/')) {
      return path;
    }
    if (path.startsWith('..')) {
      let count = 0;
      while (path.startsWith('..')) {
        path = path.replace(/\.\.\/?/, '');
        count++;
      }
      if (count >= this.dirs.length) {
        return '/' + path;
      }
      return this.dirs.slice(0, this.dirs.length - count).map(d => d.dir.name).join('/') + '/' + path;
    }
    if (path.startsWith('.')) {
      return this.dirs.map(d => d.dir.name).join('/');
    }
    return this.absPath('.') + '/' + path;
  },

  updateHistory(push) {
    const path = this.dirs.map(d => d.dir.name).join('/') + '/';
    const url = path === '/' ? '.' : `#${path}`;
    if (push) {
      window.history.pushState(path, '', url);
    } else {
      window.history.replaceState(path, '', url);
    }
  },

  async enumFiles(path='.') {
    return await enumFiles(this.absPath(path));
  },

  // mkdir creates a new directory into the current directory.
  async mkdir(name) {
    await this.currDir.getDirectoryHandle(name, { create: true });
    await this.render();
  },

  async ls() {
    const entries = [];
    for await (const [name, handle] of this.currDir.entries()) {
      entries.push([name, handle]);
    }
    entries.sort((a, b) => a[0].localeCompare(b[0], undefined, { numeric: true, sensitivity: 'base' }));
    return entries;
  },

  // cd moves to a child directory in the current directory.
  async cd(name) {
    try {
      if (name == '..') {
        this.popDir(1);
      } else {
        const dir = await this.currDir.getDirectoryHandle(name);
        this.pushDir(dir);
      }
      this.updateHistory(false);
      await this.render();
    } catch (err) {
      alert(err);
    }
  },

  async loadAs(name, file) {
    if (await this.exist(name)) {
      if (!confirm(`Are you sure to overwrite a file with the same name, "${name}"?`)) {
        return false;
      }
    }
    const fileHandle = await this.currDir.getFileHandle(name, { create: true });
    const writable = await fileHandle.createWritable();
    await writable.write(file);
    await writable.close();
    await this.render();
    alert(`Uploaded "${file.name}" as "${name}" to OPFS successfully.`);
    return true;
  },

  async exist(name) {
    try {
      const handle = await this.currDir.getDirectoryHandle(name);
      return true;
    } catch (err) {
      if (!err instanceof DOMException && !['TypeMismatchError', 'NotFoundError'].includes(err.name)) {
        throw err;
      }
    }
    try {
      const handle = await this.currDir.getFileHandle(name);
      return true;
    } catch (err) {
      if (!err instanceof DOMException && err.name != 'NotFoundError') {
        throw err;
      }
    }
    return false;
  },

  // touch creates a new file or updates 'Modified At' of the existing file,
  // and append the contents to the file.
  async touch(name, contents) {
    try {
      const fileHandle = await this.currDir.getFileHandle(name, { create: true });
      const file = await fileHandle.getFile();
      const writable = await fileHandle.createWritable();
      if (contents) {
        await writable.write(contents ? contents : '');
      }
      await writable.close();
      await this.render();
    } catch (err) {
      this.alertErr(err);
    }
  },

  async alertErr(err) {
    if (err instanceof DOMException && err.name === 'NoModificationAllowedError') {
      alert('FILE MAY BE LOCKED.\n\nIt looks like the file is open somewhere else and locked. Try closing any tabs that might be using it, like the DuckDB Shell.');
    } else {
      alert(err);
    }
  },

  // rm removes an entry from the current directory.
  async rm(name, recursive = false) {
    try {
      await this.currDir.removeEntry(name, { recursive: recursive });
      await this.render();
    } catch (err) {
      this.alertErr(err);
    }
  },

  // clearAll deletes all entries from the current directory.
  async clearAll() {
    try {
      path = this.currPath
      if (!confirm(`Are you sure to delete all contents in the current direcotry?\n\n${path}`)) {
        return;
      }
      for await (const [name, handle] of this.currDir.entries()) {
        await this.currDir.removeEntry(name, { recursive: true });
      }
      await this.render();
    } catch (err) {
      this.alertErr(err);
    }
  },

  // Update screen
  async render() {
    const path = this.dirs.map(d => d.dir.name).join('/') + '/';
    document.querySelector('title').innerText = 'OPFS: ' + path;
    await this.renderHeader();
    await this.renderEntries();
  },

  async renderHeader() {
    const breadcrumbs = [];
    const last = this.dirs.length - 1;
    for (const i in this.dirs) {
      const dir = this.dirs[i];
      if (breadcrumbs.length != 0) {
        breadcrumbs.push(m('span', ' / '));
      }
      const n = last - i;
      breadcrumbs.push(n == 0 ?
        m('span', dir.name) :
        m('a', { onclick: () => this.actCd(n) }, dir.name));
    }
    m.render(document.querySelector('#header'), breadcrumbs);
  },

  async renderEntries() {
    const rows = [ m('div.grid-header',
        m('div',
          m('input', {
            type: 'checkbox',
            id: 'toggle-selection-all',
            onclick: () => this.selectionToggleAll(),
          }),
          m('span', 'Name')),
        m('div', 'Type'),
        m('div', 'Size'),
        m('div', 'Modified At'),
        m('div', 'Actions'),
    )];
    for await (const [name, handle] of this.currDir.entries()) {
      const cols = [];
      const acts = [];
      if (handle instanceof FileSystemFileHandle) {
        const file = await handle.getFile();
        cols.push(
          m('label.name',
            m('input.selectedFile', {
              type: 'checkbox',
              name: name,
              onchange: () => this.selectionChanged(),
            }), ' ',
            name),
          m('div', 'file'),
          m('div.size', file.size),
          m('div', new Date(file.lastModified).toLocaleString()),
        );
        // Compose actions for file
        acts.push(
          m('a', { onclick: () => this.actSave(name) }, icon('download'), 'Save as'),
        );
        if (file.size < 64*1024) {
          acts.push(' ', m('a', { onclick: () => this.actEdit(name) }, icon('edit'), 'Edit'));
        }
      } else {
        const displayName = name + '/';
        cols.push(
          m('label.name',
            m('input.selectedFile', {
              type: 'checkbox',
              name: displayName,
              onchange: () => this.selectionChanged(),
            }), ' ',
            m('a', {
              onclick: (e) => { e.preventDefault(); this.actCd(name) },
            }, displayName)),
          m('div', 'dir'),
          m('div.size', '(N/A)'),
          m('div', '(N/A)'),
        );
      }
      cols.push(m('div', acts));
      rows.push(m('div.grid-row', cols));
    }
    m.render(document.querySelector('#main > .directory'), rows);
    this.selectionChanged();
  },

  // Selection

  selectedFiles() {
    return Array.from(document.querySelectorAll('input.selectedFile:checked')).map(e => e.name);
  },

  selectionChanged() {
    const all = document.querySelectorAll('input.selectedFile');
    const selected = this.selectedFiles();

    this.dispatchEvent(new CustomEvent('selection:change', { detail: selected.length }));

    const toggle = document.querySelector('#toggle-selection-all');
    toggle.checked = selected.length > 0 && selected.length == all.length;
    toggle.indeterminate = selected.length > 0 && selected.length < all.length
  },

  async selectionToggleAll() {
    const toggle = document.querySelector('#toggle-selection-all');
    for (const checkbox of document.querySelectorAll('input.selectedFile')) {
      checkbox.checked = toggle.checked;
    }
    this.selectionChanged()
  },

  async unselectAll() {
    for (const checkbox of document.querySelectorAll('input.selectedFile')) {
      checkbox.checked = false;
    }
    const toggle = document.querySelector('#toggle-selection-all');
    toggle.checked = false;
    toggle.indeterminate = false;
    this.selectionChanged()
  },

  // Actions

  async actDeleteSelectedFiles() {
    const files = this.selectedFiles();
    if (!confirm(`Are you sure you want to delete the following files/directories and its contents?\n\n- ${files.join('\n- ')}`)) {
      return;
    }
    try {
      for (const name of files) {
        if (name.endsWith('/')) {
          await this.currDir.removeEntry(name.slice(0, name.length - 1), { recursive: true });
        } else {
          await this.currDir.removeEntry(name, { recursive: false });
        };
      }
      await this.unselectAll();
    } catch (err) {
      this.alertErr(err);
    }
    await this.render()
  },

  async actDuckDBWithSelectedFiles() {
    const files = this.selectedFiles();
    let paths = [];
    for (const file of files) {
      if (file.endsWith('/')) {
        paths = paths.concat(await this.enumFiles(file));
      } else {
        paths.push(this.absPath(file));
      }
    }
    // Strip a leading '/' slash from all paths.
    paths = paths.map(e => e.slice(1));
    openWithDuckDB(paths);
  },

  async actCd(nameOrCount) {
    await this.unselectAll();
    if (typeof nameOrCount === 'number') {
      this.popDir(nameOrCount);
    } else {
      const dir = await this.currDir.getDirectoryHandle(nameOrCount);
      this.pushDir(dir);
    }
    this.updateHistory(true);
    await this.render();
  },

  async actSave(name) {
    try {
      const dst = await window.showSaveFilePicker({suggestedName: name})
      const fileHandle = await this.currDir.getFileHandle(name);
      const file = await fileHandle.getFile();
      const writable = await dst.createWritable();
      await writable.write(file);
      await writable.close();
      alert(`File ${name} in OPFS is saved as ${dst.name} in local successfully.`);
    } catch (err) {
      // Ignore file picker's cancellation.
      if (err instanceof DOMException && err.name == 'AbortError') {
        return;
      }
      this.alertErr(err);
    }
  },

  async actEdit(name) {
    const fileHandle = await this.currDir.getFileHandle(name);
    const file = await fileHandle.getFile();
    const el0 = document.querySelector('#editor-name');
    const el1 = document.querySelector('#editor-edit');
    el0.value = name;
    el1.value = await file.text();
  },
}

function swapchars(str) {
  let newStr = '';
  for (let i = 0; i < str.length; i++) {
    const char = str[i];
    switch (char) {
      case ' ':
        newStr += '-';
        break;
      case '-':
        newStr += ' ';
        break;
      case ';':
        newStr += '~';
        break;
      case '~':
        newStr += ';';
        break;
      default:
        newStr += char;
    }
  }
  return newStr;
}

function icon(name) {
  return m('span.material-symbols', name);
}

function makehash(queries) {
  return queries.map(swapchars).join(',');
}

function openWithDuckDB(paths) {
  const qparams = paths.map(v => 'opfs=' + encodeURIComponent(v)).join('&');
  const queries = paths.filter(v => supportedByDuckDB(v)).map((v, i) => `CREATE VIEW opfs${i} AS SELECT * FROM 'opfs://${v}';`);
  queries.push('SHOW TABLES;');
  const url = `${origin}/duckdb/?${qparams}#,${makehash(queries)}`;
  window.open(url, '_blank');
}

function supportedByDuckDB(name) {
  const supportedExtensions = [
    '.csv', '.csv.gz', '.csv.zst',
    '.tsv', '.tsv.gz', '.tsv.zst',
    // See https://github.com/koron/nvgd/issues/157 why stop to create views for xlsx files.
    //'.xlsx',
    '.json',
    '.parquet',
  ];
  for (const ext of supportedExtensions) {
    if (name.endsWith(ext)) {
      return true;
    }
  }
  return false;
}

async function openDir(path='/') {
  let dir = await navigator.storage.getDirectory();
  const entries = path.replace(/^\/|\/$/g, '').split('/').filter(e => e !== '');
  for (const entry of entries) {
    dir = await dir.getDirectoryHandle(entry);
  }
  return dir
}

async function enumFiles(root='/') {
  root = root.replace(/\/+$/, '');
  const dir = await openDir(root);
  let files = [];
  for await (const [name, handle] of dir.entries()) {
    const path = root + '/' + name;
    if (handle instanceof FileSystemDirectoryHandle) {
      files = files.concat(await enumFiles(path));
    } else {
      files.push(path)
    }
  }
  return files;
}

//////////////////////////////////////////////////////////////////////////////
// Command pallet

(function() {
  // Reload
  const elReload = document.getElementById('command-reload');
  async function reload() {
    await opfs.renderEntries();
  }
  elReload.addEventListener('click', reload);

  // Delete
  const elDelete = document.getElementById('command-delete');
  elDelete.addEventListener('click', () => opfs.actDeleteSelectedFiles());
  opfs.addEventListener('selection:change', ev => elDelete.disabled = ev.detail == 0);

  // DuckDB
  const elDuckDB = document.getElementById('command-duckdb');
  elDuckDB.addEventListener('click', () => opfs.actDuckDBWithSelectedFiles());
  opfs.addEventListener('selection:change', ev => elDuckDB.disabled = ev.detail == 0);
})();

//////////////////////////////////////////////////////////////////////////////
// Make directory

(function() {
  const elName = document.getElementById('mkdir-name');
  const elMkdir = document.getElementById('mkdir-mkdir');

  async function mkdir() {
    try {
      const name = elName.value.trim();
      if (name == '') {
        alert('Need directory name');
        return;
      }
      await opfs.mkdir(name);
      elName.value = '';
    } catch (err) {
      opfs.alertErr(`Failed to create "${name}" directory, because of:\n\n${err}`);
    }
  }

  elMkdir.addEventListener('click', mkdir);
})();

//////////////////////////////////////////////////////////////////////////////
// Upload a local file

(function() {
  const elFile = document.getElementById('upload-file');
  const elName = document.getElementById('upload-name');
  const elUpload = document.getElementById('upload-upload');

  async function fileChanged() {
    if (elFile.files.length == 0) {
      elName.value = '';
      elUpload.disabled = true;
      return;
    }
    const name = elFile.files[0].name;
    elName.value = name;
    elUpload.disabled = false;
  }

  async function upload() {
    const name = elName.value;
    const file = elFile.files[0];
    if (await opfs.loadAs(name, file)) {
      elName.value = '';
      elFile.value = '';
      elUpload.disabled = true;
    }
  }

  elFile.addEventListener('change', fileChanged);
  elUpload.addEventListener('click', upload);
})();

//////////////////////////////////////////////////////////////////////////////
// Upload a local file

(function() {
  const elName = document.getElementById('editor-name');
  const elSave = document.getElementById('editor-save');
  const elClear = document.getElementById('editor-clear');
  const elEdit = document.getElementById('editor-edit');

  async function save() {
    try {
      const name = elName.value;
      const body = elEdit.value;
      if (name == '') {
        alert('Need file name');
        return;
      }
      await opfs.touch(name, body);
      clear();
    } catch (err) {
      alert(`Failed to create "${name}" file, because of:\n\n${err}`);
    }
  }

  function clear() {
    elName.value = '';
    elEdit.value = '';
  }

  // Pressing the Tab key inserts a tab character.
  function handleTabKey(e) {
    if (e.keyCode == 9) {
      e.preventDefault();
      document.execCommand('insertText', false, '\t');
    }
  }

  elSave.addEventListener('click', save);
  elClear.addEventListener('click', clear);
  elEdit.addEventListener('keydown', handleTabKey);
})();

//////////////////////////////////////////////////////////////////////////////
// Download URL

(function() {
  const elURL = document.getElementById('download-url');
  const elAs = document.getElementById('download-as');
  const elDownload = document.getElementById('download-download');
  const elClear = document.getElementById('download-clear');

  function updateButtons() {
    const url = elURL.value;
      const as = elAs.value;
    elDownload.disabled = !((url.startsWith('http:') || url.startsWith('https:')) && as.length > 0);
    elClear.disabled = !(url.length > 0 || as.length > 0);
  }

  function clearInputs() {
    elURL.value = '';
    elAs.value = '';
  }

  async function download() {
    try {
      const url = elURL.value;
      const as = elAs.value;
      if (await opfs.exist(as)) {
        if (!confirm(`Are you sure to overwrite a file with the same name, "${as}"?`)) {
          return false;
        }
      }

      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`Failed to fetch data from ${url}: ${response.statusText}`);
      }
      const blob = await response.blob();
      const fileHandle = await opfs.currDir.getFileHandle(as, { create: true });
      const writable = await fileHandle.createWritable();
      await writable.write(blob);
      await writable.close();
      await opfs.render()

      return true;
    } catch (err) {
      await opfs.alertErr(err);
    }
  }

  elURL.addEventListener('input', updateButtons);
  elAs.addEventListener('input', updateButtons);
  elClear.addEventListener('click', clearInputs);
  elDownload.addEventListener('click', download);
})();

// Initialize the application.

async function init() {
  onpopstate = async (ev) => {
    await opfs.setCurrPath(ev.state);
    await opfs.unselectAll();
    await opfs.render();
  };

  const path = window.location.hash.substring(1);
  await opfs.setCurrPath(path);
  opfs.updateHistory(false);

  await opfs.render();
}

init();
