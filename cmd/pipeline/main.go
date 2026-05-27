package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
)

// notas sobre a sintaxe: a struct em minusculo indica é privado, ou seja, não pode ser acessado de fora do pacote.
// criamos a struct Job para encapsula os dados de cada tarefa que os workers vão processar.
type Job struct {
  SourcePath string `json:"source_path"`
  DestPath string `json:"dest_path"`
}

// CONSUMIDOR
// nossa função worker vai receber canal (ch <-chan Job), processa cada Job que o PRODUTOR até que o canal feche.
func consume(id int, ch <-chan Job) {
    for job := range ch {
        log.Print("[CONSUMIDOR %d] Iniciando: %s", id, filepath.Base(job.SourcePath))
        // chama a função que processa o Job
        processImage(job)
        log.Print("[CONSUMIDOR %d] Concluído: %s", id, filepath.Base(job.SourcePath))
    }
}

// PRODUTOR
// essa função (produceJobs) funciona como o produtor. Ela vai percorrer o diretório "origin" encontrar as imagens e enviar as tarefas (jobs) para o canal
func produceJobs(sourceDirectory string, targetDirectory string, ch chan<- Job) {

    err := filepath.Walk(sourceDirectory, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
    
    if !info.IsDir() {
        // Calcula o caminho relativo para replicar a estrutura de pastas no destino
        rel, _ := filepath.Rel(sourceDirectory, path)
        newDestination := filepath.Join(targetDirectory, rel)

        Job := Job{
            SourcePath: path,
            DestPath: newDestination,
            
        }
        // Envia o job estruturado para o canal, bloqueando caso o canal esteja cheio
        ch <- Job
    }
    return nil
    })

    if err != nil {
        log.Printf("Erro ao percorrer o diretório %s: %v", sourceDirectory, err)
    }
}

// processImage realiza a computação pesada (CPU-bound) para redimensiona e converte para tons de cinza as imagens.
func processImage(job Job) {
    image, err := imaging.Open(job.SourcePath)
    if err != nil {
        log.Printf("Erro ao abrir a imagem %s: %v", job.SourcePath, err)
        return 
    }
    // Processamento com algoritmo Lanczos para manter a qualidade
    resizeImage := imaging.Resize(image, 48, 48, imaging.Lanczos)
    resizeImage = imaging.Grayscale(resizeImage)

    // Garante que a pasta de destino exista antes de salvar
    os.MkdirAll(filepath.Dir(job.DestPath), os.ModePerm)

    err = imaging.Save(resizeImage, job.DestPath)
    if err != nil {
        log.Printf("Erro ao salvar a imagem %s: %v", job.DestPath, err)
    }
}


func main() {
    // Cria um canal unbuffered para comunicação segura entre as Goroutines
    ch := make(chan Job)
    var wg sync.WaitGroup

    // Inicializa um pool de 5 workers concorrentes (Goroutines)
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            consume(id, ch)
        }(i)
    }

    // Executa o Produtor em uma Goroutine separada (background).
	// Isso evita o Deadlock, pois permite que os workers comecem a consumir 
	// as imagens enquanto o diretório ainda está sendo varrido.
    go  func(){
        log.Println("[Produtor] Iniciando a busca por arquivos em 'origin'...")
        produceJobs("data/origin", "data/destination", ch)
        log.Println("[Produtor] Todos os arquivos foram listados. Fechando canal.")
        close(ch) // Fecha o canal para avisar os workers que a produção acabou.
    }() // Nota de sintaxe: parênteses obrigatórios para invocar a função anônima imediatamente.

  
    // Aguarda todas as Goroutines do pool terminarem o processamento
    wg.Wait()
    log.Println("Processamento concluído com sucesso!")
}
