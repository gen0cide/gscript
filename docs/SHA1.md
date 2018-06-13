# `SHA1(data)`

Calculate SHA1 hash of data.

## Argument List

 * `data` ([]byte) - Data to hash.

## Return Type

 * `obj.value` (string) - SHA1 hash of data.

## Example

```js
function BeforeDeploy() {
  console.log("SHA1 test starting")
  return true;
}

function Deploy() {
  sha1res = SHA1("test")
  console.log("SHA1 of \"test\": " + sha1res.value)
  return true;
}

function AfterDeploy() {
  console.log("SHA1 test complete");
  return true;
}
```

