# Programação Concorrente e Distribuída

**Trabalho Final**  
Este programa implementa um pipeline utilizando o padrão Produtor-Consumidor, suportado pela concorrência nativa do Go (Goroutines e Channels). 

---

##  Equipa
*   **Adriane Gonçalves** - 240000004
*   **Bruno Hortelão** - 240001083
---
---

## Checklist de Implementação

- [✅]. Go Pipeline
- [✅]. Implemntação de Teste (Benchmark e profiling)

---
# Instruções de Uso

### 1. Preparação do Código
* Faça o **download** do código para a sua máquina.

### 2. Configuração do Arquivo `.gitignore`
Antes de rodar os scripts, certifique-se de que possui um arquivo chamado `.gitignore` na raiz do projeto com o seguinte conteúdo para evitar o envio acidental de arquivos pesados ao GitHub:

```gitignore
# Ignorar pastas de dados geradas localmente
pipeline-go/data/origin/
pipeline-go/data/destination/

# Executáveis compilados
*.exe
pipeline-go/pipeline-go
```

### 2. Configuração das Pastas
* Crie uma pasta chamada `origin` dentro do diretório `pipeline-go/data`

Rode o comando abaixo para fazer download das imagens

```
python script-download.py
```

### 3. Execução
* Rode o projeto utilizando o comando:

```
go run main.go
```

### 4. Benchmark

Cenário A: Rodar apenas o Teste Rápido (Uma única imagem)

```
go test -bench=BenchmarkProcessImage -benchmem
```

Cenário B: Rodar o Teste de Concorrência (Pipeline Completo)

```
go test -bench=BenchmarkPipelineNWorkers -benchmem -benchtime=1x
```

---
---

## Tecnologias Utilizadas e Fontes

[Golang com struct, funções e metodos](https://tomelin-tech.medium.com/golang-com-struct-funções-e-metodos-698a25b6221a)
[Go by Example: Channels](https://gobyexample.com/channels)
[filepath](https://pkg.go.dev/path/filepath#example-Base)
[imaging](https://pkg.go.dev/github.com/disintegration/imaging#section-readme)
[Golang: Desmistificando channels - Buffered Channels](https://dev.to/igormelo/golang-desmistificando-channels-buffered-channels-16d0)
[como-fazer-benchmark-do-seu-codigo](https://ttemporin.dev/como-fazer-benchmark-do-seu-codigo/)
---

