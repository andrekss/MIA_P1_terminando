package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"main.go/build"
	"main.go/structs"
)

func MKDisk() { // go run *.go mkdisk -size=30 -fit=FF -unit=K
	var aviso bool = true

	size := flag.Int("size", 0, "Tamaño del disco")
	fit := flag.String("fit", "FF", "Ajuste")                             //opcional
	unit := flag.String("unit", "M", "Unidad del tamaño (opcional, K/M)") // opcional

	flag.CommandLine.Parse(os.Args[2:])
	Ajuste := [3]string{"BF", "FF", "WF"}
	// Verificar si se proporciona el tamaño
	if *size <= 0 {
		fmt.Println("Error: Debes proporcionar un tamaño válido para el disco.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, valor := range Ajuste {
		if *fit == valor {
			aviso = false
			break
		}
	}
	if aviso {
		fmt.Println("Error no existe este ajuste")
		flag.PrintDefaults()
		os.Exit(1)
	}

	*size = build.Conversion(*unit, *size)

	fmt.Printf("Creando disco con tamaño: %d bytes\n", *size)
	// Aquí puedes agregar lógica adicional para crear el disco con el tamaño y la unidad proporcionados.

	Arch := build.CrearArchivo(*size)

	var fitSlice []byte
	fitSlice = append(fitSlice, (*fit)[0])

	var tecnica structs.MBR
	//tecnica.Fecha = time.Now()
	//tecnica.fit = (*fit)[0]

	tecnica.Tamaño = int32(*size)
	tecnica.Signature = int32(build.Letra)

	copy(tecnica.Fit[:], fitSlice)
	copy(tecnica.Fecha[:], time.Now().String())

	// Write object in bin file
	if err := build.Escribir(Arch, tecnica, 0); err != nil {
		return
	}
	defer Arch.Close()

}

func RMDisk() { // go run *.go rmdisk -driveletter=A

	driveletter := flag.String("driveletter", "", "Borrar disco")

	flag.CommandLine.Parse(os.Args[2:])
	var nombre string = "MIA/P1/" + *driveletter + ".dsk"
	var Verificación string

	fmt.Print("¿Quiere confirmar esta accion? (tabule V para confirmar, si no tabule cualquier otro): ")
	fmt.Scan(&Verificación)

	if strings.EqualFold(Verificación, "V") {
		build.Existencia((*driveletter)[0])
		build.Eliminar(nombre)
	} else {
		fmt.Print("No se eliminó ningun archivo")
	}
}

func FDisk() { //go run *.go fdisk -size=300 -driveletter=A -name=Particion1 -unit=B

	size := flag.Int("size", 0, "Tamaño de la partición")
	driveletter := flag.String("driveletter", "", "Archivo a elegir")
	name := flag.String("name", "", "Nombre de la partición")
	// opcional
	unit := flag.String("unit", "K", "Unidad del tamaño (opcional, B/K/M)")
	types := flag.String("type", "P", "Tipo de partición (opcional, P/E/L)")
	fit := flag.String("fit", "WF", "Ajuste Partición")       // BF FF WF
	delete := flag.String("delete", "", "Eliminar partición") // Se usa junto a name, size y driveletter
	add := flag.Int("add", 0, "Agregar espacio")              // agregar o quitar espacio en una partición se usa con unit driveletter y name

	flag.CommandLine.Parse(os.Args[2:])
	ruta := "./MIA/P1/" + strings.ToUpper(*driveletter) + ".dsk"
	// Verificar si se proporciona el tamaño
	if *size <= 0 {
		fmt.Println("Error: Debes proporcionar un tamaño válido para el disco.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	build.Existencia((*driveletter)[0]) // verificación

	*size = build.Conversion(*unit, *size)

	// partición nueva
	var Partición structs.Partition
	Partición.Size = int32(*size)

	var name1 [16]byte
	copy(name1[:], *name)
	copy(Partición.Name[:], name1[:])

	tip := []byte(*types)
	copy(Partición.Tipo[:], tip)

	fits := []byte(*fit)
	copy(Partición.Fit[:], fits)
	// fin del llenado de la nueva partición

	if strings.EqualFold(*delete, "Full") { // eiminar partición
		build.EliminarParticiones(*driveletter, name1)
	} else if *add == 0 { // agregar partición

		Arch, err := build.AbrirArchivo(ruta)
		if err != nil {
			return
		}

		var Editable structs.MBR
		//fmt.Print(Editable.partitions[0])
		build.LeerArchivo(Arch, &Editable, 0)

		//build.Escribir(Arch, Partición, 1)

		// 4 particiones
		for i := 0; i < 4; i++ {
			if Editable.Partitions[i].Size == 0 {
				Editable.Partitions[i] = Partición // agregamos la partición

				fmt.Println(Editable.Partitions[i])
			}
		}
		build.Escribir(Arch, Editable, 0) // sobrescribimos
		defer Arch.Close()                // cerramos el archivo para todo

	} else { // agregar o quitar espacio
		*add = build.Conversion(*unit, *add)

		Arch, err := build.AbrirArchivo(ruta)
		if err != nil {
			return
		}
		var Editable structs.MBR
		build.LeerArchivo(Arch, &Editable, 0)
		for i := 0; i < 4; i++ {
			if Editable.Partitions[i].Name == *&name1 {
				Editable.Partitions[i].Size += int32(*add)
			}
		}
		build.Escribir(Arch, Editable, 0) // sobrescribimos
		defer Arch.Close()                // cerramos el archivo para todo
	}

	fmt.Printf("Creando una partición con tamaño: %d bytes\n", *size)
}

func Mount() {
	driveletter := flag.String("driveletter", "", "Archivo a elegir")
	name := flag.String("name", "", "Nombre de la partición")
	flag.CommandLine.Parse(os.Args[2:])

	// Abrir archivo
	ruta := "./MIA/P1/" + strings.ToUpper(*driveletter) + ".dsk"
	file, err := build.AbrirArchivo(ruta)
	if err != nil {
		return
	}

	var TempMBR structs.MBR
	// Read object from bin file
	if err := build.LeerArchivo(file, &TempMBR, 0); err != nil {
		return
	}

	var index int = -1
	var count = 0
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			count++
			if strings.Contains(string(TempMBR.Partitions[i].Name[:]), *name) {
				index = i
				break
			}
		}
	}

	// id = DriveLetter + Correlative + 19

	id := strings.ToUpper(*driveletter) + strconv.Itoa(count) + "19"

	copy(TempMBR.Partitions[index].Status[:], "1")
	copy(TempMBR.Partitions[index].Id[:], id)

	// Overwrite the MBR
	if err := build.Escribir(file, TempMBR, 0); err != nil {
		return
	}

	var TempMBR2 structs.MBR
	// Read object from bin file
	if err := build.LeerArchivo(file, &TempMBR2, 0); err != nil {
		return
	}

	// Close bin file
	defer file.Close()
}
