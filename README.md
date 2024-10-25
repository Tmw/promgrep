# promgrep

Prometheus aware grep

## Getting started

```console
make build
```

## Example

```console
$ cat testdata/example_scrape.txt | promgrep {instance=server-1}

cpu_usage_percentage{core="0", instance="server-1"} 75.5
cpu_usage_percentage{core="1", instance="server-1"} 65.2
```

## License

[MIT](./LICENSE)
