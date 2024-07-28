package hw10programoptimization

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, _ := os.Open("testdata/benchmark_test.txt")

		_, err := GetDomainStat(data, "biz")
		if err != nil {
			fmt.Println(err.Error())
		}
		data.Close()
	}
}
