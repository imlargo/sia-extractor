package main

import "github.com/go-rod/rod"

/*
type Horario struct {
	inicio: string;
	fin: string;
	dia: string;
}

 type Grupo struct {
	grupo: string
	cupos: number
	profesor: string
	duracion: string
	jornada: string
	horarios: Horario[]
}
*/

type Asignatura struct {
	// codigo    string
	nombre string
	// tipologia string
	// creditos  string
	/*
		facultad: string;
		carrera: string;
		fechaExtraccion: string;
		cuposDisponibles: number;

		grupos: Grupo[];
	*/
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

func main() {
	codigo := Codigo{
		nivel:     "Pregrado",
		sede:      "1102 SEDE MEDELLÍN",
		facultad:  "3068 FACULTAD DE MINAS",
		carrera:   "3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA",
		tipologia: "TODAS MENOS  LIBRE ELECCIÓN",
	}

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
	println("Nivel")
	page.MustWaitStable().MustElement(paths.sede).MustClick().MustSelect(codigo.sede)
	println("Sede")
	page.MustWaitStable().MustElement(paths.facultad).MustClick().MustSelect(codigo.facultad)
	println("Facultad")
	page.MustWaitStable().MustElement(paths.carrera).MustClick().MustSelect(codigo.carrera)
	println("Carrera")
	page.MustWaitStable().MustElement(paths.tipologia).MustClick().MustSelect(codigo.tipologia)
	println("Tipologia")

	// select all checkboxes
	checkboxesDias := page.MustWaitStable().MustElements(".af_selectBooleanCheckbox_native-input")
	for _, checkbox := range checkboxesDias {
		checkbox.MustClick()
	}
	println("Dias")

	// Hacer clic en el botón para ejecutar la búsqueda
	page.MustWaitStable().MustElement(".af_button_link").MustClick()
	println("Buscar ejecutado")

	size := len(page.MustWaitStable().MustElement(".af_table_data-table-VH-lines").MustElement("tbody").MustElements("tr"))

	println("Asignaturas encontradas: ", size)
	println()
	println()

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
		println(data.nombre)

		// Regresar
		backButton := page.MustElement(".af_button")
		backButton.MustClick()

		page.MustWaitStable()
		println()

	}

	println("")
	println("Finalizado")
}

func procesarMateria(page *rod.Page) Asignatura {

	rawName := page.MustElement("h2").MustText()

	return Asignatura{
		nombre: rawName,
	}

}
