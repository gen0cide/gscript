function BeforeDeploy() {
  LogInfo("Getting started");
  Halt();
  return true;
}

function Deploy() {
  LogInfo("Deploy()");
  return true;
}

function AfterDeploy() {
  LogInfo("AfterDeploy()");
  return true;
}

function OnError() {
  LogError("Shouldn't read this");
}
