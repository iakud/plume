package lua

/*
typedef struct lua_State lua_State;

// The gateway function
int lua_go_log_info(lua_State*);
int lua_go_time_now(lua_State*);
*/
import "C"

import (
	"runtime"
	"time"

	"github.com/iakud/plume/log"
)

func FuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	return runtime.FuncForPC(pc[0]).Name()
}

func OpenPlume(L *Stack) {
	lua_register(L)
}

func lua_register(L *Stack) {
	L.BeginModule("")
	{
		lua_register_go(L)
	}
	L.EndModule()
}

func lua_register_go(L *Stack) {
	L.Module("plume")
	L.BeginModule("plume") // _G.go
	{
		lua_register_go_log(L)
		lua_register_go_time(L)
	}
	L.EndModule()
}

func lua_register_go_log(L *Stack) {
	L.Module("log")
	L.BeginModule("log") // _G.go.log
	{
		L.Function("info", (CFunction)(C.lua_go_log_info))
	}
	L.EndModule()
}

//export lua_go_log_info
func lua_go_log_info(l *C.lua_State) C.int {
	L := (*Stack)(l)
	argc := L.GetTop()
	if argc == 1 {
		info := L.ToLString(-1)
		log.Infof(info)
		return 1
	}
	return C.int(L.Error("'%s' has wrong number of arguments: %d, was expecting %d \n", FuncName(), argc, 1))
}

func lua_register_go_time(L *Stack) {
	L.Module("time")
	L.BeginModule("time") // _G.go.time
	{
		L.Function("now", (CFunction)(C.lua_go_time_now))
	}
	L.EndModule()
}

//export lua_go_time_now
func lua_go_time_now(l *C.lua_State) C.int {
	L := (*Stack)(l)
	argc := L.GetTop()
	if argc == 0 {
		L.PushInt64(time.Now().UnixNano())
		return 1
	}
	return C.int(L.Error("'%s' has wrong number of arguments: %d, was expecting %d \n", FuncName(), argc, 0))
}