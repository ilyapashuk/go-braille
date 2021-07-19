package main

import "os"
import "braille"
import "strings"

func main() {
indata,_ := os.ReadFile(os.Args[1])
instr := strings.ReplaceAll(string(indata), "\r", "")
inpage,_ := braille.PageFromUnicode(instr)
outdata := inpage.ToBrf()
os.WriteFile(os.Args[2], outdata, 0644)
}