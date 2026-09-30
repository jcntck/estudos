# Laboratório de SRP — Single Responsibility Principle

## Objetivo

Praticar o Princípio da Responsabilidade Única (SRP) em Go: um módulo deve ter um único motivo para mudar. Responsabilidade, aqui, está ligada às regras e às pessoas ou áreas que solicitam mudanças; não significa simplesmente ter uma única função ou poucas linhas de código.

O laboratório contém o enunciado e testes automatizados de integração. Você vai implementar a v1, experimentar uma mudança de requisito na v2 ainda com o código acoplado e só depois refatorar com SRP. Nenhuma implementação do programa é fornecida.

## Organização

- `problema/`: implemente a v1 e depois evolua esse mesmo código acoplado para a v2.
- `solucao/`: desenvolva a versão refatorada aplicando SRP.

- `integracao_test.go`: suíte que executa o programa real, verifica a saída e lê os arquivos gravados, sem mocks.

As pastas de implementação começam vazias. Os testes ficam na raiz do módulo e podem avaliar qualquer uma delas.

## Enunciado: fechamento de vendas

Uma pequena loja precisa automatizar o fechamento das vendas de um vendedor. Hoje, uma única rotina recebe as vendas, valida os dados, calcula a comissão, salva o fechamento em disco e produz um relatório textual para o gerente.

Essa rotina atende a três áreas:

- **Comercial:** define quais vendas entram no fechamento e como a comissão é calculada.
- **Operações:** define o formato e o local de armazenamento dos fechamentos.
- **Gestão:** define quais informações aparecem no relatório e como são apresentadas.

O problema surge quando pedidos dessas áreas exigem alterações na mesma rotina. Uma mudança na apresentação do relatório, por exemplo, obriga a mexer no mesmo trecho que calcula os valores financeiros.

Seu desafio é tornar essas mudanças independentes usando SRP, sem alterar as regras iniciais.

## Regras da primeira versão

### Entrada e cálculo

Cada fechamento recebe um identificador, o nome do vendedor e uma lista de vendas. Cada venda possui um identificador, um valor inteiro em centavos e um status: `confirmada` ou `cancelada`.

- O identificador do fechamento e o nome do vendedor são obrigatórios e não podem conter apenas espaços.
- Cada venda deve ter identificador não vazio, valor maior que zero e um dos dois status permitidos.
- Identificadores de vendas não podem se repetir dentro do mesmo fechamento.
- Valide todas as vendas, inclusive as canceladas. Uma entrada inválida deve encerrar o processamento com erro, sem salvar ou emitir relatório de sucesso.
- Apenas vendas confirmadas entram no total e na quantidade de vendas do fechamento.
- A comissão é de 5% sobre o total das vendas confirmadas. Descarte eventual fração de centavo depois de calcular a comissão sobre o total.
- Uma lista vazia ou composta apenas por vendas canceladas gera um fechamento válido com quantidade, total e comissão iguais a zero.

### Armazenamento

Salve um arquivo JSON em um diretório informado pelo chamador, usando o identificador do fechamento como nome do arquivo. Para este exercício, considere identificadores de fechamento compostos apenas por letras, números e hífens.

O arquivo deve conter: identificador do fechamento, vendedor, quantidade de vendas confirmadas, total em centavos e comissão em centavos. Crie o diretório se necessário. Se o arquivo já existir, retorne erro sem sobrescrevê-lo.

Uma falha ao salvar deve ser comunicada ao chamador e impedir a emissão do relatório de sucesso.

### Relatório

Depois de salvar com sucesso, produza um texto com uma informação por linha, nesta ordem: identificador do fechamento, vendedor, quantidade de vendas confirmadas, total vendido e comissão. Apresente os valores monetários em reais, com duas casas decimais e vírgula como separador decimal.

Não é necessário enviar e-mail, criar uma API ou integrar serviços externos. Para permitir os testes de integração, exponha o programa pelo contrato de terminal descrito abaixo.

