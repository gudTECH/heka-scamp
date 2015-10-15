local cjson = require "cjson"

-- code "borrowed" from https://summitroute.com/blog/2015/06/14/shipping_windows_events_to_heka_and_elasticsearch/#get-events-into-elasticsearch

function process_message()
  local payload = read_message("Payload")

  -- check length (cjson.encode will crash if payload is > 11500 characters)
  if #payload > 11500 then
     return -1
  end

  -- We just need to be sure the payload is valid JSON
  local ok, json = pcall(cjson.decode, payload)

  if not ok then
    return -1
  end

  for k,v in pairs(json) do
    if type(v) == "table" then
        write_message("Fields[" .. k .. "]", cjson.encode(v))
    else
        write_message("Fields[" .. k .. "]", v)
    end
  end

  return 0
end