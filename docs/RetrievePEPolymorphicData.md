# `RetrievePEPolymorphicData(peFile)`

Retrieve data stored within uninitalized space at the end of the gscript binary.

## Argument List

 * `peFile` (string) - Path to PE file.

## Return Type

 * `obj.data` ([]byte) - Data after the end of the specified PE file.
 * `obj.runtimeError (error) - Error.

## Example

```js
var obj = RetrievePEPolymorphicData("C:\test.exe")
// obj.data
// obj.runtimeError
```

