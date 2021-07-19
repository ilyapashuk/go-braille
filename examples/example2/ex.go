package main

import "github.com/ilyapashuk/go-braille"
import "os"
import "bufio"
func main() {
inf,_ := os.Open(os.Args[1])
of,_ := os.Create(os.Args[2])
ifr := bufio.NewReader(inf)
defer inf.Close()
defer of.Close()
for {
r,_,err := ifr.ReadRune()
if err != nil {
return
}
if r == '\r' || r == '\n' {
continue
}
f,err := braille.FieldFromUnicode(r)
if err != nil {
continue
}
s := f.String()
s += "\n"
of.Write([]byte(s))
}
}