TODO: Intro

## New Features!

### CLI/Options: Add `--tag` flag and `tags` option to set test-wide tags (#553)

You can now specify any number of tags on the command line using the `--tag NAME=VALUE` flag. You can also use the `tags` option to the set tags in the code.

The specified tags will be applied across all metrics. However if you have set a tag with the same name on a request, check or custom metric in the code that tag value will have precedence.

Thanks to @antekresic for their work on this!

**Docs**: [Test wide tags](https://docs.k6.io/v1.0/docs/tags-and-groups#section-test-wide-tags) and [Options](https://docs.k6.io/v1.0/docs/options#section-available-options)

### k6/http: Support for HTTP NTLM Authentication (#556)

```js
import http from "k6/http";
import { check } from "k6";

export default function() {
    // Passing username and password as part of URL plus the auth option will authenticate using HTTP Digest authentication
    let res = http.get("http://user:passwd@example.com/path", {auth: "ntlm"});

    // Verify response
    check(res, {
        "status is 200": (r) => r.status === 200
    });
}
```

**Docs**: [HTTP Params](http://k6.readme.io/docs/params-k6http)

### HAR converter: Add support for correlating JSON values (#516)

There is now support for correlating JSON values in recordings, replacing recorded request values with references to the previous response.

Thanks to @cyberw for their work on this!

### InfluxDB output: Add support for sending certain sample tags as fields (#585)

Since InfluxDB indexes tags, highly variable information like `vu`, `iter` or even `url` may lead to high memory usage. The InfluxDB documentation [recommends](https://docs.influxdata.com/influxdb/v1.5/concepts/schema_and_data_layout/#encouraged-schema-design) to use fields in that case, which is what k6 does now. There is a new `INFLUXDB_TAGS_AS_FIELDS` option (`collectors.influxdb.tagsAsFields` in the global k6 JSON config) that specifies which of the tags k6 emits will be sent as fields to InfluxDB. By default that's only `url` (but not `name`), `vu` and `iter` (if enabled).

Thanks to @danron for their work on this!

### Configurable setup and teardown timeouts (#602)

Previously the `setup()` and `teardown()` functions timed out after 10 seconds. Now that period is configurable via the `setupTimeout` and `teardownTimeout` script options or the `K6_SETUP_TIMEOUT` and `K6_TEARDOWN_TIMEOUT` environment variables. The default timeouts are still 10 seconds and at this time there are no CLI options for changing them to avoid clutter.

### In-client aggregation for metrics streamed to the cloud (#600)

Metrics streamed to the Load Impact cloud can be partially aggregated to reduce bandwidth usage and processing times. Outlier metrics are automatically detected and excluded from that aggregation.

**Docs**: [Load Impact Insights Aggregation](https://docs.k6.io/docs/load-impact-insights#section-aggregation)

### Remote IP address as an optional system metric tag (#616)

It's now possible to add the remote server's IP address to the tags for HTTP and WebSocket metrics. The `ip` [system tag](https://docs.k6.io/docs/tags-and-groups#section-system-tags) is not included by default, but it could easily be enabled by modifying the `systemTags` [option](https://docs.k6.io/docs/options).

### Raw log format (#634)
There is a new log format called `raw`. When used, it will print only the log message without adding any debug information like, date or the log level. It should be useful for debuging scripts when printing a HTML response for example.

```
$ k6 run --log-format raw ~/script.js
```

### Option to output metrics to Apache Kafka (#617)

There is now support for outputing metrics to Apache Kafka! You can configure a Kafka broker (or multiple ones), topic and message format directly from the command line like this:

`k6 --out kafka=brokers={broker1,broker2},topic=k6,format=json`

The default `format` is `json`, but you can also use the [InfluxDB line protocol](https://docs.influxdata.com/influxdb/v1.5/write_protocols/line_protocol_tutorial/) for direct ingestion by InfluxDB:

`k6 --out kafka=brokers=my_broker_host,topic=k6metrics,format=influxdb`

You can even specify format options such as the [`tagsAsFields` option](#influxdb-output-add-support-for-sending-certain-sample-tags-as-fields-585) for InfluxDB:

`k6 --out kafka=brokers=someBroker,topic=someTopic,format=influxdb,influxdb.tagsAsFields={url,name,myCustomTag}`

**Docs**: [Apache Kafka output](https://docs.k6.io/docs/results-output#section-apache-kafka-output)

Thanks to @jmccann for their work on this!


### Multiple outputs (#624)

It's now possible to simultaneously send the emitted metrics to several outputs by using the CLI `--out` flag multiple times, for example:
`k6 run --out json=test.json --out influxdb=http://localhost:8086/k6`

Thanks to @jmccann for their work on this!

## UX

* Clearer error message when using `open` function outside init context (#563)
* Better error message when a script or module can't be found (#565). Thanks to @antekresic for their work on this!

## Internals

* Removed all httpbin.org usage in tests, now a local transient HTTP server is used instead (#555). Thanks to @mccutchen for the great [go-httpbin](https://github.com/mccutchen/go-httpbin) library!
* Fixed various data races and enabled automated testing with `-race` (#564)

## Bugs
* Archive: archives generated on Windows can now run on *nix and vice versa. (#566)
* Submetrics are being tagged properly now. (#609)
* HTML: fixed the `Selection.each(fn)` function, which was returning only the first element. (#610)
* Invalid Stages option won't keep k6 running indefinitely. (#615)
* the `--no-color` option is now being repected for the logs. (#634)

## Breaking changes
* The Load Impact cloud configuration options `no_compress` and `project_id` and the `payload_size` InfluxDB option have been renamed to `noCompress`, `projectID` and `payloadSize` respectively, to match the other JS option names.
