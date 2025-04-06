wrk.method = "GET"
local headers = {}
headers["Content-Type"] = "application/json"
headers["Cookie"] =
"shortugo=MTc0Mzg3NjM3NXwxc254dHo2b3NRd2t2OXB3R1JGeVUzUGt5Mk80Z2F3ZWgyVVlsclZKVmtzRmhZOHl8DmVfNXB015eCtBY4Rmf2elFRLKy1hYV5mblMS6JmJ0I=; Path=/; HttpOnly;"

function request()
    local body = ""
    return wrk.format("GET", "/api/user/urls", headers, body)
end
