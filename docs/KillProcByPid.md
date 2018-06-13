# `KillProcByPid(pid)`

Kill a process.

## Argument List

 * `pid` (in64) - PID to kill.

## Return Type

 * `obj.dead` (bool) - true if dead, false if still alive.
 * `obj.runtimeError (error) - Error.

## Example

```js
var obj = KillProcByPid(1337)
// obj.dead
// obj.runtimeError
```

