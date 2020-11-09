# Simple Helm2CSAR

A (very) simple tool for turning a Helm chart into a CSAR for use with VMware Telco Cloud Automation (TCA).

## Usage

```text
Simple Helm-2-CSAR generator

Usage:
  h2c [command]

Available Commands:
  generate    Generate
  help        Help about any command

Flags:
  -h, --help   help for h2c

Use "h2c [command] --help" for more information about a command.
```

### Generate

`./h2c generate <path_to_chart_dir>` will create a CSAR package in the same directory by reading information in the chart provided. By default, the `provider` is set to `VMware` although this can be overridden with the `--provider` flag.

## Development

Run `make all`. This will repackage the static assets / templates and build the binary.

Please file issues before PRs for any changes more than a few lines of code :)