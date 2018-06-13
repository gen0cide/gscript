# `InjectIntoProc(shellcode, processID)`

Inject shellcode into `processID`.

## Argument List

 * `shellcode` (string) - Shellcode
 * `processID` (int64) - PID to inject `shellcode` into.

## Return Type

 * `obj.runtimeError` (error) - Error.

## Example

```js
var obj = InjectIntoProc(shellcode, processID)
// obj.runtimeError
```

