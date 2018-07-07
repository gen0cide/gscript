//go_import:github.com/gen0cide/gscript/xdeprecated/cmd/gcomp/testlib as testlib
//go_import:$IMPORT_KEY as $NAMESPACE

function Deploy() {
  param1 = "hello";
  param2 = "world";
  param3 = "this is longggg";
  ret = testlib.Test1("http://google.com/search");
  console.log(ret[0].Host);
}


// dan eats piggy shit