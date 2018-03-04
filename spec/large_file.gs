//import:/Users/flint/Downloads/e200.tar

function BeforeDeploy() {
  LogInfo("BeforeDeploy()");
  return true;
}

function Deploy() {
  WriteFile("/Users/flint/Downloads/e301.tar", Asset("e200.tar"), "0644");
  return true;
}

function AfterDeploy() {
  LogInfo("AfterDeploy()");
  return true;
}
