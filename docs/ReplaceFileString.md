# `ReplaceFileString(file, match, replacement)`

Search `file` for `match`, replacing with `replacement`.

## Argument List

 * `file` (string) - Path to file.
 * `match` (string) - String to find.
 * `replacement` (string) - String to replace `match` with.

## Return Type

 * `obj.stringsReplaced` (int) - Number of strings replaced.
 * `obj.fileError` (error) - Error.

## Example

```js
var obj = ReplaceFileString("/tmp/file", "foo", "bar")
// obj.stringsReplaced
// obj.fileError
```

