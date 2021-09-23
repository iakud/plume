#ifndef _STACK_H_
#define _STACK_H_

#include <tolua/tolua.h>
#include <tolua/tolua_call.h>
#include <tolua/tolua_function.h>
#include <string.h>

void stack_addPackagePath(lua_State *L, const char* path) {
	lua_getglobal(L, LUA_LOADLIBNAME);
	lua_getfield(L, -1, "path");
	lua_pushfstring(L, "%s;%s/?.lua", lua_tostring(L, -1), path);
	lua_setfield(L, -3, "path");
	lua_pop(L, 2);
}

void unload(lua_State *L, const char* modname) {
	if (modname == NULL || strlen(modname) == 0) {
		return;
	}

	lua_getglobal(L, LUA_LOADLIBNAME);
	lua_getfield(L, -1, "loaded");
	lua_pushstring(L, modname);
	lua_gettable(L, -2);
	if (!lua_isnil(L, -1)) {
		lua_pushstring(L, modname);
		lua_pushnil(L);
		lua_settable(L, -4);
	}
	lua_pop(L, 3);
}

int stack_isnil(lua_State *L, int idx) {
	return lua_isnil(L, idx);
}

lua_Integer stack_tointeger(lua_State *L, int idx) {
	return lua_tointeger(L, idx);
}

lua_Number stack_tonumber(lua_State *L, int idx) {
	return lua_tonumber(L, idx)	;
}

void stack_insert(lua_State *L, int idx) {
	lua_insert(L, idx);
}

void stack_pop(lua_State *L, int n) {
	lua_pop(L, n);
}

#endif // _STACK_H_