// this code is licensed under mit license, see license file in repository root
// Copyright 2021 ilyapashuk


// this package provides very simple table driven translator and backtranslator
// only plain char to dots tables, computer braille tables, are supported

package translation
import "strings"
import "fmt"
import "errors"
import . "braille"

type UnmappableCharacterError struct {
Rune rune
}
func (c UnmappableCharacterError) Error() string {
return fmt.Sprintf("mapping for character %s is missing in the table", string(c.Rune))
}

type UnmappableDotsError struct {
Dots BrailleField
}
func (c UnmappableDotsError) Error() string {
return fmt.Sprintf("mapping for dots %v is missing in the table", c.Dots)
}

type TranslationRule struct {
Rune rune
Dots BrailleField
}
type RuleList []TranslationRule
// rule list format is simmeler with liblouis table format, but much more simple and don't supports any backslash escaping, so not compatible with it
func ParseRuleList(td []string) (RuleList, error) {
var res RuleList
for i,l := range td {
if l == "" || strings.HasPrefix(l, "#") {
// skip it as an comment
continue
}
w := strings.Split(l, " ")
if len(w) < 3 {
return res,InTextError{i, errors.New("invalid rule format")}
}
charlist := w[1]
dots := w[2]
f,err := FieldFromString(dots)
if err != nil {
return res,InTextError{i, err}
}
for _,char := range charlist {
res = append(res, TranslationRule{Dots:f,Rune:char})
}
}
return res,nil
}
func (c RuleList) ToForwardTable() ForwardTable {
res := make(ForwardTable)
Add := func(r rune, b BrailleField) {
if _,ok := res[r]; ! ok {
res[r] = b
}
}
for _,v := range c {
Add(v.Rune, v.Dots)
}
return res
}
func (c RuleList) ToBackTable() BackTable {
res := make(BackTable)
Add := func(f BrailleField, r rune) {
if _,ok := res[f]; ! ok {
res[f] = r
}
}
for _,v := range c {
Add(v.Dots, v.Rune)
}
return res
}
type ForwardTable map[rune]BrailleField
func (c ForwardTable) Translate(s string) (BrailleRow, error) {
ss := []rune(s)
res := make(BrailleRow, len(ss))
for i,v := range ss {
if v == ' ' {
res[i] = BrailleField(0)
continue
}
if IsBrailleUnicode(v) {
f,_ := FieldFromUnicode(v)
res[i] = f
} else {
f,ok := c[v]
if ! ok {
return res,InRowError{i, UnmappableCharacterError{v}}
}
res[i] = f
}
}
return res,nil
}
func (c ForwardTable) TranslateText(s string) (BraillePage, error) {
var res BraillePage
ss := strings.Split(s, "\n")
for i,v := range ss {
row,err := c.Translate(v)
if err != nil {
return res,InTextError{i,err}
}
res = append(res, row)
}
return res,nil
}

type BackTable map[BrailleField]rune
func (c BackTable) Translate(br BrailleRow) (string, error) {
sb := new(strings.Builder)
for i,v := range br {
if v == BrailleField(0) {
sb.WriteRune(' ')
continue
}
r,ok := c[v]
if ! ok {
return sb.String(),InRowError{i, UnmappableDotsError{v}}
}
sb.WriteRune(r)
}
return sb.String(),nil
}
func (c BackTable) TranslateText(bp BraillePage) (string, error) {
sb := new(strings.Builder)
for i,br := range bp {
s,err := c.Translate(br)
if err != nil {
return sb.String(),InTextError{i, err}
}
sb.WriteString(s)
if i != len(bp)-1 {
sb.WriteRune('\n')
}
}
return sb.String(),nil
}