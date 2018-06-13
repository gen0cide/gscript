# `ForkExecuteCommand(baseCmd, cmdArgs)`

Fork and execute a system command.

## Argument List

 * `baseCmd` (string) - Command to execute.
 * `cmdArgs` ([]string) - Command arguments.

## Return Type

 * `obj.pid` (int) - PID of forked process.
 * `obj.execError (error) - Error.

## Example

```js
var obj = ForkExecuteCommand(baseCmd, cmdArgs)
// obj.pid
// obj.execError
```

