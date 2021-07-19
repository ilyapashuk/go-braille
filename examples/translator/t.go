package main
import "os"
import "braille/translation"
import "strings"

func fatal(err error) {
if err == nil {
return
}
os.Stderr.Write([]byte("error\n"))
os.Stderr.Write([]byte(err.Error() + "\n"))
os.Exit(1)
}
func main() {
tn := os.Args[1]
tb,err := os.ReadFile(tn)
fatal(err)
ts := string(tb)
ts = strings.ReplaceAll(ts, "\r", "")
td := strings.Split(ts, "\n")
rl,err := translation.ParseRuleList(td)
fatal(err)
t := rl.ToForwardTable()
infile := os.Args[2]
indata,err := os.ReadFile(infile)
fatal(err)
indata = []byte(strings.ReplaceAll(string(indata), "\r", ""))
inpage,err := t.TranslateText(string(indata))
fatal(err)
outfile := os.Args[3]
of,_ := os.Create(outfile)
defer of.Close()
of.Write([]byte(inpage.ToUnicode()))
}