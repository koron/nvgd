# NVGD DuckDB Protocol Documentation

## Introduction

The NVGD DuckDB protocol integrates the power of DuckDB directly into your web browser. Its key characteristic is **client-side execution**: the DuckDB database engine runs entirely within your browser using WebAssembly (WASM). This means that when you interact with data through this protocol, all query processing and data manipulation happen locally on your machine, not on the NVGD server. NVGD's role is primarily to serve the static assets (HTML, JavaScript) that make this in-browser database environment possible.

There are two main ways to use the NVGD DuckDB protocol:

*   **Interactive Shell:** A command-line like interface in your browser for running SQL queries.
*   **'Show as View' Utility:** A quick way to load data from a link, such as the `indexhtml` filter, directly into the interactive shell for inspection and querying.

For more information about DuckDB WASM, see: <https://duckdb.org/docs/stable/clients/wasm/overview.html>

## Accessing the Interactive Shell

The NVGD DuckDB protocol provides a powerful interactive SQL shell that runs entirely within your web browser, powered by DuckDB-WASM. This allows you to directly query data sources accessible to your NVGD instance and perform complex data analysis without needing to install DuckDB locally.

To access the interactive shell:
1.  Ensure your NVGD instance is running.
2.  Open your web browser and navigate to the `/duckdb/` path on your NVGD instance. For example, if NVGD is running on `http://localhost:9280`, you would go to `http://localhost:9280/duckdb/`.

Upon loading, you will be presented with an interface resembling a command-line terminal. This is the DuckDB SQL prompt. You can type any valid DuckDB SQL query directly into this prompt and press `Enter` to execute it. The results of your query will be displayed within the same interface.

For example, you can create tables, load data from URLs, and run analytical queries:

```sql
-- Create a table from a CSV file accessible via a URL
CREATE VIEW items AS SELECT * FROM 'https://shell.duckdb.org/data/tpch/0_01/parquet/lineitem.parquet';

-- Perform an aggregation
SELECT l_shipmode, AVG(l_extendedprice) FROM items GROUP BY l_shipmode;
```

It's important to remember that this DuckDB instance is running entirely in your browser. The database state, including any tables you create or data you load, is **session-specific**. This means it will be lost if you close the browser tab or refresh the page, unless you explicitly export the data (e.g., using `EXPORT DATABASE` commands if needed). All processing and data handling occur client-side.

## Using the 'Show as View' Feature

The 'Show as View' feature is a convenient utility within the NVGD DuckDB protocol that allows you to quickly load data from an external source (especially, the `indexhtml` filter) directly into the interactive DuckDB shell. This is particularly useful for inspecting files or data streams accessible via links without manually typing `CREATE VIEW` statements.

This feature is specifically intended for use by the `indexhtml` filter, which uses this feature to provide links to Parquet, CSV and Excel files that can be opened in DuckDB.

**URL Structure:**

To use this feature, you construct a URL in the following format:

`http://<your-nvgd-host>/duckdb/show-as-view.html?t=<PATH_ON_NVGD>`

Let's break this down:

*   `<your-nvgd-host>`: This is the hostname and port where your NVGD instance is running (e.g., `localhost:9280`).
*   `/duckdb/show-as-view.html`: This is the fixed path to the utility page.
*   `?t=<PATH_ON_NVGD>`: This is the crucial part. The `t` query parameter takes a path on NVGD as its value. This path should point directly to the raw data you want to load.

**The `t` Parameter:**

The value of the `t` parameter (`<PATH_ON_NVGD>`) can be:

*   A path to another resource served by your NVGD instance (e.g., `/s3obj://my-bucket/data.parquet`, `/file:///path/to/local/file.csv`).
<!--*   A public URL to a raw data file (e.g., `https://example.com/data.csv`, `https://some-data-api.com/data.json`). DuckDB's WASM environment will attempt to fetch and process this URL. Ensure the URL points to the raw data, not an HTML page displaying the data.-->

**How it Works and Querying the Data:**

When you navigate to the 'Show as View' URL:

