package main

import "braille"
import "fmt"
import "os"
func main() {
s := os.Args[1]
bf,_ := braille.FieldFromString(s)
fmt.Println(string(bf.ToUnicode()))
}