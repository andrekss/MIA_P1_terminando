package build

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"strings"

	"main.go/structs"
)

func AbrirArchivo(ruta string) (*os.File, error) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir", err)
		return nil, err
	}

	return file, nil
}

func LeerArchivo(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data) // archivo, almacenamiento de datos leídos,objeto a almacenar
	if err != nil {
		//position += 1
		//LeerArchivo(file, data, position)
		fmt.Println("Error al leer archivo", err)
		return err
	}
	return nil
}

func EliminarParticiones(driveletter string, name [16]byte) {

	// Abrir bin file
	ruta := "./MIA/P1" + strings.ToUpper(driveletter) + ".dsk"
	arch, err := AbrirArchivo(ruta)
	if err != nil {
		return
	}
	var Mbr structs.MBR
	// Leer
	if err := LeerArchivo(arch, &Mbr, 0); err != nil { // posición 0 del mbr
		return
	}
	var vacio structs.Partition

	Mbr.Partitions[IndiceByName(Mbr, name)] = vacio
	Escribir(arch, Mbr, 0)
	defer arch.Close()

}

func IndiceByName(par structs.MBR, name [16]byte) int {

	for i := 0; i < 4; i++ {
		if par.Partitions[i].Name == name {
			return i
		}
	}
	return 0
}

func Conversion(unit string, add int) int {
	switch unit { // conversion
	case "K":
		return add * 1024
	case "B":
		return add
	case "M":
		return add * 1024 * 1024
	default:
		fmt.Println("Error: Unidad no válida. Utiliza 'K' o 'M'.")
		flag.PrintDefaults()

		os.Exit(1)
		return 0
	}
}

func Funcionalidades(driveletter string, delete string, name1 [16]byte, Partición structs.Partition, ruta string, add int, unit string) {

	if strings.EqualFold(delete, "Full") { // eiminar partición
		EliminarParticiones(driveletter, name1)
	} else if add == 0 { // agregar partición

		Arch, err := AbrirArchivo(ruta)
		if err != nil {
			return
		}

		var Editable structs.MBR
		//fmt.Print(Editable.partitions[0])
		LeerArchivo(Arch, &Editable, 0)

		//Escribir(Arch, Partición, 1)

		// 4 particiones
		for i := 0; i < 4; i++ {
			if Editable.Partitions[i].Size == 0 {
				Editable.Partitions[i] = Partición // agregamos la partición

				fmt.Println(Editable.Partitions[i])
			}
		}
		Escribir(Arch, Editable, 0) // sobrescribimos

		defer Arch.Close() // cerramos el archivo para todo

	} else { // agregar o quitar espacio
		add = Conversion(unit, add)

		Arch, err := AbrirArchivo(ruta)
		if err != nil {
			return
		}
		var Editable structs.MBR
		LeerArchivo(Arch, &Editable, 0)
		for i := 0; i < 4; i++ {
			if Editable.Partitions[i].Name == *&name1 {
				Editable.Partitions[i].Size += int32(add)
			}
		}
		Escribir(Arch, Editable, 0) // sobrescribimos
		defer Arch.Close()          // cerramos el archivo para todo
	}
}

func Part_Extendida_EBR() {

}
