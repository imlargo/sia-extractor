package core

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
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

func ExtraerGrupo(indexGrupo int) map[string]*[]Asignatura {

	var listadoGrupos [][]map[string]string
	bytesGrupos, _ := os.ReadFile(Path_Grupos)
	json.Unmarshal(bytesGrupos, &listadoGrupos)
	grupo := listadoGrupos[indexGrupo-1]

	chanAsignaturas := make(chan *[]Asignatura, len(grupo))
	browser := rod.New().MustConnect()

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
			asignaturas := GetAsignaturasCarrera(browser, codigo)
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

func GetAsignaturasCarrera(browser *rod.Browser, codigo Codigo) *[]Asignatura {

	jSExtractorFunctionContent = LoadJSExtractor()

	page, _ := LoadPageCarrera(browser, codigo)
	defer page.MustClose()

	println("Campos seleccionados...ejecutando búsqueda")

	// Hacer clic en el botón para ejecutar la búsqueda
	page.MustWaitStable().MustElement(".af_button_link").MustClick()
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

	return &data
}

func getTable(page *rod.Page) rod.Elements {

	var rows rod.Elements

	for {
		table := page.MustElement(".af_table_data-table-VH-lines")
		if table == nil {
			continue
		}

		tbody := table.MustElement("tbody")
		if tbody == nil {
			continue
		}

		rows = tbody.MustElements("tr")
		if rows == nil || len(rows) > 100 {
			continue
		}

		break
	}

	return rows

}