1.  NVGD serves the `show-as-view.html` page.
2.  A script on this page takes the URL from the `t` parameter.
3.  It then redirects your browser to the main DuckDB interactive shell (`/duckdb/`).
4.  As part of this redirection, it instructs the shell to automatically execute two commands:
    *   `CREATE VIEW t AS SELECT * FROM '<NVGD_ORIGIN>/<PATH_ON_NVGD>';`
    *   `SHOW t;`

This means the data from your specified URL is automatically loaded into a view named `t` within the DuckDB shell, and the structure of this view (columns and data types) is displayed.

Once loaded, you can immediately start querying this data using standard SQL commands against the view `t`. For example:

```sql
-- Assuming the data source has columns 'product_name' and 'price'
SELECT * FROM t WHERE price > 100;

SELECT product_name, COUNT(*) AS item_count
FROM t
GROUP BY product_name
ORDER BY item_count DESC;
```

This feature provides a quick and easy way to begin exploring datasets accessible via URLs without needing to manually type the `CREATE VIEW` commands in the shell. Remember, like the rest of the interactive shell, this view `t` exists only for your current browser session.

## Examples

Here are a few examples to illustrate how to use the NVGD DuckDB protocol. For these examples, assume your NVGD instance is running on `http://localhost:9280`.

### Example 1: Accessing the Interactive Shell

1.  **Open the shell:** Navigate your browser to `http://localhost:9280/duckdb/`.
2.  **Run a query:** Once the shell prompt appears, type the following SQL query and press Enter:

    ```sql
    SELECT 'Hello DuckDB from NVGD!' AS message;
    ```

    You should see the result displayed in the shell:

    ```
    ┌─────────────────────────┐
    │ message                 │
    ╞═════════════════════════╡
    │ Hello DuckDB from NVGD! │
    └─────────────────────────┘
    ```

### Example 2: Using 'Show as View' with a Local File via NVGD's `file` Protocol

Let's say you have a CSV file named `my_local_data.csv` located at `/shared/data/my_local_data.csv` on the same machine where NVGD is running, and you have NVGD's `file:` protocol configured such that `/file:///shared/data/my_local_data.csv` maps to this path.

1.  **Construct the 'Show as View' URL:**
    To load this CSV file into the DuckDB shell, you would use the following URL in your browser:

    `http://localhost:9280/duckdb/show-as-view.html?t=/file:///shared/data/my_local_data.csv`

2.  **Query the data:**
    After the browser redirects to the interactive shell, the contents of `my_local_data.csv` will be available in a view named `t`. You can then query it. For instance, to count the number of rows:

    ```sql
    SELECT COUNT(*) AS total_rows FROM t;
    ```

    Or, if your CSV has a column named `status`:

    ```sql
    SELECT status, COUNT(*) FROM t GROUP BY status;
    ```

## APPENDIX

### HTTP range requests

DuckDB supports HTTP range requests, so queries against Parquet files can be performed by loading only a portion of the file, depending on its content. For more information, see <https://duckdb.org/docs/stable/core_extensions/httpfs/https#partial-reading>

NVGD supports HTTP range requests over the `file` and `s3obj` protocols.

### How to use OPFS

As mentioned above, operations performed with NVDG's DuckDB shell are generally lost when the browser tab is closed. However, they can be made persistent by using a local file system called OPFS. For more information about OPFS, see the following URL: <https://developer.mozilla.org/en-US/docs/Web/API/File_System_API/Origin_private_file_system>

To save the table `t` in Parquet format on OPFS, execute the following.

```
duckdb> .files register opfs://output.parquet
Registering OPFS file handle for: opfs://output.parquet

duckdb> COPY (SELECT * FROM t) TO 'opfs://output.parquet';
┌───────┐
│ Count │
╞═══════╡
│ 60175 │
└───────┘
```

**WARINNG** If you omit or type the wrong `opfs://` when doing `.files register`, the DuckDB shell will hang. You will then need to reload your browser.

NVGD does not have the ability to browse OPFS currently, but there are plans to add this in the future. In the meantime, you can use the [Chrome extension: OPFS Explorer](https://chromewebstore.google.com/detail/opfs-explorer/acndjpgkpaclldomagafnognkcgjignd) to manage OPFS.
