# `Chmod(path, perms)`

Change permissions of a file.

## Argument List

 * `path` (string) - Path to file to chmod.
 * `perms (int64) - Permissions to apply to `path`.

## Return Type

 * `nil`
 * `obj.osError` (error) - Error.

## Example

```js
var obj = Chmod("/tmp/foo", 0755)
// obj.osError
```

