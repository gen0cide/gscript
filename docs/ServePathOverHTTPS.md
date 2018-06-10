# `ServePathOverHTTPS(port, path, timeout)`

Serve `path` over HTTPS for `timeout` seconds. Default port is 443.

## Argument List

 * `port` (string) - Listen port.
 * `path` (string) - Path to serve.
 * `timeout` (int64) - Timeout (in seconds).

## Return Type

 * `nil`
 * `obj.runtimeError` (error) - Error

## Example

```js
var obj = ServePathOverHTTPS(port, path, timeout)
// obj.runtimeError
```

