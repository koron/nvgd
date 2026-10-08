# "grep" filter

Output lines which matches against regular expression.

As default, matching is made for whole line.  But when valid option `field` is
given, then matching is made for specified a field, which is splitted by
`delim` character.

`grep` command equivalent.

* filter\_name: `grep`
* options
    * `re` - regular expression used for match.
    * `match` - output when match or not match.  default is true.
    * `only` - output only matching texts. default is false.
      When `only` is enabled, `match=false` doesn't print any lines, and
      `context` is ignored.
    * `field` - a match target N'th field counted from 1.
        default is none (whole line).
    * `delim` - field delimiter string (default: TAB character).
    * `context` - show a few lines before and after the matched line.
        default is `0` (no contexts).
    * `number` - prefix each line of output with the 1-based line number.
        when `true`
