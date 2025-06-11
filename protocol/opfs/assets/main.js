const opfs = {
  // OPFS root directory handler (FileSystemDirectoryHandle)
  root: undefined,
  // Current directory level ([FileSystemDirectoryHandle])
  dirs: [],

  // Change current dir with adding the last of hierarchy.
  pushDir(dir, name) {
    this.dirs.push({name: name ? name : dir.name, dir: dir});
  },

  // Get current direcotry (FileSystemDirectoryHandle)
  get currDir() {
    return this.dirs[this.dirs.length - 1].dir;
  },

  // mkdir creates a new directory into the current directory.
  async mkdir(name) {
    console.log('mkdir', name);
    await this.currDir.getDirectoryHandle(name, { create: true });
    await this.render();
  },

  // up moves to the parent directory of the current directory.
  async up(n = 1) {
    if (this.dirs.length <= n) {
      alert('No parent directory');
      return;
    }
    while (n > 0) {
      this.dirs.pop();
      n--;
    }
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
    if (name == '..') {
      return await this.up();
    }
    try {
      const dir = await this.currDir.getDirectoryHandle(name);
      this.pushDir(dir);
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
      alert(err);
    }
  },

  // rm removes an entry from the current directory.
  async rm(name, recursive = false) {
    try {
      await this.currDir.removeEntry(name, { recursive: recursive });
      await this.render();
    } catch (err) {
      alert(err);
    }
  },

  // Update screen
  async render() {
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
        m('a', { onclick: () => this.up(n) }, dir.name));
    }
    m.render(document.querySelector('#header'), breadcrumbs);
  },

  async renderEntries() {
    const rows = [ m('tr',
        m('th', 'Name'),
        m('th', 'Type'),
        m('th', 'Size'),
        m('th', 'Modified At'),
        m('th', 'Actions'),
    )];
    for await (const [name, handle] of this.currDir.entries()) {
      const cols = [];
      const acts = [];
      if (handle instanceof FileSystemFileHandle) {
        const file = await handle.getFile();
        cols.push(
          m('td', name),
          m('td', 'file'),
          m('td', file.size),
          m('td', new Date(file.lastModified).toLocaleString()),
        );
        // Compose actions for file
        acts.push(
          m('a', { onclick: () => this.actRm(name, false) }, 'Remove'),
          ' ',
          m('a', { onclick: () => this.actSave(name) }, 'Save as'),
        );
        if (file.size < 64*1024) {
          acts.push(' ', m('a', { onclick: () => this.actEdit(name) }, 'Edit'));
        }
        if (supportedByDuckDB(name)) {
          acts.push(' ', m('a', { onclick: () => this.actDuckDB(name) }, 'DuckDB'));
        }
      } else {
        cols.push(
          m('td',
            m('a', { onclick: () => this.cd(name), }, name + '/')),
          m('td', 'dir'),
          m('td', '(N/A)'),
          m('td', '(N/A)'),
        );
        // Compose actions for directory
        acts.push(m('a', { onclick: () => this.actRm(name, true) }, 'Remove'));
      }
      cols.push(m('td', acts));
      rows.push(m('tr', cols));
    }
    m.render(document.querySelector('#main > table'), rows);
  },

  async actRm(name, recursive) {
    if (confirm(`Are you sure you want to delete the following file/directory and its contents?\n\n${name}`)) {
      this.rm(name, recursive);
    }
  },

  async actSave(name) {
    const dst = await window.showSaveFilePicker({suggestedName: name})
    const fileHandle = await this.currDir.getFileHandle(name);
    const file = await fileHandle.getFile();
    const writable = await dst.createWritable();
    await writable.write(file);
    await writable.close();
    alert(`File ${name} in OPFS is saved as ${dst.name} in local successfully.`);
  },

  async actEdit(name) {
    const fileHandle = await this.currDir.getFileHandle(name);
    const file = await fileHandle.getFile();
    const el0 = document.querySelector('#touch-name');
    const el1 = document.querySelector('#touch-body');
    el0.value = name;
    el1.value = await file.text();

  },

  async actDuckDB(name) {
    const dirs = this.dirs.slice(1).map((e) => e.name)
    dirs.push(name);
    openWithDuckDB(dirs.join('/'));
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

function makehash(queries) {
  return queries.map(swapchars).join(',');
}

function openWithDuckDB(path) {
  const queries = [
    `CREATE VIEW opfs AS SELECT * FROM 'opfs://${path}';`,
    `SHOW opfs;`,
  ];
  const url = `${origin}/duckdb/?opfs=${path}#,${makehash(queries)}`;
  window.open(url, '_blank');
}

function supportedByDuckDB(name) {
  const supportedExtensions = ['.csv', '.xlsx', '.json', '.parquet'];
  const lastDotIndex = name.lastIndexOf('.');
  if (lastDotIndex === -1) {
    return false;
  }
  const ext = name.substring(lastDotIndex).toLowerCase();
  return supportedExtensions.includes(ext);
}

async function init() {
  const root = await navigator.storage.getDirectory();
  opfs.pushDir(root, '(Root)');
  opfs.render();
}

init();
