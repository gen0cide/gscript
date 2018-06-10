# `ReadFile(path)`

Read contents of file located at `path`

## Argument List

 * `path` (String) - Path to file you would like to read.

## Return Type

 * `obj.fileBytes` ([]byte) - Contents of target file.
 * `obj.fileError` (error) - Error.
 
## Example

```js
var obj = ReadFile("/tmp/foo")
// obj.fileBytes
// obj.fileError
```
