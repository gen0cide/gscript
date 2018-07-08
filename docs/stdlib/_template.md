# Package: asset

Handles embedded assets.

## Function Index

* Foo 
* Bar

## Details

### GetAssetAsString

**Author:** gen0cide

**Description:** fizzbuzz

**Method Signature:**

```
GetAssetAsString(name string) (string, error)
```

**Arguments:**

| Label   | Type         | Description                  |
|---------|--------------|------------------------------|
| `name`  | `string`     | name of the asset            |

**Returns:**

| Position | Type         | Description                  |
|----------|--------------|------------------------------|
| `0`      | `string`     | Asset bytes as a string      |
| `1`      | `error`      | (optional) function error    |

**Example Usage:**

```
// when an asset was embedded into a binary
var asset = G.asset.GetAssetAsString("real.txt");
console.log(asset[0]);
// => "contents of real.txt"
console.log(asset[1]);
// => null

// if you call the wrong asset name
var asset2 = G.asset.GetAssetAsString("notreal.txt");
console.log(asset[0]);
// => ""
console.log(asset[1]);
// => "error msg..."
```

-
