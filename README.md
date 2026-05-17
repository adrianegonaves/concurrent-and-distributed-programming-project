# Programação Concorrente e Distribuída

**Trabalho Final**  
Este programa implementa um pipeline utilizando o padrãoProdutor-Consumidor, suportado pela concorrência nativa do Go (Goroutines e Channels). O objetivo principal é acelerar a preparação do dataset FER-2013 antes do treinamento do modelo de Inteligência Artificial (projeto da UC de Introdução à Inteligência Artificial)

---

## 👥 Equipa
*   **Adriane Gonçalves** - 240000004
*   **Bruno Hortelão** - 240001083
---
---

## Checklist de Implementação

- [✅]. Go Pipeline 
- [ ]. Implemntação de Teste

---
# Instruções de Uso

### 1. Preparação do Código
* Faça o **download** do código para a sua máquina.
* Crie um arquivo chamado `.gitignore` na raiz do projeto.

### 2. Configuração das Pastas
* Crie uma pasta chamada `origin` dentro do diretório `pipeline-go`.
* Coloque as fotos que deseja processar dentro dessa pasta `origin`.

### 3. Execução
* Rode o projeto utilizando o comando `go run main.go`.

---

### Configuração do arquivo `.gitignore`
Adicione o seguinte conteúdo no seu arquivo `.gitignore` para evitar o envio de arquivos pesados ao GitHub:

```
# Ignorar pastas de dados do FER-2013
origin/
destination/
*.keras

```
---

## Tecnologias Utilizadas e Fontes

[Golang com struct, funções e metodos](https://tomelin-tech.medium.com/golang-com-struct-funções-e-metodos-698a25b6221a)
[Go by Example: Channels](https://gobyexample.com/channels)
(filepath)[https://pkg.go.dev/path/filepath#example-Base]
(imaging)[https://pkg.go.dev/github.com/disintegration/imaging#section-readme]
(Golang: Desmistificando channels - Buffered Channels)[https://dev.to/igormelo/golang-desmistificando-channels-buffered-channels-16d0]
---

