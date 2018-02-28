function BeforeDeploy() {
  LogInfo("Lets begin");
  return true;
}

function Deploy() {
  Sleep(5);
  return true;
}

function AfterDeploy() {
  LogInfo("Should fire...");
  return true;
}