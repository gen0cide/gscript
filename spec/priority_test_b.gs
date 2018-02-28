//priority:30

function sleep(seconds){
    var waitUntil = new Date().getTime() + seconds*1000;
    while(new Date().getTime() < waitUntil) true;
}

function BeforeDeploy() { return true; }

function Deploy() { console.log("b"); sleep(5); return true; }

function AfterDeploy() { return true; }
