function BeforeDeploy() {
  LogInfo("Getting started");
  Halt();
  return true;
}

function Deploy() {
  return true;
}

function AfterDeploy() {
  return true;
}

function OnError() {
  LogError("Shouldn't read this");
}
