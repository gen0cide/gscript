# `FindProcByName(procName)`

Find the PID of `procName`.

## Argument List

 * `procName` (string) - Name of process.

## Return Type

 * `obj.pid` (int) - PID of process.
 * `obj.procError (error) - Error.

## Example

```js
var obj = FindProcByName("httpd")
// obj.pid
// obj.procError
```

