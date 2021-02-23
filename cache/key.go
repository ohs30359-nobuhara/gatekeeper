package customCache

import (
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func CreateKeyFromRequest(r *http.Request) string {
	u, e := url.Parse(r.RequestURI)

	if e != nil {
		return ""
	}

	a := strings.Split(u.RequestURI(), "?")

	if len(a) == 2 {
		// return r.Host + a[0] + "?" +sortQuery(u.Query())
		 return a[0] + "?" +sortQuery(u.Query())
	}

	return u.RequestURI()
}

func sortQuery(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		vs := v[k]
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}
	return buf.String()
}