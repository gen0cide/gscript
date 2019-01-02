// Example gscript template
// Title: Multi Platform Write File Example
// Author: ahhh
// Purpose: Testing an asset and write file
// Gscript version: 1.0.0
// ATT&CK: 

//import:/go/src/github.com/gen0cide/gscript/stdlib/file/file_test.gs
//import:/go/src/github.com/gen0cide/gscript/stdlib/file/file_test.go

//priority:150
//timeout:150

function Deploy() {

    console.log("Starting Write First File form Bytes");
    var myBin = GetAssetAsBytes("file_test.gs");
    errors = G.file.WriteFileFromBytes("/go/src/github.com/gen0cide/gscript/stdlib/file/file_test", myBin[0]);
    console.log("errors: "+errors);

    console.log("Running ReadFileAsString on file_test 1");
    var readFile = G.file.ReadFileAsString("/go/src/github.com/gen0cide/gscript/stdlib/file/file_test");
    console.log("errors: "+ readFile[1]);
    console.log("contains:\n"+readFile[0]);

    console.log("Starting Write Second File form Bytes");
    var myBin2 = GetAssetAsBytes("file_test.go");
    errors = G.file.WriteFileFromBytes("/go/src/github.com/gen0cide/gscript/stdlib/file/file_test2", myBin2[0]);
    console.log("errors: "+errors);

    console.log("Running ReadFileAsString on file_test 2");
    var readFile2 = G.file.ReadFileAsString("/go/src/github.com/gen0cide/gscript/stdlib/file/file_test2");
    console.log("errors: "+ readFile2[1]);
    console.log("contains:\n"+readFile2[0]);

    return true;
}
