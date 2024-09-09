// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// ------- Voce deve ter feito o Ex4 desta serie para entender e continuar aqui ------

// PROBLEMA:
//          como proximo passo, implemente que, durante a inundação de ida, a rota vai sendo
//          gravada.  A rota é a sequencia de nodos por onde a mensagem passa.
//          Pode ser uma pilha de inteiros.  Cada nodo antes de repassar, empilha seu id.
//          Desta forma, a resposta pode ser enviada somente pela rota de retorno.
//          Ou seja, a mensagem trafega pela rota reversa.  Basta que cada nodo intermediario
//          desempilhe o identificador do proximo e repasse a mensagem para este.
// ATENCAO
//          o codigo abaixo NAO apresenta a solucao.  é uma cópia do anterior.
//          VOCE DEVE DESENVOLVER ESTA SOLUCAO.
//          este desenvolvimento pode ser com seu grupo.
//          deverá ser entregue em data marcada.

package main

import (
	"fmt"
)

// nro de nodos
const N = 10

// tamanho do buffer de cada canal de entrada
const channelBufferSize = 5

// Topologia é uma matriz NxN onde 1 em [i][j] indica presenca da aresta do nodo i para o j
type Topology [N][N]int

// O que é enviado entre nodos, pode adicionar campos nesta estrutura ...
// O que é enviado entre nodos
type Message struct {
	id       int // identificador da mensagem - um nro sequencial ...
	source   int // nodo de origem
	receiver int // nodo de destino
	data     string
	route    []int // pilha de rotas (sequência de IDs)
}

// um canal de entrada para cada nodo i
type inputChan [N]chan Message

type nodeStruct struct {
	id               int
	topo             Topology
	inCh             inputChan
	received         map[int]Message // repassadas
	receivedMessages []Message       // destino
}

// difusão ou broadcast - um nodo manda para TODOS seus vizinhos do grafo
// nodo origin, conforme topology, usando canais do vetor inCh, manda message para todos eles
func (n *nodeStruct) broadcast(m Message) { // broadcast(origin int, topo Topology, inCh inputChan, m Message) {
	for j := 0; j < N; j++ { // para todo vizinho j em N
		if n.topo[n.id][j] == 1 { //  a matriz em [origin][j] diz se origin conectado com j
			n.inCh[j] <- m // escreve m no canal de j
			// fmt.Println("      ", n.id, " -> ", j, " final ", m.receiver)
		}
	}
}

// cada nodo recebe toda matriz de conectividade e os canais de entrada de todos processos
// cada nodo le o seu canal de entrada e escreve a mensagem em todos canais de saida
// (dele para outros nodos usando a funcao send)
func (n *nodeStruct) nodo() {
	fmt.Println(n.id, "ativo!")
	for {
		m := <-n.inCh[n.id] // espera entrada

		if m.receiver == n.id { // a mensagem é para mim
			n.receivedMessages = append(n.receivedMessages, m)
			if m.id > 0 { // se mensagem de ida, responde. senão não. mensagem resposta tem id negativo
				fmt.Println("                                   ", n.id, " recebe de ", m.source, "msg ", m.id, "  ", m.data)
				// enviar resposta pela rota inversa
				if len(m.route) > 0 {
					nextNode := m.route[len(m.route)-1] // próximo nodo na rota (retorno) primeiro nodo da pilha
					m.route = m.route[:len(m.route)-1]  // remove o primeiro nodo da pilha
					go func() {
						n.inCh[nextNode] <- Message{
							id:       -m.id,
							source:   n.id,
							receiver: m.source,
							data:     "resp to msg",
							route:    m.route,
						}
					}()
				}
			} else {
				fmt.Println("                                                                      ", n.id, " recebe de ", m.source, "msg ", m.id, "  ", m.data)
			}
		} else { // não é para mim ... repassar se for a primeira vez
			if _, achou := n.received[m.id]; !achou { // não achou = não recebi a msg antes
				n.received[m.id] = m // guarda para saber no futuro

				// Empilha o ID do nodo atual na rota antes de repassar
				m.route = append(m.route, n.id)

				go n.broadcast(m) // repassa a primeira vez
			}
		}
	}
}

// ----
// carga nos nodos

func carga(nodoInicial int, inCh chan Message) {
	for i := 1; i < 10; i++ { // gera mensagem de teste a cada segundo
		inCh <- Message{
			id:       (nodoInicial * 1000) + i,
			source:   nodoInicial,
			receiver: i,
			data:     "req",
			route:    []int{},
		}
		// time.Sleep(20 * time.Millisecond)
	}
}

// ------------------------------------------------------------------------------------------------
// no main: montagem da topologia, criacao de canais, inicializacao de nodos e geracao de mensagens
// ------------------------------------------------------------------------------------------------
func main() {
	var topo Topology
	//  se [i,j]==1, então o nodo i pode enviar para o nodo j pelo canal j.
	//  para alterar a topologia basta adicionar 1s.  cada 1 é uma aresta direcional.
	//  para modelar comunicação em ambas direções entre i e j, então [i,j] e [j,i] devem ser 1
	topo = [N][N]int{
		// conforme algoritmo na função "nodo"
		//  0  1  2  3  4  5  6  7  8  9       aresta de    para - BIDIRECIONAIS
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 0           0 - 1
		{1, 0, 1, 0, 0, 0, 0, 0, 0, 0}, // 1           1 - 2
		{0, 1, 0, 1, 0, 0, 0, 0, 0, 0}, // 2           2 - 3
		{0, 0, 1, 0, 1, 0, 0, 0, 0, 0}, // 3           3 - 4
		{0, 0, 0, 1, 0, 1, 0, 0, 0, 1}, // 4           4 - 5 e  4 - 9
		{0, 0, 0, 0, 1, 0, 1, 0, 0, 0}, // 5           5 - 6
		{0, 0, 0, 0, 0, 1, 0, 1, 0, 0}, // 6           6 - 7
		{0, 0, 0, 0, 0, 0, 1, 0, 1, 0}, // 7           7 - 8
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 1}, // 8           8 - 9
		{0, 0, 0, 0, 0, 1, 0, 0, 1, 0}, // 9           9 - 4
	}

	var inCh inputChan // cada nodo i tem um canal de entrada, chamado inCh[i]
	for i := 0; i < N; i++ {
		inCh[i] = make(chan Message, channelBufferSize) // criando cada um dos canais
	}

	// lança todos os nodos
	for id := 0; id < N; id++ {
		n := nodeStruct{id, topo, inCh, make(map[int]Message), []Message{}}
		go n.nodo() // por simplicidade todos nodos têm acesso à mesma topologia (que não muda)
	}

	// carga de mensagens para que sejam "inundadas" na rede
	go carga(0, inCh[0])
	carga(5, inCh[5])

	<-make(chan struct{}) // bloqueia senão nodos acabam
}
