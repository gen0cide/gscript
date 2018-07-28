# Assets

## Embedding Assets
As previously covered in **Macros**, you can use the `//import:` macro to include arbitrary files into your final binary. These files are addressable in the gscripts they are called in. These assets have weak static encryption applied to them while embedded, and are not addressable from other gscripts. 
### Asset embedding example:
`//import:/path/to/binary/file.exe`

## Asset retrieval functions
### GetAssetAsBytes("embedded_asset")
This function returns two objects, a byte array (pos 0) and an error object (pos 1), this takes a string as a variable which is referenced by the filename of the embedded asset.
### GetAssetAsBytes example:
`file = GetAssetAsBytes("file.exe");`
`errors = G.file.WriteFileFromBytes("example_test", file[0]);`

### GetAssetAsString("embedded_asset")
This function returns two objects, a string (pos 0) and an error object (pos 1), this takes a string as a variable which is referenced by the filename of the embedded asset. A useful workaround in javascript to avoid type errors is converting most things to a string, then bringing them back to their original type in GoLang.
### GetAssetAsString example:
`file = GetAssetAsString("file.exe");`
`errors = G.file.WriteFileFromString("example_test", file[0]);`
