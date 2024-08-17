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
	checkboxesDias := page.MustElements(".af_selectBooleanCheckbox_native-input")
	for _, checkbox := range checkboxesDias {
		checkbox.MustClick()
	}
	println("Campos seleccionados")

	// Hacer clic en el botón para ejecutar la búsqueda
	page.MustElement(".af_button_link").MustClick()
	size := len(page.MustWaitStable().MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr"))

	println("Asignaturas encontradas: ", size)
	println()
	println()

	var dataAsignaturas []Asignatura = make([]Asignatura, size)

	var tiemposTotales []time.Duration = make([]time.Duration, size)
	var tiemposCarga []time.Duration = make([]time.Duration, size)
	var tiemposExtraccion []time.Duration = make([]time.Duration, size)

	startTime := time.Now()

	// Recorrer asignaturas
	for i := 0; i < size; i++ {
		println(i, " / ", size)

		timeTotal := time.Now()

		timeLoad := time.Now()

		timeListadAsignaturas := time.Now()
		asignaturas := page.MustWaitStable().MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr") // Delay
		fmt.Printf("Tiempo listado asignaturas: %s\n", time.Since(timeListadAsignaturas))

		// Cargar asignatura
		asignatura := asignaturas[i]

		timeNave := time.Now()

		link := asignatura.MustElement(".af_commandLink")
		link.MustClick()

		page.MustWaitStable() // Delay
		fmt.Printf("Tiempo navegacion: %s\n", time.Since(timeNave))

		timefinLoad := time.Since(timeLoad)

		// Extraer datos
		timeExtraccion := time.Now()
		data := procesarMateria(page)
		timefinExtraccion := time.Since(timeExtraccion)
		dataAsignaturas[i] = data

		// Regresar
		backButton := page.MustElement(".af_button")
		backButton.MustClick()

		timefinTotal := time.Since(timeTotal)

		println(data.Nombre, data.Codigo)

		fmt.Printf("Tiempo carga: %s\n", timefinLoad)
		fmt.Printf("Tiempo extraccion: %s\n", timefinExtraccion)
		fmt.Printf("Tiempo total: %s\n", timefinTotal)

		tiemposTotales[i] = timefinTotal
		tiemposCarga[i] = timefinLoad
		tiemposExtraccion[i] = timefinExtraccion

		println("")

	}

	elapsedTime := time.Since(startTime)

	promedioTotal := 0.0
	promedioCarga := 0.0
	promedioExtraccion := 0.0

	for i := 0; i < size; i++ {
		promedioTotal += (tiemposTotales[i].Seconds())
		promedioCarga += (tiemposCarga[i].Seconds())
		promedioExtraccion += (tiemposExtraccion[i].Seconds())
	}

	println(".............................")

	fmt.Printf("Tiempo de ejecución: %s\n", elapsedTime)
	println("--- Tiempos promedios ---")
	println("Promedio total: ", (promedioTotal / float64(size)))
	println("Promedio carga: ", (promedioCarga / float64(size)))
	println("Promedio extraccion: ", (promedioExtraccion / float64(size)))

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
