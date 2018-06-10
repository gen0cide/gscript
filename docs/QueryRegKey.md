# `QueryRegKey(registryString, path, key)`

Retrieve a registry key.

## Argument List

 * `registryString` (string) - Registry string.
 * `path` (string) - Path.
 * `key` (string) - Key to retrieve

## Return Type

 * `obj.keyObj` (RegistryRetValue) - The contents of `key`.
 * `obj.runtimeError (error) - Error.

## Example

```js
var obj = QueryRegKey(registryString, path, key)
// obj.keyObj
// obj.runtimeError
```

