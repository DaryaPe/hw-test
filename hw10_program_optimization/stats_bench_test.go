//go:build bench
// +build bench

package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

// go test -bench=. -benchmem > new.txt -run=BenchmarkGetDomainStat -tags=bench
// benchstat old.txt new.txt

func BenchmarkGetDomainStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r, err := zip.OpenReader("testdata/users.dat.zip")
		require.NoError(b, err)
		defer r.Close()

		require.Equal(b, 1, len(r.File))

		data, err := r.File[0].Open()
		require.NoError(b, err)

		if _, err := GetDomainStat(data, "info"); err != nil {
			b.Fatal(err)
		}
	}
}
