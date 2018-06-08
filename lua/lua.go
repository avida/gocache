package lua

import (
  "github.com/yuin/gopher-lua"
  "log"
  "errors"
  "strings"
)

type LuaVM struct {
  vm *lua.LState
  functions map[string]lua.LValue
  variables map[string]lua.LValue
}

var global *LuaVM

func InitVM() error {
  if global != nil {
    return errors.New("Lua vm already initialized")
  }
  global = &LuaVM{ vm: lua.NewState(),
                    functions: make(map[string] lua.LValue),
                    variables: make(map[string] lua.LValue)}
  return nil
}

func LoadScript(filename string) error {
  fn, err := global.vm.LoadFile(filename)
  if err != nil {
    return err
  }
  global.vm.Push(fn)
  return global.vm.PCall(0, lua.MultRet, nil)
}

func LoadVars() {
  tb := global.vm.ToTable(lua.GlobalsIndex)
  tb.ForEach(func(k lua.LValue, val lua.LValue){
      key := k.String()
      if val.Type() == lua.LTFunction {
        if strings.Index(key, "do") == 0{
          global.functions[key] = val
        }
      } else
      {
        if strings.Index(key, "g_") == 0 {
          global.variables[key] = val
        }
      }
    })
}

func PrintVMInfo() {
  tb := global.vm.ToTable(lua.GlobalsIndex)
  tb.ForEach(func(key lua.LValue, val lua.LValue){
    log.Printf("var is : %v - %v", key.String() , val.Type() )
  })
}

func Print(t *lua.LTable, indent int) {
  indt := strings.Repeat("  ", indent) 
  t.ForEach(func(key lua.LValue, val lua.LValue){
      if val.Type() == lua.LTTable{
        log.Printf("%stable %s", indt, key)
        Print(val.(*lua.LTable), indent + 1)
      } else {
      log.Printf("%s%s: %s ", indt, key, val )
      }
  })
}

func PrintGlobalVars() {
  for name, val := range(global.variables) {
      log.Printf("var %s has type of %s, val: %s ", name , val.Type(), val.String() )
  }
}

func CallFunc(name string) error {
  f, found := global.functions[name]
  if ! found {
    return errors.New("Function not found")
  }
  err := global.vm.CallByParam(lua.P{
    Fn: f,
    NRet: 1,
    Protect: true,})
  if err != nil {
    return err
  }
  global.vm.Pop(1)
  return err
}

func GetVar(name string) (string, error) {
    if _, found := global.variables[name]; found == true{
      //return val.String(), nil
      return global.vm.GetGlobal(name).String(), nil
    } else
    {
      return "", errors.New("Variable not found")
    }

}
