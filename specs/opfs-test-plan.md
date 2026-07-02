# OPFS (Origin Private File System) Protocol - Comprehensive E2E Test Plan

## Overview

This test plan covers the OPFS protocol implementation in NVGD, including:
- **Main OPFS Protocol** (`/opfs/`): Directory navigation, file operations, multi-selection, DuckDB integration, and download URL to OPFS
- **toopfs Filter**: Downloading files from NVGD URLs to OPFS

**Base URL**: `http://localhost:9280`  
**Assumption**: NVGD server is running and accessible at the base URL.

---

## Suite 1: OPFS Page Loading and Initial State

### Test 1.1: OPFS Page Loads and Displays Root Directory Listing

**File**: `tests/opfs-core/opfs-page-loads.spec.ts`

**Steps**:
1. Navigate to `http://localhost:9280/opfs/`
2. Verify the page title is `OPFS: (Root)`
3. Verify the directory listing displays the root contents
4. Verify the breadcrumb shows `(Root)`

**Expected**:
- Page loads successfully with no JavaScript errors in console
- Title and breadcrumb match the root directory

---

### Test 1.2: OPFS Page Structure and Elements

**File**: `tests/opfs-core/opfs-page-structure.spec.ts`

**Steps**:
1. Navigate to `http://localhost:9280/opfs/`
2. Verify the page contains:
   - A header section with breadcrumb trail
   - A directory listing area (grid table)
   - Footer controls: Reload, Delete, DuckDB buttons
   - Make directory input and button
   - Upload file input and button
   - Simple Editor section (name, textarea, save, clear buttons)
   - Download URL section (URL input, name input, download button)

**Expected**:
- All UI elements are present and visible
- Delete and DuckDB buttons are initially disabled
- Download button is initially disabled

---

## Suite 2: Directory Navigation

### Test 2.1: Click into Subdirectory

**File**: `tests/opfs-core/opfs-directory-navigation.spec.ts`

**Steps**:
1. Create a directory named `testdir`
2. Create a subdirectory `testdir/subdir`
3. Navigate into `testdir` by clicking on it
4. Verify the breadcrumb trail shows `(Root) / testdir`
5. Verify the directory listing shows the contents of `testdir`
6. Verify the page title updates to `OPFS: testdir`

**Expected**:
- Navigation into directory works correctly
- Breadcrumb trail updates with clickable parent directories
- Page title reflects the current directory path

---

### Test 2.2: Navigate Up with ..

**File**: `tests/opfs-core/opfs-navigate-up.spec.ts`

**Steps**:
1. Create a directory `testdir`
2. Create a subdirectory `testdir/subdir`
3. Navigate into `testdir/subdir`
4. Click the `..` link in the breadcrumb
5. Verify the listing shows the contents of `testdir`
6. Verify the breadcrumb trail reflects the parent directory
7. Attempt to navigate up from the root directory
8. Verify an alert `No parent directory` is shown

**Expected**:
- `..` navigates up one level
- Root directory navigation up shows error alert
- Breadcrumb trail correctly reflects the current path

---

### Test 2.3: Breadcrumb Navigation

**File**: `tests/opfs-core/opfs-breadcrumb.spec.ts`

**Steps**:
1. Create nested directories: `a/b/c/`
2. Navigate into `a/b/c/`
3. Click on `b` in the breadcrumb trail
4. Verify the listing shows the contents of `a/b/`
5. Click on `a` in the breadcrumb trail
6. Verify the listing shows the contents of `a/`

**Expected**:
- Clicking any breadcrumb navigates directly to that level
- Breadcrumb trail correctly represents the full path hierarchy

---

### Test 2.4: History Back/Forward

**File**: `tests/opfs-core/opfs-history.spec.ts`

**Steps**:
1. Navigate to `http://localhost:9280/opfs/`
2. Create a directory `testdir`
3. Navigate into `testdir`
4. Use browser Back button
5. Verify the OPFS UI restores the previous directory view
6. Verify the browser history URL updates with hash fragments

**Expected**:
- Back button navigates to the previous OPFS directory
- URL hash fragments update correctly for OPFS paths

---

## Suite 3: Create Directory Operations

### Test 3.1: Create Directory

**File**: `tests/opfs-core/opfs-mkdir.spec.ts`

**Steps**:
1. Enter `mydir` in the "Directory name to create" input
2. Click "Create new directory" button
3. Verify a new directory `mydir/` appears in the listing
4. Click into `mydir` and verify it is empty
5. Enter empty name and click "Create new directory"
6. Verify an alert `Need directory name` is shown
7. Enter an existing directory name and click "Create new directory"
8. Verify an error is thrown

**Expected**:
- Directory creation succeeds and appears in listing
- Empty name shows error alert
- Duplicate name throws an error

---

### Test 3.2: Create Nested Directories

**File**: `tests/opfs-core/opfs-nested-mkdir.spec.ts`

**Steps**:
1. Create directory `project`
2. Navigate into `project`
3. Create directory `src`
4. Navigate into `src`
5. Create directory `main`
6. Navigate back up and verify the full nested structure
7. Verify breadcrumb trail correctly represents the nested path

