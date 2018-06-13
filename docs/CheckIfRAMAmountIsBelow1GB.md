# `CheckIfRAMAmountIsBelow1GB()`

Checks for low amounts of RAM, which is common for sandboxed environments.

## Argument List

 * None.

## Return Type

 * `obj.areWeInASandbox` (bool) - true if below 1GB, otherwise false.

## Example

```js
var obj = CheckIfRAMAmountIsBelow1GB()
// obj.areWeInASandbox
```

