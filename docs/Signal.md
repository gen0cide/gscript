# `Signal(pid, signal)`

Send a signal to a process.

## Argument List

 * `pid` (String) - Process ID you wish to signal.
 * `signal` (Integer) - Signal you wish to send (9, 15, etc.)

## Return Type

- `obj.runtimeError` (error) - Error.

## Example

```js
var obj = Signal(1337, 9)
// obj.runtimeError
```
