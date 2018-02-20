// genesis script

function BeforeDeploy() {
    console.log("testing GetHostname!");
    host = GetHost();
    if (host.toUpperCase() == "TEQUILABOOMBOOM") {
        console.log("Sandbox detected, exiting");
        return false;
    }
    if (host.toUpperCase() == "SANDBOX") {
        console.log("Sandbox detected, exiting");
        return false;
    }
    if (host.toUpperCase() == "VIRUS") {
        console.log("Sandbox detected, exiting");
        return false;
    }
    if (host.toUpperCase() == "MALWARE") {
        console.log("Sandbox detected, exiting");
        return false;
    }
    if (host.toUpperCase() == "MALTEST") {
        console.log("Sandbox detected, exiting");
        return false;
    }

    console.log("testing RegKeys for Windows!");
    var VMWare_val = QueryRegKey("LM", "Software\\Microsoft\\Windows\\CurrentVersion\\Run", "VMWare User Process");
    if (VMWare_val == "\"C:\\Program Files\\VMware\\VMware Tools\\vmtoolsd.exe\" -n vmusr" ){
        console.log("VM detected, exiting");
        return false;
    }

    return true; 
}

function Deploy() {
  console.log("No Sandboxes / VMs Detected!");
}

function AfterDeploy() {
  return true;
}
