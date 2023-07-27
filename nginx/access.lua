local cjson = require "cjson"
local httpc = require("resty.http").new()

ngx.log(ngx.ERR, ngx.req.get_headers()["accessToken"])

if not ngx.req.get_headers()["accessToken"] then
    ngx.exit(ngx.HTTP_UNAUTHORIZED)
    return
end
local res, err = httpc:request_uri("http://app:8000/validate", {
    method = "POST",
    body = "{\"accessToken\":\"" .. ngx.req.get_headers()["accessToken"] .. "\"}",
    headers = {
        ["Content-Type"] = "application/json",
    },
})
if (not res or res.status ~= 200) then
    ngx.exit(ngx.HTTP_UNAUTHORIZED)
    return
end
local body = cjson.decode(res.body)
ngx.req.set_header("identityToken", body["identityToken"])