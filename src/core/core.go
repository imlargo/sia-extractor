package core

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/ysmood/gson"
)

type Horario struct {
	Inicio string `json:"inicio"`
	Fin    string `json:"fin"`
	Dia    string `json:"dia"`
}

type Grupo struct {
	Grupo    string    `json:"grupo"`
	Cupos    int       `json:"cupos"`
	Profesor string    `json:"profesor"`
	Duracion string    `json:"duracion"`
	Jornada  string    `json:"jornada"`
	Horarios []Horario `json:"horarios"`
}

type Asignatura struct {
	Nombre           string  `json:"nombre"`
	Codigo           string  `json:"codigo"`
	Tipologia        string  `json:"tipologia"`
	Creditos         string  `json:"creditos"`
	Facultad         string  `json:"facultad"`
	FechaExtraccion  string  `json:"fechaExtraccion"`
	CuposDisponibles string  `json:"cuposDisponibles"`
	Grupos           []Grupo `json:"grupos"`
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
	SIA_URL          string = "https://sia.unal.edu.co/Catalogo/facespublico/public/servicioPublico.jsf?taskflowId=task-flow-AC_CatalogoAsignaturas"
	ValueNivel       string = "Pregrado"
	ValueSede        string = "1102 SEDE MEDELLÍN"
	Path_Carreras    string = "carreras.json"
	Tipologia_All    string = "TODAS MENOS  LIBRE ELECCIÓN"
	SizeGrupo        int    = 3
	Path_Grupos      string = "grupos.json"
	Path_JsExtractor string = "src/core/getData.js"
)

var Paths = Path{
	Nivel:     "#pt1\\:r1\\:0\\:soc1\\:\\:content",
	Sede:      "#pt1\\:r1\\:0\\:soc9\\:\\:content",
	Facultad:  "#pt1\\:r1\\:0\\:soc2\\:\\:content",
	Carrera:   "#pt1\\:r1\\:0\\:soc3\\:\\:content",
	Tipologia: "#pt1\\:r1\\:0\\:soc4\\:\\:content",
}

func parseAsignatura(rawData *gson.JSON) Asignatura {
	rawGrupos := rawData.Get("grupos").Arr()

	var grupos []Grupo = make([]Grupo, len(rawGrupos))

	for i, rawGrupo := range rawGrupos {
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
		Nombre:           rawData.Get("nombre").Str(),
		Codigo:           rawData.Get("codigo").Str(),
		Tipologia:        rawData.Get("tipologia").Str(),
		Creditos:         rawData.Get("creditos").Str(),
		Facultad:         rawData.Get("facultad").Str(),
		FechaExtraccion:  rawData.Get("fechaExtraccion").Str(),
		CuposDisponibles: rawData.Get("cuposDisponibles").Str(),
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

	var dataAsignaturas []Asignatura = make([]Asignatura, size)

	var tiemposTotales []time.Duration = make([]time.Duration, size)
	var tiemposCarga []time.Duration = make([]time.Duration, size)
	var tiemposExtraccion []time.Duration = make([]time.Duration, size)

	startTime := time.Now()

	// Recorrer asignaturas
	for i := 0; i < size; i++ {
		timeTotal := time.Now()
		timeLoad := time.Now()

		asignaturas := page.MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr") // Delay

		// Cargar asignatura
		asignatura := asignaturas[i]

		link := asignatura.MustElement(".af_commandLink")
		link.MustClick()

		// page.MustWaitStable() // Delay
		page.MustElement(".af_showDetailHeader_content0")

		timefinLoad := time.Since(timeLoad)
		timeExtraccion := time.Now()

		// Extraer datos
		rawData := page.MustEval(jSExtractorFunctionContent)
		var dataAsignatura Asignatura = parseAsignatura(&rawData)

		timefinExtraccion := time.Since(timeExtraccion)
		dataAsignaturas[i] = dataAsignatura

		// Regresar
		backButton := page.MustElement(".af_button")
		backButton.MustClick()

		timefinTotal := time.Since(timeTotal)
		tiemposTotales[i] = timefinTotal
		tiemposCarga[i] = timefinLoad
		tiemposExtraccion[i] = timefinExtraccion

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

	fmt.Printf("Tiempo de ejecución, %s %s\n", codigo.Carrera, elapsedTime)
	println("--- Tiempos promedios ---")
	println("Promedio total: ", (promedioTotal / float64(size)))
	println("Promedio carga: ", (promedioCarga / float64(size)))
	println("Promedio extraccion: ", (promedioExtraccion / float64(size)))

	println("")

	return dataAsignaturas

}

func LoadJSExtractor() string {
	content, _ := os.ReadFile(Path_JsExtractor)
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

func GenerarGruposCarreras() {
	var listadoCarreras []map[string]string

	contentCarreras, _ := os.ReadFile(Path_Carreras)
	json.Unmarshal(contentCarreras, &listadoCarreras)

	stacks := int(math.Ceil(float64(len(listadoCarreras)) / float64(SizeGrupo)))

	println("Cantidad de stacks: ", (stacks))

	var grupos [][]map[string]string
	for i := 0; i < stacks; i++ {
		var grupo []map[string]string

		for j := 0; j < SizeGrupo; j++ {
			if (i*SizeGrupo)+j < len(listadoCarreras) {
				grupo = append(grupo, listadoCarreras[(i*SizeGrupo)+j])
			}
		}

		grupos = append(grupos, grupo)
	}

	dataGruposJSON, _ := json.Marshal(grupos)
	os.WriteFile(Path_Grupos, dataGruposJSON, 0644)

}
