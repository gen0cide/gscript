// a new script
//go_import:github.com/tdewolff/minify as minify
//go_import:github.com/tdewolff/minify/js as js
var foo = "bar";

function poop() {
  return null;
}

function Deploy() {
  poop();
  minfier = minify.New()
  minifier.AddFunc("text/javascript", js.Minify)
  return true;
}

function AfterDeploy() {
  return true;
}