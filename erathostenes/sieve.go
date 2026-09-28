package sieve

// Retornar uma lista vazia quando limit < 2.
// Criar uma marcação para os números de 0 até limit.
// Considerar inicialmente os números como possíveis primos.
// Para cada primo encontrado, marcar seus múltiplos como compostos.
// Retornar somente os números que permanecerem marcados como primos.

type SieveOfEratosthenes interface {
	PrimesUpTo(limit int) []int
}
