# geth_exporter

`geth_exporter` is a metrics exporter for [Prometheus](https://github.com/prometheus/prometheus).

## Usage

```
go build && \
  ./geth_exporter -ipc node/data/path/geth.ipc -filter="whisper_*" -filter="les_*"
```
