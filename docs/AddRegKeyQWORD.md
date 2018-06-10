# `AddRegKeyQWORD(registryString, path, name, value)`

Add a QWORD registry key.

## Argument List

 * `registryString` (string) - Contents of registry key.
 * `path` (string) - Path to registry key.
 * `name` (string) - Registry key name.
 * `value` (uint42) - Contents of registry key.

## Return Type

 * `obj.runtimeError` (error) - Error

## Example

```js
var obj = AddRegKeyQWORD(registryString, path, name, value)
// obj.runtimeError
```

