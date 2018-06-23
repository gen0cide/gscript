### `CopyFile(srcPath, dstPath)`

Copy file from `srcPath` to `dstPath`.

##### Argument List

 * `srcPath` (String) - Path to source file.
 * `dstFile` (String) - Path to destination file.

##### Return Type

`boolean` (true = success, false = error)

# Example

```
    var file_1 = "/etc/passwd";
    var file_2 = "/tmp/rock";
    var return_value = CopyFile(file_1, file_2);
```

This will Copy a file.