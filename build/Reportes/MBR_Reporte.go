package Reportes

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"main.go/build"
	"main.go/structs"
)

func CreandoEstructura(NombreReporte, id string) [][]string {
	var tabla [][]string
	switch {
	case strings.EqualFold(NombreReporte, "MBR"):

		tabla = append(tabla, []string{"Reporte MBR"})
		ruta := "./MIA/P1/" + string(id[0]) + ".dsk"
		var revision structs.MBR
		file, err := build.AbrirArchivo(ruta)
		if err != nil {
			break
		}
		if err := build.LeerArchivo(file, &revision, 0); err != nil {
			break
		}

		tabla = append(tabla, []string{"Tamaño", strconv.FormatInt(int64(revision.Tamaño), 10)})
		tabla = append(tabla, []string{"Fecha de creación", string(revision.Fecha[:])})
		tabla = append(tabla, []string{"Signature", strconv.FormatInt(int64(revision.Signature), 10)})

		for i := 0; i < 4; i++ { // 4 particiones
			nameBytes := bytes.TrimRight(revision.Partitions[i].Name[:], "\x00")

			tabla = append(tabla, []string{"----Partición----"})
			tabla = append(tabla, []string{"Status", string(revision.Partitions[i].Status[0])})
			tabla = append(tabla, []string{"Tipo", string(revision.Partitions[i].Tipo[0])})
			tabla = append(tabla, []string{"Fit", string(revision.Partitions[i].Fit[0])})
			tabla = append(tabla, []string{"Start", strconv.FormatInt(int64(revision.Partitions[i].Start), 10)})
			tabla = append(tabla, []string{"Size", strconv.FormatInt(int64(revision.Partitions[i].Size), 10)})
			tabla = append(tabla, []string{"Nombre", string(nameBytes)})
		}

		return tabla
	}
	return tabla
}

func Reportes(tableData [][]string, path string) {
	// Crear una cadena DOT que representa la tabla
	dot := "digraph G {\n"
	dot += "rankdir=\"LR\";\n"
	dot += "node [shape=plaintext]\n"
	dot += "node [fontname=\"Arial\"]\n"
	dot += "edge [style=invis]\n"
	dot += "tbl [label=<<table border=\"1\" cellspacing=\"0\">\n"

	// Agregar filas a la tabla
	for _, row := range tableData {
		dot += "<tr>"
		for _, cell := range row {
			dot += fmt.Sprintf("<td>%s</td>", cell)
		}
		dot += "</tr>\n"
	}
	dot += "</table>>, ];\n"
	dot += "}\n"

	// Generar el archivo DOT
	err := EscribirArchivo(".dot", dot)
	if err != nil {
		log.Fatal("Error al escribir el archivo DOT:", err)
	}

	// Ejecutar Graphviz para generar la imagen PNG
	err = GenerarPNG(path, path) // generación png ej: reporte.png
	if err != nil {
		log.Fatal("Error al generar la imagen PNG:", err)
	}

	fmt.Println("Tabla generada exitosamente. Archivo PNG: tabla.png")
}

func EscribirArchivo(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func GenerarPNG(inputFile, outputFile string) error {
	cmd := exec.Command("dot", "-Tpng", inputFile, "-o", outputFile)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