**Expected**:
- Nested directory creation works correctly
- Breadcrumb trail accurately reflects the full path hierarchy

---

## Suite 4: File Upload Operations

### Test 4.1: Upload Local File

**File**: `tests/opfs-core/opfs-upload-file.spec.ts`

**Steps**:
1. Select a local file using the file picker
2. Verify the file name populates the "File name to upload" field
3. Click "Upload" button
4. Verify success alert: `Uploaded "filename" as "filename" to OPFS successfully.`
5. Verify the uploaded file appears in the directory listing with correct size
6. Verify the file shows the correct last modified date
7. Upload an empty file and verify size is 0
8. Upload a file with a name that already exists
9. Verify a confirmation dialog appears before overwriting
10. After confirming overwrite, verify the file listing shows updated size/date
11. Upload a file with an empty name
12. Verify an alert `Need file name` is shown

**Expected**:
- File upload succeeds and appears in listing
- Empty file upload shows size 0
- Overwrite confirmation dialog appears for existing files
- Empty name shows error alert

---

### Test 4.2: Upload File with Custom Name

**File**: `tests/opfs-core/opfs-upload-custom-name.spec.ts`

**Steps**:
1. Select a local file `source.txt`
2. Change the file name in the "File name to upload" field to `dest.txt`
3. Click "Upload" button
4. Verify the file appears under its new name `dest.txt`
5. Verify the original file `source.txt` does not appear in the listing

**Expected**:
- File is uploaded with the custom name
- Original file name is not created

---

## Suite 5: File Editor Operations

### Test 5.1: Create New File via Editor

**File**: `tests/opfs-core/opfs-editor-create.spec.ts`

**Steps**:
1. Enter `test.txt` in the "File name to create or update" field
2. Enter file content in the editor textarea
3. Click "Create or update a file" button
4. Verify success alert: `Uploaded "test.txt" as "test.txt" to OPFS successfully.`
5. Verify the new file appears in the directory listing with entered content
6. Verify the file size reflects the content entered

**Expected**:
- New file creation via editor succeeds
- File content is correctly stored in OPFS

---

### Test 5.2: Update Existing File via Editor

**File**: `tests/opfs-core/opfs-editor-update.spec.ts`

**Steps**:
1. Create a file with initial content using the editor
2. Enter the same file name in the editor
3. Modify the content in the editor textarea
4. Click "Create or update a file" button
5. Verify the file is updated in place (no duplicate created)
6. Verify the file size reflects the updated content

**Expected**:
- Existing file is updated in place
- File size reflects the updated content

---

### Test 5.3: Overwrite Confirmation

**File**: `tests/opfs-core/opfs-editor-overwrite.spec.ts`

**Steps**:
1. Create a file with initial content
2. Enter the same file name in the editor
3. Click "Create or update a file" without modifying content
4. Verify a confirmation dialog appears asking to overwrite
5. Cancel the confirmation - verify the file remains unchanged
6. Re-enter the file name and click "Create or update a file"
7. Confirm the overwrite - verify the file is updated

**Expected**:
- Overwrite confirmation dialog appears
- Cancel preserves original content
- Confirm updates the file

---

### Test 5.4: File Size Limit (64KiB)

**File**: `tests/opfs-core/opfs-editor-size-limit.spec.ts`

**Steps**:
1. Create a file with content larger than 64KiB (65536 bytes)
2. Verify the "Edit" action is NOT available for files >= 64KiB
3. Verify the file is still visible in the listing with correct size
4. Verify the file is still downloadable via "Save as" action

**Expected**:
- Files >= 64KiB cannot be edited via the editor
- Files still show in listing and are downloadable

---

### Test 5.5: Tab Key Handling

**File**: `tests/opfs-core/opfs-editor-tab-key.spec.ts`

**Steps**:
1. Focus the editor textarea
2. Press the Tab key while in the editor
3. Verify a tab character is inserted into the editor content
4. Save the file and verify the tab character is preserved

**Expected**:
- Tab key inserts a tab character in the editor
- Tab character is preserved when saving the file

---

### Test 5.6: Clear Button

**File**: `tests/opfs-core/opfs-editor-clear.spec.ts`

**Steps**:
1. Enter a file name and content in the editor
2. Click the "Clear" button
3. Verify the file name field is cleared
4. Verify the editor textarea is cleared
5. Verify the previously saved file still exists in OPFS

**Expected**:
- Clear button resets editor fields
- Previously saved file remains in OPFS

---

## Suite 6: File Download Operations

### Test 6.1: Save as - Download File from OPFS to Local

**File**: `tests/opfs-core/opfs-save-as.spec.ts`

**Steps**:
1. Create a file with known content in OPFS
2. Click the "Save as" action (download icon) next to the file
3. Verify a file picker dialog appears with the file name pre-filled
4. Select a destination location and save the file
5. Verify a success alert confirming the file was saved
6. Verify the saved file contains the expected content
7. Verify the file is saved with the original OPFS file name

**Expected**:
- File download to local disk works correctly
- Saved file contains the expected content

---

### Test 6.2: Cancel File Download

**File**: `tests/opfs-core/opfs-save-as-cancel.spec.ts`

