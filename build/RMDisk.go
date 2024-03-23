package build

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func Eliminar(nombreArchivo string) {

	// Intenta eliminar el archivo
	err := os.Remove(nombreArchivo)
	if err != nil {
		fmt.Println("Error al eliminar el archivo:", err)
		return
	}

	fmt.Println("Archivo eliminado:", nombreArchivo)
}

func Existencia(letra byte) {

	carpeta := "./MIA/P1"
	var aviso bool = true
	// Lee la lista de archivos en la carpeta
	archivos, err := ioutil.ReadDir(carpeta)
	if err != nil {
		fmt.Println("Error al leer la carpeta:", err)
		return
	}

	for _, archivo := range archivos {

		if archivo.Name()[0] == letra {
			aviso = false
			break
		}
	}
	if aviso {
		fmt.Println("Error no existe este archivo")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
