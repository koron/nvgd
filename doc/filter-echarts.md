# "echarts" filter

echarts filter provides chart rendering.
Which uses [Apache Echarts][echarts].

[echarts]:https://echarts.apache.org/en/index.html

## Parameters

Name    | Type    | Requirements | Description
--------|---------|--------------|-------------
`t`     | String  | Optional     | [Chart type](#chart-types-supported) (default is `line`)
`d`     | String  | Optional     | Direction of series, `column` or `row` (default is `column`)
`titleOpts` | JSON| Optional     | [Title settings](#title-settings)
`legendOpts`| JSON| Optional     | [Legend settings](#legend-settings)

### Chart types supported

* `line`
* `bar`
* `pie`
* `scatter`

### Title settings

Title settings are given by JSON.  Main properties are:

* `text` - Main title (string)
* `subtext` - Sub title (string)

Example to show title and sub title:

```
titleOpts:{"text":"Awesome chart","subtext":"awful results"}
```

See <https://echarts.apache.org/en/option.html#title> for other properties.

### Legend settings

Legend settings are given by JSON.  Main properties are:

* `show` - Show legends (boolean)

Example to show legends:

```
legendOpts:{"show":true}
```

See <https://echarts.apache.org/en/option.html#legend> for other properties.

## Data examples

### for line, bar, and scatter

```csv
date,A series,B series,C series
1/7,1.0,3.5,0.0
1/14,1.5,3.1,2.0
1/21,1.9,2.8,3.0
1/28,2.2,2.6,3.5
2/4,2.4,2.5,3.7
2/11,2.5,2.5,3.8
2/18,2.5,2.4,3.7
3/4,2.6,2.2,3.5
3/11,2.8,1.9,3.0
3/18,3.1,1.5,2.0
3/25,3.5,1.0,0.0
```

### for pie

```csv
PC OS,Share
Windows,71.4
macOS,11.74
Linux,0.8
Unknown,15.44
```
