package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

//easyjson:json
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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (result users, err error) {

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		user := User{}
		if err = easyjson.Unmarshal([]byte(scanner.Text()), &user); err != nil {
			return nil, err
		}
		result = append(result, user)
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.Contains(user.Email, "."+domain) {
			subDomain := strings.ToLower(user.Email[strings.Index(user.Email, "@")+1:])
			num := result[subDomain]
			num++
			result[subDomain] = num
		}
	}
	return result, nil
}
