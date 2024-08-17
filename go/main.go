package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
)

/*
type Horario struct {
	inicio string
	fin    string
	dia    string
}

type Grupo struct {
	grupo    string
	cupos    int
	profesor string
	duracion string
	jornada  string
	horarios []Horario
}
*/

type Asignatura struct {
	Nombre           string `json:"nombre"`
	Codigo           string `json:"codigo"`
	Tipologia        string `json:"tipologia"`
	Creditos         string `json:"creditos"`
	Facultad         string `json:"facultad"`
	FechaExtraccion  string `json:"fechaExtraccion"`
	CuposDisponibles string `json:"cuposDisponibles"`
	// grupos           []Grupo
}

type Codigo struct {
	nivel     string
	sede      string
	facultad  string
	carrera   string
	tipologia string
}

type Paths struct {
	nivel     string
	sede      string
	facultad  string
	carrera   string
	tipologia string
}

var JSFunction string = ""

func main() {
	codigo := Codigo{
		nivel:     "Pregrado",
		sede:      "1102 SEDE MEDELLÍN",
		facultad:  "3068 FACULTAD DE MINAS",
		carrera:   "3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA",
		tipologia: "TODAS MENOS  LIBRE ELECCIÓN",
	}

	content, _ := os.ReadFile("./getData.js")
	JSFunction = string(content)

	getAsignaturasCarrera(codigo)
}

func getAsignaturasCarrera(codigo Codigo) {
	url := "https://sia.unal.edu.co/Catalogo/facespublico/public/servicioPublico.jsf?taskflowId=task-flow-AC_CatalogoAsignaturas"

	paths := Paths{
		nivel:     "#pt1\\:r1\\:0\\:soc1\\:\\:content",
		sede:      "#pt1\\:r1\\:0\\:soc9\\:\\:content",
		facultad:  "#pt1\\:r1\\:0\\:soc2\\:\\:content",
		carrera:   "#pt1\\:r1\\:0\\:soc3\\:\\:content",
		tipologia: "#pt1\\:r1\\:0\\:soc4\\:\\:content",
	}

	println("Iniciando...")
	page := rod.New().MustConnect().MustPage(url).MustWaitStable()
	println("Cargado. ok")
	println("")

	page.MustWaitStable().MustElement(paths.nivel).MustClick().MustSelect(codigo.nivel)
	page.MustWaitStable().MustElement(paths.sede).MustClick().MustSelect(codigo.sede)
	page.MustWaitStable().MustElement(paths.facultad).MustClick().MustSelect(codigo.facultad)
	page.MustWaitStable().MustElement(paths.carrera).MustClick().MustSelect(codigo.carrera)
	page.MustWaitStable().MustElement(paths.tipologia).MustClick().MustSelect(codigo.tipologia)

	// select all checkboxes
	checkboxesDias := page.MustWaitStable().MustElements(".af_selectBooleanCheckbox_native-input")
	for _, checkbox := range checkboxesDias {
		checkbox.MustClick()
	}
	println("Campos seleccionados")

	// Hacer clic en el botón para ejecutar la búsqueda
	page.MustWaitStable().MustElement(".af_button_link").MustClick()
	size := len(page.MustWaitStable().MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr"))

	println("Asignaturas encontradas: ", size)
	println()
	println()

	var dataAsignaturas []Asignatura = make([]Asignatura, size)

	startTime := time.Now()

	// Recorrer asignaturas
	for i := 0; i < size; i++ {
		println(i, " / ", size)

		asignaturas := page.MustWaitStable().MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr")

		// Cargar asignatura
		asignatura := asignaturas[i]
		link := asignatura.MustElement(".af_commandLink")
		link.MustClick()

		page.MustWaitStable()

		// Extraer datos
		data := procesarMateria(page)
		dataAsignaturas[i] = data
		println(data.Nombre, data.Codigo)

		// Regresar
		backButton := page.MustElement(".af_button")
		backButton.MustClick()

		page.MustWaitStable()
		println()

	}

	elapsedTime := time.Since(startTime)

	fmt.Printf("Tiempo de ejecución: %s\n", elapsedTime)

	for _, asignatura := range dataAsignaturas {
		println(asignatura.Nombre, asignatura.Codigo)
	}

	// Guardar datos de asignaturas en archivo json
	dataAsignaturasJSON, _ := json.Marshal(dataAsignaturas)
	os.WriteFile("asignaturas.json", dataAsignaturasJSON, 0644)

	println("")
	println("Finalizado")
}

func procesarMateria(page *rod.Page) Asignatura {

	dataAsignatura := page.MustEval(JSFunction)

	return Asignatura{
		Nombre:           dataAsignatura.Get("nombre").Str(),
		Codigo:           dataAsignatura.Get("codigo").Str(),
		Tipologia:        dataAsignatura.Get("tipologia").Str(),
		Creditos:         dataAsignatura.Get("creditos").Str(),
		Facultad:         dataAsignatura.Get("facultad").Str(),
		FechaExtraccion:  dataAsignatura.Get("fechaExtraccion").Str(),
		CuposDisponibles: dataAsignatura.Get("cuposDisponibles").Str(),
	}

}
