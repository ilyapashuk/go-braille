// this code is licensed under mit license, see license file in repository root
// Copyright 2021 ilyapashuk


package braille

//this file contains code for brf support

import "bytes"
import _ "embed"
import "io"
import "strings"
//go:embed brf.dat
var brfdata []byte
var brfback map[rune]BrailleField
var brfforw map[BrailleField]rune

func init() {
brfback = make(map[rune]BrailleField)
brfforw = make(map[BrailleField]rune)
brfforw[BrailleField(0)] = ' '
brfback[' '] = BrailleField(0)
brfr := bytes.NewReader(brfdata)
for {
r,size,err := brfr.ReadRune()
if err != nil {
if err == io.EOF {
break
} else {
panic(err)
}
}
if size != 1 {
panic("invalid internal brf table format")
}
bb,err := brfr.ReadByte()
if err != nil {
panic(err)
}
brfback[r] = BrailleField(bb)
}
for i,v := range brfback {
brfforw[v] = i
}
}

// gets braille ascii representation of this braille field. dots 7 and 8 will be stripped.
func (f BrailleField) ToBrf() rune {
ff := f
if ff.Dot7() {
ff = ff ^ BrailleField(Dot7)
}
if ff.Dot8() {
ff = ff ^ BrailleField(Dot8)
}
return brfforw[ff]
}

//converts this braille row to brf
func (c BrailleRow) ToBrf() []byte {
sb := new(strings.Builder)
for _,f := range c {
sb.WriteRune(f.ToBrf())
}
return []byte(sb.String())
}

func (c BraillePage) ToBrf() []byte {
sb := new(strings.Builder)
for _,row := range c {
sb.Write(row.ToBrf())
sb.WriteRune('\n')
}
return []byte(sb.String())
}

// creates braille field from it's braille ascii representation
func FieldFromBrf(r rune) (BrailleField, error) {
if bf,ok := brfback[r]; ok {
return bf,nil
} else {
return BrailleField(0),InvalidBrailleError
}
}


func RowFromBrf(bdata []byte) (BrailleRow, error) {
res := make(BrailleRow, len(bdata))
for i,v := range bdata {
bf,err := FieldFromBrf(rune(v))
if err != nil {
return res,InRowError{i, err}
}
res[i] = bf
}
return res,nil
}

func PageFromBrf(bdata []byte) (BraillePage, error) {
bl := strings.Split(string(bdata), "\n")
res := make(BraillePage, len(bl))
for i,v := range bl {
row,err := RowFromBrf([]byte(v))
if err != nil {
return res,InTextError{i, err}
}
res[i] = row
}
return res,nil
}