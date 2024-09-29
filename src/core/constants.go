package core

const (
	SIA_URL          string = "https://sia.unal.edu.co/Catalogo/facespublico/public/servicioPublico.jsf?taskflowId=task-flow-AC_CatalogoAsignaturas"
	ValueNivel       string = "Pregrado"
	ValueSede        string = "1102 SEDE MEDELLÍN"
	Tipologia_All    string = "TODAS MENOS  LIBRE ELECCIÓN"
	SizeGrupo        int    = 1
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
