package pion

import "strings"

func FindHostInCandidate(raw string) string {
	split := strings.Fields(raw)
	// Foundation not specified: not RFC 8445 compliant but seen in the wild
	if len(raw) != 0 && raw[0] == ' ' {
		split = append([]string{" "}, split...)
	}
	if len(split) < 8 {
		return ""
	}

	address := split[4]

	return address

}
