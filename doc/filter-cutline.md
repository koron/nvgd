# "cutline" filter

This cutline is a filter for NVGD that extracts a specific range of lines
from a text stream. The start and end of the range you want to extract can
be specified using regular expression patterns, allowing for flexible data
extraction. For example, this is useful for extracting only specific error
sections from a log file, or extracting parts of a configuration file.

## Parameters

Name    |Type   |Requirements |Description
--------|-------|-------------|------------
`start` |Regexp |Optional     |A regexp pattern that matches the first line in the range
`end`   |Regexp |Optional     |A regexp pattern that matches the end line in the range

## Description

This cutline filter displays the range from the line that matches `start` to
the line that matches `end`.  The line matching `end` are included to display.
If neither is specified, all lines are displayed.

If only `start` is specified, all lines after the matching line are displayed.
If only `end` is specified, all lines from the beginning to the matching line
are displayed.
