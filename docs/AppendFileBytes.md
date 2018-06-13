# `AppendFileBytes(path, fileData)`

Add a byte array to the end of a file.

## Argument List

 * `path` (string) - Path to file
 * `fileData` ([]byte) - Data to write.

## Return Type

 * `nil`
 * `obj.fileError` (error) - Error.

## Example

```js
var obj = AppendFileBytes("/tmp/foo", variable)
// obj.fileError
```

