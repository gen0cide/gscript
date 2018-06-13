# `XorBytes(aByteArray, bByteArray)`

XOR two byte arrays.

### Limitations:

 * Byte arrays must be equal in size.
 * Limited to 20 bytes in length.

## Argument List

 * `aByteArray` ([]byte) - First byte array.
 * `bByteArray` ([]byte) - Second byte array.

## Return Type

 * `obj.value` ([20]byte) - Array containing XOR results.

## Example

```js
function BeforeDeploy() {
  console.log("Testing XorBytes");
  return true; 
}

function Deploy() {
  var string1 = "AAAA"
  var string2 = "BBBB"
  var xorresult = XorBytes(string1, string2)
  console.log("Result: " + xorresult.value)
  return true;
}

function AfterDeploy() {
  console.log("Test complete");
  return true;
}
```
