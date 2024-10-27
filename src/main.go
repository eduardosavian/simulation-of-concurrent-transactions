package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Recursos compartilhados
type Resource struct {
	sync.Mutex
	name string
}

// Thread que representa uma transação
type Transaction struct {
	id       int
	timestamp time.Time
}

// Recursos X e Y compartilhados
var (
	X = Resource{name: "X"}
	Y = Resource{name: "Y"}
)

// Função para simular transações concorrentes tentando acessar X e Y
func transactionRoutine(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Thread T(%d) entrou em execução.\n", transaction.id)

	// Atraso randômico para simular solicitação de bloqueio em momentos aleatórios
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	// Tenta acessar os recursos X e Y
	if tryAccessResources(transaction) {
		fmt.Printf("Thread T(%d) finalizou sua execução.\n", transaction.id)
	} else {
		fmt.Printf("Thread T(%d) foi finalizada devido a deadlock.\n", transaction.id)
	}
}

// Tenta acessar os recursos X e Y utilizando a técnica "wait-die"
func tryAccessResources(transaction *Transaction) bool {
	// Gerando um tempo limite para evitar deadlock (timestamp como "idade" da thread)
	deadlockDetected := false
	for !deadlockDetected {
		// Primeiro tenta obter o bloqueio no recurso X
		if lockResource(&X) {
			fmt.Printf("Thread T(%d) obteve o bloqueio no recurso X.\n", transaction.id)
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

			// Depois tenta obter o bloqueio no recurso Y
			if lockResource(&Y) {
				fmt.Printf("Thread T(%d) obteve o bloqueio no recurso Y.\n", transaction.id)
				// Libera os recursos após o uso
				unlockResource(&Y, transaction)
				unlockResource(&X, transaction)
				return true
			} else {
				// Falha ao obter o recurso Y, libera o recurso X
				unlockResource(&X, transaction)
				fmt.Printf("Thread T(%d) falhou ao obter o recurso Y e liberou o recurso X.\n", transaction.id)
			}
		} else {
			fmt.Printf("Thread T(%d) está esperando pelo recurso X.\n", transaction.id)
		}

		// Espera aleatória antes de tentar novamente
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		// Simula deadlock detection (opcionalmente poderia haver uma lógica mais complexa aqui)
		deadlockDetected = time.Since(transaction.timestamp) > 2*time.Second
	}

	return false
}

// Tenta obter o bloqueio no recurso usando a técnica wait-die
func lockResource(resource *Resource) bool {
	// Técnica "wait-die": a thread mais nova espera; se for mais velha, interrompe a execução
	if resource.TryLock() {
		return true
	} else {
		// Se o recurso já está bloqueado por outra transação, verifica se há deadlock
		// Aqui estamos usando uma lógica simples onde a transação mais nova espera
		return false
	}
}

// Libera o recurso
func unlockResource(resource *Resource, transaction *Transaction) {
	resource.Unlock()
	fmt.Printf("Thread T(%d) liberou o bloqueio no recurso %s.\n", transaction.id, resource.name)
}

func main() {
	// Inicializa o grupo de espera para as threads
	var wg sync.WaitGroup

	// Cria 5 threads concorrentes
	for i := 0; i < 5; i++ {
		wg.Add(1)
		transaction := &Transaction{
			id:        i + 1,
			timestamp: time.Now(),
		}
		go transactionRoutine(transaction, &wg)
	}

	// Espera todas as threads concluírem
	wg.Wait()
	fmt.Println("Simulação finalizada.")
}
