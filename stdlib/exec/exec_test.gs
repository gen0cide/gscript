// Example gscript template
// Title: Cross Platform Execute Example
// Author: ahhh
// Purpose: executes something that should work on multiple platforms 
// Gscript version: 1.0.0

//priority:150
//timeout:150
//go_import:github.com/gen0cide/gscript/stdlib/exec as exec

function Deploy() {  
    console.log("Starting Exec Command");
    var response = G.exec.ExecuteCommand("netstat", ["-a"]);
    console.log("Pid: "+response[0]);
    console.log("stdout: "+response[1])
    console.log("stderr: "+response[2])
    console.log("exit code: "+response[3])
    console.log("go errors: "+response[4])
    return true;
  }

