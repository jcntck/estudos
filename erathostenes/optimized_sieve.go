package sieve

import "math"

type OptimizedSieveOfEratosthenes struct {
	crossedOut []bool
	result     []int
}

func NewOptimizedSieveOfEratosthenes() *OptimizedSieveOfEratosthenes {
	return &OptimizedSieveOfEratosthenes{}
}

func (e *OptimizedSieveOfEratosthenes) PrimesUpTo(limit int) []int {
	if limit < 2 {
		return []int{}
	}
	e.uncrossIntegersUpTo(limit)
	e.crossOutMultiples()
	e.putUncrossedIntegersIntoResult()
	return e.result
}

func (e *OptimizedSieveOfEratosthenes) uncrossIntegersUpTo(limit int) {
	e.crossedOut = make([]bool, limit+1)
	for i := 2; i < len(e.crossedOut); i++ {
		e.crossedOut[i] = false
	}
}

func (e *OptimizedSieveOfEratosthenes) crossOutMultiples() {
	limit := e.determineIterationLimit()
	for i := 2; i <= limit; i++ {
		if !e.crossedOut[i] {
			e.crossOutMultiplesOf(i)
		}
	}
}

func (e *OptimizedSieveOfEratosthenes) determineIterationLimit() int {
	iterationLimit := math.Sqrt(float64(len(e.crossedOut)))
	return int(iterationLimit)
}

func (e *OptimizedSieveOfEratosthenes) crossOutMultiplesOf(i int) {
	for multiple := 2 * i; multiple < len(e.crossedOut); multiple += i {
		e.crossedOut[multiple] = true
	}
}

func (e *OptimizedSieveOfEratosthenes) putUncrossedIntegersIntoResult() {
	e.result = []int{}
	for i := 2; i < len(e.crossedOut); i++ {
		if !e.crossedOut[i] {
			e.result = append(e.result, i)
		}
	}
}
