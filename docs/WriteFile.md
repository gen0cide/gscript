# `WriteFile(path, bytes, perms)`

Write `bytes` to `path` and set perms to `perms`.

## Argument List

 * `path` (string) - Path to file you wish to write.
 * `bytes` ([]byte) - Array of bytes you wish to write to the `path` location.
 * `perms` (string) - Unix permissions represented as a string. ie: `0777`.

## Return Type

 * `obj.bytesWritten` (int) - Number of bytes writen to `path`.
 * `obj.fileError` (error) - Error. 

# Example

```js
function BeforeDeploy() {
  console.log("WriteFile test starting")
  return true;
}

function Deploy() {
  filedata = "lol"
  WriteFile("/tmp/lol", filedata, 0644)
  return true;
}

function AfterDeploy() {
  console.log("WriteFile test complete");
  return true;
}
```

