wrk.method = "POST"
local headers = {}
headers["Content-Type"] = "application/json"
headers["Cookie"] =
"shortugo=MTc0Mzg3NjM3NXwxc254dHo2b3NRd2t2OXB3R1JGeVUzUGt5Mk80Z2F3ZWgyVVlsclZKVmtzRmhZOHl8DmVfNXB015eCtBY4Rmf2elFRLKy1hYV5mblMS6JmJ0I=; Path=/; HttpOnly;"

local counter = 1400000

function request()
    n1 = counter + 1
    n2 = counter + 2
    n3 = counter + 3
    counter = counter + 3
    local body = string.format('["%d","%d","%d"]', n1, n2, n3)
    return wrk.format("DELETE", "/api/user/urls", headers, body)
end
