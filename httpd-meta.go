/*
GoGet
GoBuildNull
GoBuild
GoRun
ListenAddr=:8080 GoRun
*/
package main
import (
	"fmt"
	"io"
	"net/http"
	"os"
)
const (
	ListenAddrDef=":80"
	SP=" "
	TAB="\t"
	NL="\n"
)
var (
	ListenAddr string
	F=fmt.Sprintf
	EF=fmt.Errorf
	pout=fmt.Print
)
func perr(msg string) (int, error) { return fmt.Fprint(os.Stderr, msg+NL) }
func funchttp(w http.ResponseWriter, req *http.Request) {
	var err error
	var reqbody []byte
	reqbody,err=io.ReadAll(req.Body)
	if err!=nil { perr(F("ERROR request body ReadAll %v", err)) }
	defer req.Body.Close()
	// https://pkg.go.dev/net/http#Request
	// https://pkg.go.dev/net/url#URL.Query
	reqmeta:=F(
		"proto[%s]"+SP+"method[%s]"+SP+"path[%s]"+NL+
		"query[%s]"+NL+
		"headers{ "+NL+
			"%s"+
		" }"+NL+
		"body[-"+NL+
			"%s"+NL+
		"-]",
		req.Proto, req.Method, req.URL.Path,
		req.URL.RawQuery,
		func(m map[string][]string)(s string){
			for k, vv := range m {
				s += "["+k+"] ("+NL
				for _, v := range vv {
					s += TAB+"["+v+"]"+NL
				}
				s += ")"+NL
			}
			return s
		}(req.Header),
		string(reqbody),
	)
	perr(reqmeta)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	reqq:=req.URL.Query()
	hubchallenge:=reqq.Get("hub.challenge")
	perr(F(
		"hub.mode [%s]"+SP+"hub.verify_token [%s]"+SP+"hub.challenge [%s]",
		reqq.Get("hub.mode"), 
		reqq.Get("hub.verify_token"), 
		hubchallenge,
	))
	if hubchallenge!="" {
		io.WriteString(w, hubchallenge)
		return
	}
	io.WriteString(w, reqmeta+NL+NL)
	return
}
func main() {
	ListenAddr=ListenAddrDef
	if la:=os.Getenv("ListenAddr"); la!="" { ListenAddr=la }
	perr(F("ListenAddr [%s]", ListenAddr))
	http.HandleFunc("/", funchttp)
	if err:=http.ListenAndServe(ListenAddr, nil); err!=nil {
		perr(F("ERROR ListenAndServe [%s] %v", ListenAddr, err))
		os.Exit(1)
	}
}

