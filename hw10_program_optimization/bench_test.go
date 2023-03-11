package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func Benchmark(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		return
	}
	defer r.Close()
	data, err := r.File[0].Open()
	if err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		_, err = GetDomainStat(data, "biz")
		if err != nil {
			return
		}
	}
}
