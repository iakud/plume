#ifndef _STACK_H_
#define _STACK_H_

#include <tolua/tolua.h>
#include <tolua/tolua_call.h>
#include <tolua/tolua_function.h>
#include <string.h>

void AddPackagePath(lua_State *L, const char* path) {
	lua_getglobal(L, LUA_LOADLIBNAME);
	lua_getfield(L, -1, "path");
	lua_pushfstring(L, "%s;%s/?.lua", lua_tostring(L, -1), path);
	lua_setfield(L, -3, "path");
	lua_pop(L, 2);
}

void Unload(lua_State *L, const char* modname) {
	lua_getglobal(L, LUA_LOADLIBNAME);
	lua_getfield(L, -1, "loaded");
	lua_getfield(L, -1, modname);
	if (!lua_isnil(L, -1)) {
		lua_pushnil(L);
		lua_setfield(L, -3, modname);
	}
	lua_pop(L, 3);
}

lua_Integer Lua_tointeger(lua_State *L, int idx) {
	return lua_tointeger(L, idx);
}

lua_Number Lua_tonumber(lua_State *L, int idx) {
	return lua_tonumber(L, idx)	;
}

void Lua_insert(lua_State *L, int idx) {
	lua_insert(L, idx);
}

void Lua_pop(lua_State *L, int idx) {
	lua_pop(L, idx);
}

int LuaL_error(lua_State *L, const char *s) {
	return luaL_error(L, s);
}

#endif // _STACK_H_