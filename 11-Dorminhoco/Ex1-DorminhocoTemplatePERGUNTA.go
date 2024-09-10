// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// PROBLEMA:
//   o dorminhoco especificado no arquivo Ex1-ExplanacaoDoDorminhoco.pdf nesta pasta
// ESTE ARQUIVO
//   Um template para criar um anel generico.
//   Adapte para o problema do dorminhoco.
//   Nada está dito sobre como funciona a ordem de processos que batem.
//   O ultimo leva a rolhada ...
//   ESTE  PROGRAMA NAO FUNCIONA.    É UM RASCUNHO COM DICAS.

// POR João Pedro Demari, Athos Endele Puna e Thales Xavier

package main

import (
	"fmt"
	"math/rand"
)

const NJ = 5 // numero de jogadores
const M = 4  // numero de cartas

type carta string // carta é um strirng

var ch [NJ]chan carta // NJ canais de itens tipo carta

func remove(slice []carta, s int) []carta {
	return append(slice[:s], slice[s+1:]...)
}

func jogador(id int, in chan carta, out chan carta, cartasIniciais []carta, bateuCh chan int) {
	mao := cartasIniciais   // estado local - as cartas na mao do jogador
	nroDeCartas := len(mao) // quantas cartas ele tem
	// carta recebida é vazia
	acabou := false

	var cartaRecebida carta = " "

	for !acabou {
		if nroDeCartas == (M + 1) { // tenho que jogar
			// fmt.Println(id, " joga") // escreve seu identificador

			// escolhe alguma carta para passar adiante ...
			cartaEscolhidaIndex := rand.Intn(nroDeCartas)

			cartaParaSair := mao[cartaEscolhidaIndex]

			// NÃO PODE PASSAR O CORINGA ASSIM QUE RECEBE
			for cartaParaSair == "@" && cartaRecebida == "@" {
				cartaEscolhidaIndex = rand.Intn(nroDeCartas)

				cartaParaSair = mao[cartaEscolhidaIndex]
			}

			fmt.Printf("jogador %d escolheu carta "+string(cartaParaSair)+" para passar\n", id)

			mao = remove(mao, cartaEscolhidaIndex)

			nroDeCartas = len(mao)
			// fmt.Printf("Eu sou o jogador %d e antes de passar tenho %d cartas na mão\n", id, len(mao))

			bateu := true
			for i := 1; i < len(mao); i++ {
				if mao[0] != mao[i] {
					bateu = false
				}
			}

			if bateu {
				fmt.Printf("MÃO DO JOGADOR %d VENCEDOR\n", id)
				for i := 0; i < len(mao); i++ {
					fmt.Printf(string(mao[i]) + " - ")
				}
				fmt.Print("\n")
				bateuCh <- id
				for i := 0; i < NJ; i++ {
					if i != id {
						ch[i] <- "BATI"
					}
				}
				acabou = true
				break
			} else {
				// manda carta escolhida o proximo
				out <- cartaParaSair
			}

		} else {
			// ...
			// fmt.Printf("Eu sou o jogador %d e estou esperando a carta recebida\n", id)
			cartaRecebida = <-in // recebe carta na entrada
			// fmt.Printf("Eu sou o jogador %d e recebi uma carta\n", id)

			// e se alguem bate ?

			if cartaRecebida == "BATI" {
				acabou = true
				bateuCh <- id
			} else {
				mao = append(mao, cartaRecebida)
				// fmt.Printf("Eu sou o jogador %d e tenho %d cartas na mão\n", id, len(mao))
				nroDeCartas = len(mao)
			}

			// guarda carta que entrou

		}
	}

}

func main() {
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan carta)
	}

	bateu := make(chan int)

	cartasDisponiveis := []carta{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	baralho := [NJ*M + 1]carta{}

	// cria um baralho com NJ*M cartas

	fmt.Printf("---- BARALHO DISPONÍVEL ----")

	for nj := 0; nj < NJ; nj++ {
		cartaASerAdicionada := cartasDisponiveis[nj]
		for m := 0; m < M; m++ {
			baralho[nj*M+m] = cartaASerAdicionada
			fmt.Printf("Carta %d: "+string(cartaASerAdicionada)+"\n", nj*M+m)
		}
	}

	baralho[NJ*M] = "@"
	fmt.Printf("Carta %d: "+string(baralho[NJ*M])+"\n", NJ*M)

	for i := 0; i < NJ; i++ {
		// escolhe aleatoriamente (tira) cartas do baralho, passa cartas para jogador
		cartasEscolhidas := []carta{}

		for m := 0; m < M; m++ {

			cartaASerDadaIndex := rand.Intn(NJ*M + 1)
			var cartaASerDada carta = "#"
			for cartaASerDada == "#" {
				cartaASerDada = baralho[cartaASerDadaIndex]
				if cartaASerDada != "#" {
					cartasEscolhidas = append(cartasEscolhidas, cartaASerDada)
					baralho[cartaASerDadaIndex] = "#"
				} else {
					cartaASerDadaIndex = rand.Intn(NJ*M + 1)
				}
			}
			if i == NJ-1 && m == M-1 {
				cartaASerDadaIndex := rand.Intn(NJ*M + 1)
				var cartaASerDada carta = "#"
				for cartaASerDada == "#" {
					cartaASerDada = baralho[cartaASerDadaIndex]
					if cartaASerDada != "#" {
						cartasEscolhidas = append(cartasEscolhidas, cartaASerDada)
						baralho[cartaASerDadaIndex] = "#"
					} else {
						cartaASerDadaIndex = rand.Intn(NJ*M + 1)
					}
				}
			}
		}

		fmt.Printf("----CARTAS DO JOGADOR %d ----\n", i)

		for m := 0; m < len(cartasEscolhidas); m++ {
			fmt.Printf("Carta %d: "+string(cartasEscolhidas[m])+"\n", m)
		}
		go jogador(i, ch[i], ch[(i+1)%NJ], cartasEscolhidas, bateu) // cria processos conectados circularmente
	}

	for i := 0; i < NJ; i++ {
		jogador := <-bateu
		if i == 0 {
			fmt.Println("ACABOU O JOGO")
		}
		fmt.Printf("%dº lugar: Jogador %d \n", i+1, jogador)
	}

}