**Steps**:
1. Create a file in OPFS
2. Click the "Save as" action on the file
3. Cancel the file picker dialog without saving
4. Verify the file still exists in OPFS unchanged

**Expected**:
- Canceling the file picker does not affect the OPFS file

---

## Suite 7: File Deletion Operations

### Test 7.1: Delete Selected Files

**File**: `tests/opfs-core/opfs-delete-files.spec.ts`

**Steps**:
1. Create two files: `file1.txt` and `file2.txt`
2. Check the checkbox next to `file1.txt`
3. Check the checkbox next to `file2.txt`
4. Verify the "Delete" button becomes enabled
5. Click the "Delete" button
6. Verify a confirmation dialog appears listing the files to delete
7. Confirm the deletion
8. Verify both files are removed from the directory listing
9. Verify the "Delete" button becomes disabled again
10. Delete a file with an empty name
11. Verify an alert is shown

**Expected**:
- Selected files are deleted after confirmation
- Delete button is enabled only when files are selected
- Delete button becomes disabled after deletion

---

### Test 7.2: Delete Directory Recursively

**File**: `tests/opfs-core/opfs-delete-directory.spec.ts`

**Steps**:
1. Create a directory `mydir` with files inside (`file1.txt`, `file2.txt`)
2. Check the checkbox next to `mydir/`
3. Click the "Delete" button
4. Confirm the deletion
5. Verify the directory and all its contents are removed
6. Verify the directory does not appear in the listing
7. Delete an empty directory
8. Verify the deletion succeeds

**Expected**:
- Directory and all contents are deleted recursively
- Empty directory deletion succeeds

---

### Test 7.3: Delete Single File

**File**: `tests/opfs-core/opfs-delete-single-file.spec.ts`

**Steps**:
1. Create a file `single.txt` in OPFS
2. Check the checkbox next to the file
3. Click the "Delete" button
4. Confirm the deletion
5. Verify the file is removed from the listing
6. Verify other files in the same directory remain unaffected

**Expected**:
- Single file deletion works correctly
- Other files in the directory are unaffected

---

## Suite 8: File Selection Operations

### Test 8.1: Select All

**File**: `tests/opfs-core/opfs-selection-all.spec.ts`

**Steps**:
1. Create three files in OPFS
2. Click the "Select all" checkbox (top-left)
3. Verify all files are checked
4. Verify the "Delete" and "DuckDB" buttons become enabled
5. Uncheck one file
6. Verify the "Select all" checkbox becomes indeterminate
7. Click the "Select all" checkbox again
8. Verify all files are selected again

**Expected**:
- Select all checks all files
- Action buttons become enabled when files are selected
- Indeterminate state is shown for partial selection

---

### Test 8.2: Select Individual Files

**File**: `tests/opfs-core/opfs-selection-individual.spec.ts`

**Steps**:
1. Create two files in OPFS
2. Check the checkbox for the first file only
3. Verify the "Select all" checkbox is unchecked
4. Check the checkbox for the second file
5. Verify the "Select all" checkbox becomes indeterminate
6. Uncheck the second file
7. Verify the "Select all" checkbox becomes unchecked

**Expected**:
- Individual file selection works correctly
- Select all checkbox reflects the selection state

---

### Test 8.3: Mixed Files and Directories

**File**: `tests/opfs-core/opfs-selection-mixed.spec.ts`

**Steps**:
1. Create a directory and a file in OPFS
2. Check the checkbox for the file
3. Check the checkbox for the directory
4. Verify both are selected simultaneously
5. Verify the "Delete" button is enabled
6. Verify the "DuckDB" button is enabled

**Expected**:
- Files and directories can be selected simultaneously
- Action buttons are enabled when any items are selected

---

## Suite 9: Download URL to OPFS

### Test 9.1: Download URL to OPFS

**File**: `tests/opfs-core/opfs-download-url.spec.ts`

**Steps**:
1. Navigate to an NVGD URL that returns content (e.g., `/help/` or `/version/` or a `file://` URL)
2. Enter the URL in the "URL to be downloaded" field
3. Enter a file name in the "Name to save" field
4. Verify the "Download" button becomes enabled
5. Click the "Download" button
6. Verify the file appears in the OPFS directory listing
7. Verify the file contains the expected content from the URL
8. Download a file with a name that already exists
9. Verify a confirmation dialog appears
10. After confirming overwrite, verify the file is updated
11. Enter an invalid URL
12. Verify an error alert is shown
13. Enter a URL without a filename
14. Verify the Download button is disabled
15. Enter a URL without `http://` or `https://`
16. Verify the Download button is disabled
17. A URL that returns an error
18. Verify an error alert is shown

**Expected**:
- URL download to OPFS works correctly
- File content matches the source URL
- Overwrite confirmation appears for existing files
- Invalid URLs show error alerts
- Download button is disabled for invalid inputs

---

### Test 9.2: Download from NVGD Resource

**File**: `tests/opfs-core/opfs-download-nvgd-resource.spec.ts`

**Steps**:
1. Use an NVGD `file://` URL as the download source
2. Enter the NVGD URL in the download URL field
3. Enter a destination filename
4. Click Download and verify the file was saved to OPFS with correct content
5. Use the NVGD `/version/` endpoint as the download source
6. Verify the downloaded version file contains the nvgd version string

