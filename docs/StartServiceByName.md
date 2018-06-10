# `StartServiceByName(name)`

Start a system service.

## Argument List

 * `name` (string) - Name of service to start.

## Return Type

 * `nil`
 * `obj.startError` (error) - Error.

## Example

```js
var obj = StartServiceByName("httpd")
// obj.startError
```

