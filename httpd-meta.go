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
	NL="\n"
)
var (
	ListenAddr string
	F=fmt.Sprintf
	EF=fmt.Errorf
	pout=fmt.Print
)
func main() {
	
	var err error
	ListenAddr=ListenAddrDef
	if la:=os.Getenv("ListenAddr"); la!="" { ListenAddr=la }
	perr(F("ListenAddr [%s]", ListenAddr))
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var reqbody []byte
		reqbody,err=io.ReadAll(req.Body)
		if err!=nil { perr(F("ERROR request body ReadAll %v", err)) }
		defer req.Body.Close()
		// https://pkg.go.dev/net/http#Request
		// https://pkg.go.dev/net/url#URL.Query
		perr(F(
			"proto[%s]"+SP+"method[%s]"+SP+"path[%s]"+NL+
			"query[%s]"+NL+
			"headers{ %v }"+NL+
			"body[%s]",
			req.Proto, req.Method, req.URL.Path,
			req.URL.RawQuery,
			req.Header,
			string(reqbody),
			))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		reqq:=req.URL.Query()
		mode:=reqq.Get("hub.mode")
		token:=reqq.Get("hub.verify_token")
		challenge:=reqq.Get("hub.challenge")
		perr(F(
			"hub.mode [%s]"+SP+"hub.verify_token [%s]"+SP+"hub.challenge [%s]",
			mode, token, challenge,
			))
		io.WriteString(w, challenge)
	})
	if err:=http.ListenAndServe(ListenAddr, nil); err!=nil {
		perr(F("ERROR ListenAndServe [%s] %v", ListenAddr, err))
		os.Exit(1)
	}
	
}
func perr(msg string) (int, error) { return fmt.Fprint(os.Stderr, msg+NL) }

