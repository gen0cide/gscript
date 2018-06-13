# `ModifyTimestamp(path, accessTime, modifyTime)`

Change the access and modified time of `path`.

## Argument List

 * `path` (string) - Path to file to modify.
 * `accessTime` (int64) - Access time to apply to `path`.
 * `modifyTime` (int64) - Modify time to apply to `path`.

## Return Type

 * `nil`
 * `obj.fileError` (error) - Error.

## Example

```js
var obj = ModifyTimestamp("/tmp/foo", accessTime, modifyTime)
// obj.fileError
```

