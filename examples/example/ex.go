package main

import "github.com/ilyapashuk/go-braille"
import "fmt"
import "os"
func main() {
s := os.Args[1]
bf,_ := braille.FieldFromString(s)
fmt.Println(string(bf.ToUnicode()))
}