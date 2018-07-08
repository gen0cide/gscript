//priority:15

function sleep(seconds){
    var waitUntil = new Date().getTime() + seconds*1000;
    while(new Date().getTime() < waitUntil) true;
}

function BeforeDeploy() { return true; }

function Deploy() { console.log("a"); sleep(5); return true; }

function AfterDeploy() { return true; }
