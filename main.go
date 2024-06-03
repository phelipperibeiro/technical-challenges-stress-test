package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// Declaração de variáveis globais para armazenar os argumentos da linha de comando.
var (
	url         string
	requests    int
	concurrency int
)

func main() {
	// Cria o comando raiz usando o pacote cobra.
	rootCmd := &cobra.Command{
		Use:   "technical-challenges-stress-test", // Nome do comando.
		Short: "Stress testing tool",              // Descrição curta do comando.
		Run:   executeLoadTest,                    // Função a ser executada quando o comando é chamado.
	}

	// Define as flags da linha de comando para o comando raiz.
	rootCmd.Flags().StringVar(&url, "url", "", "URL of the service to test")
	rootCmd.Flags().IntVar(&requests, "requests", 1, "Total number of requests")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 1, "Number of concurrent requests")

	// Marca a flag "url" como obrigatória.
	if err := rootCmd.MarkFlagRequired("url"); err != nil {
		log.Fatalf("url flag is required: %v", err)
	}

	// Executa o comando raiz.
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("cmd.Execute() failed: %v", err)
	}
}

// Função chamada quando o comando raiz é executado.
func executeLoadTest(cmd *cobra.Command, args []string) {
	// Verifica se a concorrência é maior que o número de solicitações.
	if concurrency > requests {
		fmt.Println("Error: Concurrency cannot be greater than the number of requests.")
		os.Exit(1)
	}
	// Executa o teste de carga.
	runLoadTest(url, requests, concurrency)
}

// Função que executa o teste de carga.
func runLoadTest(url string, requests, concurrency int) {
	var waitGroup sync.WaitGroup // Define um WaitGroup para sincronizar as goroutines.
	var mutex sync.Mutex         // Define um Mutex para proteger o acesso aos recursos compartilhados.

	startTime := time.Now() // Marca o início do teste de carga.

	statusCodes := make(map[int]int) // Mapa para contar os códigos de status HTTP.
	totalRequests := 0               // Contador de solicitações totais.

	sem := make(chan struct{}, concurrency) // Canal para limitar a concorrência.
	done := make(chan struct{})             // Canal para sinalizar o término do teste.

	// Goroutine para imprimir uma mensagem de progresso.
	go func() {
		fmt.Println("Making requests, please wait....")
	}()

	// Loop para iniciar as goroutines que farão as solicitações.
	for i := 0; i < requests; i++ {
		waitGroup.Add(1)  // Incrementa o contador do WaitGroup.
		sem <- struct{}{} // Envia um sinal para o canal de concorrência.

		// Inicia uma goroutine para fazer a solicitação.
		go func() {
			defer waitGroup.Done()   // Decrementa o contador do WaitGroup ao finalizar.
			defer func() { <-sem }() // Libera um espaço no canal de concorrência.

			resp, err := http.Get(url) // Faz a solicitação HTTP.
			mutex.Lock()               // Bloqueia o Mutex para acessar os recursos compartilhados.
			defer mutex.Unlock()       // Desbloqueia o Mutex ao finalizar.

			if err != nil {
				handleError(err, statusCodes) // Trata o erro, se houver.
			} else {
				processResponse(resp, statusCodes) // Processa a resposta.
			}
			totalRequests++ // Incrementa o contador de solicitações totais.
		}()
	}

	// Goroutine para fechar o canal done quando todas as goroutines terminarem.
	go func() {
		waitGroup.Wait() // Aguarda todas as goroutines finalizarem.
		close(done)      // Fecha o canal done.
	}()

	<-done // Aguarda o sinal de término.

	totalTime := time.Since(startTime)                    // Calcula o tempo total do teste.
	generateReport(totalRequests, statusCodes, totalTime) // Gera o relatório.
}

// Função para tratar erros de solicitação.
func handleError(err error, statusCodes map[int]int) {
	if os.IsTimeout(err) {
		statusCodes[404]++ // Incrementa o contador de erros 404.
	} else {
		statusCodes[500]++ // Incrementa o contador de erros 500.
	}
}

// Função para processar a resposta HTTP.
func processResponse(resp *http.Response, statusCodes map[int]int) {
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err) // Loga o erro ao fechar o corpo da resposta.
		}
	}(resp.Body)

	statusCodes[resp.StatusCode]++ // Incrementa o contador do código de status da resposta.
}

// Função para gerar o relatório do teste de carga.
func generateReport(totalRequests int, statusCodes map[int]int, totalTime time.Duration) {
	fmt.Println("Report: ")
	fmt.Printf("Total requests: %d\n", totalRequests) // Imprime o total de solicitações.
	fmt.Printf("Time taken: %s\n", totalTime)         // Imprime o tempo total do teste.
	fmt.Println("Status code distribution:")

	// Loop para imprimir a distribuição dos códigos de status.
	for code, count := range statusCodes {
		fmt.Printf("[%d] %d requests \n", code, count)
	}
}