**Expected**:
- NVGD resources can be downloaded to OPFS
- Downloaded content is correct

---

## Suite 10: DuckDB Integration

### Test 10.1: Open Supported File Types

**File**: `tests/opfs-core/opfs-duckdb-supported.spec.ts`

**Steps**:
1. Create a CSV file in OPFS with tabular data
2. Check the checkbox next to the CSV file
3. Click the "DuckDB" button
4. Verify a DuckDB WASM shell opens in a new tab
5. Verify the DuckDB shell shows a view named `opfs0`
6. Verify the DuckDB shell shows the data from the CSV file
7. Close the DuckDB tab

**Expected**:
- DuckDB shell opens for supported file types
- Views are created and queryable

---

### Test 10.2: Open Multiple Files

**File**: `tests/opfs-core/opfs-duckdb-multiple.spec.ts`

**Steps**:
1. Create two CSV files in OPFS (`data1.csv` and `data2.csv`)
2. Check both file checkboxes
3. Click the "DuckDB" button
4. Verify the DuckDB shell shows views `opfs0` and `opfs1`
5. Verify each view contains data from the corresponding file

**Expected**:
- Multiple files create sequential views (`opfs0`, `opfs1`, etc.)
- Each view contains the correct data

---

### Test 10.3: Open Directory Recursively

**File**: `tests/opfs-core/opfs-duckdb-directory.spec.ts`

**Steps**:
1. Create a directory with multiple CSV files inside
2. Check the checkbox next to the directory
3. Click the "DuckDB" button
4. Verify the DuckDB shell creates views for all files in the directory
5. Verify all views are queryable

**Expected**:
- Directory selection recursively opens all files
- All files create DuckDB views

---

### Test 10.4: Unsupported File Types

**File**: `tests/opfs-core/opfs-duckdb-unsupported.spec.ts`

**Steps**:
1. Create a binary file (e.g., an image or executable) in OPFS
2. Check the checkbox next to the unsupported file
3. Click the "DuckDB" button
4. Verify the DuckDB shell opens but does not create a view for the unsupported file
5. Verify only supported file types generate DuckDB views

**Expected**:
- Unsupported file types do not create DuckDB views
- DuckDB shell still opens but shows no views for unsupported files

---

## Suite 11: OPFS File Listing

### Test 11.1: File Metadata Display

**File**: `tests/opfs-core/opfs-file-metadata.spec.ts`

**Steps**:
1. Create a file with known content
2. Verify the file listing displays the correct file size
3. Verify the file listing displays the correct last modified date
4. Verify the file listing shows "file" as the type for regular files
5. Verify directories show "(N/A)" for size and modified date
6. Verify directory names end with "/" in the listing

**Expected**:
- File metadata is correctly displayed in the listing

---

### Test 11.2: File Sorting

**File**: `tests/opfs-core/opfs-file-sorting.spec.ts`

**Steps**:
1. Create files with names in different orders (`b.txt`, `a.txt`, `c.txt`)
2. Verify the directory listing sorts files alphabetically (case-insensitive, numeric)
3. Verify files and directories are mixed correctly in the sorted listing

**Expected**:
- Files are sorted alphabetically with numeric sorting
- Files and directories are intermixed in the listing

---

### Test 11.3: Empty Directory

**File**: `tests/opfs-core/opfs-empty-directory.spec.ts`

**Steps**:
1. Create a new empty directory
2. Navigate into the empty directory
3. Verify the directory listing is empty (no files or directories)
4. Verify the "Select all" checkbox is unchecked
5. Verify the "Delete" and "DuckDB" buttons are disabled

**Expected**:
- Empty directory listing shows no items
- Action buttons are disabled when no items are selected

---

### Test 11.4: Large Number of Files

**File**: `tests/opfs-core/opfs-many-files.spec.ts`

**Steps**:
1. Create 50+ files in OPFS using a script or repeated uploads
2. Verify the directory listing displays all files
3. Verify the "Select all" checkbox works with many files
4. Verify the "Delete" button is enabled with many selected files
5. Delete a batch of files and verify they are all removed

**Expected**:
- Large number of files are displayed correctly
- Selection and deletion work with many files

---

### Test 11.5: Special Characters in Filenames

**File**: `tests/opfs-core/opfs-special-characters.spec.ts`

**Steps**:
1. Create a file with spaces in the name (`my file.txt`)
2. Verify the file appears correctly in the listing
3. Upload a file with hyphens in the name (`my-file.txt`)
4. Verify the file appears correctly in the listing
5. Create a file with dots in the name (`my.file.txt`)
6. Verify the file appears correctly in the listing
7. Navigate into a directory with spaces in its name
8. Verify the breadcrumb correctly displays the directory with spaces

**Expected**:
- Special characters in filenames are handled correctly
- Breadcrumb trail correctly displays directories with special characters

---

### Test 11.6: File Handle Locked by DuckDB

**File**: `tests/opfs-core/opfs-file-locked.spec.ts`

