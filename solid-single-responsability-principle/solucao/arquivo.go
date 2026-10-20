package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Arquivo struct {
	origem       *string
	destino      *string
	pastaDestino string
	arquivoSaida *os.File
}

func (a *Arquivo) CarregarPametros() {
	a.origem = flag.String("origem", "", "arquivo JSON de entrada")
	a.destino = flag.String("destino", "", "local de armazenamento do JSON de saída")
	flag.Parse()
}

func (a *Arquivo) CarregarArquivo() Entrada {
	arquivo := a.lerArquivo()
	if arquivo != os.Stdin {
		defer arquivo.Close()
	}
	var entrada Entrada
	err := json.NewDecoder(arquivo).Decode(&entrada)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro ao ler entrada:", err)
		os.Exit(1)
	}

	return entrada
}

func (a *Arquivo) lerArquivo() *os.File {
	arquivo := os.Stdin
	if *a.origem != "" {
		var err error
		arquivo, err = os.Open(*a.origem)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro ao abrir arquivo:", err)
			os.Exit(1)
		}
	}
	return arquivo
}

func (a *Arquivo) CriarArquivo(nome string) {
	a.criarPastaSeNecessario()
	caminho := filepath.Join(a.pastaDestino, nome+".json")
	arquivoSaida, err := os.OpenFile(caminho, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro ao criar arquivo:", err)
		os.Exit(1)
	}
	a.arquivoSaida = arquivoSaida
}

func (a *Arquivo) criarPastaSeNecessario() {
	a.pastaDestino = *a.destino
	if err := os.MkdirAll(a.pastaDestino, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "erro ao criar diretório:", err)
		os.Exit(1)
	}
}

func (a *Arquivo) Escrever(conteudo []byte) {
	if _, err := a.arquivoSaida.Write(conteudo); err != nil {
		a.arquivoSaida.Close()
		fmt.Fprintln(os.Stderr, "erro ao gravar arquivo:", err)
		os.Exit(1)
	}
	if err := a.arquivoSaida.Close(); err != nil {
		fmt.Fprintln(os.Stderr, "erro ao fechar arquivo:", err)
		os.Exit(1)
	}
}
