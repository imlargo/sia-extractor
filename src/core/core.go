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

type Prerequisito struct {
	Tipo        string              `json:"tipo"`
	IsTodas     bool                `json:"isTodas"`
	Cantidad    int                 `json:"cantidad"`
	Asignaturas []map[string]string `json:"asignaturas"`
}

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
	Nombre           string         `json:"nombre" bson:"nombre"`
	Codigo           string         `json:"codigo" bson:"codigo"`
	Tipologia        string         `json:"tipologia" bson:"tipologia"`
	Creditos         int            `json:"creditos" bson:"creditos"`
	Facultad         string         `json:"facultad" bson:"facultad"`
	Carrera          string         `json:"carrera" bson:"carrera"`
	FechaExtraccion  string         `json:"fechaExtraccion" bson:"fechaExtraccion"`
	CuposDisponibles int            `json:"cuposDisponibles" bson:"cuposDisponibles"`
	Prerequisitos    []Prerequisito `json:"prerequisitos" bson:"prerequisitos"`
	Grupos           []Grupo        `json:"grupos" bson:"grupos"`
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
	Tipologia_All    string = "TODAS MENOS  LIBRE ELECCIÓN"
	SizeGrupo        int    = 3
	Path_Carreras    string = "data/carreras.json"
	Path_Grupos      string = "data/grupos.json"
	Path_JsExtractor string = "src/getData.js"
)

var Paths = Path{
	Nivel:     "#pt1\\:r1\\:0\\:soc1\\:\\:content",
	Sede:      "#pt1\\:r1\\:0\\:soc9\\:\\:content",
	Facultad:  "#pt1\\:r1\\:0\\:soc2\\:\\:content",
	Carrera:   "#pt1\\:r1\\:0\\:soc3\\:\\:content",
	Tipologia: "#pt1\\:r1\\:0\\:soc4\\:\\:content",
}

func parseAsignatura(rawData *gson.JSON, codigo *Codigo) Asignatura {
	rawGrupos := rawData.Get("grupos").Arr()
	var grupos []Grupo = make([]Grupo, len(rawGrupos))

	// Agregar grupos
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

	// Agregar prerequisitos
	rawPrerequisitos := rawData.Get("prerequisitos").Arr()
	prerequisitos := make([]Prerequisito, len(rawPrerequisitos))

	for i, rawPrerequisito := range rawPrerequisitos {
		rawAsignaturas := rawPrerequisito.Get("asignaturas").Arr()
		var asignaturas []map[string]string = make([]map[string]string, len(rawAsignaturas))

		for j, rawAsignatura := range rawAsignaturas {
			asignaturas[j] = map[string]string{
				"codigo": rawAsignatura.Get("codigo").Str(),
				"nombre": rawAsignatura.Get("nombre").Str(),
			}
		}

		prerequisitos[i] = Prerequisito{
			Tipo:        rawPrerequisito.Get("tipo").Str(),
			IsTodas:     rawPrerequisito.Get("isTodas").Bool(),
			Cantidad:    rawPrerequisito.Get("cantidad").Int(),
			Asignaturas: asignaturas,
		}
	}

	return Asignatura{
		Nombre:           rawData.Get("nombre").Str(),
		Codigo:           rawData.Get("codigo").Str(),
		Tipologia:        rawData.Get("tipologia").Str(),
		Creditos:         rawData.Get("creditos").Int(),
		Facultad:         codigo.Facultad,
		Carrera:          codigo.Carrera,
		FechaExtraccion:  rawData.Get("fechaExtraccion").Str(),
		CuposDisponibles: rawData.Get("cuposDisponibles").Int(),
		Prerequisitos:    prerequisitos,
		Grupos:           grupos,
	}
}