**Steps**:
1. Create a CSV file in OPFS
2. Open the CSV file in DuckDB shell (creates a view)
3. Try to edit the file using the OPFS editor
4. Verify an alert appears: `FILE MAY BE LOCKED...`
5. Try to delete the file while it is open in DuckDB
6. Verify an error is shown for the locked file
7. Close the DuckDB tab
8. After closing DuckDB, verify the file is editable again

**Expected**:
- Locked files show appropriate error messages
- Files become editable after closing the DuckDB tab

---

## Suite 12: toopfs Filter - Download Files from NVGD to OPFS

### Test 12.1: Download Files from NVGD to OPFS

**File**: `tests/opfs-core/toopfs-filter.spec.ts`

**Steps**:
1. Navigate to an NVGD page with a `file://` source and `?toopfs` filter
2. Verify the "Download to OPFS" UI is displayed
3. Verify all files are pre-selected by default
4. Verify the file count and total size are displayed correctly
5. Unselect some files
6. Verify the file count and total size update
7. Verify the "Download" button is disabled when no files are selected
8. Select files and click "Download"
9. Verify a confirmation dialog appears
10. Confirm the download
11. Verify the progress overlay appears with download steps
12. After download completes, verify the OPFS destination directory is opened in a new tab
13. Verify the downloaded files are visible in the OPFS destination directory

**Expected**:
- toopfs filter UI is displayed correctly
- File selection, count, and size are accurate
- Download process works end-to-end

---

### Test 12.2: Nested Directory Structure

**File**: `tests/opfs-core/toopfs-nested-dir.spec.ts`

**Steps**:
1. Navigate to an NVGD page listing files in a nested directory structure
2. Select files from a nested directory listing
3. Set a destination directory path in the "OPFS directory for downloads" field
4. Download the selected files
5. Verify the files are saved in the correct subdirectory of OPFS
6. Verify the directory structure is preserved in OPFS

**Expected**:
- Nested directory structure is preserved in OPFS
- Files are saved in the correct subdirectory

---

### Test 12.3: Select All/Unselect All

**File**: `tests/opfs-core/toopfs-select-all.spec.ts`

**Steps**:
1. Navigate to an NVGD page with `?toopfs` filter
2. Click the "Select/unselect all" checkbox
3. Verify all files are selected
4. Click the "Select/unselect all" checkbox again
5. Verify all files are unselected
6. Verify the "Download" button is disabled when nothing is selected

**Expected**:
- Select/unselect all works correctly
- Download button state reflects selection state

---

### Test 12.4: Download Progress Display

**File**: `tests/opfs-core/toopfs-progress.spec.ts`

**Steps**:
1. Navigate to an NVGD page with multiple files and `?toopfs` filter
2. Select files and click Download
3. Verify the progress bar appears and updates
4. Verify the download message shows step numbers (e.g., `#1/2 downloading...`)
5. After completion, verify the message shows `completed.`
6. Verify the progress bar is removed after completion

**Expected**:
- Progress bar updates during download
- Step numbers are displayed correctly
- Completion message is shown

---

### Test 12.5: Clear Destination Directory

**File**: `tests/opfs-core/toopfs-clear-destdir.spec.ts`

**Steps**:
1. Navigate to an NVGD page with `?toopfs` filter
2. Enter a destination directory path
3. Click the "Clear" button next to the destination field
4. Verify the destination directory field is cleared
5. Verify the field is focused after clearing

**Expected**:
- Clear button resets the destination directory field
- Field is focused after clearing

---

### Test 12.6: Download Single File

**File**: `tests/opfs-core/toopfs-single-file.spec.ts`

**Steps**:
1. Navigate to an NVGD page with a single file and `?toopfs` filter
2. Select the file
3. Set a destination directory
4. Download the file
5. Verify the file appears in the OPFS destination directory
6. Verify the file content matches the source

**Expected**:
- Single file download works correctly
- File content is preserved

---

### Test 12.7: Download with file:// Protocol

**File**: `tests/opfs-core/toopfs-file-protocol.spec.ts`

**Steps**:
1. Navigate to an NVGD page with `file://` source and `?toopfs` filter
2. Select a file
3. Download the file to OPFS
4. Verify the file is downloaded and uploaded to OPFS successfully

**Expected**:
- file:// protocol files can be downloaded to OPFS

---

### Test 12.8: Download with https:// Protocol

**File**: `tests/opfs-core/toopfs-https-protocol.spec.ts`

**Steps**:
1. Navigate to an NVGD page with `https://` source and `?toopfs` filter
2. Select a file
3. Download the file to OPFS
4. Verify the file is downloaded and uploaded to OPFS successfully

**Expected**:
- https:// protocol files can be downloaded to OPFS

---

### Test 12.9: Download Error Handling

**File**: `tests/opfs-core/toopfs-error-handling.spec.ts`

**Steps**:
1. Navigate to an NVGD page with a broken source URL and `?toopfs` filter
2. Select a file and attempt to download
3. Verify an error message appears indicating the download failed
4. Verify the progress overlay is removed after error
5. Verify the files that were successfully downloaded remain in OPFS

**Expected**:
- Download errors are handled gracefully
- Partially downloaded files remain in OPFS

---

### Test 12.10: Download to Nested OPFS Path