## Fase 1 — Implementar a v1 acoplada

Em `problema/`, crie um programa Go executável (`package main`). Concentre em uma única rotina a validação, o cálculo, a criação do JSON, a gravação e a composição do relatório. O propósito é observar decisões que mudam por razões diferentes convivendo no mesmo fluxo.

Faça todos os testes de integração da v1 passarem. A taxa nessa versão é sempre de 5%, inclusive para fechamentos de R$ 1.000,00 ou mais. Antes de seguir, preserve essa versão em um commit, se quiser poder revisitá-la.

## Fase 2 — A mudança que gera a v2

O Comercial pediu uma nova regra:

- Total confirmado abaixo de R$ 1.000,00: comissão de 5%.
- Total confirmado igual ou superior a R$ 1.000,00: comissão de 7% **sobre o total inteiro**, não apenas sobre o excedente.
- Vendas canceladas não participam da soma nem da escolha da taxa.
- Continue descartando a fração de centavo somente depois de calcular a comissão sobre o total.

Todo o restante permanece como na v1: validações, armazenamento, formato do relatório e tratamento de erros.

Primeiro rode os testes da v2 contra sua v1: os cenários que alcançam o limite devem falhar. Depois altere a implementação em `problema/`, ainda acoplada, até passar na suíte v2. Observe quais trechos precisou entender e modificar e o risco de afetar o armazenamento ou o relatório ao mudar uma regra comercial.

As suítes representam requisitos diferentes: não é esperado que a v2 passe nos testes da v1 para valores a partir do limite. Não crie uma opção de versão no programa para satisfazer ambas ao mesmo tempo.

## Depois das duas fases — Refatorar com SRP

Em `solucao/`, refatore a **v2** para que regras comerciais, armazenamento e apresentação possam evoluir independentemente. Você decide os nomes, os tipos e os limites dos componentes. Mantenha o contrato externo e faça a mesma suíte v2 passar nessa pasta.

Evite apenas distribuir trechos em arquivos diferentes mantendo as mesmas dependências entre decisões. Não é obrigatório criar uma interface para cada componente.

Após refatorar, explique onde alteraria o programa caso Operações pedisse outro meio de armazenamento ou a Gestão pedisse um relatório HTML. Esses dois pedidos são apenas perguntas de reflexão; não precisam ser implementados.

## Contrato externo usado pelos testes

Use o módulo Go já existente. Cada pasta de implementação deve poder ser compilada separadamente com `go build ./problema` ou `go build ./solucao`. Os testes compilam o alvo e executam seu binário uma vez por cenário.

O programa deve:

1. Receber o diretório de saída pelo argumento `-destino <diretorio>`.
2. Ler um objeto JSON da entrada padrão (`stdin`) com este formato:

```json
{
  "id": "F-001",
  "vendedor": "Ana",
  "vendas": [
    {"id": "V1", "valor_centavos": 10000, "status": "confirmada"},
    {"id": "V2", "valor_centavos": 25000, "status": "confirmada"},
    {"id": "V3", "valor_centavos": 8000, "status": "cancelada"}
  ]
}
```

3. Em caso de sucesso, salvar `<diretorio>/F-001.json` com os campos abaixo. A ordem dos campos e a indentação são livres; os cinco campos devem existir mesmo quando os valores forem zero.

```json
{
  "id": "F-001",
  "vendedor": "Ana",
  "quantidade_vendas_confirmadas": 2,
  "total_centavos": 35000,
  "comissao_centavos": 1750
}
```

4. Após salvar, escrever exatamente este formato de relatório em `stdout`, incluindo uma quebra de linha após a última linha, e encerrar com código zero. Não use separador de milhares.

```text
Fechamento: F-001
Vendedor: Ana
Vendas confirmadas: 2
Total vendido: R$ 350,00
Comissão: R$ 17,50
```

