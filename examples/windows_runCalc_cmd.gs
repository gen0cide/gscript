// genesis script

function BeforeDeploy() {
    console.log("Entered BeforeDeploy!");
    return true; 
}

function Deploy() {
    console.log("Entered Deploy!");
    ExecuteCommand("calc.exe", Array(""))
    return true;
}

function AfterDeploy() {
    console.log("Entered AfterDeploy!");
    return true;
}


