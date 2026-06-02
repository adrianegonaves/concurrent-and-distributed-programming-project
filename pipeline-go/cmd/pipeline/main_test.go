package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func BenchmarkProcessImage(b *testing.B) {
	// 1. Descobre a pasta real de origem
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	originDir := filepath.Clean(filepath.Join(currentDir, "../../data/origin"))

	// 2. Procura a primeira imagem válida dentro da vossa pasta real
	var sampleImage string
	err := filepath.Walk(originDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Se não for pasta e for um ficheiro (ex: .jpg, .png), escolhemos este
		if !info.IsDir() && (filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".png" || filepath.Ext(path) == ".jpeg") {
			sampleImage = path
			return filepath.SkipDir // Encontrou uma, pode parar de procurar
		}
		return nil
	})

	// Proteção: Se a vossa pasta estiver vazia, o teste avisa
	if err != nil || sampleImage == "" {
		b.Fatalf("Não foram encontradas imagens válidas na pasta %s para correr o benchmark!", originDir)
	}

	// 3. Cria o Job com uma imagem real do vosso dataset
	job := Job{
		SourcePath: sampleImage,
		DestPath:   filepath.Join(os.TempDir(), "teste_out" + filepath.Ext(sampleImage)),
	}

	// Executa o benchmark b.N vezes com a imagem real
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processImage(job)
	}
}