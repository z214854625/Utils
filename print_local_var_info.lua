
--[[
	通过debug.getlocal(1, 1) 可以打印局部变量的详细信息 可以重载一下全局的tostring()函数将table详细信息打印出来
	也可以统计 BasicAndDiamondGift中子表的节点数和层级
]]--
local function show1()
    local BasicAndDiamondGift = require "BasicAndDiamondGift"
    local val1, name1 = debug.getlocal(1, 1)
    print("first--", val1, type(name1))
    for k1, v1 in pairs(name1) do
        print("sub--", k1, type(k1), v1, type(v1))
    end
end
show1()