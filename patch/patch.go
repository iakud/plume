package patch

import (
	"errors"
	"plugin"

	gomonkey "github.com/agiledragon/gomonkey/v2"
)

const kPatchName = "Patch"

type Patch interface {
	Apply()
}

var patches = gomonkey.NewPatches()

func Open(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}
	patchSymbol, err := p.Lookup(kPatchName)
	if err != nil {
		return err
	}
	patch, ok := patchSymbol.(Patch)
	if !ok {
		return errors.New("patch error")
	}
	patch.Apply()
	return nil
}

func Reset() {
	patches.Reset()
}

func ApplyFunc(target, double interface{}) {
	patches.ApplyFunc(target, double)
}

func ApplyMethod(target interface{}, methodName string, double interface{}) {
	patches.ApplyMethod(target, methodName, double)
}

func ApplyMethodFunc(target interface{}, methodName string, doubleFunc interface{}) {
	patches.ApplyMethodFunc(target, methodName, doubleFunc)
}

func ApplyPrivateMethod(target interface{}, methodName string, double interface{}) {
	patches.ApplyPrivateMethod(target, methodName, double)
}

func ApplyGlobalVar(target, double interface{}) {
	patches.ApplyGlobalVar(target, double)
}

func ApplyFuncVar(target, double interface{}) {
	patches.ApplyFuncVar(target, double)
}