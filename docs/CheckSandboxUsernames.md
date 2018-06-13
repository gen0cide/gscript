# `CheckSandboxUsernames()`

Check for common sandbox system usernames. Currently checks for SANDBOX, VIRUS,
MALWARE, malware, virus, sandbox, .bin, .elf, and .exe usernames.

## Argument List

 * None.

## Return Type

 * `obj.areWeInASandbox` (bool) - true if users are present, otherwise false
 * `obj.runtimeError` (error) - Error.

## Example

```js
var obj = CheckSandboxUsernames()
// obj.areWeInASandbox
// obj.runtimeError
```

