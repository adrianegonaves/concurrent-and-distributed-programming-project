package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
)

// o struct em minusculo indica é privado, ou seja, não pode ser acessado de fora do pacote
// Exemplo de como calcular o destino dinamicamente:


type Job struct {
  ImageURLOrigen string `json:"image_url"`
  ImageURLDestino string `json:"image_url_destino"`
}

// nossa função worker vai receber o job como parametro
func worker(id int, ch <-chan Job) {
    for job := range ch {
        // chamamos a função no worker para ela pegar um job
        processImage(job)
    }
    
}

// essa função funciona como o produtor 
func pathToFile(sourceDir string, targetDir string, ch chan<- Job) {
    err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
    
    if !info.IsDir() {
        rel, _ := filepath.Rel(sourceDir, path)
        novoDestino := filepath.Join(targetDir, rel)
        Job := Job{
            ImageURLOrigen: rel,
            ImageURLDestino: novoDestino,
            
        }
        // vamos enviar o job no canal ( recebemos o canal como parametro da função)
        ch <- Job
    }
    return nil
    })
}

//  criamos a função processImage para processar as imagens para ficar com o tamanho 48x48 e na escala cinza
func processImage(job Job) {
    image, err := imaging.Open(job.ImageURLOrigen)
    if err != nil {
        log.Printf("Erro ao abrir a imagem %s: %v", job.ImageURLOrigen, err)
        return 
    }
    imageRedimensiona := imaging.Resize(image, 48, 48, imaging.Lanczos)
    imageRedimensiona = imaging.Grayscale(imageRedimensiona)
    os.MkdirAll(filepath.Dir(job.ImageURLDestino), os.ModePerm)
    err = imaging.Save(imageRedimensiona, job.ImageURLDestino)
    if err != nil {
        log.Printf("Erro ao salvar a imagem %s: %v", job.ImageURLDestino, err)
    }
}


func main() {
    // o canal é criado para enviar e receber mensagens do tipo Job
    image := make(chan Job)
    
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)
        // incia a goroutine para cada worker
        go func(id int) {
            defer wg.Done()
            worker(id, image)
        }(i)
    }

    // funciona como um produtor para enviar o job
    pathToFile('origem', 'destino', image)

    // fechamos o canal quando os workers terminarem de processar os jobs, para evitar deadlocks
    close(image)

    wg.Wait()

}
