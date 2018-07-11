// Example gscript template
// Title: Multi Platform Delete File Example
// Author: ahhh
// Purpose: Testing a delete file on different platforms
// Gscript version: 1.0.0
// ATT&CK:

//priority:150
//timeout:150
//go_import:os as os

function Deploy() {

  console.log("Starting Delete File");
  var errors = os.Remove("test_file.txt");
  console.log("errors: " + errors);

  return true;
}

// //import:/tmp/testasset.txt

// function Deploy() {
//   a = G.time.GetUnix();
//   console.log("CURRENT UNIX TIME: " + a);
//   b = GetAssetAsString("testasset.txt");
//   console.log(b);
//   return true;
// }