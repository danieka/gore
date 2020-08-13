# Gore

Gore is a reporting server that makes it easy to write and manange you data reports.

Reports are a written in a single `.rpt` file and contains both a query for getting data from an input source and one or many formatters for outputting the data. Supported input formats are MySQL.

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
name: test
</info>
<source sql>
SELECT email FROM user
</source>
```

Save this file as `test.rpt` in the repository. Now you can start Gore with `go run *.go`. Gore will now start a web server and your report will be accesible on `http://localhost:16772/`.

The `info` section contains metadata about the report.