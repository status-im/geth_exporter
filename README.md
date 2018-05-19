# geth_exporter

`geth_exporter` is a metrics exporter for [Prometheus](https://github.com/prometheus/prometheus).

## Usage

```
go build && \
  ./geth_exporter -ipc node/data/path/geth.ipc -filter="whisper_*" -filter="les_*"
```

## Docker example

```
cd docker-example
docker-compose up
```

At `http://localhost:9200/metrics` you will have the `geth_exporter`, it should send a response similar to this:

```
discv5_inboundTraffic_avgRate01Min 2981
discv5_inboundTraffic_avgRate05Min 2981
discv5_inboundTraffic_avgRate15Min 2981
discv5_inboundTraffic_meanRate 3591.687045952474
discv5_inboundTraffic_overall 18213
discv5_outboundTraffic_avgRate01Min 1733.4
discv5_outboundTraffic_avgRate05Min 1733.4
discv5_outboundTraffic_avgRate15Min 1733.4
discv5_outboundTraffic_meanRate 2356.6951261606164
....
```

At `http://localhost:9090/graph` you can use Prometheus to query the exporter.

At `http://localhost:3000/` you can log in to grafana with username `admin` and password `admin`.
