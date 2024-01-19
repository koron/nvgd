# Testcases for issue:88 (serve root contents)

*   Disable `root_contents_file` (empty)

    Check it doens't work

*   Enable `root_contents_file` with invalid path)

    Check it doens't work

*   Enable `root_contents_file` with valid path)
    *   GET: type, size, modtime, body
    *   HEAD: type, size, modtime
