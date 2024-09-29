package core

import (
	"time"

	"github.com/go-rod/rod"
)

func LoadPageCarrera(browser *rod.Browser, codigo Codigo) (*rod.Page, *rod.Browser) {

	var page *rod.Page
	timeoutLoad := 15 * time.Second
	timeoutSelect := 10 * time.Second
	intentos := 0

	for {

		err := rod.Try(func() {
			page = browser.MustIncognito().MustPage(SIA_URL).Timeout(timeoutLoad).MustWaitStable().CancelTimeout()

			println("Selecionando...")
			Sel(page, Paths.Nivel, codigo.Nivel, timeoutSelect)
			println("Nivel seleccionado...")
			Sel(page, Paths.Sede, codigo.Sede, timeoutSelect)
			println("Sede seleccionada...")
			Sel(page, Paths.Facultad, codigo.Facultad, timeoutSelect)
			println("Facultad seleccionada...")
			Sel(page, Paths.Carrera, codigo.Carrera, timeoutSelect)
			println("Carrera seleccionada...")
			Sel(page, Paths.Tipologia, codigo.Tipologia, timeoutSelect)
			println("Campos seleccionados...")
		})

		if err == nil {
			break
		}

		intentos++
		println("Pooling again...")
		page.MustClose()

		if intentos > 3 {
			panic("Error al cargar la pagina, timeout")
		}

	}

	// select all checkboxes
	checkboxesDias := page.MustElements(".af_selectBooleanCheckbox_native-input")
	for _, checkbox := range checkboxesDias {
		checkbox.MustClick()
	}

	return page, browser
}

func Sel(page *rod.Page, path string, value string, t1 time.Duration) {
	// Wait for done

	page.Timeout(t1).MustWaitStable().CancelTimeout()

	// Get element
	el := page.MustElement(path)

	// Click and select

	el.MustClick().MustSelect(value)
}
