package build

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

/*
seek(tamaño, posición)
si tamaño = 0 <--- se sobrescriba
*/

var TamañoMBR int = 159
var TamañoEbr int = 30

var Alfabeto = [26]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
	'Y', 'Z'}

var Letra int = 0

var NombreArchivo string = "MIA/P1/" + string(Alfabeto[Letra]) + ".dsk"

func Analizar(file *os.File) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		comando := strings.TrimSpace(scanner.Text())
		if comando[0] != '#' {
			comando = "go run *.go" + comando
			// Ejecutar el comando en la consola
			cmd := exec.Command("bash", "-c", comando)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	}

	if err := scanner.Err(); err != nil {
		return
	}
}
