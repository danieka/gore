# Gore

Gore is a reporting server that makes it easy to write and manange you data reports.

Reports are a written in a single `gore` file and contains both a query for getting data from an input source and one or many formatters for outputting the data. The outputs of these reports are then exposed through a REST API. So Gore handles all the shuffling of data, you just decleratively define what data you want formatted how. Gore handles fetching data, rendering outputs and serving the finished reports.

Supported input formats are MySQL.

## Installation

Right now the only way to install Gore is to clone the repo, download Go and then `go run *.go` in the repository.

## Getting started

To get started you need to configure a data source and write your first report. The format for configuring the data source is:

```
sources:
  default:
    type: mysql
    host: 127.0.0.1
    port: 3306
    database: test
    username: user
    password: passwd
```

Edit this and save it as config.yaml in the repository. Now you need to create your first report. The format for reports is:

```
<info>
id: test
</info>
<source sql>
SELECT id, description FROM article
</source>
<output json>
[
    {{ range .Rows }}
    {
        "id": "{{ .id }}",
        "description": "{{ .description }}",
    },
    {{ end }}
]
</output>
```

Save this file as `test.gore` in the repository. Now you can start Gore with `go run *.go`. Gore will now start a web server and your report will be accesible on `http://localhost:16772/report/{id}`.

The `info` section contains metadata about the report. The `id` attribute is used when accessing the report through the rest API.

The `source` section contains the SQL query for fetching data. If no specific source is given it uses the `default` source.

The `output` section contains a template detailing how the data should be formatted when outputted. Supported formats right now are `json`, `xlsx` and `csv` but also expect HTML and PDF in the future. Output uses the Go templating language. The variable `.Rows` is the rows resulting from the query. Also available is the `.Cols` variable which is a `[]string` of the column names from the database.

If you you leave the `output` tag empty Gore will use a default template. This template should be usable for most cases. XLXS only has a default format, so any custom formatting you enter will not be taken into account.

## Gore command

```
gore [flags]
```

```
  -i	Run server in interactive mode (suitable for development)
  -s string
    	filesystem path where reports sources are
  -w	Watch .gore files for changes and reload reports on changes
```