**File**: `tests/opfs-core/toopfs-nested-path.spec.ts`

**Steps**:
1. Navigate to an NVGD page with `?toopfs` filter
2. Set a nested destination path (e.g., `output/data/`)
3. Select and download files
4. Verify the files are saved in the nested OPFS path
5. Navigate to the OPFS destination and verify the directory structure

**Expected**:
- Nested OPFS paths are created correctly
- Files are saved in the correct nested path

---

### Test 12.11: Download After Navigation

**File**: `tests/opfs-core/toopfs-after-navigation.spec.ts`

**Steps**:
1. Navigate to an NVGD page with `?toopfs` filter
2. Navigate away from the page (e.g., to `/opfs/`)
3. Navigate back to the original page
4. Verify the toopfs UI is still functional
5. Select and download files as before

**Expected**:
- toopfs UI remains functional after navigation
- Download works after returning to the page

---

### Test 12.12: Download Multiple Files in Sequence

**File**: `tests/opfs-core/toopfs-multiple-sequence.spec.ts`

**Steps**:
1. Navigate to an NVGD page with multiple files and `?toopfs` filter
2. Select all files
3. Click Download
4. Verify each file is downloaded and uploaded sequentially
5. Verify the progress shows the correct step numbers
6. Verify all files appear in the OPFS destination directory

**Expected**:
- Multiple files are downloaded sequentially
- Progress accurately reflects the download steps

---

### Test 12.13: Download with Empty Selection

**File**: `tests/opfs-core/toopfs-empty-selection.spec.ts`

**Steps**:
1. Navigate to an NVGD page with `?toopfs` filter
2. Unselect all files
3. Verify the "Download" button is disabled
4. Attempt to click Download
5. Verify no action is triggered

**Expected**:
- Download is disabled when no files are selected
- Clicking Download with empty selection does nothing

---

### Test 12.14: File Size Display Formatting

**File**: `tests/opfs-core/toopfs-file-size-format.spec.ts`

**Steps**:
1. Navigate to an NVGD page with files of various sizes and `?toopfs` filter
2. Verify the file sizes are displayed with locale formatting (commas as thousand separators)
3. Verify the total size is calculated correctly across all selected files

**Expected**:
- File sizes are formatted with locale separators
- Total size is calculated correctly

---

### Test 12.15: Download and Verify File Integrity

**File**: `tests/opfs-core/toopfs-file-integrity.spec.ts`

**Steps**:
1. Navigate to an NVGD page with a `file://` source containing known content and `?toopfs` filter
2. Download the file to OPFS
3. Open the OPFS file and verify the content matches the source
4. Verify the file is not corrupted during the download/upload process

**Expected**:
- File integrity is preserved during download/upload
- Downloaded content matches the source

---

### Test 12.16: Download with Special Characters in Paths

**File**: `tests/opfs-core/toopfs-special-chars.spec.ts`

**Steps**:
1. Navigate to an NVGD page with files containing special characters in paths and `?toopfs` filter
2. Select and download the files
3. Verify the files are saved correctly with special characters in the OPFS path
4. Navigate to the OPFS destination and verify the files exist

**Expected**:
- Special characters in paths are handled correctly
- Files are saved in the correct paths

---

### Test 12.17: Download After Page Reload

**File**: `tests/opfs-core/toopfs-after-reload.spec.ts`

**Steps**:
1. Navigate to an NVGD page with `?toopfs` filter
2. Set up a destination directory and select files
3. Reload the page
4. Verify the toopfs UI is restored
5. Select and download files again
6. Verify the files are downloaded successfully

**Expected**:
- toopfs UI is restored after page reload
- Download works after reload

---

### Test 12.18: Download with Indeterminate Selection

**File**: `tests/opfs-core/toopfs-indeterminate.spec.ts`

**Steps**:
1. Navigate to an NVGD page with multiple files and `?toopfs` filter
2. Select only some files (not all)
3. Verify the "Select/unselect all" checkbox is indeterminate
4. Verify the file count and total size reflect only the selected files
5. Verify the "Download" button is enabled

**Expected**:
- Indeterminate state is shown for partial selection
- File count and size reflect only selected files

---

### Test 12.19: Download with Zero-Size Files

**File**: `tests/opfs-core/toopfs-zero-size.spec.ts`

**Steps**:
1. Navigate to an NVGD page with a zero-size file and `?toopfs` filter
2. Download the zero-size file
3. Verify the file is saved in OPFS with size 0

**Expected**:
- Zero-size files can be downloaded to OPFS

---

### Test 12.20: Download with Large Total Size

**File**: `tests/opfs-core/toopfs-large-total-size.spec.ts`

**Steps**:
1. Navigate to an NVGD page with multiple large files and `?toopfs` filter
2. Select all files
3. Verify the total size is displayed correctly with locale formatting
4. Download all files
5. Verify the download completes successfully

**Expected**:
- Large total sizes are displayed correctly
- Download completes successfully for large files

---

### Test 12.21: Download with Special Characters in Names

**File**: `tests/opfs-core/toopfs-special-name-chars.spec.ts`

**Steps**:
1. Navigate to an NVGD page with files containing special characters in names and `?toopfs` filter
2. Select and download the files
3. Verify the files are saved correctly in OPFS
4. Verify the file names are preserved in OPFS

