package deploy

import (
	"sia-extractor/src/core"
	"slices"
)

func getTipologiasUnicas(asignaturas []core.Asignatura) []string {
	tipologiasUnicas := make([]string, 0)

	for _, asignatura := range asignaturas {

		if slices.Contains(tipologiasUnicas, asignatura.Tipologia) {
			continue
		}

		tipologiasUnicas = append(tipologiasUnicas, asignatura.Tipologia)
	}

	return tipologiasUnicas
}

func CreateDocumentCarrera(carrera string, facultad string, asignaturas []core.Asignatura) DocumentCarrera {
	return DocumentCarrera{
		ID:          carrera,
		Facultad:    facultad,
		Carrera:     carrera,
		Asignaturas: asignaturas,
	}
}
