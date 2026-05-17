package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
)

// notas sobre a sintaxe: o struct em minusculo indica é privado, ou seja, não pode ser acessado de fora do pacote.
// criamos a struct Job para representar o trabalho que cada worker vai processar.
type Job struct {
  ImageURLOrigin string `json:"image_url"`
  ImageURLDestination string `json:"image_url_destino"`
}

// nossa função worker vai receber o job como parametro ( consumidor) e vai consumir o job.
func worker( ch <-chan Job) {
    for job := range ch {
        // chamamos a função no worker para ela pegar um job
        processImage(job)
    }
    
}

// essa função funciona como o produtor ( recebe diretorio de origem e destino e o canal para enviar os jobs) e vai percorrer o diretório de origem, encontrar as imagens e enviar os jobs para os workers processarem.
func pathToFile(sourceDirectory string, targetDirectory string, ch chan<- Job) {
    err := filepath.Walk(sourceDirectory, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
    
    if !info.IsDir() {
        // calculamos o caminho da imagens
        rel, _ := filepath.Rel(sourceDirectory, path)
        newDestination := filepath.Join(targetDirectory, rel)
        Job := Job{
            ImageURLOrigin: path,
            ImageURLDestination: newDestination,
            
        }
        // vamos enviar o job no canal ( recebemos o canal como parametro da função)
        ch <- Job
    }
    return nil
    })

    if err != nil {
        log.Printf("Erro ao percorrer o diretório %s: %v", sourceDirectory, err)
    }
}

//  criamos a função processImage para processar as imagens para ficar com o tamanho 48x48 e na escala cinza
func processImage(job Job) {
    image, err := imaging.Open(job.ImageURLOrigin)
    if err != nil {
        log.Printf("Erro ao abrir a imagem %s: %v", job.ImageURLOrigin, err)
        return 
    }
    ResizeImage := imaging.Resize(image, 48, 48, imaging.Lanczos)
    ResizeImage = imaging.Grayscale(ResizeImage)

    os.MkdirAll(filepath.Dir(job.ImageURLDestination), os.ModePerm)

    err = imaging.Save(ResizeImage, job.ImageURLDestination)
    if err != nil {
        log.Printf("Erro ao salvar a imagem %s: %v", job.ImageURLDestination, err)
    }
}


func main() {
    // o canal é criado para enviar e receber mensagens do tipo Job ( Job é a struct que criamos)
    ch := make(chan Job)
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            worker(ch)
        }(i)
    }

    // funciona como um produtor para enviar o job, essa função ocorre em segundo plano, ou seja, a função main não espera ela terminar para continuar a execução do código. Isso para não trovar a main.
    go  func(){
        pathToFile("origin", "destination", ch)
        close(ch)
    }() // notas sobre a sintaxe: não esquecer de colocar os parenteses para chamar a função anônima, caso contrário ela não será executada.

    // fechamos o canal quando os workers terminarem de processar os jobs, para evitar deadlocks
    

    wg.Wait()
    log.Println("Processamento concluído com sucesso!")

}
