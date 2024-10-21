package core

import (
	"os"
	"sia-extractor/src/utils"

	"github.com/ysmood/gson"
)

func LoadJSExtractor() string {
	content, _ := os.ReadFile(Path_JsExtractor)
	JSExtractorFunctionContent := string(content)

	return JSExtractorFunctionContent
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
				Lugar:  rawHorario.Get("lugar").Str(),
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

func ConstructCodigo(facultad string, carrera string) Codigo {
	return Codigo{
		Nivel:     ValueNivel,
		Sede:      ValueSede,
		Facultad:  facultad,
		Carrera:   carrera,
		Tipologia: Tipologia_All,
	}
}

func ConstructCodigoElectiva(facultad string, carrera string) PathElectiva {
	return PathElectiva{
		Por:         ValuesElectiva.Por,
		SedePor:     ValuesElectiva.SedePor,
		FacultadPor: facultad,
		CarreraPor:  carrera,
	}
}

func loadListadoGrupos() [][]map[string]string {
	var listadoGrupos [][]map[string]string
	utils.LoadJsonFromFile(&listadoGrupos, Path_Grupos)

	return listadoGrupos
}

func GetCodigoFromGrupo(indexGrupo int) *Codigo {

	listadoGrupos := loadListadoGrupos()

	if indexGrupo > len(listadoGrupos) {
		println("El grupo seleccionado no existe")
		return nil
	}

	carrera := listadoGrupos[indexGrupo-1]

	codigo := ConstructCodigo(carrera[0]["facultad"], carrera[0]["carrera"])

	return &codigo
}
