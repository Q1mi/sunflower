
-- KEYS[1]: 用户签到key
-- ARGV[1]: 当前日偏移量(从年初开始的天数)
-- ARGV[2]: 提醒阈值

local offset = tonumber(ARGV[1])
local threshold = tonumber(ARGV[2])

-- 检查今天是否已签到(最新位)
local today = redis.call('GETBIT', KEYS[1], offset)
if today == 1 then return 0 end -- 已签到无需提醒

-- 检查最近threshold天的签到情况
local continuous = true
for i = 1, threshold do
    local bit = redis.call('GETBIT', KEYS[1], offset - i)
    if bit ~= 1 then
        continuous = false
        break
    end
end

-- 返回结果: 1需要提醒 0不需要
return continuous and 1 or 0
