package sieve

import (
	"reflect"
	"testing"
)

func TestEratosthenesPrimesUpTo(t *testing.T) {
	implementations := []struct {
		name string
		new  func() SieveOfEratosthenes
	}{
		{
			name: "regular",
			new: func() SieveOfEratosthenes {
				return NewRegularSieveOfEratosthenes()
			},
		},
		{
			name: "optimized",
			new: func() SieveOfEratosthenes {
				return NewOptimizedSieveOfEratosthenes()
			},
		},
	}

	tests := []struct {
		name     string
		limit    int
		expected []int
	}{
		{name: "limite menor que dois", limit: 1, expected: []int{}},
		{name: "limite igual a dois", limit: 2, expected: []int{2}},
		{name: "limite igual a dez", limit: 10, expected: []int{2, 3, 5, 7}},
		{
			name:     "limite igual a trinta",
			limit:    30,
			expected: []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29},
		},
	}

	for _, implementation := range implementations {
		t.Run(implementation.name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					sieve := implementation.new()
					result := sieve.PrimesUpTo(test.limit)

					if !reflect.DeepEqual(result, test.expected) {
						t.Fatalf("resultado = %v; esperado = %v", result, test.expected)
					}
				})
			}
		})
	}
}

func TestImplementationsStress(t *testing.T) {
	implementations := []struct {
		name string
		new  func() SieveOfEratosthenes
	}{
		{
			name: "regular",
			new: func() SieveOfEratosthenes {
				return NewRegularSieveOfEratosthenes()
			},
		},
		{
			name: "optimized",
			new: func() SieveOfEratosthenes {
				return NewOptimizedSieveOfEratosthenes()
			},
		},
	}

	for _, implementation := range implementations {
		t.Run(implementation.name, func(t *testing.T) {
			primes := implementation.new().PrimesUpTo(1_000_000)

			if len(primes) != 78_498 {
				t.Fatalf("quantidade de primos = %d; esperado = 78498", len(primes))
			}
			if last := primes[len(primes)-1]; last != 999_983 {
				t.Fatalf("último primo = %d; esperado = 999983", last)
			}
		})
	}
}
