package obfuscator

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"math/rand"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/gen0cide/gscript/compiler/computil"
	"github.com/gen0cide/gscript/engine"
)

const (
	GENERICMATCH = `[[:word:]\.\/\(\)\*\[\]\&\%\{\}\$ ]*`
)

var (
	pathTest = regexp.MustCompile(`/`)

	reservedWords = map[string]bool{
		"init":         true,
		"execute":      true,
		"__PRELOAD":    true,
		"__ENTRYPOINT": true,
		"root":         true,
		"home":         true,
		"":             true,
		"Users":        true,
	}

	importantMatches = []string{
		`github\.com/gen0cide/gscript/[[:word:]\.\/\(\)\*\[\]\&\%\{\}\$ ]*`,
		`github\.com/robertkrimen/otto/[[:word:]\.\/\(\)\*\[\]\&\%\{\}\$ ]*`,
		`github\.com/gen0cide/gscript`,
		`github\.com/robertkrimen/otto`,
		`github\.com/gen0cide/`,
		`github\.com/robertkrimen/`,
		`/home/[[:word:]\.\\/ ]*`,
		`/Users/[[:word:]\.\\/ ]*`,
		`/root/[[:word:]\.\\/ ]*`,
		`/tmp/[[:word:]\.\\/ ]*`,
		`/usr/local/[[:word:]\.\\/ ]*`,
		`gopkg[[:word:]\.\\/ ]*`,
		`UPX[[:word:]\.\\/ ]*`,
		`$Info[[:word:]\.\\/ ]*`,
		`github[[:word:]\.\\/ ]*`,
		`google[[:word:]\.\\/ ]*`,
	}

	softDefaults = strings.Split(string(computil.MustAsset("soft_reserved")), "\n")
	hardDefaults = strings.Split(string(computil.MustAsset("hard_reserved")), "\n")
)

type Mordor struct {
	sync.RWMutex
	Horde  map[string]*Orc
	Dead   map[string]bool
	Logger engine.Logger
}

type Orc struct {
	Name   string
	Hits   int
	Filter *regexp.Regexp
}

func NewMordor(l engine.Logger) *Mordor {
	m := &Mordor{
		Horde:  map[string]*Orc{},
		Dead:   map[string]bool{},
		Logger: l,
	}
	for _, im := range importantMatches {
		m.AddSingleGhostLiteral(im)
	}
	libDir, err := computil.ResolveEngineLibrary()
	if err == nil {
		ghosts, err := WalkEngineLibForGhosts(libDir, "engine")
		if err == nil {
			m.AddGhosts(ghosts)
		}
	}
	for _, d := range hardDefaults {
		m.AddSingleGhostLiteral(d)
	}
	m.AddGhosts(softDefaults)
	return m
}

func buildFilter(s string) string {
	return strings.Join([]string{s, GENERICMATCH}, "")
}

func WalkEngineLibForGhosts(dirpath string, pkgName string) ([]string, error) {
	ghosts := []string{}
	fs := token.NewFileSet()
	pkgGlob, err := parser.ParseDir(fs, dirpath, nil, parser.ParseComments)
	if err != nil {
		return ghosts, err
	}
	pkg := pkgGlob[pkgName]
	if pkg == nil {
		return ghosts, fmt.Errorf("no package named %s located within that path", pkgName)
	}
	exists := ast.PackageExports(pkg)
	if exists != true {
		return ghosts, nil
	}
	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			if gd, ok := d.(*ast.GenDecl); ok && gd.Tok == token.TYPE {
				for _, s := range gd.Specs {
					if ts, ok := s.(*ast.TypeSpec); ok {
						tn := ts.Name.Name
						ghosts = append(ghosts, fmt.Sprintf("%s\\.%s", pkgName, tn))
						ghosts = append(ghosts, fmt.Sprintf("\\*%s\\.%s", pkgName, tn))
					}
				}
			}
		}
	}
	return ghosts, nil
}

func (m *Mordor) AddGhosts(g []string) {
	for _, e := range g {
		if m.Dead[e] == true || m.Horde[e] != nil {
			continue
		}
		if err := m.AddSingleGhost(e); err != nil {
			m.Logger.Errorf("Error creating ghost %s: %v", e, err)
			continue
		}
		if pathTest.MatchString(e) {
			b := path.Base(e)
			if err := m.AddSingleGhost(b); err != nil {
				continue
			}
		}
	}
}

func (m *Mordor) AddSingleGhostLiteral(g string) error {
	if reservedWords[g] == true {
		m.Lock()
		m.Dead[g] = true
		m.Unlock()
		return errors.New("reserved word cannot be used as ghost")
	}
	r, err := regexp.Compile(g)
	if err != nil {
		m.Lock()
		m.Dead[g] = true
		m.Unlock()
		return err
	}
	orc := &Orc{
		Name:   g,
		Hits:   0,
		Filter: r,
	}
	m.Lock()
	m.Horde[g] = orc
	m.Unlock()
	return nil
}

func (m *Mordor) AddSingleGhost(g string) error {
	if reservedWords[g] == true {
		m.Lock()
		m.Dead[g] = true
		m.Unlock()
		return errors.New("reserved word cannot be used as ghost")
	}
	fn := buildFilter(g)
	r, err := regexp.Compile(fn)
	if err != nil {
		m.Lock()
		m.Dead[g] = true
		m.Unlock()
		return err
	}
	orc := &Orc{
		Name:   g,
		Hits:   0,
		Filter: r,
	}
	m.Lock()
	m.Horde[fn] = orc
	m.Unlock()
	return nil
}

func (m *Mordor) PrintStats() {
	for n, o := range m.Horde {
		if o.Hits == 0 {
			continue
		}
		m.Logger.Infof("Mordor Count: %d - %s", o.Hits, n)
	}
}

func (m *Mordor) Assault(srcFile string) error {
	data, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}
	for k, o := range m.Horde {
		data = o.Filter.ReplaceAllFunc(data, func(b []byte) []byte {
			m.Horde[k].Hits++
			for i := range b {
				b[i] = byte(rand.Int() % 256)
			}
			return b
		})
	}
	err = ioutil.WriteFile(srcFile, data, 0755)
	return err
}
