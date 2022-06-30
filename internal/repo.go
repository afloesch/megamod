package internal

import "strings"

type Repo string

func (r Repo) Organization() string {
	return strings.Split(r.String(), "/")[0]
}

func (r Repo) Name() string {
	parts := strings.Split(r.String(), "/")
	if len(parts) > 1 {
		return parts[1]
	}
	return r.Name()
}

func (r Repo) String() string {
	return string(r)
}
