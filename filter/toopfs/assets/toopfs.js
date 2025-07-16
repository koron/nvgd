const elSelectAll = document.getElementById('select-all');
const elFileCount = document.getElementById('file-count');
const elTotalSize = document.getElementById('total-size');
const elDestdir = document.getElementById('destdir');
const elClearDestdir = document.getElementById('clear-destdir');
const elDownload = document.getElementById('download');
const elDownloadCover = document.getElementById('download-cover');
const elDownloadMessage = document.getElementById('download-message');
const elDownloadProgress = document.getElementById('download-progress');

async function sleep(msec) {
  await new Promise((resolve) => setTimeout(resolve, msec));
}

function selectedFilesAll(base=document) {
  return base.querySelectorAll('input[type="checkbox"][data-isfile="true"]:checked');
}

function updateTotalSize() {
  const selected = selectedFilesAll();
  let count = 0;
  let totalSize = 0;
  selected.forEach((e, i) => {
    totalSize += e.dataset.size - 0;
    count++
  });
  elFileCount.innerText = count;
  elTotalSize.innerText = totalSize.toLocaleString();
  elDownload.disabled = selected.length == 0;
}

async function uploadOPFS(name, blob) {
  // Create directories if necessary
  let destdirHandle = await navigator.storage.getDirectory();
  let destName = name;
  while (true) {
    const index = destName.indexOf('/');
    if (index < 0) {
      break;
    }
    destdirHandle = await destdirHandle.getDirectoryHandle(destName.slice(0, index), { create: true });
    destName = destName.slice(index+1);
  }

  const fileHandle = await destdirHandle.getFileHandle(destName, { create: true });
  const writable = await fileHandle.createWritable();
  await writable.write(blob);
  await writable.close();
}

async function downloadFiles() {
  // Normalize destdir
  let destdir = elDestdir.value;
  if (destdir !== '') {
    destdir = destdir.replace(/(^\/+|\/+$)/, '') + '/';
  }
  const targets = selectedFilesAll()
  const count = targets.length;
  let step = 0;
  elDownloadProgress.max = count * 2;
  elDownloadProgress.value = step;
  for (file of targets) {
    const src = file.dataset.link;
    const dst = destdir + src.slice(window.location.pathname.length);
    // Download a file
    elDownloadMessage.innerText = `#${step+1}/${count*2} downloading from ${src} ...`;
    const isFileProtocol = src.startsWith('/files/') || src.startsWith('/file:///');
    const response = await fetch(src + (isFileProtocol ? '?keepcompress&all' : '?all'));
    if (!response.ok) {
      throw new Error(`Failed to fetch data from ${src}: ${response.statusText}`);
    }
    elDownloadProgress.value = ++step;
    const blob = await response.blob();
    // Upload a file to OPFS
    elDownloadMessage.innerText = `#${step+1}/${count*2} uploading to opfs://${dst} ...`;
    await uploadOPFS(dst, blob);
    elDownloadProgress.value = ++step;
  }
  elDownloadMessage.innerText = 'completed.';
}

function on_click_cleanDestdir() {
  elDestdir.value = '';
  elDestdir.focus();
}

async function on_click_download() {
  elDownloadCover.style.display = '';
  try {
    const links = selectedFilesAll()
    if (!confirm(`Download ${selectedFilesAll().length} files, total size ${elTotalSize.innerText} bytes?`)) {
      return;
    }
    await downloadFiles();
    await sleep(100);
    // Open the OPFS destination directory.
    if (confirm('Download completed.\n\nOpen the OPFS destination directory in a new tab?')) {
      const destdir = elDestdir.value.replace(/(^\/+|\/+$)/, '');
      const hashPath = destdir != '' ? '#/' + destdir + '/' : '';
      const url = `/opfs/${hashPath}`;
      window.open(url, '_blank');
    }
  } catch (err) {
    alert(`Problem occurred while downloading: ${err.message}`);
    throw err;
  } finally {
    elDownloadCover.style.display = 'none';
    elDownloadMessage.innerText = '';
    elDownloadProgress.removeAttribute('max');
    elDownloadProgress.removeAttribute('value');
  }
}

function associateCheckboxes(parentCheckbox, ul) {
  const children = Array.from(ul.querySelectorAll('ul > li > label > input[type="checkbox"]'));

  parentCheckbox.addEventListener('change', ev => {
    for (const child of children) {
      child.checked = parentCheckbox.checked;
      child.indeterminate = false;
      child.dispatchEvent(new Event('propagateParentChange'));
    }
    updateTotalSize();
  });
  parentCheckbox.addEventListener('childChange', ev => {
    const total = children.length;
    const checked = children.filter(e => e.checked).length;
    parentCheckbox.checked = checked > 0 && checked == total;
    parentCheckbox.indeterminate = checked > 0 && checked < total;
    parentCheckbox.dispatchEvent(new Event('propagateChildChange'));
  });
  parentCheckbox.addEventListener('parentChange', ev => {
    for (const child of children) {
      child.checked = parentCheckbox.checked;
      child.indeterminate = false;
      child.dispatchEvent(new Event('propagateParentChange'));
    }
  });

  for (const child of children) {
    child.addEventListener('change', ev => {
      updateTotalSize();
      parentCheckbox.dispatchEvent(new Event('childChange'));
    });
    child.addEventListener('propagateChildChange', ev => {
      parentCheckbox.dispatchEvent(new Event('childChange'));
    });
    child.addEventListener('propagateParentChange', ev => {
      child.dispatchEvent(new Event('parentChange'));
    });
  }

  updateTotalSize();
}

function rewriteNumbers(el) {
  el.querySelectorAll('.number').forEach(v => {
    v.innerText = new Number(v.innerText).toLocaleString();
  });
}

// Events

elClearDestdir.addEventListener('click', on_click_cleanDestdir);
elDownload.addEventListener('click', on_click_download);

document.body.addEventListener('htmx:afterSwap', ev => {
  const parentCheckbox = ev.target.parentElement.querySelector('li > label > input[type="checkbox"]');
  associateCheckboxes(parentCheckbox, ev.target);
  rewriteNumbers(ev.target);
});

const list = document.querySelector('#input-section > ul.file-selection');
associateCheckboxes(elSelectAll, list);
rewriteNumbers(list);
