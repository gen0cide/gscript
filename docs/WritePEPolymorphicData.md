# `WritePEPolymorphicData(peFile, data)`

Write `data` to the uninitalized space at the end of `peFile`.

## Argument List

 * `peFile` (string) - Path to PE file.
 * `data` ([]byte) - Data to write to end of `peFile`.

## Return Type

 * `nil`
 * `obj.runtimeError` (error)

## Example

```js
var obj = WritePEPolymorphicData("C:\test.exe", "datahere")
// obj.runtimeError
```

