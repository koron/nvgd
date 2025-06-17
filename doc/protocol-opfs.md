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

The NVGD OPFS interface is typically available at the `/opfs/` URL path of your NVGD instance. For example, if your NVGD instance is running at `http://localhost:9280`, you would access the OPFS interface by navigating your browser to `http://localhost:9280/opfs/`.

## OPFS User Interface and Functionalities

The NVGD OPFS interface provides a way to manage files and directories stored within the browser's Origin Private File System.

### 1. Directory Navigation

The UI allows users to navigate the OPFS directory structure.

*   **File and Directory Listing:**
    *   The main view displays a list of files and directories within the current OPFS directory.
    *   Each entry shows its name, type (file/directory), size (for files), and last modified date.
*   **Creating Directories:**
    *   A "Create Directory" button or input field allows users to create new subdirectories within the current directory.
*   **Moving Up and Down the Directory Tree:**
    *   Clicking on a directory name in the list navigates into that directory.
    *   A "Parent Directory" button (often represented as ".." or an up arrow icon) allows navigation to the directory containing the current one.
    *   A breadcrumb trail might also be present, showing the current path and allowing quick navigation to any parent directory in the path.

### 2. File Operations

Various operations can be performed on individual files and directories.

*   **Uploading Local Files:**
    *   An "Upload File(s)" button opens a system file dialog, allowing users to select one or more files from their local computer to be uploaded into the current OPFS directory.
*   **Creating or Editing Files Directly:**
    *   A "Create New File" button might allow users to create a blank text file or a file of a specific type.
    *   An "Edit File" option (often available when selecting a file, especially text-based files) would open a simple in-browser text editor to modify the file's content. Changes are saved back to OPFS.
*   **Removing Files/Directories:**
    *   A "Remove" or "Delete" option (often an icon or a button available next to each file/directory, or after selecting items) allows users to delete selected files or directories.
    *   A confirmation prompt is typically shown before permanent deletion. Deleting a directory will also delete all its contents.
*   **Saving Files from OPFS to Local Disk (Download):**
    *   A "Download" option (often an icon or a button available next to each file, or after selecting a file) allows users to save a copy of the selected file from OPFS to their computer's local file system.

### 3. Actions for Multiple Selected Files

The UI supports operations on multiple selected items.

*   **Selection:** Users can typically select multiple files and/or directories using checkboxes next to each item or by using standard keyboard shortcuts (e.g., Ctrl+Click, Shift+Click).
*   **"Open multiple files with DuckDB":**
    *   When multiple files (e.g., Parquet, CSV) are selected, an action button like "Open selected with DuckDB" or "Load into DuckDB" becomes available.
    *   Clicking this button will typically instruct the integrated DuckDB WASM instance to create views or tables based on the selected files, allowing users to query them collectively. The exact mechanism might involve creating separate views for each file or attempting to union compatible files.
*   **Other Batch Actions:** Depending on the implementation, other batch actions like "Download Selected" or "Delete Selected" might be available.

### 4. "Delete all contents" Button

*   **Functionality:** A prominent button, often labeled "Delete all contents in this directory," "Empty Directory," or similar, is provided.
*   **Purpose:** This allows for the quick removal of all files and subdirectories within the currently viewed OPFS directory.
*   **Confirmation:** Due to its destructive nature, clicking this button will always trigger a confirmation dialog asking the user to verify they indeed want to delete all content in the current directory. This helps prevent accidental data loss.

This UI aims to provide a familiar file explorer-like experience for managing data within the browser's OPFS, integrated with the data processing capabilities of tools like DuckDB.

## DuckDB Integration

NVGD's OPFS interface offers seamless integration with DuckDB, allowing you to directly query files stored in OPFS using DuckDB's powerful SQL capabilities running in your browser via WebAssembly (WASM).

### Supported File Types

The DuckDB integration primarily supports the following file types for direct querying from OPFS:

*   **CSV (Comma Separated Values):** `.csv` files.
*   **XLSX (Microsoft Excel Open XML Format Spreadsheet):** `.xlsx` files. DuckDB can read data from sheets within these files.
*   **JSON (JavaScript Object Notation):** `.json` files. This includes standard JSON files and newline-delimited JSON (NDJSON).
*   **Parquet:** `.parquet` files. Parquet is a columnar storage format optimized for analytics.

### Opening Single Files with DuckDB

When you choose to open a single supported file from the OPFS interface with DuckDB (e.g., via a context menu option "Open with DuckDB" or a dedicated button):

*   NVGD will instruct the in-browser DuckDB instance to create a view for that file.
*   This view is typically named `opfs`.
*   You can then query this view directly in the DuckDB SQL shell. For example, if you opened `my_data.csv`:
    ```sql
    SELECT * FROM opfs;
    SELECT column_name, COUNT(*) FROM opfs GROUP BY column_name;
    ```

### Opening Multiple Files with DuckDB

The interface also allows you to select and open multiple supported files with DuckDB simultaneously:

*   When multiple files are selected and opened with DuckDB:
    *   NVGD will instruct DuckDB to create a separate view for each selected file.
    *   These views are typically named sequentially, such as `opfs0`, `opfs1`, `opfs2`, and so on. The numbering corresponds to the order in which the files were selected or processed.
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
