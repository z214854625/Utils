
-- 创建一个基类 test
local test = {}

-- 不设置 __index
--test.__index = test  -- 这行被注释掉了

function test:new(o)
    -- 设置元表
    --return setmetatable(o or {}, test)
,return setmetatable(o or {}, {__index = self})
end

function test:sayHello()
    print("Hello from test!")
end

-- 使用 test 类创建实例
local obj = test:new()

-- 尝试调用 sayHello 方法
obj:sayHello()  -- 这将引发一个错误，因为 sayHello 在 obj 中无法找到