**Expected**:
- Special characters in file names are preserved
- Files are saved correctly in OPFS

---

### Test 12.22: Download and Edit

**File**: `tests/opfs-core/toopfs-edit-after-download.spec.ts`

**Steps**:
1. Navigate to an NVGD page with a text file and `?toopfs` filter
2. Download the file to OPFS
3. Open the OPFS file in the editor
4. Modify the content
5. Save the changes
6. Verify the modified content is persisted in OPFS

**Expected**:
- Downloaded files can be edited and saved
- Modified content is persisted in OPFS

---

### Test 12.23: Download Then Delete

**File**: `tests/opfs-core/toopfs-delete-after-download.spec.ts`

**Steps**:
1. Navigate to an NVGD page with files and `?toopfs` filter
2. Download files to OPFS
3. Navigate to the OPFS destination directory
4. Select and delete the downloaded files
5. Verify the files are removed from OPFS

**Expected**:
- Downloaded files can be deleted from OPFS

---

### Test 12.24: Download with Different Destination Directories

**File**: `tests/opfs-core/toopfs-different-destdirs.spec.ts`

**Steps**:
1. Navigate to an NVGD page with files and `?toopfs` filter
2. Download files to the root OPFS directory
3. Navigate away and come back
4. Download the same files to a different destination directory (e.g., `backup/`)
5. Verify both sets of files exist in their respective directories

**Expected**:
- Files can be downloaded to different destination directories
- Both sets of files exist in their respective directories

---

### Test 12.25: Download with Existing Files in Destination

**File**: `tests/opfs-core/toopfs-existing-files.spec.ts`

**Steps**:
1. Create some files in OPFS first
2. Navigate to an NVGD page with `?toopfs` filter
3. Download files that have the same names as existing OPFS files
4. Verify a confirmation dialog appears for each conflicting file
5. After confirming, verify the files are overwritten
6. Verify the overwritten files contain the new content

**Expected**:
- Overwrite confirmation appears for conflicting files
- Overwritten files contain the new content

---

### Test 12.26: Download Progress Cancellation

**File**: `tests/opfs-core/toopfs-cancel-progress.spec.ts`

**Steps**:
1. Navigate to an NVGD page with multiple files and `?toopfs` filter
2. Select files and click Download
3. While the download is in progress, verify the progress overlay is visible
4. Verify the download completes or fails based on the source availability

**Expected**:
- Progress overlay is visible during download
- Download completes or fails gracefully

---

### Test 12.27: Download with No Destination Set

**File**: `tests/opfs-core/toopfs-no-destdir.spec.ts`

**Steps**:
1. Navigate to an NVGD page with files and `?toopfs` filter
2. Select files without setting a destination directory
3. Download the files
4. Verify the files are saved in the root OPFS directory

**Expected**:
- Files are saved in the root OPFS directory when no destination is set

---

### Test 12.28: Download with file:// Protocol keepcompress

**File**: `tests/opfs-core/toopfs-keepcompress.spec.ts`

**Steps**:
1. Navigate to an NVGD page with `file://` source containing compressed data and `?toopfs` filter
2. Download the file to OPFS
3. Verify the compressed file is saved correctly in OPFS

**Expected**:
- Compressed files are saved correctly in OPFS

---

### Test 12.29: Download with https:// Protocol all Parameter

**File**: `tests/opfs-core/toopfs-https-all.spec.ts`

**Steps**:
1. Navigate to an NVGD page with `https://` source and `?toopfs&all` filter
2. Select and download files
3. Verify the files are downloaded correctly with the `?all` parameter

**Expected**:
- Files are downloaded correctly with the `?all` parameter

---

### Test 12.30: Download with file:// Protocol all Parameter

**File**: `tests/opfs-core/toopfs-file-all.spec.ts`

**Steps**:
1. Navigate to an NVGD page with `file://` source and `?toopfs&all` filter
2. Select and download files
3. Verify the files are downloaded correctly with the `?all` parameter

**Expected**:
- Files are downloaded correctly with the `?all` parameter

---

## Suite 13: Browser Compatibility Considerations

### Test 13.1: OPFS Operations in Secure Context

**File**: `tests/opfs-core/opfs-secure-context.spec.ts`

**Steps**:
1. Verify the OPFS page loads in a secure context (HTTPS or localhost)
2. Verify all OPFS operations work correctly in a secure context
3. Verify file upload, download, and editor operations work
4. Verify DuckDB integration works in a secure context

**Expected**:
- All OPFS operations work in a secure context
- OPFS API is available in a secure context

---

### Test 13.2: OPFS Operations in Non-Secure Context (if applicable)

**File**: `tests/opfs-core/opfs-non-secure-context.spec.ts`

**Steps**:
1. If testing in a non-secure context (non-localhost, non-HTTPS)
2. Verify that OPFS operations are blocked or show appropriate errors
3. Verify that file upload, download, and editor operations are disabled

**Expected**:
- OPFS operations are blocked in non-secure contexts
- Appropriate error messages are shown

---

### Test 13.3: OPFS Operations with Different User Agents

