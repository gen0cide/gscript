# `HTTPGetFile(url)`

Retrieve `url` over HTTP(s).

## Argument List

 * `url` (string) - URL to retrieve.

## Return Type

 * `obj.statusCode` (int) - HTTP response.
 * `obj.file` ([]byte) - Contents of URL.
 * `obj.runtimeError` (error) - Error.

## Example

```js
var obj = HTTPGetFile("https://github.com")
// obj.statusCode
// obj.file
// obj.runtimeError
```

