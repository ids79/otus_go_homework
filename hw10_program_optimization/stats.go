package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jlexer "github.com/mailru/easyjson/jlexer"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domain = "." + domain
	var user User
	scanner := bufio.NewScanner(r)
	j := &jlexer.Lexer{}
	for scanner.Scan() {
		// if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
		//	 return nil, err
		// }
		// if err := user.UnmarshalJSON(scanner.Bytes()); err != nil {
		// 	return nil, err
		// }
		*j = jlexer.Lexer{}
		j.Data = scanner.Bytes()
		user.UnmarshalEasyJSON(j)
		if strings.Contains(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
