# OPFS Protocol — Test Plan

## Setup

NVGD must be running before test execution:

```powershell
go build -o nvgd.exe . && .\nvgd.exe
```

Tests assume NVGD is listening on `http://localhost:9280`.

Implementation files live in `tests/opfs/` and are listed by category below.

---

## Category: Rendering

**File:** `tests/opfs/rendering.spec.ts`

| # | Test | Key steps |
|---|------|-----------|
| 1 | Page title | Navigate to `/opfs/`. Assert `document.title` is `"OPFS: (Root)"` |
| 2 | All UI sections | Assert breadcrumb area, grid header (Name, Type, Size, Modified At, Actions), footer sections (action buttons, mkdir, upload, editor, download URL) are present |
| 3 | Initial button states | Assert Delete and DuckDB buttons are initially disabled |

---

## Category: Directory Management

**File:** `tests/opfs/directory.spec.ts`

| # | Test | Key steps |
|---|------|-----------|
| 4 | Create directory | Type name in `#mkdir-name`, click `#mkdir-mkdir`. Assert directory row appears with type `"dir"`, size and modified `"(N/A)"` |
| 5 | Navigate in | Click directory name. Assert breadcrumb updates, title becomes `"OPFS: (Root)/<dirname>"` |
| 6 | Navigate back | Click parent breadcrumb link. Assert listing returns to root |
| 7 | Reload | Click `#command-reload`. Assert entry count unchanged |
| 28 | Deep nesting | Create `a/b/c/d/e`, navigate to deepest. Assert breadcrumb shows full path |

---

## Category: File I/O (Editor + Upload)

**File:** `tests/opfs/file-io.spec.ts`

| # | Test | Key steps |
|---|------|-----------|
| 8 | Create file | Fill `#editor-name` and `#editor-edit`, click `#editor-save`. Assert file row appears with type `"file"`, size > 0. Assert editor fields cleared |
| 9 | Load file for editing | Click "Edit" action link on file. Assert `#editor-name` and `#editor-edit` populated with file name and content |
| 10 | Modify and save | Change `#editor-edit` content, click save. Assert file size updates |
| 11 | Tab in textarea | Press Tab. Assert tab character inserted (no focus loss) |
| 12 | Clear button | Click `#editor-clear`. Assert both fields cleared |
| 13 | Select file auto-fills name | Set `#upload-file` via `page.setInputFiles()`. Assert `#upload-name` auto-populated, `#upload-upload` becomes enabled |
| 14 | Upload file | After setting file, click Upload. Assert file row appears with correct name and size |
| 15 | Upload with custom name | Change `#upload-name` before upload. Assert file appears with custom name |

---

## Category: Download URL

**File:** `tests/opfs/download-url.spec.ts`

| # | Test | Key steps |
|---|------|-----------|
| 16 | Button state | Assert Download disabled until URL starts with `http:`/`https:` AND download-as is non-empty |
| 17 | Clear button | Assert clear enabled when either input has content |
| 18 | Download from NVGD URL | Enter URL to `/version://` and name, click Download. Assert file appears with content matching version response |

---

## Category: Bulk Operations (Selection + Delete + DuckDB)

**File:** `tests/opfs/operations.spec.ts`

| # | Test | Key steps |
|---|------|-----------|
| 19 | Single file delete | Select file checkbox, click Delete, accept confirm dialog. Assert file removed |
| 20 | Cancel delete | Select file, click Delete, dismiss dialog. Assert file remains |
| 21 | Multi-select delete | Create 3 files, select all via toggle-all, click Delete, confirm. Assert all removed |
| 22 | Toggle-all | With mix checked/unchecked, toggle-all checks/unchecks all. Assert indeterminate state works |
| 23 | Directory delete | Create nested dir with child file, delete directory. Assert both removed |
| 24 | DuckDB button state | Assert DuckDB disabled when no selection; enabled when file(s) selected |
| 25 | Open with DuckDB | Select `.csv` file, click DuckDB. Assert new page opens at `/duckdb/` with `opfs=` param and `CREATE VIEW` in hash |

---

## Category: Edge Cases

**File:** `tests/opfs/edge-cases.spec.ts`

| # | Test | Key steps |
|---|------|-----------|
| 26 | Empty dir name | Click mkdir with empty input. Assert alert: `"Need directory name"` |
| 27 | Empty file name | Click editor save with empty name. Assert alert: `"Need file name"` |
| 29 | Overwrite prompt (upload) | Upload same file twice. Assert confirm dialog shown |
| 30 | Overwrite prompt (download) | Download URL with name matching existing file. Assert confirm dialog shown |

---

## Implementation notes

- Each test gets an isolated browser context, so OPFS state is fresh per test.
- Use `page.on('dialog')` to handle alert/confirm dialogs.
- For file upload, use `page.setInputFiles()` with `{ name, mimeType, buffer }` inline objects.
- For DuckDB popup, use `context.waitForEvent('page')` and verify URL.
- The "Save as" action (`showSaveFilePicker()`) may not work in headless mode — test `npx playwright test --headed` for those.
- Run all OPFS tests: `npx playwright test tests/opfs/`
- Run a single category: `npx playwright test tests/opfs/directory.spec.ts`
- Run with visible browser: `npx playwright test tests/opfs/ --headed`
