package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	StatusVendaConfirmada StatusVenda = "confirmada"
	StatusVendaCancelada  StatusVenda = "cancelada"
)

func (s StatusVenda) Valido() bool {
	switch s {
	case StatusVendaConfirmada, StatusVendaCancelada:
		return true
	default:
		return false
	}
}

func ProcessarFechamento() {
	arquivo := Arquivo{}
	arquivo.CarregarPametros()
	entrada := arquivo.CarregarArquivo()

	vistos := map[string]struct{}{}
	total := int64(0)
	qtdVendas := 0
	comissao := int64(0)
	for _, venda := range entrada.Vendas {
		if strings.TrimSpace(venda.ID) == "" || venda.Valor <= 0 || !venda.Status.Valido() {
			fmt.Fprintln(os.Stderr, "ID, valor ou status estão invalidos")
			os.Exit(1)
		}

		if _, ok := vistos[venda.ID]; ok {
			fmt.Fprintln(os.Stderr, "ID já informado anteriormente:", venda.ID)
			os.Exit(1)
		}
		vistos[venda.ID] = struct{}{}

		if venda.Status == StatusVendaConfirmada {
			total += venda.Valor
			qtdVendas++
		}
	}

	comissao = total/100*5 + total%100*5/100

	fechamento := Fechamento{
		ID:         entrada.ID,
		Vendedor:   entrada.Vendedor,
		Quantidade: qtdVendas,
		Total:      total,
		Comissao:   comissao,
	}

	conteudo, err := json.MarshalIndent(fechamento, "", " ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro ao criar JSON:", err)
		os.Exit(1)
	}

	arquivo.CriarArquivo(fechamento.ID)
	arquivo.Escrever(conteudo)

	fmt.Printf("Fechamento: %s\n", fechamento.ID)
	fmt.Printf("Vendedor: %s\n", fechamento.Vendedor)
	fmt.Printf("Vendas confirmadas: %d\n", fechamento.Quantidade)
	fmt.Printf("Total vendido: %s\n", formatarReais(fechamento.Total))
	fmt.Printf("Comissão: %s\n", formatarReais(fechamento.Comissao))
}

func formatarReais(centavos int64) string {
	return fmt.Sprintf("R$ %d,%02d", centavos/100, centavos%100)
}

func main() {
	ProcessarFechamento()
}