**File**: `tests/opfs-core/opfs-user-agents.spec.ts`

**Steps**:
1. Test OPFS operations with different user agents (if browser supports)
2. Verify that OPFS operations work correctly with different user agents
3. Verify that DuckDB integration works with different user agents

**Expected**:
- OPFS operations work with different user agents
- DuckDB integration works with different user agents

---

### Test 13.4: OPFS Operations with Different Screen Sizes

**File**: `tests/opfs-core/opfs-screen-sizes.spec.ts`

**Steps**:
1. Test OPFS operations at different screen sizes (desktop, tablet, mobile)
2. Verify that the OPFS UI is responsive and functional at different screen sizes
3. Verify that all OPFS operations work correctly at different screen sizes

**Expected**:
- OPFS UI is responsive at different screen sizes
- All OPFS operations work at different screen sizes

---

## Suite 14: Edge Cases and Error Handling

### Test 14.1: OPFS Operations with Very Large Files

**File**: `tests/opfs-core/opfs-large-files.spec.ts`

**Steps**:
1. Create a file with content larger than 1MB in OPFS
2. Verify the file appears in the listing with correct size
3. Verify the file can be downloaded via "Save as" action
4. Verify the file can be uploaded via the file picker

**Expected**:
- Large files can be created, listed, and downloaded
- Files larger than 64KiB cannot be edited via the editor

---

### Test 14.2: OPFS Operations with Unicode Filenames

**File**: `tests/opfs-core/opfs-unicode-filenames.spec.ts`

**Steps**:
1. Create a file with a Unicode name (e.g., `日本語.txt`)
2. Verify the file appears correctly in the listing
3. Upload a file with a Unicode name
4. Verify the file appears correctly in the listing

**Expected**:
- Unicode filenames are handled correctly
- Files with Unicode names appear correctly in the listing

---

### Test 14.3: OPFS Operations with Very Long Filenames

**File**: `tests/opfs-core/opfs-long-filenames.spec.ts`

**Steps**:
1. Create a file with a very long name (e.g., 255 characters)
2. Verify the file appears correctly in the listing
3. Verify the file can be downloaded via "Save as" action

**Expected**:
- Very long filenames are handled correctly
- Files with very long names can be downloaded

---

### Test 14.4: OPFS Operations with Special Characters in Paths

**File**: `tests/opfs-core/opfs-special-chars-paths.spec.ts`

**Steps**:
1. Create a directory with a name containing special characters (e.g., `my dir/`)
2. Create a file inside the directory
3. Verify the file appears correctly in the listing
4. Verify the breadcrumb correctly displays the directory with special characters

**Expected**:
- Special characters in directory names are handled correctly
- Breadcrumb trail correctly displays directories with special characters

---

### Test 14.5: OPFS Operations with Concurrent Access

**File**: `tests/opfs-core/opfs-concurrent-access.spec.ts`

**Steps**:
1. Create a file in OPFS
2. Open the file in DuckDB shell
3. Try to edit the file in OPFS
4. Verify an error is shown for the locked file
5. Close the DuckDB tab
6. Verify the file is editable again

**Expected**:
- Concurrent access is handled correctly
- Locked files show appropriate error messages

---

### Test 14.6: OPFS Operations with Network Errors

**File**: `tests/opfs-core/opfs-network-errors.spec.ts`

**Steps**:
1. Navigate to an NVGD page with a broken source URL and `?toopfs` filter
2. Select a file and attempt to download
3. Verify an error message appears indicating the download failed
4. Verify the progress overlay is removed after error

**Expected**:
- Network errors are handled gracefully
- Error messages are displayed appropriately

---

### Test 14.7: OPFS Operations with Timeout

**File**: `tests/opfs-core/opfs-timeout.spec.ts`

**Steps**:
1. Navigate to an NVGD page with a slow source URL and `?toopfs` filter
2. Select a file and attempt to download
3. Verify the download completes or times out
4. Verify the progress overlay is removed after timeout

**Expected**:
- Timeout is handled gracefully
- Progress overlay is removed after timeout

---

### Test 14.8: OPFS Operations with Invalid Input

**File**: `tests/opfs-core/opfs-invalid-input.spec.ts`

**Steps**:
1. Enter an invalid directory name in the "Directory name to create" input
2. Verify an error is shown
3. Enter an invalid file name in the "File name to create or update" field
4. Verify an error is shown
5. Enter an invalid URL in the "URL to be downloaded" field
6. Verify an error alert is shown

**Expected**:
- Invalid input is handled gracefully
- Error messages are displayed appropriately

---

## Summary

This test plan covers:
- **14 suites** with **80+ test cases**
- **Core OPFS operations**: Page loading, directory navigation, file operations, multi-selection
- **DuckDB integration**: Supported and unsupported file types, multiple files, directory recursion
- **Download URL to OPFS**: NVGD resources, error handling, file integrity
- **toopfs filter**: Download files from NVGD URLs to OPFS, progress display, error handling
- **Browser compatibility**: Secure context, different user agents, screen sizes
- **Edge cases and error handling**: Large files, Unicode filenames, concurrent access, network errors

**Total estimated test time**: 4-6 hours (depending on browser setup and server configuration)
