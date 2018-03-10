function BeforeDeploy() {
  LogInfo("BeforeDeploy()");
  return true;
}

function Deploy() {
  ts = Timestamp();
  LogInfo("Deploy()");
  LogDebug("Timestamp: " + ts.value);
  return true;
}

function AfterDeploy() {
  LogInfo("AfterDeploy()");
  return true;
}
