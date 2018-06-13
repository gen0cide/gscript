# `AppendFileString(path, addString)`

Append a string to the end of a file.

## Argument List

 * `path` (string) - Path to file.
 * `addString` (string) - String to append to file.

## Return Type

 * `nil`
 * `obj.fileError` (error) - Error.

## Example

```js
var obj = AppendFileString("/tmp/foo", "International House of Burgers")
// obj.fileError
```

