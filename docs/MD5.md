# `MD5(data)`

Calculate an MD5 hash of data.

## Argument List

 * `data` ([]byte) - Data to hash.

## Return Type

 * `obj.value` (string) - Hash of data.

## Example

```js
function BeforeDeploy() {
  console.log("MD5 test starting")
  return true;
}

function Deploy() {
  md5res = MD5("test")
  console.log("MD5 of \"test\": " + md5res.value)
  return true;
}

function AfterDeploy() {
  console.log("MD5 test complete");
  return true;
}
```

