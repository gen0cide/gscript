// Example gscript template
// Title: Multi Platform Write File Example
// Author: ahhh
// Purpose: Testing an asset and write file on different platforms
// Gscript version: 1.0.0
// ATT&CK: 

//import:/private/tmp/example.bin

//priority:150
//timeout:150

function Deploy() {

    console.log("Starting Write file from String");
    var writeStringErrors = G.file.WriteFileFromString("/go/src/github.com/gen0cide/gscript/stdlib/file/file_test4.txt", "Testing some stuff\n");
    console.log("errors: "+ writeStringErrors)

};
