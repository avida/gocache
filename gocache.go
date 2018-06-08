package main

import (
  "log"
  "flag"
  "./lua"
  "fmt"
  "time"
)

type FilesList []string

func (list *FilesList) String() string {
  return ""
}
func (list *FilesList) Set(value string) error{
  *list = append(*list, value)
  return nil
}

func main() {
  var scripts FilesList
  flag.Var(&scripts, "f", "Lua script to executes")
  flag.Parse()
  if len(scripts) == 0 {
    fmt.Println("Please specify script to load")
    return
  }
  lua.InitVM()
  for _, script := range(scripts){
    if err := lua.LoadScript(script); err != nil {
      log.Printf("Error loading script %s: %s", script, err)
      return
    }
  }
  lua.LoadVars()
  lua.PrintGlobalVars()
  start := time.Now()
  for i:= 0; i < 10000; i++ {
    if err := lua.CallFunc("doStuff"); err != nil {
      log.Printf("error calling func %v", err)
    }
  }
  elapsed := time.Since(start)
  log.Printf("Call took %s", elapsed)
  if val, err := lua.GetVar("g_stuff"); err != nil{
    log.Printf("Error occured: %s", err)
  } else {
    log.Printf("val is: %s", val)
  }

}
