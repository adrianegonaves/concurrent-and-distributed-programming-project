package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

func findSampleImage(originDir string) (string, error) {
	var found string
	err := filepath.Walk(originDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || found != "" {
			return err
		}
		ext := filepath.Ext(path)
		if !info.IsDir() && (ext == ".jpg" || ext == ".jpeg" || ext == ".png") {
			found = path
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if found == "" {
		return "", fmt.Errorf("nenhuma imagem encontrada em %s", originDir)
	}
	return found, nil
}

func originDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "../../data/origin"))
}

// BenchmarkProcessImage mede a performance processando uma única imagem repetidamente.
// O objetivo é testar a eficiência da CPU/algoritmo, isolando variações de tamanho de arquivo.
func BenchmarkProcessImage(b *testing.B) {
	// Busca a imagem de amostra no diretório de origem
	sampleImage, err := findSampleImage(originDir())
	if err != nil {
		b.Fatalf("Pré-condição falhou: %v", err)
	}

	// Cria o cenário do Job com a imagem fixa
	job := Job{
		SourcePath: sampleImage,
		DestPath:   filepath.Join(os.TempDir(), "bench_single"+filepath.Ext(sampleImage)),
		Size:       48,
	}

	b.ResetTimer() // Ignora o tempo de carregamento inicial no resultado do benchmark
	// Executa o processamento da mesma imagem N vezes para calcular a média de tempo
	for i := 0; i < b.N; i++ {
		if err := processImage(job, job.Size); err != nil {
			b.Errorf("processImage falhou: %v", err)
		}
	}
}

func BenchmarkPipelineNWorkers(b *testing.B) {
	src := originDir()
	if _, err := os.Stat(src); os.IsNotExist(err) {
		b.Skipf("Pasta de origem não encontrada (%s) — a saltar benchmark de pipeline", src)
	}

	workerCounts := []int{1, 2, 4, runtime.NumCPU()}
	seen := map[int]bool{}
	unique := []int{}
	for _, n := range workerCounts {
		if !seen[n] {
			seen[n] = true
			unique = append(unique, n)
		}
	}

	for _, n := range unique {
		n := n
		b.Run(fmt.Sprintf("workers=%d", n), func(b *testing.B) {
			destDir := filepath.Join(os.TempDir(), fmt.Sprintf("bench_pipeline_%d_workers", n))
			defer os.RemoveAll(destDir)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
					b.Fatal(err)
				}

				ch := make(chan Job, n*2)
				var wg sync.WaitGroup
				stats := &Stats{}

				for id := 1; id <= n; id++ {
					wg.Add(1)
					go func(id int) {
						defer wg.Done()
						consume(id, ch, stats)
					}(id)
				}

				
				total := produceJobs(src, destDir, 48, ch)
				stats.total.Store(int64(total))
				close(ch)
				wg.Wait()

				os.RemoveAll(destDir)
			}
		})
	}
}