package srp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type venda struct {
	ID     string `json:"id"`
	Valor  int64  `json:"valor_centavos"`
	Status string `json:"status"`
}

type entrada struct {
	ID       string  `json:"id"`
	Vendedor string  `json:"vendedor"`
	Vendas   []venda `json:"vendas"`
}

type fechamento struct {
	ID         string `json:"id"`
	Vendedor   string `json:"vendedor"`
	Quantidade int    `json:"quantidade_vendas_confirmadas"`
	Total      int64  `json:"total_centavos"`
	Comissao   int64  `json:"comissao_centavos"`
}

func TestIntegracaoV1(t *testing.T) { testarVersao(t, false) }
func TestIntegracaoV2(t *testing.T) { testarVersao(t, true) }

func compilar(t *testing.T) string {
	t.Helper()
	alvo := os.Getenv("SRP_IMPL")
	if alvo == "" {
		alvo = "problema"
	}
	if alvo != "problema" && alvo != "solucao" {
		t.Fatal("SRP_IMPL deve ser problema ou solucao")
	}
	binario := filepath.Join(t.TempDir(), "fechamento")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-o", binario, "./"+alvo)
	if saida, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("implemente um programa executável em %s/ conforme o README; compilação falhou: %v\n%s", alvo, err, saida)
	}
	return binario
}

func executar(t *testing.T, binario, diretorio string, dados entrada) (string, string, error) {
	t.Helper()
	payload, err := json.Marshal(dados)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binario, "-destino", diretorio)
	cmd.Stdin = bytes.NewReader(payload)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err = cmd.Run()
	if ctx.Err() != nil {
		t.Fatalf("programa excedeu o limite de 5 segundos: %v", ctx.Err())
	}
	return stdout.String(), stderr.String(), err
}

func conferirErro(t *testing.T, stdout, stderr string, err error) {
	t.Helper()
	if err == nil {
		t.Error("esperava término com código diferente de zero")
	}
	if stdout != "" {
		t.Errorf("não deve emitir relatório em caso de erro; stdout: %q", stdout)
	}
	if strings.TrimSpace(stderr) == "" {
		t.Error("esperava mensagem de erro em stderr")
	}
}

func reais(centavos int64) string {
	return fmt.Sprintf("R$ %d,%02d", centavos/100, centavos%100)
}

