/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cl

import (
	"go/token"
	"go/types"

	"github.com/goplus/gogen"
)

// -----------------------------------------------------------------------------

func initMathBig(_ *gogen.Package, conf *gogen.Config, big gogen.PkgRef) {
	conf.UntypedBigInt = big.Ref("UntypedBigint").Type().(*types.Named)
	conf.UntypedBigRat = big.Ref("UntypedBigrat").Type().(*types.Named)
	conf.UntypedBigFloat = big.Ref("UntypedBigfloat").Type().(*types.Named)
}

type Builtin struct {
	types.Object
	name string
	pkg  string
	sym  string
}

func (t *Builtin) Parent() *types.Scope {
	return Universe
}
func (t *Builtin) Pos() token.Pos {
	return 0
}
func (t *Builtin) Pkg() *types.Package {
	return nil
}
func (t *Builtin) Name() string {
	return t.name
}
func (t *Builtin) Type() types.Type {
	return types.Typ[types.Invalid]
}
func (t *Builtin) Exported() bool {
	return false
}
func (t *Builtin) Id() string {
	return "_." + t.name
}
func (t *Builtin) String() string {
	return "builtin_" + t.name
}
func (t *Builtin) Sym() string {
	return t.pkg + "." + t.sym
}

var (
	Universe *types.Scope
)

var builtinDefs = [...]struct {
	name string
	pkg  string
	sym  string
}{
	{"bigint", "github.com/goplus/gop/builtin/ng", ""},
	{"bigrat", "github.com/goplus/gop/builtin/ng", ""},
	{"bigfloat", "github.com/goplus/gop/builtin/ng", ""},
	{"int128", "github.com/goplus/gop/builtin/ng", ""},
	{"uint128", "github.com/goplus/gop/builtin/ng", ""},
	{"lines", "github.com/goplus/gop/builtin/iox", ""},
	{"blines", "github.com/goplus/gop/builtin/iox", "BLines"},
	{"newRange", "github.com/goplus/gop/builtin", "NewRange__0"},
	{"echo", "fmt", "Println"},
	{"print", "fmt", ""},
	{"println", "fmt", ""},
	{"printf", "fmt", ""},
	{"errorf", "fmt", ""},
	{"fprint", "fmt", ""},
	{"fprintln", "fmt", ""},
	{"sprint", "fmt", ""},
	{"sprintln", "fmt", ""},
	{"sprintf", "fmt", ""},
	{"open", "os", ""},
	{"create", "os", ""},
	{"type", "reflect", "TypeOf"},
}

type defSym struct {
	name string
	sym  string
}

var (
	builtinSym map[string][]defSym
)

func init() {
	Universe = types.NewScope(nil, 0, 0, "universe")
	builtinSym = make(map[string][]defSym)
	for _, def := range builtinDefs {
		if def.sym == "" {
			def.sym = string(def.name[0]-('a'-'A')) + def.name[1:]
		}
		builtinSym[def.pkg] = append(builtinSym[def.pkg], defSym{name: def.name, sym: def.sym})
		obj := &Builtin{name: def.name, pkg: def.pkg, sym: def.sym}
		Universe.Insert(obj)
	}
}

func initBuiltin(pkg *gogen.Package, builtin *types.Package) {
	scope := builtin.Scope()
	for im, defs := range builtinSym {
		if p := pkg.TryImport(im); p.Types != nil {
			for _, def := range defs {
				obj := p.Ref(def.sym)
				if _, ok := obj.Type().(*types.Named); ok {
					scope.Insert(types.NewTypeName(token.NoPos, builtin, def.name, obj.Type()))
				} else {
					scope.Insert(gogen.NewOverloadFunc(token.NoPos, builtin, def.name, obj))
				}
			}
		}
	}
	scope.Insert(types.NewTypeName(token.NoPos, builtin, "any", gogen.TyEmptyInterface))
}

func newBuiltinDefault(pkg *gogen.Package, conf *gogen.Config) *types.Package {
	builtin := types.NewPackage("", "")
	ng := pkg.TryImport("github.com/goplus/gop/builtin/ng")
	strx := pkg.TryImport("github.com/qiniu/x/stringutil")
	stringslice := pkg.TryImport("github.com/goplus/gop/builtin/stringslice")
	pkg.TryImport("strconv")
	pkg.TryImport("strings")
	if ng.Types != nil {
		initMathBig(pkg, conf, ng)
	}
	initBuiltin(pkg, builtin)
	gogen.InitBuiltin(pkg, builtin, conf)
	if strx.Types != nil {
		ti := pkg.BuiltinTI(types.Typ[types.String])
		ti.AddMethods(
			&gogen.BuiltinMethod{Name: "Capitalize", Fn: strx.Ref("Capitalize")},
		)
	}
	if stringslice.Types != nil {
		ti := pkg.BuiltinTI(types.NewSlice(types.Typ[types.String]))
		ti.AddMethods(
			&gogen.BuiltinMethod{Name: "Capitalize", Fn: stringslice.Ref("Capitalize")},
			&gogen.BuiltinMethod{Name: "ToTitle", Fn: stringslice.Ref("ToTitle")},
			&gogen.BuiltinMethod{Name: "ToUpper", Fn: stringslice.Ref("ToUpper")},
			&gogen.BuiltinMethod{Name: "ToLower", Fn: stringslice.Ref("ToLower")},
			&gogen.BuiltinMethod{Name: "Repeat", Fn: stringslice.Ref("Repeat")},
			&gogen.BuiltinMethod{Name: "Replace", Fn: stringslice.Ref("Replace")},
			&gogen.BuiltinMethod{Name: "ReplaceAll", Fn: stringslice.Ref("ReplaceAll")},
			&gogen.BuiltinMethod{Name: "Trim", Fn: stringslice.Ref("Trim")},
			&gogen.BuiltinMethod{Name: "TrimSpace", Fn: stringslice.Ref("TrimSpace")},
			&gogen.BuiltinMethod{Name: "TrimLeft", Fn: stringslice.Ref("TrimLeft")},
			&gogen.BuiltinMethod{Name: "TrimRight", Fn: stringslice.Ref("TrimRight")},
			&gogen.BuiltinMethod{Name: "TrimPrefix", Fn: stringslice.Ref("TrimPrefix")},
			&gogen.BuiltinMethod{Name: "TrimSuffix", Fn: stringslice.Ref("TrimSuffix")},
		)
	}
	return builtin
}

// -----------------------------------------------------------------------------
