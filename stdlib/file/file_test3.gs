// Example gscript template
// Title: Multi Platform Write File Example
// Author: ahhh
// Purpose: Testing an asset and write file
// Gscript version: 1.0.0
// ATT&CK: 

//import:/go/src/github.com/gen0cide/gscript/stdlib/file/file.go

//priority:150
//timeout:150

function Deploy() {

    console.log("Starting Write Third File form Bytes");
    var myBin = GetAssetAsBytes("file.go");
    errors = G.file.WriteFileFromBytes("/go/src/github.com/gen0cide/gscript/stdlib/file/file_test3", myBin[0]);
    console.log("errors: "+errors);

    console.log("Running ReadFileAsString on file_test3");
    var readFile = G.file.ReadFileAsString("/go/src/github.com/gen0cide/gscript/stdlib/file/file_test3");
    console.log("errors: "+ readFile[1]);
    console.log("contains:\n"+readFile[0]);


    return true;
}
