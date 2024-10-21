package core

import (
	"fmt"
	"math"
	"sia-extractor/src/utils"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

var jSExtractorFunctionContent string = ""

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

	err := utils.SaveJsonToFile(listadoCarrerasSede, "carreras.json")
	if err != nil {
		fmt.Println("Error al guardar archivo: ", err)
	}

	println("Finalizado!!! :D")
}

func GenerarGruposCarreras() {
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

func ExtraerElectivas(codigo Codigo) *[]Asignatura {

	codigo.Facultad = "3068 FACULTAD DE MINAS"
	codigo.Carrera = "3520 INGENIERÍA DE SISTEMAS E INFORMÁTICA"
	codigo.Tipologia = Tipologia_Electiva

	jSExtractorFunctionContent = LoadJSExtractor()

	driver := NewDriver()
	page := driver.LoadPageCarrera(codigo)

	driver.SelectElectivas(codigo, ConstructCodigoElectiva(ValuesElectiva.FacultadPor, ValuesElectiva.CarreraPor))
	defer page.MustClose()

	println("Campos seleccionados...ejecutando búsqueda", codigo.Carrera)

	codigoCopy := codigo
	codigoCopy.Facultad = ValuesElectiva.FacultadPor
	codigoCopy.Carrera = ValuesElectiva.CarreraPor

	asignaturas := extraerAsignaturas(codigoCopy, page)

	return &asignaturas
}

func ExtraerGrupo(indexGrupo int) map[string]*[]Asignatura {

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

func GetAsignaturasCarrera(codigo Codigo) *[]Asignatura {

	jSExtractorFunctionContent = LoadJSExtractor()

	driver := NewDriver()
	page := driver.LoadPageCarrera(codigo)
	defer page.MustClose()

	println("Campos seleccionados...ejecutando búsqueda", codigo.Carrera)

	asignaturas := extraerAsignaturas(codigo, page)

	return &asignaturas
}

func extraerAsignaturas(codigo Codigo, page *rod.Page) []Asignatura {

	// Hacer clic en el botón para ejecutar la búsqueda
	page.MustWaitStable().MustWaitIdle().MustWaitDOMStable()
	page.MustElement(".af_button_link").MustClick()
	page.MustWaitStable().MustWaitIdle().MustWaitDOMStable()

	asignaturas := getTable(page)
	size := len(asignaturas)
	println("Asignaturas encontradas: ", size)

	data := make([]Asignatura, size)

	// Recorrer asignaturas
	for i := 0; i < size; i++ {

		asignaturas = page.MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr")

		// Cargar link
		asignaturas[i].MustElement(".af_commandLink").MustClick()

		page.MustElement(".af_showDetailHeader_content0")

		// Extraer datos

		rawData := page.MustEval(jSExtractorFunctionContent)
		data[i] = parseAsignatura(&rawData, &codigo)
		println(i, "/", size, data[i].Nombre)

		// Regresar
		page.MustElement(".af_button").MustClick()
	}

	println("Finalizado...")

	return data
}

func getTable(page *rod.Page) rod.Elements {

	var rows rod.Elements

	for {
		println("Buscando tabla...")

		table := page.MustElement(".af_table_data-table-VH-lines")
		time.Sleep(3 * time.Second)

		if table == nil {
			continue
		}

		if !table.MustHas("tbody") {
			break
		}

		tbody := table.MustElement("tbody")
		if tbody == nil {
			continue
		}

		rows = tbody.MustElements("tr")
		if rows == nil || len(rows) > 200 {
			continue
		}

		break
	}

	return rows

}
