# `GetProcName(pid)`

Return the name of a PID.

## Argument List

 * `pid` (int) - PID of process to look up process name.

## Return Type

 * `obj.procName` (string) - Name of process.
 * `obj.runtimeError` (error) - Error.

## Example

```js
var obj = GetProcName(1337)
// obj.procName
// obj.runtimeError
```

