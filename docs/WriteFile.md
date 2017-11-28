### `WriteFile(path, bytes, perms)`

Write `bytes` to `path` and set perms to `perms`.

##### Argument List

 * `path` (String) - Path to file you wish to write.
 * `bytes` (Array) - Array of bytes you wish to write to the `path` location.
 * `perms` (String) - Octal unix permissions represented as a string. ie: `0777`.

##### Return Type

`boolean` (true = success, false = error)

#Example

This will write a file.