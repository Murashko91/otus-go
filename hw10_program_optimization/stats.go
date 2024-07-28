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
	Email    string `json:"email"`
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

type users [100_000]User

func getUsers(r io.Reader) (users, error) {
	result := users{}

	// 2. Replaced ReadAll with  bufio.NewReader(benchmark_results/benchstat_2.txt)
	scanner := bufio.NewReader(r)
	i := 0
	for {
		line, isPrefix, readErr := scanner.ReadLine()
		if isPrefix {
			break
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return users{}, readErr
		}

		user := User{}

		// 3. Updated Unmarshal to easyjson (benchmark_results/benchstat_3.txt)
		if err := easyjson.Unmarshal((line), &user); err != nil {
			return users{}, err
		}
		result[i] = user
		i++
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		// 1. Updated Regexp  to strings.Contains (benchmark_results/benchstat_1.txt)
		if strings.Contains(user.Email, "."+domain) {
			fullDomain := strings.ToLower(user.Email[strings.Index(user.Email, "@")+1:])
			num := result[fullDomain]
			num++
			result[fullDomain] = num
		}
	}
	return result, nil
}
