package core

import (
	"fmt"
	"math"
	"sia-extractor/src/utils"
	"time"

	"github.com/go-rod/rod"
)

type Extractor struct {
	Driver          *Driver
	JsExtractorFunc string
}

func NewExtractor() *Extractor {
	return &Extractor{
		Driver: NewDriver(),
	}
}

func (extractor *Extractor) LoadJSFunc() {
	extractor.JsExtractorFunc = LoadJSExtractorFunc()
}

func (extractor *Extractor) ExtraerAsignaturas(codigo Codigo) []Asignatura {

	// Hacer clic en el botón para ejecutar la búsqueda
	extractor.Driver.Page.MustWaitStable().MustWaitIdle().MustWaitDOMStable()
	extractor.Driver.Page.MustElement(".af_button_link").MustClick()
	extractor.Driver.Page.MustWaitStable().MustWaitIdle().MustWaitDOMStable()

	asignaturas := extractor.Driver.GetTable()
	size := len(asignaturas)
	println("Asignaturas encontradas: ", size)

	data := make([]Asignatura, size)

	// Recorrer asignaturas
	for i := 0; i < size; i++ {

		asignaturas = extractor.Driver.Page.MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr")

		// Cargar link
		asignaturas[i].MustElement(".af_commandLink").MustClick()

		extractor.Driver.Page.MustElement(".af_showDetailHeader_content0")

		// Extraer datos

		rawData := extractor.Driver.Page.MustEval(extractor.JsExtractorFunc)
		data[i] = parseAsignatura(&rawData, &codigo)
		println(i, "/", size, data[i].Nombre)

		// Regresar
		extractor.Driver.Page.MustElement(".af_button").MustClick()
	}

	println("Finalizado...")

	return data
}

func (extractor *Extractor) CreatePathsCarreras() {

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

	err := utils.SaveJsonToFile(listadoCarrerasSede, "carreras.json")
	if err != nil {
		fmt.Println("Error al guardar archivo: ", err)
	}

	println("Finalizado!!! :D")
}

func (extractor *Extractor) GenerarGruposCarreras() {
	var listadoCarreras []map[string]string
	utils.LoadJsonFromFile(&listadoCarreras, Path_Carreras)

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

	err := utils.SaveJsonToFile(grupos, "grupos.json")
	if err != nil {
		fmt.Println("Error al guardar archivo: ", err)
	}

}

func (extractor *Extractor) ExtraerElectivas(codigo Codigo) *[]Asignatura {

	codigo.Facultad = "3068 FACULTAD DE MINAS"
	codigo.Carrera = "3520 INGENIERÍA DE SISTEMAS E INFORMÁTICA"
	codigo.Tipologia = Tipologia_Electiva

	extractor.LoadJSFunc()

	driver := NewDriver()
	page := driver.LoadPageCarrera(codigo)

	driver.SelectElectivas(codigo, ConstructCodigoElectiva(ValuesElectiva.FacultadPor, ValuesElectiva.CarreraPor))
	defer page.MustClose()

	println("Campos seleccionados...ejecutando búsqueda", codigo.Carrera)

	codigoCopy := codigo
	codigoCopy.Facultad = ValuesElectiva.FacultadPor
	codigoCopy.Carrera = ValuesElectiva.CarreraPor

	asignaturas := extractor.ExtraerAsignaturas(codigoCopy)

	return &asignaturas
}

/*
func (extractor *Extractor) ExtraerGrupo(indexGrupo int) map[string]*[]Asignatura {

	var listadoGrupos [][]map[string]string
	utils.LoadJsonFromFile(&listadoGrupos, Path_Grupos)

	if indexGrupo > len(listadoGrupos) {
		println("El grupo seleccionado no existe")
		return nil
	}

	grupo := listadoGrupos[indexGrupo-1]

	chanAsignaturas := make(chan *[]Asignatura, len(grupo))

	var wg sync.WaitGroup
	for _, carrera := range grupo {

		wg.Add(1)

		go func(carrera map[string]string) {
			defer wg.Done()

			codigo := Codigo{
				Nivel:     ValueNivel,
				Sede:      ValueSede,
				Facultad:  carrera["facultad"],
				Carrera:   carrera["carrera"],
				Tipologia: Tipologia_All,
			}

			println("Iniciando: ", codigo.Carrera)
			asignaturas := GetAsignaturasCarrera(codigo)
			println("Finalizado: ", codigo.Carrera)

			chanAsignaturas <- asignaturas

		}(carrera)
	}

	go func() {
		wg.Wait()
		close(chanAsignaturas)
	}()

	data := make(map[string]*[]Asignatura)
	for asignaturas := range chanAsignaturas {
		carrera := (*asignaturas)[0].Carrera
		data[carrera] = asignaturas
	}

	return data
}
*/

func (extractor *Extractor) GetAsignaturasCarrera(codigo Codigo) *[]Asignatura {

	extractor.LoadJSFunc()

	driver := NewDriver()
	page := driver.LoadPageCarrera(codigo)
	defer page.MustClose()

	println("Campos seleccionados...ejecutando búsqueda", codigo.Carrera)

	asignaturas := extractor.ExtraerAsignaturas(codigo)

	return &asignaturas
}
