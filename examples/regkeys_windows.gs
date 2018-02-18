// genesis script

function BeforeDeploy() {
    console.log("testing GetHostname!");
    host = GetHost();
    console.log(host);
    return true; 
}

function Deploy() {
    console.log("testing registry key functions");

    console.log("reading a VMWare key");
    var key_val = QueryRegKey("LM", "Software\\Microsoft\\Windows\\CurrentVersion\\Run", "VMWare User Process");
    console.log(key_val);

    console.log("reading test key");
    var key_value = QueryRegKey("CU", "Environment", "TEMP");
    console.log(key_value);

    console.log("Adding a reg key for current user run");
    AddRegKey("CU", "Software\\Microsoft\\Windows\\CurrentVersion\\Run", "testExe", "C:\\test.exe")
   
    console.log("Done adding key, reading key");
    var key_value = QueryRegKey("CU", "Software\\Microsoft\\Windows\\CurrentVersion\\Run", "testExe");
    console.log(key_value);

    console.log("Deleting the key now");
    DelRegKey("CU", "Software\\Microsoft\\Windows\\CurrentVersion\\Run", "testExe");
    return true;
}

function AfterDeploy() {
    console.log("All done!");
    return true;
}
