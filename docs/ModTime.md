# `ModTime(path)`

Get last moditified time of a file.

## Argument List

 * `path` (string) - Path to file.

## Return Type

 * `obj.modTime` (int64) - Last modified time.
 * `obj.fileError` (error) - Error. 

## Example

```js
var obj = ModTime("/tmp/foo")
// obj.modTime
// obj.fileError
```

