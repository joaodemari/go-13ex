// João Pedro Tiellet Demari, Athos Endele Puna, Thales Xavier 
// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// EXERCÍCIO:  dado o programa abaixo
//    1) quantos processos concorrentes são gerados ?
// 		São gerados N processos, que no contexto desse código, são 40.
//    2) execute e observe: que se pode supor sobre a velocidade relativa dos mesmos?
//    Não se pode supor nada sobre a velocidade relativa, tendo em vista que é aleatório a ordem
// dos processos
// OBSERVACAO:o sleep no método main serve para este nao acabar, o que acabaria todos processos em execucao.
//     mais adiante veremos outras formas de sincronizar isto

package main

import (
	"fmt"
	"time"
)

var N int = 40

func funcaoA(id int, s string) {
	for {
		fmt.Println(s, id)
	}
}

func geraNespacos(n int) string {
	s := "  "
	for j := 0; j < n; j++ {
		s = s + "   "
	}
	return s
}

func main() {
	for i := 0; i < N; i++ {
		go funcaoA(i, geraNespacos(i))
	}
	for true {
		time.Sleep(100 * time.Millisecond)
	}
}
