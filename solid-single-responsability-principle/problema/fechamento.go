package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type StatusVenda string

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

type Venda struct {
	ID     string      `json:"id"`
	Valor  int64       `json:"valor_centavos"`
	Status StatusVenda `json:"status"`
}

type Entrada struct {
	ID       string  `json:"id"`
	Vendedor string  `json:"vendedor"`
	Vendas   []Venda `json:"vendas"`
}

type Fechamento struct {
	ID         string `json:"id"`
	Vendedor   string `json:"vendedor"`
	Quantidade int    `json:"quantidade_vendas_confirmadas"`
	Total      int64  `json:"total_centavos"`
	Comissao   int64  `json:"comissao_centavos"`
}

func ProcessarFechamento() {
	origem := flag.String("origem", "", "arquivo JSON de entrada")
	destino := flag.String("destino", "", "local de armazenamento do JSON de saída")
	flag.Parse()

	arquivo := os.Stdin

	if *origem != "" {
		var err error
		arquivo, err = os.Open(*origem)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro ao abrir arquivo:", err)
			os.Exit(1)
		}
		defer arquivo.Close()
	}

	var entrada Entrada
	err := json.NewDecoder(arquivo).Decode(&entrada)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro ao ler entrada:", err)
		os.Exit(1)
	}

	if strings.TrimSpace(entrada.ID) == "" || strings.TrimSpace(entrada.Vendedor) == "" {
		fmt.Fprintln(os.Stderr, "ID e/ou vendedor invalidos")
		os.Exit(1)
	}

	vistos := map[string]bool{}
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
		vistos[venda.ID] = true

		if venda.Status == StatusVendaConfirmada {
			total += int64(venda.Valor)
			qtdVendas++
		}
	}

	if total > 0 {
		comissao = int64(float64(total) * 0.05)
	}

	pastaDestino := *destino
	if err := os.MkdirAll(pastaDestino, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "erro ao criar diretório:", err)
		os.Exit(1)
	}

	fechamento := Fechamento{
		ID:         entrada.ID,
		Vendedor:   entrada.Vendedor,
		Quantidade: qtdVendas,
		Total:      total,
		Comissao:   int64(comissao),
	}

	conteudo, err := json.MarshalIndent(fechamento, "", " ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro ao criar JSON:", err)
		os.Exit(1)
	}

	caminho := filepath.Join(pastaDestino, fechamento.ID+".json")
	arquivoSaida, err := os.OpenFile(caminho, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro ao criar arquivo:", err)
		os.Exit(1)
	}
	defer arquivoSaida.Close()

	if _, err := arquivoSaida.Write(conteudo); err != nil {
		fmt.Fprintln(os.Stderr, "erro ao gravar arquivo:", err)
		os.Exit(1)
	}

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
