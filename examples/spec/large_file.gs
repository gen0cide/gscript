//import:/Users/flint/Downloads/e200.tar

function BeforeDeploy() {
  LogInfo("BeforeDeploy()");
  return true;
}

function Deploy() {
  return WriteFile("Z:/Public/e200.tar", Asset("e200.tar"), "0644");
}

function AfterDeploy() {
  LogInfo("AfterDeploy()");
  return true;
}
