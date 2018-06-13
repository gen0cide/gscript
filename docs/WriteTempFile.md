# `WriteTempFile(name, fileData)`

Writes data from a byte array to a temporary file and returns the file's path.

## Argument List

 * `name` (string) - temporary file name.
 * `fileData` ([]byte) - data to write to file.

## Return Type

 * `obj.fillpath` (string) - fill path to temporary file.
 * `obj.fileError` (error) - Error.

## Example

```js
function BeforeDeploy() {
  console.log("WriteTempFile test starting")
  return true;
}

function Deploy() {
  tempfile = WriteTempFile("lollerskates", "roflcopter")
  console.log("Temporary file: " + tempfile.fullpath)
  return true;
}

function AfterDeploy() {
  console.log("WriteTempFile test complete");
  return true;
}
```

