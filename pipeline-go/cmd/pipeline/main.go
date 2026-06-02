package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/disintegration/imaging"
)

// notas sobre a sintaxe: a struct em minusculo indica é privado, ou seja, não pode ser acessado de fora do pacote.
// Job encapsula os dados de cada tarefa que os workers vão processar.
type Job struct {
  SourcePath string `json:"source_path"`
  DestPath string `json:"dest_path"`
  Size int
}

// Stats agrupa as métricas de execução do pipeline (processado/falho/total)
type Stats struct {
	processed atomic.Int64
	failed    atomic.Int64
	total     atomic.Int64
}

// CONSUMIDOR
// nossa função worker vai receber canal (ch <-chan Job), processa cada Job que o PRODUTOR até que o canal feche.
func consume(id int, ch <-chan Job, stats *Stats) {
	for job := range ch {
		log.Printf("[WORKER %d] Iniciando: %s", id, filepath.Base(job.SourcePath))
		if err := processImage(job, job.Size); err != nil {
			log.Printf("[WORKER %d] Erro em %s: %v", id, filepath.Base(job.SourcePath), err)
			stats.failed.Add(1)
		} else {
			stats.processed.Add(1)
			log.Printf("[WORKER %d] Concluído: %s", id, filepath.Base(job.SourcePath))
		}
	}
}

// PRODUTOR
// essa função (produceJobs) funciona como o produtor. Ela vai percorrer o diretório "origin" encontrar as imagens e enviar as tarefas (jobs) para o canal

func produceJobs(sourceDirectory string, targetDirectory string, targetSize int, ch chan<- Job) int {
	count := 0
	err := filepath.Walk(sourceDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
 
		ext := filepath.Ext(path)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".bmp" {
			return nil
		}
 
		// Replica a estrutura de subpastas no destino
        // Calcula o caminho relativo para replicar a estrutura de pastas no destino, tbm garante que a estrutura original de subpastas seja clonada no final.
        rel, _ := filepath.Rel(sourceDirectory, path) // descobre o caminho relativo do arquivo
        newDestination := filepath.Join(targetDirectory, rel) // Junta o caminho relativo à pasta de destino.
	
		// depois que descobrir onde as imagens estão e para onde vão motamos a estrutura Job contendo o caminho de origem (SourcePath) e o de destino (DestPath).
        Job := Job{
            SourcePath: path,
            DestPath: newDestination,
            Size: targetSize,
        }
        // Envia o job estruturado para o canal, bloqueando caso o canal esteja cheio
        ch <- Job

		count++
		return nil
	})
 
	if err != nil {
		log.Printf("[PRODUTOR] Erro ao percorrer %s: %v", sourceDirectory, err)
	}
	return count
}


// processImage realiza a computação pesada (CPU-bound) para redimensiona e converte para tons de cinza as imagens.
// retorna erro para realizar tramento de erro 
func processImage(job Job, size int) error {
    image, err := imaging.Open(job.SourcePath)
    if err != nil {
       return fmt.Errorf("Erro ao abrir a imagem: %w", err)
    }
    // Processamento com algoritmo Lanczos para manter a qualidade, isso foi feito antes de melhorar a performace 
    //resizeImage := imaging.Resize(image, size, size, imaging.Lanczos)
    //Para melhorar a perfomace usamos o CatmullRom
    resizeImage := imaging.Resize(image, size,size, imaging.CatmullRom)
    grayImage := imaging.Grayscale(resizeImage)

    // Garante que a pasta de destino exista antes de salvar e faz tratamento de erro
    if err := os.MkdirAll(filepath.Dir(job.DestPath), os.ModePerm); err != nil {
		return fmt.Errorf("criar diretório de destino: %w", err)
	}
 
	if err := imaging.Save(grayImage, job.DestPath); err != nil {
		return fmt.Errorf("salvar imagem: %w", err)
	}
	return nil
}


func main() {
	originDir := flag.String("origin", "../../data/origin", "pasta com as imagens originais")
	destDir := flag.String("dest", "../../data/destination", "pasta de destino para as imagens processadas")
	numWorkers := flag.Int("workers", runtime.NumCPU(), "número de workers concorrentes")
	imgSize := flag.Int("size", 48, "tamanho de saída (largura e altura em px)")
	flag.Parse()

    if err := os.MkdirAll(*destDir, os.ModePerm); err != nil {
		log.Fatalf("Erro ao criar diretório de destino: %v", err)
	}

    // Cria um canal com buffered  para comunicação segura entre as Goroutines, com o buffered canal pode fazer numWorkers*2 envios sem bloquear a goroutine.
    ch := make(chan Job, *numWorkers*2)
    var wg sync.WaitGroup
    stats := &Stats{}

    // Inicia os workers (consumidores)
    // O número de workers é pego de forma dinâmica com base em quantos núcleos lógicos a CPU possui.
    for i := 1; i <= *numWorkers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            consume(id, ch, stats)
        }(i)
    }

    // Executa o produtor em goroutine separada para evitar deadlock
	// Isso evita o Deadlock, pois permite que os workers comecem a consumir 
	// as imagens enquanto o diretório ainda está sendo varrido.
    start := time.Now()
    go  func(){
        total := produceJobs(*originDir, *destDir,*imgSize, ch)
		stats.total.Store(int64(total))
        close(ch) // Fecha o canal para avisar os workers que a produção acabou.
    }() // Nota de sintaxe: parênteses obrigatórios para invocar a função anônima imediatamente.

  
    // Aguarda todas as Goroutines do pool terminarem o processamento
    wg.Wait()
    

    elapsed := time.Since(start)
	total := stats.total.Load()
	processed := stats.processed.Load()
	failed := stats.failed.Load()
 
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("Processamento concluído em %s", elapsed.Round(time.Millisecond))
	log.Printf("Total encontrado : %d imagens", total)
	log.Printf("Processadas      : %d", processed)
	log.Printf("Falhas           : %d", failed)
	if total > 0 {
		log.Printf("Throughput       : %.1f imagens/s", float64(processed)/elapsed.Seconds())
	}
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

