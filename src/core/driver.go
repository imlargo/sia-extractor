package core

import "github.com/go-rod/rod"

func LoadPageCarrera(codigo *Codigo) (*rod.Page, *rod.Browser) {

	browser := rod.New().MustConnect()
	page := browser.MustIncognito().MustPage(SIA_URL)

	println("Selecionando...")
	page.MustWaitStable().MustElement(Paths.Nivel).MustClick().MustSelect(codigo.Nivel)
	page.MustWaitStable().MustElement(Paths.Sede).MustClick().MustSelect(codigo.Sede)
	page.MustWaitStable().MustElement(Paths.Facultad).MustClick().MustSelect(codigo.Facultad)
	page.MustWaitStable().MustElement(Paths.Carrera).MustClick().MustSelect(codigo.Carrera)
	page.MustWaitStable().MustElement(Paths.Tipologia).MustClick().MustSelect(codigo.Tipologia)
	println("Campos seleccionados...")

	// select all checkboxes
	checkboxesDias := page.MustElements(".af_selectBooleanCheckbox_native-input")
	for _, checkbox := range checkboxesDias {
		checkbox.MustClick()
	}

	return page, browser
}
