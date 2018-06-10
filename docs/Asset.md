# `Asset(assetName)`

Retrieves the packed bytes of file `assetName` from the GENESIS packed archive.

## Argument List

 * `file` (string) - Name of the file you imported via your gscript.

## Return Type

 * `obj.fileData` ([]byte) - Asset data.
 * `obj.err` (error) - An error message.

## Example
```js
//import:/bin/ls

function BeforeDeploy() {
  console.log("Start of Asset() test")
  return true; 
}

function Deploy() {
  data = Asset("ls").fileData
  WriteFile("/tmp/ls", data, 0755);
  console.log("Wrote file to /tmp/ls");
  return true;
}

function AfterDeploy() {
  console.log("End of Asset() test")
  return true;
}
```

