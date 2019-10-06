# TCC

This project evaluates the efficiency of different data formats and their serialization libraries over an HTTP(S) channel. The intent is to, besides evaluating each format, consider their impact in a real communication scenario. The scenarios for this project are run locally and in the cloud. Formats evaluated:  
* JSON
* Protocol Buffers
* Apache Avro
* MessagePack

The project has three different programs:  
1. **consumer** is the passive one, being an HTTP server which receives marshalled data, unmarshals it and responds with the time measures taken from receiving the request to sending the response.
2. **producer** is the active one, being a consumer client, marshaling some generated data (the time for generating the data is not considered), sending it to the consumer and grouping its own measures (marshaling + requesting) with the ones coming from the consumer response. The measures are logged to *stderr*.
3. **parser** parses the measures from a log file (either from the producer or the consumer).

## Measuring

For running measures, first you need a running **consumer**, then you may run the **producer** against that consumer endpoint (`c` parameter). The producer takes several different parameters explained further below. Once you run the producer, you may use the **parser** over the producer stderr output. E.g.:

```sh
producer -c "http://localhost:9000" &> p.log
cat p.log | parser
```

### Request parameters

The single producer run may generate a number of requests `r` with a given interval `i` in milliseconds. The requests are started concurrently, unless `i` is `0` (its default value), which then generates sequential requests without an additional interval between them.

E.g.:
```sh
# produces 100 requests with 100ms in between one another
producer -r 100 -i 100
```

### Payload parameters

The payload can be customized in format and size. To set a format, use the `f` parameter (defaults to JSON) in the producer; the consumer will automatically detected the format according to the *Content-Type* header. The `f` value follows a MIME-like types string (not all formats have a standard MIME type):
* JSON: `application/json`
* ProtoBuff: `application/x-protobuf`
* Avro: `application/x-avro`
* MsgPack: `application/x-msgpack`

The type of the data is set via the `t` parameter. Each request sends an array of length `l` of that given type. Each type may have its own custom parameters. Supported types are:

1. `string`: go string type with a set of random characters (a-Z) of size `strlen`.
2. `int`: go int64 type ranging from abs(`intmin`) to abs(`intmax`). The absolute value is then set as negative value with 50% chance.
3. `float`: go float64 type.
4. `object`: go struct type with the given fields:
    ```
    I int64
	F float64
	T bool
	S string
	B []byte
    ```
    B is subject to the same string randomization (which is converted to bytes afterwards).