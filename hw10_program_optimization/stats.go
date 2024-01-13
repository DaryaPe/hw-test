package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
)

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string `json:"Email"`
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	user := &User{}
	scan := bufio.NewScanner(r)
	result := make(DomainStat)
	var err error
	for scan.Scan() {
		*user = User{}
		if err = user.UnmarshalJSON(scan.Bytes()); err != nil {
			return nil, err
		}

		if strings.Contains(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
