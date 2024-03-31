package build

import (
	"encoding/binary"
	"fmt"
	"os"
)

func IndiceAlfabeto(letra string) int {
	for i := 0; i < 26; i++ {
		if Alfabeto[i] == byte(letra[0]) {
			return i
		}
	}
	return 0

}

func CrearArchivo(TamanhoArchivo int) *os.File {
	for { // cliclo sirve para rotular con letras los discos
		// Verificar si el archivo ya existe
		_, err := os.Stat(NombreArchivo)
		if !os.IsNotExist(err) {
			Letra += 1
			NombreArchivo = "MIA/P1/" + string(Alfabeto[Letra]) + ".dsk"
			continue
		}

		Arch, err := os.Create(NombreArchivo) // valor ya existente en el main
		if err != nil {
			fmt.Println(err)
			return Arch
		}

		// Llenar el archivo con la cantidad de 0s especificada
		for i := 0; i < TamanhoArchivo; i++ {
			err := Escribir(Arch, byte(0), int64(i))
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
		//ceros := make([]byte, TamanhoArchivo)
		//arch.Write(ceros)

		fmt.Printf("Archivo creado con %d bytes\n", int(TamanhoArchivo))
		return Arch

	}
}

/*
func Escribir(nuevo interface{}) { // Interface permite qu ecepte cualquier objeto

	// abrir archivo
	file, err := os.OpenFile(NombreArchivo, os.O_RDWR, 0644)

	if err != nil {
		fmt.Println(err)
		return
	}

	file.Seek(0, 0)

	binary.Write(file, binary.LittleEndian, &nuevo)

	defer file.Close()

	fmt.Println("Se ha escrito en el archivo")

}
*/
func Escribir(file *os.File, data interface{}, position int64) error {

	//file.Seek(0, io.SeekCurrent) // Obtener la posición inicial
	//imprimirPos(*file)
	file.Seek(position, 0)

	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err objeto", err)
		return err
	}

	//imprimirPos(*file)
	//file.Seek(0, io.SeekCurrent) // Obtener la posición final

	return nil
}
