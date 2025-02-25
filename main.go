package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
)

type Url struct {
	Scheme     string        `json:"scheme"`
	Opaque     string        `json:"opaque"`      // encoded opaque data
	User       *url.Userinfo `json:"user"`        // username and password information
	Host       string        `json:"host"`        // host or host:port
	Path       string        `json:"path"`        // path (relative paths may omit leading slash)
	RawPath    string        `json:"raw_path"`    // encoded path hint (see EscapedPath method)
	ForceQuery bool          `json:"force_query"` // append a query ('?') even if RawQuery is empty
	// RawQuery    string        `json:"raw_query"`    // encoded query values, without '?'
	Fragment    string     `json:"fragment"`     // fragment for references, without '#'
	RawFragment string     `json:"raw_fragment"` // encoded fragment hint (see EscapedFragment method)
	Query       url.Values `json:"query"`
	// Query       map[string]interface{} `json:"query"`
}

var reverseFlag = flag.Bool("r", false, "json -> url")

func urlToJson(s string) {
	u, err := url.Parse(s)

	if err != nil {
		fmt.Fprintln(os.Stderr, "bad url", err)
		os.Exit(1)
	}

	o := &Url{
		Scheme:     u.Scheme,
		Opaque:     u.Opaque,
		User:       u.User,
		Host:       u.Host,
		Path:       u.Path,
		RawPath:    u.RawPath,
		ForceQuery: u.ForceQuery,
		// RawQuery:    u.RawQuery,
		Fragment:    u.Fragment,
		RawFragment: u.RawFragment,
		Query:       u.Query(),
	}
	e := json.NewEncoder(os.Stdout)
	e.Encode(o)
}

func jsonToUrl(o Url) {
	u := url.URL{
		Scheme:      o.Scheme,
		Opaque:      o.Opaque,
		User:        o.User,
		Host:        o.Host,
		Path:        o.Path,
		RawPath:     o.RawPath,
		ForceQuery:  o.ForceQuery,
		RawQuery:    o.Query.Encode(),
		Fragment:    o.Fragment,
		RawFragment: o.RawFragment,
	}
	fmt.Println(u.String())
}

func main() {

	flag.Parse()

	if !*reverseFlag {
		if flag.NArg() > 0 {
			for _, u := range flag.Args() {
				urlToJson(u)
			}
		} else {
			r := bufio.NewScanner(os.Stdin)
			buf := make([]byte, 0, 1024*1024)
			r.Buffer(buf, cap(buf))
			for r.Scan() {
				urlToJson(r.Text())
			}
			if err := r.Err(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	} else {
		if flag.NArg() > 0 {
			for _, u := range flag.Args() {
				var o Url
				if err := json.Unmarshal([]byte(u), &o); err != nil {
					fmt.Fprintln(os.Stderr, "bad json url", err)
					os.Exit(1)
				}
				jsonToUrl(o)
			}
		} else {
			var o Url
			r := json.NewDecoder(os.Stdin)
			for {
				err := r.Decode(&o)
				if err == io.EOF {
					return
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, "bad json url", err)
					os.Exit(1)
				}

				jsonToUrl(o)
			}
		}
	}

}