func testarVersao(t *testing.T, v2 bool) {
	binario := compilar(t)
	casos := []struct {
		nome                          string
		vendas                        []venda
		quantidade                    int
		total, comissaoV1, comissaoV2 int64
	}{
		{"exemplo_do_enunciado", []venda{{"V1", 10000, "confirmada"}, {"V2", 25000, "confirmada"}, {"V3", 8000, "cancelada"}}, 2, 35000, 1750, 1750},
		{"lista_vazia", []venda{}, 0, 0, 0, 0},
		{"apenas_canceladas", []venda{{"V1", 150000, "cancelada"}}, 0, 0, 0, 0},
		{"fracao_de_centavo", []venda{{"V1", 199, "confirmada"}}, 1, 199, 9, 9},
		{"arredondar_somente_apos_somar", []venda{{"V1", 199, "confirmada"}, {"V2", 199, "confirmada"}}, 2, 398, 19, 19},
		{"abaixo_do_limite", []venda{{"V1", 99999, "confirmada"}}, 1, 99999, 4999, 4999},
		{"exatamente_no_limite", []venda{{"V1", 100000, "confirmada"}}, 1, 100000, 5000, 7000},
		{"acima_do_limite", []venda{{"V1", 150099, "confirmada"}}, 1, 150099, 7504, 10506},
		{"limite_pela_soma", []venda{{"V1", 60000, "confirmada"}, {"V2", 40000, "confirmada"}}, 2, 100000, 5000, 7000},
		{"canceladas_nao_ativam_taxa_maior", []venda{{"V1", 90000, "confirmada"}, {"V2", 20000, "cancelada"}}, 1, 90000, 4500, 4500},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			destino := filepath.Join(t.TempDir(), "novo", "fechamentos")
			dados := entrada{"F-001", "Ana", caso.vendas}
			comissao := caso.comissaoV1
			if v2 {
				comissao = caso.comissaoV2
			}
			stdout, stderr, err := executar(t, binario, destino, dados)
			if err != nil {
				t.Fatalf("esperava sucesso: %v; stderr: %s", err, stderr)
			}
			if stderr != "" {
				t.Errorf("stderr inesperado: %q", stderr)
			}
			esperado := fechamento{dados.ID, dados.Vendedor, caso.quantidade, caso.total, comissao}
			conteudo, err := os.ReadFile(filepath.Join(destino, dados.ID+".json"))
			if err != nil {
				t.Fatal(err)
			}
			var recebido fechamento
			if err := json.Unmarshal(conteudo, &recebido); err != nil {
				t.Fatalf("JSON inválido: %v", err)
			}
			var campos map[string]json.RawMessage
			if err := json.Unmarshal(conteudo, &campos); err != nil {
				t.Fatal(err)
			}
			for _, campo := range []string{"id", "vendedor", "quantidade_vendas_confirmadas", "total_centavos", "comissao_centavos"} {
				if valor, ok := campos[campo]; !ok || string(valor) == "null" {
					t.Errorf("campo obrigatório ausente ou nulo: %s", campo)
				}
			}
			if recebido != esperado {
				t.Errorf("JSON: recebido %+v; esperado %+v", recebido, esperado)
			}
			relatorio := fmt.Sprintf("Fechamento: %s\nVendedor: %s\nVendas confirmadas: %d\nTotal vendido: %s\nComissão: %s\n", dados.ID, dados.Vendedor, caso.quantidade, reais(caso.total), reais(comissao))
			if stdout != relatorio {
				t.Errorf("relatório recebido: %q\nesperado: %q", stdout, relatorio)
			}
		})
	}

	t.Run("entradas_invalidas", func(t *testing.T) {
		casos := []struct {
			nome    string
			alterar func(*entrada)
		}{
			{"fechamento_vazio", func(e *entrada) { e.ID = "" }},
			{"fechamento_em_branco", func(e *entrada) { e.ID = "   " }},
			{"vendedor_vazio", func(e *entrada) { e.Vendedor = "" }},
			{"vendedor_em_branco", func(e *entrada) { e.Vendedor = " \t " }},
			{"venda_sem_id", func(e *entrada) { e.Vendas[1].ID = "" }},
			{"valor_zero", func(e *entrada) { e.Vendas[1].Valor = 0 }},
			{"valor_negativo", func(e *entrada) { e.Vendas[1].Valor = -1 }},
			{"status_desconhecido", func(e *entrada) { e.Vendas[1].Status = "pendente" }},
			{"status_vazio", func(e *entrada) { e.Vendas[1].Status = "" }},
			{"id_duplicado", func(e *entrada) { e.Vendas[1].ID = "V1" }},
			{"cancelada_invalida", func(e *entrada) { e.Vendas[1].Status = "cancelada"; e.Vendas[1].Valor = 0 }},
			{"cancelada_sem_id", func(e *entrada) { e.Vendas[1].Status = "cancelada"; e.Vendas[1].ID = "" }},
			{"id_duplicado_em_cancelada", func(e *entrada) { e.Vendas[1].Status = "cancelada"; e.Vendas[1].ID = "V1" }},
		}
		for _, caso := range casos {
			t.Run(caso.nome, func(t *testing.T) {
				destino := t.TempDir()
				dados := entrada{"F-001", "Ana", []venda{{"V1", 10000, "confirmada"}, {"V2", 20000, "confirmada"}}}
				caso.alterar(&dados)
				stdout, stderr, err := executar(t, binario, destino, dados)
				conferirErro(t, stdout, stderr, err)
				arquivos, err := os.ReadDir(destino)
				if err != nil {
					t.Fatal(err)
				}
				if len(arquivos) != 0 {
					t.Errorf("entrada inválida deixou arquivos no destino: %v", arquivos)
				}
			})
		}
	})

	t.Run("nao_sobrescrever_arquivo", func(t *testing.T) {
		destino := t.TempDir()
		arquivo := filepath.Join(destino, "F-001.json")
		original := []byte("conteudo preexistente que deve permanecer intacto\n")
		if err := os.WriteFile(arquivo, original, 0600); err != nil {
			t.Fatal(err)
		}
		dados := entrada{"F-001", "Ana", []venda{{"V1", 100000, "confirmada"}}}
		stdout, stderr, err := executar(t, binario, destino, dados)
		conferirErro(t, stdout, stderr, err)
		depois, err := os.ReadFile(arquivo)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(original, depois) {
			t.Error("arquivo existente foi alterado")
		}
	})

	t.Run("falha_ao_gravar", func(t *testing.T) {
		// Um arquivo no lugar de um diretório falha também quando o teste roda como root.
		destino := filepath.Join(t.TempDir(), "arquivo")
		original := []byte("isto nao e um diretorio")
		if err := os.WriteFile(destino, original, 0600); err != nil {
			t.Fatal(err)
		}
		dados := entrada{"F-001", "Ana", []venda{{"V1", 100000, "confirmada"}}}
		stdout, stderr, err := executar(t, binario, destino, dados)
		conferirErro(t, stdout, stderr, err)
		depois, err := os.ReadFile(destino)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(original, depois) {
			t.Error("arquivo que bloqueava o destino foi alterado")
		}
	})
}
