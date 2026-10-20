package main

type StatusVenda string

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
