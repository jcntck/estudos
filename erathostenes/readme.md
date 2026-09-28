# Crivo de Eratóstenes

Este estudo implementa o **Crivo de Eratóstenes** em Go para explorar a geração de números primos e comparar o desempenho de duas versões do algoritmo.

## Fonte do estudo

A referência para este exercício é o livro **Clean Code**, de Robert C. Martin, no capítulo sobre **Formatação**.

## Motivação

Meu interesse é visualizar, na prática, como limitar parte das iterações à raiz quadrada pode melhorar o desempenho de um algoritmo. Para isso, o estudo compara uma versão regular do crivo com uma versão que reduz o alcance da iteração responsável por marcar os múltiplos.

O exercício também é uma oportunidade de praticar Go, uma linguagem que estou estudando. A implementação, os testes e os benchmarks ajudam a exercitar a linguagem enquanto observo o comportamento do algoritmo.

## Por que escolhi o crivo?

O Crivo de Eratóstenes é desafiador o suficiente para estimular o raciocínio, mas simples o bastante para servir como um exercício básico de Go. Ele permite trabalhar com laços, slices e organização do código, além de oferecer uma otimização cujo efeito pode ser medido e compreendido.

## Como o crivo funciona

Um número primo é um inteiro maior que 1 que só possui dois divisores positivos: 1 e ele mesmo. O crivo encontra esses números eliminando os compostos:

1. Reserve uma marcação para cada número de 0 até o limite desejado.
2. Comece pelo 2, o primeiro primo, e marque seus múltiplos a partir de 4 como compostos.
3. Avance para o próximo número ainda não marcado e repita a marcação dos múltiplos.
4. Reúna os números a partir de 2 que permaneceram sem marcação.

Por exemplo, até 10, o 2 elimina 4, 6, 8 e 10; o 3 elimina 6 e 9. O resultado é `[2, 3, 5, 7]`.

No código, `crossedOut[i] == true` significa que o número foi marcado como composto. Os índices 0 e 1 ficam fora da coleta do resultado.

## Onde entra a raiz quadrada?

Se um número composto é escrito como `a × b`, pelo menos um desses fatores é menor ou igual à sua raiz quadrada. Caso ambos fossem maiores, o produto ultrapassaria o próprio número. Por isso, basta marcar os múltiplos dos primos até essa região para eliminar os compostos até o limite desejado.

A versão regular percorre os candidatos até o limite informado. A versão otimizada reduz essa iteração usando a raiz quadrada do tamanho do vetor, calculada como `sqrt(limit + 1)`. Nas duas versões, a coleta final continua percorrendo os números até o limite para reunir todos os primos encontrados.

Para acompanhar o estudo, comece pela [implementação regular](regular_sieve.go) e compare a etapa de marcação com a [implementação otimizada](optimized_sieve.go). Depois, consulte os [resultados dos benchmarks](benchmark.txt) para observar o efeito dessa mudança no tempo de execução.

[Voltar ao índice de estudos](../readme.md)
