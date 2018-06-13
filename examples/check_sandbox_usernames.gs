function BeforeDeploy() {
  console.log("CheckSandboxUsernames test starting")
  return true
}

function Deploy() {
  obj = CheckSandboxUsernames()
  console.log("areWeInASandbox?: " + obj.areWeInASandbox)
  return true
}

function AfterDeploy() {
  console.log("CheckSandboxUsernames test complete")
  return true
}

