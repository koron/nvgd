# Origin Private File System (OPFS) Protocol

## Introduction

The Origin Private File System (OPFS) is a storage space private to the browser tab/window and origin. It allows web applications to store and manage files directly on the user's device, offering persistent storage that is not cleared when the browser is closed (unless explicitly deleted by the user or the application).

In the context of this project (NVGD), OPFS serves several key purposes:

*   **Persistent Client-Side Storage:** Enables users to save data, query results, or entire database states generated within browser-based tools (like the DuckDB WASM shell) directly to their local machine. This data persists across browser sessions.
*   **Local Data Workspace:** Provides a sandboxed file system within the browser, allowing users to manage project files, upload local data for analysis, and save outputs without needing to constantly transfer files to/from a remote server or their main local file system.
*   **Improved Performance for Large Files:** For applications dealing with large datasets, OPFS can offer faster access and processing compared to repeatedly fetching data from remote sources, as the data is stored locally.
*   **Offline Capabilities (Future Potential):** While not fully implemented in all browsers or applications, OPFS lays the groundwork for web applications to function with local data even when offline.

This document describes the user interface (UI) provided by NVGD for interacting with OPFS.

## Accessing the OPFS Interface

The NVGD OPFS interface is available at the `/opfs/` URL path of your NVGD instance. For example, if your NVGD instance is running at `http://localhost:9280`, you would access the OPFS interface by navigating your browser to `http://localhost:9280/opfs/`.

## OPFS User Interface and Functionalities

The NVGD OPFS interface provides a way to manage files and directories stored within the browser's Origin Private File System.

### 1. Directory Navigation

The UI allows users to navigate the OPFS directory structure.

*   **File and Directory Listing:**
    *   The main view displays a list of files and directories within the current OPFS directory.
    *   Each entry shows its name and type (file/directory). For files, it also displays size and last modified date. For directories, these are shown as "(N/A)".
*   **Creating Directories:**
    *   A "Create Directory" button or input field allows users to create new subdirectories within the current directory.
*   **Moving Up and Down the Directory Tree:**
    *   Clicking on a directory name in the list navigates into that directory.
    *   A breadcrumb trail is also present, showing the current path and allowing quick navigation to any parent directory in the path.
*   **Reload Directory List:**
    *   Click the "Reload" button to refresh the directory list. This is useful for reflecting any changes to the OPFS contents.

### 2. File Operations

Various operations can be performed on individual files and directories.


*   **Uploading a Local File:**
    *   In the Upload a local file section, you can select a file in the file picker, change the name of the file to upload if necessary, and then click the Upload button to upload the local file to the current directory in OPFS.
*   **Editor:**
    *   You can create new files or update existing files.
    *   Specify a file name, enter file contents, and press the "Create or update a file" button to create or update a file with that file name and contents.
    *   If the file is less than 64KiB, clicking the "Edit" action (represented by an icon) in the list will load the file's name and contents into the Editor input fields, making it easier to update existing files.
*   **Saving Files from OPFS to Local Disk (Download):**
    *   A "Save as" option (represented by a download icon) next to each file allows users to save a copy of the selected file from OPFS to their computer's local file system.

### 3. Actions for Multiple Selected Files

The UI supports operations on multiple selected items.

*   **Selection:** Users can select multiple files and/or directories using checkboxes next to each item.
*   **Delete selected Files/Directories:**
    *   When you check the checkboxes for one or more files or directories, the "Delete" action button becomes available.
    *   This button allows users to delete the checked files or directories.
    *   A confirmation prompt is shown before permanent deletion.
    *   Deleting a directory will also delete all its contents.
*   **"Open selected files with DuckDB":**
    *   When you check the checkboxes for one or more files or directories, the "DuckDB" action button becomes available.
    *   Clicking this button will instruct the integrated DuckDB WASM instance to create separate views based on the selected files (or files within selected directories), allowing users to query them collectively.
    *   Selecting a directory will cause DuckDB to open all files contained within that directory, recursively.
    *   See "DuckDB Integration" below for more information.

## DuckDB Integration

NVGD's OPFS interface offers seamless integration with DuckDB, allowing you to directly query files stored in OPFS using DuckDB's powerful SQL capabilities running in your browser via WebAssembly (WASM).

### Supported File Types

The DuckDB integration primarily supports the following file types for direct querying from OPFS:

*   **CSV (Comma Separated Values):** `.csv`, `.csv.gz`, and `.csv.zst` files.
*   **TSV (Tab Separated Values):** `.tsv`, `.tsv.gz`, and `.tsv.zst` files.
*   **XLSX (Microsoft Excel Open XML Format Spreadsheet):** `.xlsx` files. DuckDB can read data from sheets within these files.
*   **JSON (JavaScript Object Notation):** `.json` files. This includes standard JSON files and newline-delimited JSON (NDJSON).
*   **Parquet:** `.parquet` files. Parquet is a columnar storage format optimized for analytics.

For files of these file types, the OPFS interface will automatically create views when opening them using DuckDB.

For information on DuckDB's supported file formats and their specifications, please refer to: <https://duckdb.org/docs/stable/guides/file_formats/overview>

### Opening (Multiple) Files with DuckDB

The interface allows you to select and open multiple supported files with DuckDB simultaneously:

*   When multiple files are selected and opened with DuckDB:
    *   NVGD will instruct DuckDB to create a separate view for each selected file of supported file types.
    *   These views are named sequentially, such as `opfs0`, `opfs1`, `opfs2`, and so on. The numbering corresponds to the order in which the files were selected or processed.
*   You can then query these individual views or join them in your SQL queries. For example, if you opened `data_part1.parquet` and `data_part2.parquet`:

    ```sql
    -- Query the first file
    SELECT * FROM opfs0;

    -- Query the second file
    SELECT * FROM opfs1;

    -- Example of joining data from both files (assuming compatible schemas)
    SELECT a.*, b.extra_info
    FROM opfs0 AS a
    JOIN opfs1 AS b ON a.id = b.id;
    ```

This integration allows for flexible data analysis by bringing the power of DuckDB directly to the files you manage within OPFS.
