package core

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

const timeoutSelect = 3 * time.Second

func LoadPageCarrera(codigo Codigo) (*rod.Page, *rod.Browser) {

	var page *rod.Page
	intentos := 0

	browser := rod.New().MustConnect()

	for {

		err := rod.Try(func() {
			page = getPage(browser)

			println("Selecionando...")
			page.MustElement(Paths.Nivel).MustClick().MustSelect(codigo.Nivel)
			println("Nivel seleccionado...", codigo.Carrera)

			time.Sleep(timeoutSelect)
			SelectWithRecover(page, Paths.Sede, codigo.Sede, Paths.Nivel, codigo.Nivel)
			println("Sede seleccionada...", codigo.Carrera)

			time.Sleep(timeoutSelect)
			SelectWithRecover(page, Paths.Facultad, codigo.Facultad, Paths.Sede, codigo.Sede)
			println("Facultad seleccionada...", codigo.Carrera)

			time.Sleep(timeoutSelect)
			SelectWithRecover(page, Paths.Carrera, codigo.Carrera, Paths.Facultad, codigo.Facultad)
			println("Carrera seleccionada...", codigo.Carrera)

			time.Sleep(timeoutSelect)
			SelectWithRecover(page, Paths.Tipologia, codigo.Tipologia, Paths.Carrera, codigo.Carrera)
			println("Tipologia seleccionada...", codigo.Carrera)
		})

		if err == nil {
			break
		}

		intentos++
		println("Pooling again...", codigo.Carrera)
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

func loadElectivas(codigo Codigo, codigoElectiva PathElectiva, page *rod.Page) {

	// Porque tipo
	time.Sleep(timeoutSelect)
	SelectWithRecover(page, PathsElectiva.Por, codigoElectiva.Por, Paths.Tipologia, codigo.Tipologia)
	println("Tipo electiva seleccionada...", codigo.Carrera)

	// POrque sede
	time.Sleep(timeoutSelect)
	SelectWithRecover(page, PathsElectiva.SedePor, codigoElectiva.SedePor, PathsElectiva.Por, ValuesElectiva.Por)
	println("Sede electiva seleccionada...", codigo.Carrera)

	// porque facultad
	time.Sleep(timeoutSelect)
	SelectWithRecover(page, PathsElectiva.FacultadPor, codigoElectiva.FacultadPor, PathsElectiva.SedePor, ValuesElectiva.SedePor)
	println("Facultad electiva seleccionada...", codigo.Carrera)

	// porque plan
	time.Sleep(timeoutSelect)
	SelectWithRecover(page, PathsElectiva.CarreraPor, codigoElectiva.CarreraPor, PathsElectiva.FacultadPor, ValuesElectiva.FacultadPor)
	println("Carrera electiva seleccionada...", codigo.Carrera)
}

func getPage(browser *rod.Browser) *rod.Page {
	page := browser.MustIncognito().MustPage("")

	router := page.HijackRequests()

	cancelReq := func(ctx *rod.Hijack) {
		if ctx.Request.Type() == proto.NetworkResourceTypeImage {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	}

	router.MustAdd("*.png", cancelReq)
	router.MustAdd("*.svg", cancelReq)
	router.MustAdd("*.gif", cancelReq)
	router.MustAdd("*.css", cancelReq)

	// since we are only hijacking a specific page, even using the "*" won't affect much of the performance
	go router.Run()

	page.Timeout(15 * time.Second).MustNavigate(SIA_URL).MustWaitStable().CancelTimeout()

	return page
}

func SelectWithRecover(page *rod.Page, path string, value string, prevPath string, prevValue string) {
	i := 0
	for {
		if i > 5 {
			panic("Error al cargar la pagina, timeout")
		}

		// println("Selecting element")
		selectEl := page.MustElement(path)
		options := selectEl.MustElements("option")

		// println("Loadded options")

		if len(options) != 0 {
			// println("Selecting value")
			selectEl.MustClick()
			Sel(selectEl, value)
			// println("Clicked")
			break
		}

		if len(options) == 0 {
			i++
			// println("### Pooling again ###", value)
			el2 := page.MustElement(prevPath)
			el2.MustClick()
			Sel(el2, prevValue)
			time.Sleep(timeoutSelect)
		}
	}
}

func Sel(el *rod.Element, value string) error {
	regex := fmt.Sprintf("^%s$", value)
	err := el.Select([]string{regex}, true, rod.SelectorTypeRegex)
	return err
}
