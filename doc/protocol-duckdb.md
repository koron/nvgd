**NVGD DuckDB Protocol Documentation**

**1. Introduction**

The NVGD DuckDB protocol integrates the power of DuckDB directly into your web browser. Its key characteristic is **client-side execution**: the DuckDB database engine runs entirely within your browser using WebAssembly (WASM). This means that when you interact with data through this protocol, all query processing and data manipulation happen locally on your machine, not on the NVGD server. NVGD's role is primarily to serve the static assets (HTML, JavaScript) that make this in-browser database environment possible.

There are two main ways to use the NVGD DuckDB protocol:
*   **Interactive Shell:** A command-line like interface in your browser for running SQL queries.
*   **'Show as View' Utility:** A quick way to load data from a URL directly into the interactive shell for inspection and querying.

**2. Accessing the Interactive Shell**

The NVGD DuckDB protocol provides a powerful interactive SQL shell that runs entirely within your web browser, powered by DuckDB-WASM. This allows you to directly query data sources accessible to your NVGD instance and perform complex data analysis without needing to install DuckDB locally.

To access the interactive shell:
1.  Ensure your NVGD instance is running.
2.  Open your web browser and navigate to the `/duckdb/` path on your NVGD instance. For example, if NVGD is running on `http://localhost:8080`, you would go to `http://localhost:8080/duckdb/`.

Upon loading, you will be presented with an interface resembling a command-line terminal. This is the DuckDB SQL prompt. You can type any valid DuckDB SQL query directly into this prompt and press `Enter` to execute it. The results of your query will be displayed within the same interface.

For example, you can create tables, load data from URLs, and run analytical queries:

```sql
-- Create a table from a CSV file accessible via a URL
CREATE TABLE items AS SELECT * FROM 'https://nvgd-data.s3.amazonaws.com/items.csv';

-- Perform an aggregation
SELECT type, AVG(price) AS average_price FROM items GROUP BY type;
```

It's important to remember that this DuckDB instance is running entirely in your browser. The database state, including any tables you create or data you load, is **session-specific**. This means it will be lost if you close the browser tab or refresh the page, unless you explicitly export the data (e.g., using `EXPORT DATABASE` commands if needed). All processing and data handling occur client-side.

**3. Using the 'Show as View' Feature**

The 'Show as View' feature is a convenient utility within the NVGD DuckDB protocol that allows you to quickly load data from an external source directly into the interactive DuckDB shell. This is particularly useful for inspecting files or data streams accessible via URLs without manually typing `CREATE VIEW` statements.

**URL Structure:**

To use this feature, you construct a URL in the following format:

`http://<your-nvgd-host>/duckdb/show-as-view.html?t=<URL_TO_DATA_SOURCE>`

Let's break this down:
*   `<your-nvgd-host>`: This is the hostname and port where your NVGD instance is running (e.g., `localhost:8080`).
*   `/duckdb/show-as-view.html`: This is the fixed path to the utility page.
*   `?t=<URL_TO_DATA_SOURCE>`: This is the crucial part. The `t` query parameter takes a URL as its value. This URL should point directly to the raw data you want to load.

**The `t` Parameter:**

The value of the `t` parameter (`<URL_TO_DATA_SOURCE>`) can be:
*   A path to another resource served by your NVGD instance (e.g., `/http/example.com/mydata.csv`, `/s3/my-bucket/data.parquet`, `/file/path/to/local/file.csv`).
*   A public URL to a raw data file (e.g., `https://example.com/data.csv`, `https://some-data-api.com/data.json`). DuckDB's WASM environment will attempt to fetch and process this URL. Ensure the URL points to the raw data, not an HTML page displaying the data.

**How it Works and Querying the Data:**

When you navigate to the 'Show as View' URL:
1.  NVGD serves the `show-as-view.html` page.
2.  A script on this page takes the URL from the `t` parameter.
3.  It then redirects your browser to the main DuckDB interactive shell (`/duckdb/`).
4.  As part of this redirection, it instructs the shell to automatically execute two commands:
    *   `CREATE VIEW t AS SELECT * FROM '<URL_TO_DATA_SOURCE>';`
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

**4. Examples**

Here are a few examples to illustrate how to use the NVGD DuckDB protocol. For these examples, assume your NVGD instance is running on `http://localhost:8080`.

**Example 1: Accessing the Interactive Shell**

1.  **Open the shell:** Navigate your browser to `http://localhost:8080/duckdb/`.
2.  **Run a query:** Once the shell prompt appears, type the following SQL query and press Enter:

    ```sql
    SELECT 'Hello DuckDB from NVGD!' AS message;
    ```

    You should see the result displayed in the shell:

    ```
    ┌───────────────────────────┐
    │          message          │
    │          varchar          │
    ├───────────────────────────┤
    │ Hello DuckDB from NVGD! │
    └───────────────────────────┘
    ```

**Example 2: Using 'Show as View' with a Local File via NVGD's `file:` Protocol**

Let's say you have a CSV file named `my_local_data.csv` located at `/shared/data/my_local_data.csv` on the same machine where NVGD is running, and you have NVGD's `file:` protocol configured such that `/file/data/my_local_data.csv` maps to this path.

1.  **Construct the 'Show as View' URL:**
    To load this CSV file into the DuckDB shell, you would use the following URL in your browser:

    `http://localhost:8080/duckdb/show-as-view.html?t=/file/data/my_local_data.csv`

    *(Note: The exact path after `/file/` depends on your NVGD `file:` protocol configuration and the base path it's serving from.)*

2.  **Query the data:**
    After the browser redirects to the interactive shell, the contents of `my_local_data.csv` will be available in a view named `t`. You can then query it. For instance, to count the number of rows:

    ```sql
    SELECT COUNT(*) AS total_rows FROM t;
    ```

    Or, if your CSV has a column named `status`:

    ```sql
    SELECT status, COUNT(*) FROM t GROUP BY status;
    ```

**Example 3: Using 'Show as View' with a Public URL**

You can also use 'Show as View' to directly load data from public URLs that host raw data files. For this example, we'll use a sample CSV file for the Iris dataset.

1.  **Construct the 'Show as View' URL:**
    The URL for the raw CSV data is: `https://gist.githubusercontent.com/curran/a08a1080b88344b0c8a7/raw/0e7a9b0a5d22642a06d3d5b9bcbad9890c8ee534/iris.csv`

    The URL to load this into NVGD's DuckDB shell would be:

    `http://localhost:8080/duckdb/show-as-view.html?t=https://gist.githubusercontent.com/curran/a08a1080b88344b0c8a7/raw/0e7a9b0a5d22642a06d3d5b9bcbad9890c8ee534/iris.csv`

2.  **Query the data:**
    Once the shell loads and the data is available in view `t`, you can query it. For example, to see the distinct species in the dataset:

    ```sql
    SELECT DISTINCT species FROM t;
    ```

    To find the average sepal length for each species:

    ```sql
    SELECT species, AVG(sepal_length) AS avg_sepal_length
    FROM t
    GROUP BY species;
    ```

These examples demonstrate the flexibility of the NVGD DuckDB protocol for both direct interaction and quick loading of various data sources directly into your browser.
