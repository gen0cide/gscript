//go_import:github.com/deckarep/gosx-notifier as notifier
//go_import:github.com/atotto/clipboard as clipboard
//go_import:github.com/gen0cide/gscript/ternals/linker_test/testlib as testlib
//go_import:github.com/bettercap/bettercap/network as betnet
//go_import:github.com/bettercap/bettercap/modules as betmod
//go_import:github.com/bettercap/bettercap/session as betsesh

function Deploy() {
  param1 = "hello";
  param2 = "world";
  param3 = "this is longggg";
  hotmama = notifier.NewNotification(param1 + " " + param2 + " " + param3);
  da_error = hotmama.Push();
  console.log(da_error);
  clippy = clipboard.ReadAll();
  console.log(clippy[0]);
  urlObj = testlib.Test1("https://bing.com/search");
  console.warn(urlObj[0].Host);
  DebugConsole();
}