// Example gscript template
// Title: Cross Platform Execute Async Example
// Author: ahhh
// Purpose: executes something that should work on multiple platforms
// Gscript version: 1.0.0
// ATT&CK:

//priority:150
//timeout:150

function Deploy() {
  console.log("Starting Exec Command");
  var response = G.exec.ExecuteCommandAsync("netstat", ["-a"]);
  console.log("cmd_obj: " + response[0]);
  return true;
}