g_modVar = "modValue"

g_stuff = 0
st = "1234545610000"
function doStuff()
  g_stuff = g_stuff + 1
 tostring(g_stuff)
  if string.find(st, tostring(g_stuff)) then
    print("found!", g_stuff)
  end
end
