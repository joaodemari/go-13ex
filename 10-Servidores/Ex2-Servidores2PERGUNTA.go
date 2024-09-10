// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// servidor com criacao dinamica de thread de servico
// Problema:
//   considere um servidor que recebe pedidos por um canal (representando uma conexao)
//   ao receber o pedido, sabe-se através de qual canal (conexao) responder ao cliente.
//   Abaixo uma solucao sequencial para o servidor.
// Exercicio
//   deseja-se tratar os clientes concorrentemente, e nao sequencialmente.
//   como ficaria a solucao ?
// Veja abaixo a resposta ...
//   quantos clientes podem estar sendo tratados concorrentemente ?
//
// Exercicio:
//   agora suponha que o seu servidor pode estar tratando no maximo 10 clientes concorrentemente.
//   como voce faria ?
//

// POR João Pedro Demari, Athos Endele Puna e Thales Xavier

// Como foi mostrado na pasta ./3-CanaisEBufferizacao, adicionamos um canal com buffer de tamanho 10. Assim, esse canal
// permite que apenas 10 requisições sejam tratadas concorrentemente, porque, ao tratar as 10 requisições, o buffer
// do canal fica "cheio" e precisa esperar que um espaço no buffer do canal libere.

package main

import (
	"fmt"
	"math/rand"
)

const (
	NCL  = 100
	Pool = 10
)

type Request struct {
	v      int
	ch_ret chan int
}

// ------------------------------------
// cliente
func cliente(i int, req chan Request) {
	var v, r int
	my_ch := make(chan int)
	for {
		v = rand.Intn(1000)
		req <- Request{v, my_ch}
		r = <-my_ch
		fmt.Println("cli: ", i, " req: ", v, "  resp:", r)
	}
}

// ------------------------------------
// servidor
// thread de servico calcula a resposta e manda direto pelo canal de retorno informado pelo cliente
func trataReq(id int, req Request) {
	fmt.Println("                                 trataReq ", req)
	req.ch_ret <- req.v * 2
}

func mandaReq(req Request, bloq chan struct{}) {
	trataReq(rand.Intn(100), req)
	<-bloq
}

// servidor que dispara threads de servico
func servidorConc(in chan Request) {
	bloqueio := make(chan struct{}, Pool) // Canal com buffer para limitar o número de goroutines

	for req := range in {
		bloqueio <- struct{}{} // Bloqueia se o buffer estiver cheio, liberando até 10 goroutines simultâneas
		go mandaReq(req, bloqueio)
	}
}

// ------------------------------------
// main
func main() {
	fmt.Println("------ Servidores - criacao dinamica -------")
	serv_chan := make(chan Request) // CANAL POR ONDE SERVIDOR RECEBE PEDIDOS
	go servidorConc(serv_chan)      // LANÇA PROCESSO SERVIDOR
	for i := 0; i < NCL; i++ {      // LANÇA DIVERSOS CLIENTES
		go cliente(i, serv_chan)
	}
	<-make(chan int)
}
