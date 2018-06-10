# `Timestamp()`

Get the system's current time in epoch format.

## Argument List

 * None.

## Return Type

  * `obj.value` (int64) - Timestamp in epoch format.

## Example

```js
function BeforeDeploy() {
  console.log("Timestamp test starting")
  return true;
}

function Deploy() {
  ts = Timestamp()
  console.log("Timestamp: " + ts.value)
  return true;
}

function AfterDeploy() {
  console.log("Timestamp test complete");
  return true;
}
```
