# `CopyFile(srcPath, dstPath, perms)`

Copy file from `srcPath` to `dstPath` with given `perms`.

## Argument List

 * `srcPath` (string) - Path to source file.
 * `dstFile` (string) - Path to destination file.
 * `perms` (int64) - Permissions.

## Return Type

 * `nil`
 * `obj.fileError` (error) - Error.

## Example

```js
var file_1 = "/etc/passwd";
var file_2 = "/tmp/rock";
var return_value = CopyFile(file_1, file_2);
```

