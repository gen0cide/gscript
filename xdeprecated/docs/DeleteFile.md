# `DeleteFile(path)`

Delete the file located at `path`.

##### Argument List

 * `path` (String) - Path to file you wish to delete.

##### Return Type

`boolean` (true = success, false = error)

## Example

```
var file_to_drop = "/tmp/tmpfile";

function Deploy() {
		var file_1 = file_to_drop;
    var return_value1 = DeleteFile(file_1);
}

// ....
```

This will automatically delete a file in the path 