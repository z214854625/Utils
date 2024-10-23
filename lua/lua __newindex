local t = {}
local newIndexTable = {[1] = "newindextable"}
local metatable = {
	__newindex = newIndexTable,
}
setmetatable(t, metatable)
t[1] = 2

print("1--", t[1])     -- nil
print("2--", metatable[1]) -- nil
print("3--", newIndexTable[1]) -- 2
