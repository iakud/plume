package lua

import (
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -llua -ltolua
#include <tolua/tolua.h>
#include <tolua/tolua_call.h>
#include <tolua/tolua_function.h>
#include <stack.h>
#include <stdlib.h>
 */
import "C"

type Stack = C.lua_State

func NewStack() *Stack {
	L := C.luaL_newstate()
	C.luaL_openlibs(L)
	C.tolua_open(L)
	return L
}

func (L *Stack) Close() {
	C.lua_close(L)
}

func (L *Stack) AddPackagePath(path string) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.stack_addPackagePath(L, cPath)
}

func (L *Stack) Load(modname string) {
	if len(modname) == 0 {
		return
	}
	require := fmt.Sprintf("require '%v'", modname)
	L.ExecuteString(require)
}

func (L *Stack) Unload(modname string) {
	cModname := C.CString(modname)
	defer C.free(unsafe.Pointer(cModname))
	C.unload(L, cModname)
}

func (L *Stack) Reload(modname string) {
	L.Unload(modname)
	L.Load(modname)
}

//
// push value
//
func (L *Stack) PushNil() {
	C.lua_pushnil(L)
}

func (L *Stack) PushBool(value bool) {
	if value {
		C.lua_pushboolean(L, 1)
	} else {
		C.lua_pushboolean(L, 0)
	}
}

func (L *Stack) PushInt(value int) {
	C.lua_pushinteger(L, C.lua_Integer(value))
}

func (L *Stack) PushInt32(value int32) {
	C.lua_pushinteger(L, C.lua_Integer(value))
}

func (L *Stack) PushInt64(value int64) {
	C.lua_pushnumber(L, C.lua_Number(value))
}

func (L *Stack) PushFloat32(value float32) {
	C.lua_pushnumber(L, C.lua_Number(value))
}

func (L *Stack) PushFloat64(value float64) {
	C.lua_pushnumber(L, C.lua_Number(value))
}

func (L *Stack) PushString(value string) {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	C.lua_pushstring(L, cValue)
}

func (L *Stack) PushLString(value string) {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	C.lua_pushlstring(L, cValue, C.size_t(len(value)))
}

func (L *Stack) PushLightUserdata(p unsafe.Pointer) {
	C.lua_pushlightuserdata(L, p)
}

func (L *Stack) PushUserType(p unsafe.Pointer, name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.tolua_pushusertype(L, p, cName)
}

func (L *Stack) PushFunctionRef(f *C.tolua_FunctionRef) {
	C.tolua_pushfunction_ref(L, f)
}

//
// to value
//
func (L *Stack) ToBool(index int) bool {
	return C.lua_toboolean(L, C.int(index)) != 0
}

func (L *Stack) ToInt(index int) int {
	return int(C.stack_tointeger(L, C.int(index)))
}

func (L *Stack) ToInt32(index int) int32 {
	return int32(C.stack_tointeger(L, C.int(index)))
}

func (L *Stack) ToInt64(index int) int64 {
	return int64(C.stack_tonumber(L, C.int(index)))
}

func (L *Stack) ToFloat32(index int) float32 {
	return float32(C.stack_tonumber(L, C.int(index)))
}

func (L *Stack) ToFloat64(index int) float64 {
	return float64(C.stack_tonumber(L, C.int(index)))
}

func (L *Stack) ToString(index int) string {
	return C.GoString(C.lua_tolstring(L, C.int(index), nil))
}

func (L *Stack) ToLString(index int) string {
	var l C.size_t
	cValue := C.lua_tolstring(L, C.int(index), &l)
	return C.GoStringN(cValue, C.int(l))
}

func (L *Stack) ToUserdata(index int) unsafe.Pointer {
	return C.lua_touserdata(L, C.int(index))
}

func (L *Stack) ToUserType(index int, name string) unsafe.Pointer {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.tolua_tousertype(L, C.int(index), cName)
}

func (L *Stack) ToFunctionRef(index int) *C.tolua_FunctionRef {
	return C.tolua_tofunction_ref(L, C.int(index))
}

//
// remove
//
func (L *Stack) RemoveFunctionRef(f *C.tolua_FunctionRef) {
	C.tolua_removefunction_ref(L, f)
}

//
// stack
//
func (L *Stack) GetTop() int {
	return int(C.lua_gettop(L))
}

func (L *Stack) Clean() {
	C.lua_settop(L, 0)
}

func (L *Stack) FormatIndex(index int) int {
	if index < 0 {
		return L.GetTop() + 1 + index
	} else {
		return index
	}
}

//
// excute
//
func (L *Stack) ExecuteGlobalFunction(funcname string, nargs, nresults int) {
	cFuncname := C.CString(funcname)
	defer C.free(unsafe.Pointer(cFuncname))
	C.lua_getglobal(L, cFuncname)
	if nargs > 0 {
		C.stack_insert(L, C.int(-(nargs + 1)))
	}
	L.execute(nargs, nresults)
}

func (L *Stack) ExecuteFunction(f *C.tolua_FunctionRef, nargs, nresults int) {
	C.tolua_pushfunction_ref(L, f)
	if nargs > 0 {
		C.stack_insert(L, C.int(-(nargs + 1)))
	}
	L.execute(nargs, nresults)
}

func (L *Stack) ExecuteString(codes string) {
	cCodes := C.CString(codes)
	defer C.free(unsafe.Pointer(cCodes))
	C.luaL_loadstring(L, cCodes)
	L.execute(0, 0)
}

func (L *Stack) execute(nargs, nresults int) {
	if C.tolua_docall(L, C.int(nargs), C.int(nresults)) != 0 && C.stack_isnil(L, -1) == 0 {
		err := C.lua_tolstring(L, -1, nil)
		C.stack_pop(L, 1)
		panic(err)
	}
}

//
// module
//
/*
func (L *Stack) Module(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.tolua_module(L, cName)
}

func (L *Stack) BeginModule(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.tolua_beginmodule(L, cName)
}

func (L *Stack) EndModule() {
	C.tolua_endmodule((*C.lua_State)(L))
}

func (L *Stack) Function(name string, f C.lua_CFunction) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.tolua_function(L, cName, f)
}

func (L *Stack) UserType(name string, col C.lua_CFunction) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.tolua_usertype(L, cName, col)
}

func (L *Stack) Class(lname, name, base string) {
	c_lname := C.CString(lname)
	defer C.free(unsafe.Pointer(c_lname))
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))
	c_base := C.CString(base)
	defer C.free(unsafe.Pointer(c_base))
	C.tolua_class((*C.lua_State)(L), c_lname, c_name, c_base)
}
*/