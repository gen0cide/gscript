package hades

import (
	"github.com/gen0cide/gscript"
)

func Run() {
	gse := gscript.New("")
	gse.CreateVM()
	gse.LoadScript(MustAsset("genesis_entry_point.gs"))
	gse.ExecutePlan()
}
