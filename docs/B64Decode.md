# `B64Decode(data)`

Base64 decode `data`.

## Argument List

 * `data` (string) - Data to decode.

## Return Type

 * `obj.value` ([]byte) - Decoded data.
 * `obj.execError` (error) - Error message.

## Example

```js
var obj = B64Decode("aSBjYW50IGJlbGlldmUgc29tZW9uZSB3b3VsZCBkZWNvZGUgdGhpcwo=")
// obj.value
```
