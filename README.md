# omq

A processor of the [OpenMetrics](https://openmetrics.io/) format text for the command-line, and Go library.

WIP WIP WIP

## Usage Example

### Convert OpenMetrics text to JSON

```
$ cat <<EOF | omq -j
# HELP http_response_latency_seconds the HTTP response latency for the client request
# TYPE http_response_latency_seconds histogram
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.005",status="2xx"} 4.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.01",status="2xx"} 3.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.025",status="2xx"} 2.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.05",status="2xx"} 1.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.075",status="2xx"} 1.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.1",status="2xx"} 1.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="+Inf",status="2xx"} 1.0
http_response_latency_seconds_count{http_method="GET",path="/index.html",status="2xx"} 4.0
http_response_latency_seconds_sum{http_method="GET",path="/index.html",status="2xx"} 0.001654784
EOF
[{"name":"http_response_latency_seconds","help":"the HTTP response latency for the client request","type":"HISTOGRAM","metrics":[{"type":"HISTOGRAM","labels":{"http_method":"GET","path":"/index.html","status":"2xx"},"data":{"values":{"+Inf":1,"0.005000":4,"0.010000":3,"0.025000":2,"0.050000":1,"0.075000":1,"0.100000":1},"sampleSum":0.001654784,"sampleCount":4}}]}]
```

### Convert JSON to OpenMetrics text

```
$ echo '[{"name":"http_response_latency_seconds","help":"the HTTP response latency for the client request","type":"HISTOGRAM","metrics":[{"type":"HISTOGRAM","labels":{"http_method":"GET","path":"/index.html","status":"2xx"},"data":{"values":{"+Inf":1,"0.005000":4,"0.010000":3,"0.025000":2,"0.050000":1,"0.075000":1,"0.100000":1},"sampleSum":0.001654784,"sampleCount":4}}]}]' | omq -o
# HELP http_response_latency_seconds the HTTP response latency for the client request
# TYPE http_response_latency_seconds histogram
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",status="2xx",le="+Inf"} 1
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",status="2xx",le="0.005"} 4
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",status="2xx",le="0.01"} 3
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",status="2xx",le="0.025"} 2
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",status="2xx",le="0.05"} 1
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",status="2xx",le="0.075"} 1
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",status="2xx",le="0.1"} 1
http_response_latency_seconds_sum{http_method="GET",path="/index.html",status="2xx"} 0.001654784
http_response_latency_seconds_count{http_method="GET",path="/index.html",status="2xx"} 4
```

### Query by jq equivalent query

```
$ echo -n "==> " && cat <<EOF | go run ./cmd/omq/main.go -q '.[] | select(.name == "http_response_latency_seconds") | .metrics[].data.values."0.025000"'
# HELP http_response_latency_seconds the HTTP response latency for the client request
# TYPE http_response_latency_seconds histogram
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.005",status="2xx"} 4.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.01",status="2xx"} 3.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.025",status="2xx"} 2.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.05",status="2xx"} 1.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.075",status="2xx"} 1.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="0.1",status="2xx"} 1.0
http_response_latency_seconds_bucket{http_method="GET",path="/index.html",le="+Inf",status="2xx"} 1.0
http_response_latency_seconds_count{http_method="GET",path="/index.html",status="2xx"} 4.0
http_response_latency_seconds_sum{http_method="GET",path="/index.html",status="2xx"} 0.001654784
EOF
==> 2
```

## Author

moznion (<moznion@mail.moznion.net>)

