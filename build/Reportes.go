package build

import (
	"fmt"
	"log"

	"github.com/goccy/go-graphviz/cgraph"
)

func reportes() {
	// Crear un nuevo gráfico Graphviz
	var graph cgraph.Graph
	graph = cgraph.Graph(graph)

	// Crear un nuevo nodo en el gráfico
	node, err := graph.CreateNode("table")
	if err != nil {
		log.Fatal(err)
	}

	// Especificar la forma del nodo como "plaintext" para que pueda contener HTML
	node.Set("shape", "plaintext")

	// Especificar la tabla usando la etiqueta del nodo
	node.Set("label", fmt.Sprintf(`
		<<table border="1" cellborder="1" cellspacing="0">
			<tr><td>Header 1</td><td>Header 2</td></tr>
			<tr><td>Row 1, Col 1</td><td>Row 1, Col 2</td></tr>
			<tr><td>Row 2, Col 1</td><td>Row 2, Col 2</td></tr>
		</table>>
	`))

	// Guardar el gráfico en un archivo
	if err := graph.RenderFilename("tabla.png", "png", "dot"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tabla creada correctamente en tabla.png")
}
