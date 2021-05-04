package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "syscall"
)

func main() {

    binary, lookErr := exec.LookPath(os.Args[1])
    if lookErr != nil {
        panic(lookErr)
    }

    env := os.Environ()
    var args []string

    for _, v := range os.Args[2:] {
        if strings.Contains(v, " ") {
            v = `'` + v + `'`
        }
        //fmt.Printf("DEBUG: Command arg[%d] %s\r\n", i, v)
        args = append(args, v)
    }

    //args = []string{"-c", "'echo HEY'"}
    fmt.Printf("Command: %s\n", binary)
    fmt.Printf("Args: %s\n", strings.Join(args, " "))

    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        panic(execErr)
    }
}