func GetAsignaturasCarrera(codigo Codigo) []Asignatura {

	jSExtractorFunctionContent = LoadJSExtractor()

	page, browser := LoadPageCarrera(&codigo)

	asignaturas := page.MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr")

	size := len(asignaturas)
	println("Asignaturas encontradas: ", size)

	var dataAsignaturas []Asignatura = make([]Asignatura, size)

	// Recorrer asignaturas
	for i := 0; i < size; i++ {

		asignaturas = page.MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr")

		// Cargar link
		link := asignaturas[i].MustElement(".af_commandLink")
		link.MustClick()

		page.MustElement(".af_showDetailHeader_content0")

		// Extraer datos
		rawData := page.MustEval(jSExtractorFunctionContent)
		dataAsignaturas[i] = parseAsignatura(&rawData, &codigo)

		// Regresar
		page.MustElement(".af_button").MustClick()
	}

	println("Cerrando navegador...")
	browser.MustClose()
	println("Navegador cerrado")

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
	page := rod.New().MustConnect().MustIncognito().MustPage(SIA_URL)
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

func ExtraerElectivas() []Asignatura {
	println("Holi")

	var codigo Codigo = Codigo{
		Nivel:     ValueNivel,
		Sede:      ValueSede,
		Facultad:  "3067 FACULTAD DE CIENCIAS HUMANAS  Y ECONÓMICAS",
		Carrera:   "3512 CIENCIA POLÍTICA",
		Tipologia: "LIBRE ELECCIÓN",
	}

	jSExtractorFunctionContent = LoadJSExtractor()

	println("Iniciando...")
	page := rod.New().MustConnect().MustIncognito().MustPage(SIA_URL)
	println("Cargado. ok")

	page.MustWaitStable().MustElement(Paths.Nivel).MustClick().MustSelect(codigo.Nivel)
	page.MustWaitStable().MustElement(Paths.Sede).MustClick().MustSelect(codigo.Sede)
	page.MustWaitStable().MustElement(Paths.Facultad).MustClick().MustSelect(codigo.Facultad)
	page.MustWaitStable().MustElement(Paths.Carrera).MustClick().MustSelect(codigo.Carrera)
	time.Sleep(2 * time.Second)

	err := page.MustElement(Paths.Tipologia).MustClick().Select([]string{`^LIBRE ELECCIÓN$`}, true, rod.SelectorTypeRegex)
	if err != nil {
		println("Error: ", err)
	}

	println("screen tomado")

	// Porque tipo
	page.MustWaitStable().MustElement("#pt1\\:r1\\:0\\:soc5\\:\\:content").MustClick().MustSelect("Por facultad y plan")

	println("Tipo seleccionado")

	// POrque sede
	page.MustWaitStable().MustElement("#pt1\\:r1\\:0\\:soc10\\:\\:content").MustClick().MustSelect("1102 SEDE MEDELLÍN")
	println("Sede seleccionada")

	// porque facultad
	page.MustWaitStable().MustElement("#pt1\\:r1\\:0\\:soc6\\:\\:content").MustClick().MustSelect("3 SEDE MEDELLÍN")
	println("Facultad seleccionada")

	// porque plan
	page.MustWaitStable().MustElement("#pt1\\:r1\\:0\\:soc7\\:\\:content").MustClick().MustSelect("3CLE COMPONENTE DE LIBRE ELECCIÓN")
	println("Plan seleccionado")

	// select all checkboxes
	checkboxesDias := page.MustElements(".af_selectBooleanCheckbox_native-input")
	for _, checkbox := range checkboxesDias {
		checkbox.MustClick()
	}

	println("Campos seleccionados...ejecutando búsqueda")

	// Hacer clic en el botón para ejecutar la búsqueda
	page.MustElement(".af_button_link").MustClick()

	page.MustWaitStable().MustWaitIdle().MustWaitDOMStable()

	size := len(page.MustWaitStable().MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr"))

	println("Asignaturas encontradas: ", size)

	var dataAsignaturas []Asignatura = make([]Asignatura, size)
	// Recorrer asignaturas
	for i := 0; i < size; i++ {

		asignaturas := page.MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr")

		// Cargar asignatura
		asignatura := asignaturas[i]

		link := asignatura.MustElement(".af_commandLink")
		link.MustClick()

		page.MustElement(".af_showDetailHeader_content0")

		// Extraer datos
		rawData := page.MustEval(jSExtractorFunctionContent)
		dataAsignaturas[i] = parseAsignatura(&rawData, &codigo)

		println("Asignatura: ", dataAsignaturas[i].Codigo, dataAsignaturas[i].Nombre)

		// Regresar
		backButton := page.MustElement(".af_button")
		backButton.MustClick()
	}

	return dataAsignaturas
}
