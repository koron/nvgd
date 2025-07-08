const elSelectAll = document.getElementById('select-all');
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

function dirsAll(base=document) {
  return base.querySelectorAll('input[type="checkbox"][data-isfile="false"]');
}

function filesAll(base=document) {
  return base.querySelectorAll('input[type="checkbox"][data-isfile="true"]');
}

function selectedFilesAll(base=document) {
  return base.querySelectorAll('input[type="checkbox"][data-isfile="true"]:checked');
}

function updateTotalSize() {
  const selected = selectedFilesAll();
  let totalSize = 0;
  selected.forEach((e, i) => {
    totalSize += e.dataset.size - 0;
  });
  elTotalSize.innerText = totalSize;

  const count = filesAll().length;
  elSelectAll.checked = selected.length > 0 && selected.length == count;
  elSelectAll.indeterminate = selected.length > 0 && selected.length < count;

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
    const response = await fetch(src, { headers: { "Accept-Encoding": "" } });
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

function on_change_selectAll(ev) {
  const v = ev.target.checked;
  filesAll().forEach((e, i) => e.checked = v);
  updateTotalSize();
}

function on_click_cleanDestdir() {
  elDestdir.value = '';
  elDestdir.focus();
}

function on_change_fileSelection() {
  updateTotalSize();
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
    alert(`There was a problem saving the data: ${err.message}`);
    throw err;
  } finally {
    elDownloadCover.style.display = 'none';
    elDownloadMessage.innerText = '';
    elDownloadProgress.removeAttribute('max');
    elDownloadProgress.removeAttribute('value');
  }
}

// Events

elSelectAll.addEventListener('change', on_change_selectAll);
elClearDestdir.addEventListener('click', on_click_cleanDestdir);
elDownload.addEventListener('click', on_click_download);

document.body.addEventListener('htmx:afterSwap', ev => {
  console.log(ev.target);
  filesAll(ev.target).forEach(e => e.addEventListener('change', on_change_fileSelection));
  updateTotalSize();
});

filesAll().forEach((e) => e.addEventListener('change', on_change_fileSelection));
updateTotalSize();
