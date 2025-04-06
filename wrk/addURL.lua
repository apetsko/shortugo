wrk.method = "POST"
local headers = {}
headers["Content-Type"] = "application/json"
headers["Cookie"] = "shortugo=MTc0Mzg3NjM3NXwxc254dHo2b3NRd2t2OXB3R1JGeVUzUGt5Mk80Z2F3ZWgyVVlsclZKVmtzRmhZOHl8DmVfNXB015eCtBY4Rmf2elFRLKy1hYV5mblMS6JmJ0I=; Path=/; HttpOnly;"

local counter = 1400000

function request()
    counter = counter + 1
    local body = string.format('{"url":"http://example.com/api/%d"}', counter)
    return wrk.format("POST", "/api/shorten", headers, body)
end