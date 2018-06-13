# `IsTCPPortInUse(port)`

Check if TCP port `port` is in use.

## Argument List

 * `port` (string) - Port to check

## Return Type

 * `obj.state` (bool) - true if port is in use, otherwise false.

## Example

```js
var obj = IsTCPPortInUse("1337")
// obj.state
```

