# `PostJSON(url, json)`

Post JSON data to `url`.

## Argument List

 * `url` (string) - URL to POST data to.
 * `json` (string) - JSON to post to `url`.

## Return Type

 * `obj.statusCode` (int) - HTTP response code.
 * `obj.response` ([]byte) - HTTP data.
 * `obj.runtimeError` (error) - Error.

## Example

```js
var obj = PostJSON(url, json)
// obj.statusCode
// obj.response
// obj.runtimeError
```