5. Em caso de erro, encerrar com código diferente de zero, deixar `stdout` vazio e escrever uma mensagem não vazia em `stderr`. O texto do erro é livre. Em caso de sucesso, deixe `stderr` vazio.

O contrato define somente a fronteira do programa. Os testes não exigem nomes de funções, interfaces ou uma arquitetura interna específica. Entradas dos testes são JSON bem-formado; não é necessário implementar interatividade ou imprimir prompts.

## Como rodar os testes

No terminal, entre na pasta do laboratório:

```sh
cd solid-single-responsability-principle
```

**Fase 1**, avaliando a v1 em `problema/`:

```sh
go test -v -count=1 -run '^TestIntegracaoV1$' .
```

**Fase 2**, avaliando a evolução para v2 na mesma pasta:

```sh
go test -v -count=1 -run '^TestIntegracaoV2$' .
```

**Após refatorar**, avaliando a v2 em `solucao/`:

```sh
SRP_IMPL=solucao go test -v -count=1 -run '^TestIntegracaoV2$' .
```

Os comandos são para um terminal Linux/macOS. `SRP_IMPL` assume `problema` quando omitido; `-count=1` evita resultados do cache. Use o filtro da fase desejada: `go test ./...` executa as duas suítes, cujas expectativas de comissão são incompatíveis para o mesmo programa.

Antes de você implementar o programa, a suíte falha na compilação do alvo. Isso é esperado: não há código de aplicação ou implementação provisória. Depois que a compilação passar, os testes verificarão as regras e os efeitos reais do programa.

Cada cenário usa um diretório temporário, removido automaticamente pelo Go, e cada execução do programa tem limite de cinco segundos. A suíte cobre cálculo, limites de comissão, truncamento, cancelamentos, entradas inválidas, criação de diretórios, preservação de arquivo existente e falha de gravação.

Os testes verificam o comportamento externo; passar neles não demonstra, por si só, que SRP foi aplicado. A avaliação da separação das responsabilidades usa os critérios abaixo.

## Exemplos para conferir o comportamento

Considere um fechamento válido para a vendedora Ana, com estas vendas:

| Venda | Valor | Status |
| --- | ---: | --- |
| V1 | R$ 100,00 | confirmada |
| V2 | R$ 250,00 | confirmada |
| V3 | R$ 80,00 | cancelada |

O resultado deve conter **2 vendas confirmadas**, **R$ 350,00 de total** e **R$ 17,50 de comissão**. No JSON, total e comissão devem ser, respectivamente, `35000` e `1750` centavos.

Confira também:

- Uma venda confirmada de 199 centavos gera comissão de 9 centavos.
- Lista vazia e lista contendo apenas vendas canceladas geram valores zerados.
- Valor inválido, status desconhecido ou identificador duplicado gera erro sem arquivo nem relatório de sucesso.
- Arquivo já existente permanece intacto e o processamento retorna erro.
- Falha de gravação não produz relatório de sucesso.

Use diretórios separados ao comparar as versões, para que um arquivo produzido por uma delas não impeça a execução da outra.

## Critérios de conclusão

- A v1 passa na suíte da fase 1; a v2 acoplada passa na suíte da fase 2.
- A v2 refatorada em `solucao/` passa na mesma suíte v2, preservando o comportamento da versão acoplada.
- Na solução, é possível verificar o cálculo sem gravar arquivos ou gerar relatórios.
- Uma mudança de apresentação não exige alterar a regra de comissão.
- Uma mudança de armazenamento não exige alterar o cálculo nem a formatação do relatório.
- Existe um ponto claro que coordena a sequência do fluxo e propaga falhas.
- Você consegue explicar o motivo de mudança de cada componente e relacioná-lo à área responsável.

Ao terminar, registre uma breve comparação: quais responsabilidades estavam juntas, como você as separou e quais custos ou benefícios essa separação trouxe.

## Referência

Material de apoio indicado originalmente neste laboratório: https://www.youtube.com/watch?v=EWHTE1dQM4U
