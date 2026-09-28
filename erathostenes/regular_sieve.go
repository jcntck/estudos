package sieve

type RegularSieveOfEratosthenes struct {
	crossedOut []bool
	result     []int
}

func NewRegularSieveOfEratosthenes() *RegularSieveOfEratosthenes {
	return &RegularSieveOfEratosthenes{}
}

func (e *RegularSieveOfEratosthenes) PrimesUpTo(limit int) []int {
	if limit < 2 {
		return []int{}
	}
	e.uncrossIntegersUpTo(limit)
	e.crossOutMultiples()
	e.putUncrossedIntegersIntoResult()
	return e.result
}

func (e *RegularSieveOfEratosthenes) uncrossIntegersUpTo(limit int) {
	e.crossedOut = make([]bool, limit+1)
	for i := 2; i < len(e.crossedOut); i++ {
		e.crossedOut[i] = false
	}
}

func (e *RegularSieveOfEratosthenes) crossOutMultiples() {
	for i := 2; i < len(e.crossedOut); i++ {
		if !e.crossedOut[i] {
			e.crossOutMultiplesOf(i)
		}
	}
}

func (e *RegularSieveOfEratosthenes) crossOutMultiplesOf(i int) {
	for multiple := 2 * i; multiple < len(e.crossedOut); multiple += i {
		e.crossedOut[multiple] = true
	}
}

func (e *RegularSieveOfEratosthenes) putUncrossedIntegersIntoResult() {
	e.result = []int{}
	for i := 2; i < len(e.crossedOut); i++ {
		if !e.crossedOut[i] {
			e.result = append(e.result, i)
		}
	}
}
