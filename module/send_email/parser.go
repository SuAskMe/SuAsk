package send_email

import (
	"errors"
	"sort"
	"strings"
)

type pair struct {
	key string
	idx int
}

type cache struct {
	pairs []pair
	slice []string
}

func formatTmpl(tmpl string, list []string) (*cache, error) {
	pairs := make([]pair, 0)
	for _, v := range list {
		if !strings.Contains(tmpl, v) {
			return nil, errors.New("msgTmpl has not " + v)
		}
		tmp := tmpl
		offset := 0
		for strings.Contains(tmp, v) {
			// get first index of v in tmp
			idx := strings.Index(tmp, v)
			pairs = append(pairs, pair{key: v, idx: offset + idx})
			tmp = tmp[idx+len(v):]
			offset += idx + len(v)
		}
	}
	// sort pairs by idx
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].idx < pairs[j].idx
	})
	slice := make([]string, 0)
	st := 0
	for _, v := range pairs {
		slice = append(slice, tmpl[st:v.idx])
		st = v.idx + len(v.key)
	}
	slice = append(slice, tmpl[st:])
	if len(slice) != len(pairs)+1 {
		return nil, errors.New("parse error, check formatMsgTmpl()")
	}
	return &cache{pairs: pairs, slice: slice}, nil
}

func getAuthCode(tmpl *cache, code string) string {
	var sb strings.Builder
	for i := 0; i < len(tmpl.pairs); i++ {
		sb.WriteString(tmpl.slice[i])
		sb.WriteString(code)
	}
	sb.WriteString(tmpl.slice[len(tmpl.pairs)])
	return sb.String()
}
func getMessage(tmpl *cache, notice *Notice) string {
	var sb strings.Builder
	for i := 0; i < len(tmpl.pairs); i++ {
		sb.WriteString(tmpl.slice[i])
		switch tmpl.pairs[i].key {
		case "$user$":
			sb.WriteString(notice.User)
		case "$type$":
			sb.WriteString(notice.Type)
		case "$content$":
			sb.WriteString(notice.Content)
		case "$url$":
			sb.WriteString(notice.URL)
		}
	}
	sb.WriteString(tmpl.slice[len(tmpl.pairs)])
	return sb.String()
}