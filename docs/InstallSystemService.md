# `InstallSystemService(path, name, displayname, description)`

Install `path` binary as a system service.

## Argument List

 * `path` (string) - Path to binary.
 * `name` (string) - Name of service.
 * `displayName` (string) - Display name of service.
 * `description` (string) - Description of service.

## Return Type

 * `nil`
 * `obj.installError` (error) - Error.

## Example

```js
var obj = InstallSystemService(path, name, displayName, description)
// obj.installError
```

