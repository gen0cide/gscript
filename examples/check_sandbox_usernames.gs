function BeforeDeploy() {
  LogInfo("CheckSandboxUsernames test starting")
  return true
}

function Deploy() {
  obj = CheckSandboxUsernames()
  LogInfo("areWeInASandbox?: " + obj.areWeInASandbox)
  return true
}

function AfterDeploy() {
  LogInfo("CheckSandboxUsernames test complete")
  return true
}

