package lua

import (
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -llua -ltolua
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
	C.AddPackagePath(L, cPath)
}

func (L *Stack) Load(modname string) {
	if len(modname) == 0 {
		return
	}
	require := fmt.Sprintf("require '%v'", modname)
	L.ExecuteString(require)
}

func (L *Stack) Unload(modname string) {
	if len(modname) == 0 {
		return
	}
	cModname := C.CString(modname)
	defer C.free(unsafe.Pointer(cModname))
	C.Unload(L, cModname)
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
	return int(C.Lua_tointeger(L, C.int(index)))
}

func (L *Stack) ToInt32(index int) int32 {
	return int32(C.Lua_tointeger(L, C.int(index)))
}

func (L *Stack) ToInt64(index int) int64 {
	return int64(C.Lua_tonumber(L, C.int(index)))
}

func (L *Stack) ToFloat32(index int) float32 {
	return float32(C.Lua_tonumber(L, C.int(index)))
}

func (L *Stack) ToFloat64(index int) float64 {
	return float64(C.Lua_tonumber(L, C.int(index)))
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

func (L *Stack) Pop(n int) {
	C.Lua_pop(L, C.int(n))
}

//
// excute
//
func (L *Stack) ExecuteGlobalFunction(funcname string, nargs, nresults int) {
	cFuncname := C.CString(funcname)
	defer C.free(unsafe.Pointer(cFuncname))
	C.lua_getglobal(L, cFuncname)
	if nargs > 0 {
		C.Lua_insert(L, C.int(-(nargs + 1)))
	}
	L.execute(nargs, nresults)
}

func (L *Stack) ExecuteFunction(f *C.tolua_FunctionRef, nargs, nresults int) {
	C.tolua_pushfunction_ref(L, f)
	if nargs > 0 {
		C.Lua_insert(L, C.int(-(nargs + 1)))
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
	if C.tolua_docall(L, C.int(nargs), C.int(nresults)) != 0 && C.lua_type(L, -1) != C.LUA_TNIL {
		err := L.ToString(-1)
		L.Pop(1)
		panic(err)
	}
}

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

func (L *Stack) EndModule(name string) {
	C.tolua_endmodule(L)
}

type CFunction = C.lua_CFunction

func (L *Stack) Function(name string, f CFunction) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.tolua_function(L, cName, f)
}

func (L *Stack) UserType(name string, col CFunction) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.tolua_usertype(L, cName, col)
}

func (L *Stack) Class(lname, name, base string) {
	cLName := C.CString(lname)
	defer C.free(unsafe.Pointer(cLName))
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cBase := C.CString(base)
	defer C.free(unsafe.Pointer(cBase))
	C.tolua_class(L, cLName, cName, cBase)
}

func (L *Stack) BeginUserType(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.tolua_beginusertype(L, cName)
}

func (L *Stack) EndUserType(name string) {
	C.tolua_endusertype(L)
}

func (L *Stack) IsUserTable(index int, name string) bool {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.tolua_isusertable(L, C.int(index), cName) == 0
}