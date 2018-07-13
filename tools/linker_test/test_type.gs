//go_import:github.com/gen0cide/gscript/tools/linker_test/typelib as typelib

function Deploy() {
  val = 12345
  t = typelib.ReturnHandle(val)
  check(t)
}