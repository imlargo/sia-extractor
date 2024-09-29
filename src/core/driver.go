package core

import "github.com/go-rod/rod"

func LoadPageCarrera(browser *rod.Browser, codigo Codigo) (*rod.Page, *rod.Browser) {

	page := browser.MustIncognito().MustPage(SIA_URL).MustWaitStable().MustWaitIdle().MustWaitDOMStable()

	println("Selecionando...")
	page.MustWaitStable().MustWaitIdle().MustWaitDOMStable().MustElement(Paths.Nivel).MustClick().MustSelect(codigo.Nivel)
	page.MustWaitStable().MustWaitIdle().MustWaitDOMStable().MustElement(Paths.Sede).MustClick().MustSelect(codigo.Sede)
	page.MustWaitStable().MustWaitIdle().MustWaitDOMStable().MustElement(Paths.Facultad).MustClick().MustSelect(codigo.Facultad)
	page.MustWaitStable().MustWaitIdle().MustWaitDOMStable().MustElement(Paths.Carrera).MustClick().MustSelect(codigo.Carrera)
	page.MustWaitStable().MustWaitIdle().MustWaitDOMStable().MustElement(Paths.Tipologia).MustClick().MustSelect(codigo.Tipologia)
	println("Campos seleccionados...")

	// select all checkboxes
	checkboxesDias := page.MustElements(".af_selectBooleanCheckbox_native-input")
	for _, checkbox := range checkboxesDias {
		checkbox.MustClick()
	}

	return page, browser
}
