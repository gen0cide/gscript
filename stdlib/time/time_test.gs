// Example gscript template
// Title: Cross Platform Time Example
// Author: ahhh
// Purpose: gets the current Unix time on multiple platforms 
// Gscript version: 1.0.0

//priority:150
//timeout:150

function Deploy() {  
    console.log("Starting Time");
    var response = G.time.GetUnix();
    console.log("Time: "+response);
    return true;
  }
