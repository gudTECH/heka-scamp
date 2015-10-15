local cjson = require "cjson"
local string = require "string"

-- code "borrowed" from https://summitroute.com/blog/2015/06/14/shipping_windows_events_to_heka_and_elasticsearch/#get-events-into-elasticsearch

function process_message()
  local payload = read_message("Payload")

  -- check length (cjson.encode will crash if payload is > 11500 characters)
  -- TODO: Xavier is not sure where this potential crash info came from, couldn't
  -- find any info elsewhere
  -- if #payload > 11500 then
  --    return -1
  -- end

  -- We just need to be sure the payload is valid JSON
  local ok, json = pcall(cjson.decode, payload)

  if not ok then
    return -1
  end

  for k,v in pairs(json) do
    if type(v) == "table" then
        write_message("Fields[" .. k .. "]", cjson.encode(v))
    else
      if k == "type" or k == "Type" then
        local rawtype = v

        dotstart,_dotend = string.find(rawtype, "%.")
        if dotstart == nil then
          return -1
        end
        
        local index = rawtype:sub(0,dotstart-1)
        local type = rawtype:sub(dotstart+1)

        write_message("Type", type)
        write_message("Hostname", index) -- Had to stash it somewhere unused
      else
        write_message("Fields[" .. k .. "]", v)
      end
    end
  end

  return 0
end