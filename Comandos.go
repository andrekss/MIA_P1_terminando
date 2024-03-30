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
	types := flag.String("type", "P", "Tipo de partición (opcional, P/E/L)") // Extendida trae un EBR
	fit := flag.String("fit", "WF", "Ajuste Partición")                      // BF FF WF
	delete := flag.String("delete", "", "Eliminar partición")                // Se usa junto a name, size y driveletter
	add := flag.Int("add", 0, "Agregar espacio")                             // agregar o quitar espacio en una partición se usa con unit driveletter y name

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

	// aqui ejecutamos todo el comando
	build.Funcionalidades(*driveletter, *delete, name1, Partición, ruta, *add, *unit)

	fmt.Printf("Creando una partición con tamaño: %d bytes\n", *size)
}

func Mount() { // go run *.go mount -driveletter=A -name=Particion1
	driveletter := flag.String("driveletter", "", "Archivo a elegir")
	name := flag.String("name", "", "Nombre de la partición")
	flag.CommandLine.Parse(os.Args[2:])

	// Abrir archivo
	ruta := "./MIA/P1/" + strings.ToUpper(*driveletter) + ".dsk"
	file, err := build.AbrirArchivo(ruta)
	if err != nil {
		return
	}

	var MBR structs.MBR

	if err := build.LeerArchivo(file, &MBR, 0); err != nil {
		return
	}

	var index int = -1
	var count = 0
	// buscamos la partición especidifica
	for i := 0; i < 4; i++ {
		if MBR.Partitions[i].Size != 0 {
			count++
			if strings.Contains(string(MBR.Partitions[i].Name[:]), *name) {
				index = i
				break
			}
		}
	}

	// id = DriveLetter + Correlative + 19

	id := strings.ToUpper(*driveletter) + strconv.Itoa(count) + "80"

	copy(MBR.Partitions[index].Status[:], "1") // verifica que esta montada
	copy(MBR.Partitions[index].Id[:], id)

	if err := build.Escribir(file, MBR, 0); err != nil {
		return
	}

	/*
	      solo para imprimir
	   	var MBR2 structs.MBR
	   	if err := build.LeerArchivo(file, &MBR2, 0); err != nil {
	   		return
	   	}*/

	fmt.Print("Se a montado la partición" + *name)

	defer file.Close()
}

func Unmount() { // go run *.go Unmount -id=A180
	id := flag.String("id", "", "id de la particioón")
	flag.CommandLine.Parse(os.Args[2:])

	i := 0
	for {
		ruta := "./MIA/P1/" + string(build.Alfabeto[i]) + ".dsk"
		var revision structs.MBR
		file, err := build.AbrirArchivo(ruta)
		if err != nil {
			break
		}

		if err := build.LeerArchivo(file, &revision, 0); err != nil {
			return
		}

		for j := 0; j < 4; j++ { // recorremos las particiones
			var idd [4]byte
			copy(idd[:], *id)

			if revision.Partitions[j].Id == idd {
				copy(revision.Partitions[j].Status[:], "0") // indicamos que se desmontó
				build.Escribir(file, revision, 0)

				defer file.Close()
				fmt.Print("Se a desmontó la partición " + string(revision.Partitions[j].Name[:]))
				break
			}
		}
		i += 1
	}
}

func MKfs() {

	id := flag.String("id", "", "id partición montada")
	//types := flag.String("type", "Full", "tipo de formateo")
	//fs := flag.String("fs", "2fs", "Formateo a otro sistema")

	flag.CommandLine.Parse(os.Args[2:])
	i := 0
	for {
		ruta := "./MIA/P1/" + string(build.Alfabeto[i]) + ".dsk"
		var revision structs.MBR
		file, err := build.AbrirArchivo(ruta)
		if err != nil {
			break
		}

		if err := build.LeerArchivo(file, &revision, 0); err != nil {
			return
		}

		for j := 0; j < 4; j++ { // particiones
			var idd [4]byte
			copy(idd[:], *id)

			if revision.Partitions[j].Id == idd {
			}
		}
		i += 1
	}
}
