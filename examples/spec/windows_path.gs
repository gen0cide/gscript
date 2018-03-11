function BeforeDeploy() {
  LogInfo("BeforeDeploy()");
  return true;
}

function Deploy() {
  PathTester("C:/Windows/system32/calc.exe");
  return true;
}

function AfterDeploy() {
  LogInfo("AfterDeploy()");
  return true;
}
