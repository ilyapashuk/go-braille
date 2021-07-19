// this code is licensed under mit license, see license file in repository root
// Copyright 2021 ilyapashuk
// this package provides some basic types and functions to implement digital braille processing

package braille

import "strconv"
import "errors"
import "fmt"
import "strings"

// this constants contains dot representations for our internal format
const (
Dot1 = byte(0b00000001)
Dot2 = byte(0b00000010)
Dot3 = byte(0b00000100)
Dot4 = byte(0b00001000)
Dot5 = byte(0b00010000)
Dot6 = byte(0b00100000)
Dot7 = byte(0b01000000)
Dot8 = byte(0b10000000)
)

const UnicodeBrailleBase = rune(0x2800)

// this function check weather given rune belongs to the braille unicode range
func IsBrailleUnicode(r rune) bool {
return !(r < UnicodeBrailleBase || r > rune(0x28FF))
}

var InvalidBrailleError error = errors.New("invalid braille encoding")

// this error type wraps an other error to give information about char number in line or cell number in braille row in which this error occurred
type InRowError struct {
// count starts from zero, so 1 should be added to this number before displaying
Position int
Err error
}
func (c InRowError) Unwrap() error {
return c.Err
}
func (c InRowError) Error() string {
return fmt.Sprintf("in position %v %v", c.Position+1, c.Err)
}

// this error type wraps an other error to give information about line number in the text in which this error occurred
type InTextError struct {
// count starts from zero, so 1 should be added to this number before displaying
Line int
Err error
}
func (c InTextError) Unwrap() error {
return c.Err
}
func (c InTextError) Error() string {
return fmt.Sprintf("in line %v %s", c.Line+1, c.Err.Error())
}
// this type represents a one 8-dot braille cell
// every bit in this byte represents a one braille dot, like in braille unicode
type BrailleField byte
func (c BrailleField) Dot1() bool {
return (byte(c) & Dot1) != 0
}
func (c BrailleField) Dot2() bool {
return (byte(c) & Dot2) != 0
}
func (c BrailleField) Dot3() bool {
return (byte(c) & Dot3) != 0
}
func (c BrailleField) Dot4() bool {
return (byte(c) & Dot4) != 0
}
func (c BrailleField) Dot5() bool {
return (byte(c) & Dot5) != 0
}
func (c BrailleField) Dot6() bool {
return (byte(c) & Dot6) != 0
}
func (c BrailleField) Dot7() bool {
return (byte(c) & Dot7) != 0
}
func (c BrailleField) Dot8() bool {
return (byte(c) & Dot8) != 0
}

// convert this braille cell to string dots representation
func (c BrailleField) String() string {
var res []int
var r bool
r = c.Dot1()
if r {
res = append(res, 1)
}
r = c.Dot2()
if r {
res = append(res, 2)
}
r = c.Dot3()
if r {
res = append(res, 3)
}
r = c.Dot4()
if r {
res = append(res, 4)
}
r = c.Dot5()
if r {
res = append(res, 5)
}
r = c.Dot6()
if r {
res = append(res, 6)
}
r = c.Dot7()
if r {
res = append(res, 7)
}
r = c.Dot8()
if r {
res = append(res, 8)
}
var s string
for _,i := range res {
s += strconv.Itoa(i)
}
return s
}

// get braille unicode representation for this braille cell
func (c BrailleField) ToUnicode() rune {
return UnicodeBrailleBase + rune(c)
}

func FieldFromUnicode(r rune) (BrailleField, error) {
if r < UnicodeBrailleBase || r > rune(0x28FF) {
return BrailleField(0),InvalidBrailleError
}
return BrailleField(byte(r - UnicodeBrailleBase)),nil
}
func FieldFromString(s string) (BrailleField, error) {
var res byte
for _,r := range s {
switch r {
case '1':
res = res | Dot1

case '2':
res = res | Dot2

case '3':
res = res | Dot3

case '4':
res = res | Dot4

case '5':
res = res | Dot5

case '6':
res = res | Dot6

case '7':
res = res | Dot7

case '8':
res = res | Dot8
default:
return BrailleField(res),InvalidBrailleError
}
}
return BrailleField(res), nil
}

// this type represents one braille line
type BrailleRow []BrailleField
func RowFromUnicode(s string) (BrailleRow, error) {
ss := []rune(s)
res := make(BrailleRow, len(ss))
for i,v := range ss {
if ! IsBrailleUnicode(v) {
return res,InRowError{i, InvalidBrailleError}
}
f,_ := FieldFromUnicode(v)
res[i] = f
}
return res,nil
}
func RowFromInternal(d []byte) BrailleRow {
res := make(BrailleRow, len(d))
for i,v := range d {
res[i] = BrailleField(v)
}
return res
}
func (c BrailleRow) ToUnicode() string {
res := make([]rune, len(c))
for i,v := range c {
res[i] = v.ToUnicode()
}
return string(res)
}
func (c BrailleRow) ToInternal() []byte {
res := make([]byte, len(c))
for i,v := range c {
res[i] = byte(v)
}
return res
}

// this type represents a collection of braille lines
type BraillePage []BrailleRow
func PageFromUnicode(s string) (BraillePage, error) {
var res BraillePage
ss := strings.Split(s, "\n")
for i,v := range ss {
row,err := RowFromUnicode(v)
if err != nil {
return res,InTextError{i, err}
}
res = append(res, row)
}
return res,nil
}
func (c BraillePage) ToUnicode() string {
sb := new(strings.Builder)
for i,v := range c {
sb.WriteString(v.ToUnicode())
if i != len(c)-1 {
sb.WriteRune('\n')
}
}
return sb.String()
}