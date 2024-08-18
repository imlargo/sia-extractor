package core

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
)

type Horario struct {
	Inicio string `json:"inicio"`
	Fin    string `json:"fin"`
	Dia    string `json:"dia"`
}

type Grupo struct {
	Grupo    string `json:"grupo"`
	Cupos    int    `json:"cupos"`
	Profesor string `json:"profesor"`
	Duracion string `json:"duracion"`
	Jornada  string `json:"jornada"`
	Horarios []Horario
}

type Asignatura struct {
	Nombre           string `json:"nombre"`
	Codigo           string `json:"codigo"`
	Tipologia        string `json:"tipologia"`
	Creditos         string `json:"creditos"`
	Facultad         string `json:"facultad"`
	FechaExtraccion  string `json:"fechaExtraccion"`
	CuposDisponibles string `json:"cuposDisponibles"`
	Grupos           []Grupo
}

type Codigo struct {
	Nivel     string
	Sede      string
	Facultad  string
	Carrera   string
	Tipologia string
}

type Path struct {
	Nivel     string
	Sede      string
	Facultad  string
	Carrera   string
	Tipologia string
}

var jSExtractorFunctionContent string = ""

const (
	SIA_URL       string = "https://sia.unal.edu.co/Catalogo/facespublico/public/servicioPublico.jsf?taskflowId=task-flow-AC_CatalogoAsignaturas"
	ValueNivel    string = "Pregrado"
	ValueSede     string = "1102 SEDE MEDELLÍN"
	Path_Carreras string = "carreras.json"
	Tipologia_All string = "TODAS MENOS  LIBRE ELECCIÓN"
)

var Paths = Path{
	Nivel:     "#pt1\\:r1\\:0\\:soc1\\:\\:content",
	Sede:      "#pt1\\:r1\\:0\\:soc9\\:\\:content",
	Facultad:  "#pt1\\:r1\\:0\\:soc2\\:\\:content",
	Carrera:   "#pt1\\:r1\\:0\\:soc3\\:\\:content",
	Tipologia: "#pt1\\:r1\\:0\\:soc4\\:\\:content",
}

func procesarMateria(page *rod.Page) Asignatura {

	dataAsignatura := page.MustEval(jSExtractorFunctionContent)

	rawGrupos := dataAsignatura.Get("grupos").Arr()
	println("Grupos: ", len(rawGrupos))

	var grupos []Grupo = make([]Grupo, len(rawGrupos))

	for i, rawGrupo := range rawGrupos {
		println("Grupo: ", rawGrupo.Get("profesor").Str())

		rawHorarios := rawGrupo.Get("horarios").Arr()
		var horarios []Horario = make([]Horario, len(rawHorarios))

		for j, rawHorario := range rawHorarios {
			horarios[j] = Horario{
				Inicio: rawHorario.Get("inicio").Str(),
				Fin:    rawHorario.Get("fin").Str(),
				Dia:    rawHorario.Get("dia").Str(),
			}
		}

		grupos[i] = Grupo{
			Grupo:    rawGrupo.Get("grupo").Str(),
			Cupos:    rawGrupo.Get("cupos").Int(),
			Profesor: rawGrupo.Get("profesor").Str(),
			Duracion: rawGrupo.Get("duracion").Str(),
			Jornada:  rawGrupo.Get("jornada").Str(),
			Horarios: horarios,
		}

	}

	return Asignatura{
		Nombre:           dataAsignatura.Get("nombre").Str(),
		Codigo:           dataAsignatura.Get("codigo").Str(),
		Tipologia:        dataAsignatura.Get("tipologia").Str(),
		Creditos:         dataAsignatura.Get("creditos").Str(),
		Facultad:         dataAsignatura.Get("facultad").Str(),
		FechaExtraccion:  dataAsignatura.Get("fechaExtraccion").Str(),
		CuposDisponibles: dataAsignatura.Get("cuposDisponibles").Str(),
		Grupos:           grupos,
	}

}

func GetAsignaturasCarrera(codigo Codigo) []Asignatura {

	jSExtractorFunctionContent = LoadJSExtractor()

	println("Iniciando...")
	page := rod.New().MustConnect().MustPage(SIA_URL)
	println("Cargado. ok")
	println("")

	page.MustWaitStable().MustElement(Paths.Nivel).MustClick().MustSelect(codigo.Nivel)
	page.MustWaitStable().MustElement(Paths.Sede).MustClick().MustSelect(codigo.Sede)
	page.MustWaitStable().MustElement(Paths.Facultad).MustClick().MustSelect(codigo.Facultad)
	page.MustWaitStable().MustElement(Paths.Carrera).MustClick().MustSelect(codigo.Carrera)
	page.MustWaitStable().MustElement(Paths.Tipologia).MustClick().MustSelect(codigo.Tipologia)

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
		asignaturas := page.MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr") // Delay
		fmt.Printf("Tiempo listado asignaturas: %s\n", time.Since(timeListadAsignaturas))

		// Cargar asignatura
		asignatura := asignaturas[i]

		timeNave := time.Now()

		link := asignatura.MustElement(".af_commandLink")
		link.MustClick()

		// page.MustWaitStable() // Delay
		page.MustElement(".af_showDetailHeader_content0")
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

	println("")
	println("Finalizado")

	return dataAsignaturas

}

func LoadJSExtractor() string {
	content, _ := os.ReadFile("./getData.js")
	JSExtractorFunctionContent := string(content)

	return JSExtractorFunctionContent
}

func CreatePathsCarreras() {

	startTime := time.Now()
	println("Iniciando...")
	page := rod.New().MustConnect().MustPage(SIA_URL)
	println("Cargado. ok")
	println("")

	page.MustWaitStable().MustElement(Paths.Nivel).MustClick().MustSelect(ValueNivel)
	page.MustWaitStable().MustElement(Paths.Sede).MustClick().MustSelect(ValueSede)
	println("Datos iniciales seleccionados")

	selectFacultad := page.MustWaitStable().MustElement(Paths.Facultad)
	listadoFacultades := selectFacultad.MustElements("option")

	var listadoCarrerasSede []map[string]string

	for i, facultad := range listadoFacultades {

		if i == 0 {
			continue
		}

		println("Seleccionando: ", facultad.MustText())

		page.MustWaitStable().MustElement(Paths.Facultad).MustClick().MustSelect(facultad.MustText())

		selectCarrera := page.MustWaitStable().MustElement(Paths.Carrera)
		listadoCarreras := selectCarrera.MustElements("option")

		var carreras []map[string]string = make([]map[string]string, len(listadoCarreras)-1)

		for i, carrera := range listadoCarreras {

			if i == 0 {
				continue
			}

			println(facultad.MustText(), carrera.MustText())

			carreras[i-1] = map[string]string{
				"facultad": facultad.MustText(),
				"carrera":  carrera.MustText(),
			}
		}

		listadoCarrerasSede = append(listadoCarrerasSede, carreras...)

	}

	elapsedTime := time.Since(startTime)

	println(".............................")
	fmt.Printf("Tiempo de ejecución: %s\n", elapsedTime)

	dataCarrerasJSON, _ := json.Marshal(listadoCarrerasSede)
	os.WriteFile(Path_Carreras, dataCarrerasJSON, 0644)

	println("Finalizado!!! :D")
}
