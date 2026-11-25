---
name: nvgd-reporter
description: Query an nvgd HTTP endpoint to discover available sources, filter and aggregate data server-side, and generate structured reports with minimal context overhead.

---

# nvgd Reporter Skill

`nvgd` is a stateless HTTP data-retrieval and pipeline-processing server (`/{protocol}://{target_args}?{filters}`). Use `curl` to query `nvgd`, narrow or aggregate data on the server side to conserve context tokens, and synthesize accurate reports.

## 1. Discovery & Environment

- Base URL: `http://127.0.0.1:9280` (override with `NVGD_URL` if configured differently).
- Check active instance configuration and available capabilities first when source names, allowed paths, DBs, or command keys are unknown:

    ```bash
    BASE="${NVGD_URL:-http://127.0.0.1:9280}"
    curl -fsS "$BASE/config/"
    # /help/ is a 303 redirect to /help:/, so follow it with -L.
    curl -fsSL "$BASE/help/?markdown" | head -n 120
    ```

* For detailed filter or protocol options beyond this skill, consult `/help/` or the corresponding `doc/*.md` reference rather than guessing flag names.

## 2. URL Syntax & Encoding Rules

Requests follow this form:

```text
{NVGD_URL}/{protocol}://{target_args}?{filter1}&{filter2}...
```

* **URL-encode dynamic arguments**: Always encode query values, regex patterns, quotes, spaces, and reserved characters (`=`, `&`, `?`, `#`, `%`). `curl -G --data-urlencode` URL-encodes only the *query string*; values that must live in the URL *path* (e.g., the SQL of a `db://` request) have to be encoded by hand or via a helper, since `-G` will not place them in the path.
* **Quote shell arguments**: Always double-quote URLs and arguments in shell commands to prevent globbing or word splitting.

## 3. Core Protocols & Exploration Order

Start broad enough to identify targets, then narrow to exact records or aggregates:

* **`file://`**: Read files permitted by the instance's location/forbidden rules (e.g., logs, LTSV, CSV, text). Inspect recent lines or list/locate targets before pulling full files.
* **`command://`**: Run pre-configured server commands by registered key. Check `/config/` for available command keys first.
* **`s3list://` / `s3obj://`**: List prefixes/keys (`s3list://`) to locate relevant objects before fetching specific object contents (`s3obj://`).
* **`db://`**: Query configured SQL databases via their `database/sql` drivers. The configured DB name is the `db://` "host" and the (URL-encoded) SQL goes in the path. Discover available DB names from `/config/`, inspect schemas/tables or sample rows (`LIMIT`), and run aggregations server-side. (`trdsql` is a separate filter for querying tabular files, not a `db://` driver.)

## 4. Filtering, Aggregation & Token Conservation

Treat the `nvgd` server as the primary filter and aggregator so large raw payloads do not flood the context window:

1. **Probe before pulling large slices**: Check recent entries (`tail`), line/record counts (`count`), or small samples before retrieving broad ranges.
2. **Filter early**: Narrow down by pattern (`grep`), head/tail windows, or time/key conditions as early as possible in the pipeline.
3. **Aggregate server-side**: Prefer SQL (`GROUP BY`, `COUNT`, `SUM`, `MIN`/`MAX`) or `trdsql` aggregations over fetching thousands of raw rows and counting in the prompt.
4. **Convert to compact structured formats**: Use format filters such as `ltsvconv` or `jsonarray` / table formatters when helpful so tabular data is easy to parse cleanly. Note NVGD's LTSV dialect is `label:value` (tab-separated fields, colon-separated key/value); `ltsvconv` and directory listings follow this, so standard `key=value` LTSV files will be misparsed.
5. **Respect pipeline order and caps**: Chain filters intentionally in the query string, and tighten filters or narrow time windows if output volume is still large.

## 5. Report Generation Workflow

1. **Clarify scope**: Identify the target metric or topic, time range, data source(s), and desired report structure.
2. **Inspect instance & schema**: Check `/config/` and discover relevant files, S3 prefixes, command keys, or DB tables/columns.
3. **Run exploratory & targeted queries**: Verify field names and volumes first, then execute focused queries for counts, distributions, anomalies, and representative examples.
4. **Verify & write the report**:
* Summarize key findings clearly with structured sections or tables.
* Cite the data sources and key filter/query conditions used.
* Explicitly note any truncation, row limits (`max_rows`), sampling, or time-range boundaries that affect completeness.
* Do not copy sensitive secrets, credentials, or unnecessary PII into report output.


## 6. Reference Examples

Filter parameters use colon syntax: each `filter=` pair's value is itself a
`key:value` list, e.g. `tail=limit:200`, `grep=re:ERROR|WARN`. Omitting the key
(`tail=200`, `grep=ERROR|WARN`) is silently ignored and falls back to defaults.

```bash
BASE="${NVGD_URL:-http://127.0.0.1:9280}"

# Tail recent error/warning lines from a file.
curl -fsS -G "$BASE/file:///var/log/app/app.log" \
  --data-urlencode "tail=limit:200" \
  --data-urlencode "grep=re:ERROR|WARN"

# Convert LTSV records to CSV for compact inspection.
curl -fsS -G "$BASE/file:///var/log/app/access.ltsv" \
  --data-urlencode "tail=limit:300" \
  --data-urlencode "ltsvconv=format:csv"

# Aggregate counts server-side via DB query.
# The DB name is the db:// "host"; the SQL is the (URL-encoded) path. A
# `?query=` parameter would be treated as an unknown filter (500 error), and
# `-G --data-urlencode` cannot place the SQL in the path.
curl -fsS "$BASE/db://metrics/select%20status,%20COUNT(*)%20AS%20cnt%20FROM%20events%20GROUP%20BY%20status%20ORDER%20BY%20cnt%20DESC%20LIMIT%2020"
```
