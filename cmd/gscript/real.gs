// genesis script

var hello_world;
var hash;

function BeforeDeploy() {
 hello_world = "helloworld";
 return true; 
}

function Deploy() {
  hash = MD5(hello_world);
  console.log(hash);
  return true;
}

function AfterDeploy() {
  return false;
}