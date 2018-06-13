# `FileExists(path)`

Check if a file exists.

## Argument List

 * `path` (string) - Path to file.

## Return Type

 * `obj.fileExists` (bool) - `true` if file exists, otherwise `false`.
 * `obj.fileError` (error) - Error.

## Example

```js
var obj = FileExists("/tmp/foo")
// obj.fileExists
// obj.fileError
```

