# `AddRegKeyString(registryString, path, name, value)`

Add a string registry key.

## Argument List

 * `registryString` (string) - Contents of registry key.
 * `path` (string) - Path to registry key.
 * `name` (string) - Registry key name.
 * `value` (string) - Contents of registry key.

## Return Type

 * `obj.runtimeError` (error) - Error

## Example

```js
var obj = AddRegKeyString(registryString, path, name, value)
// obj.runtimeError
```

