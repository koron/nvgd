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

  // touch creates a new file or updates "Modified At" of the existing file,
  // and append the contents to the file.
  async touch(name, contents) {
    try {
      const fileHandle = await this.currDir.getFileHandle(name, { create: true });
      const file = await fileHandle.getFile();
      const writable = await fileHandle.createWritable();
      if (contents) {
        await writable.seek(file.size);
        await writable.write(contents ? contents : '');
      } else {
        await writable.truncate(file.size);
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
      if (handle instanceof FileSystemFileHandle) {
        const file = await handle.getFile();
        rows.push(m('tr',
          m('td', name),
          m('td', 'file'),
          m('td', file.size),
          m('td', new Date(file.lastModified).toLocaleString()),
          m('td', ''), // TODO: Actions
        ));
      } else {
        rows.push(m('tr',
          m('td',
            m('a', { onclick: () => this.cd(name), }, name + '/')),
          m('td', 'dir'),
          m('td', '(N/A)'),
          m('td', '(N/A)'),
          m('td', ''), // TODO: Actions
        ));
      }
    }
    m.render(document.querySelector('#main > table'), rows);
  },
};

async function init() {
  const root = await navigator.storage.getDirectory();
  opfs.pushDir(root, '(Root)');
  opfs.render();
}

init();
