package core

import (
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

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

			time.Sleep(5 * time.Second)
			SelectWithRecover(page, Paths.Sede, codigo.Sede, Paths.Nivel, codigo.Nivel)
			println("Sede seleccionada...", codigo.Carrera)

			time.Sleep(5 * time.Second)
			SelectWithRecover(page, Paths.Facultad, codigo.Facultad, Paths.Sede, codigo.Sede)
			println("Facultad seleccionada...", codigo.Carrera)

			time.Sleep(5 * time.Second)
			SelectWithRecover(page, Paths.Carrera, codigo.Carrera, Paths.Facultad, codigo.Facultad)
			println("Carrera seleccionada...", codigo.Carrera)

			time.Sleep(5 * time.Second)
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

		println("Selecting element")
		selectEl := page.MustElement(path)
		options := selectEl.MustElements("option")

		println("Loadded options")

		if len(options) != 0 {
			println("Selecting value")
			selectEl.MustClick()
			selectEl.MustSelect(value)
			println("Clicked")
			break
		}

		if len(options) == 0 {
			i++
			println("### Pooling again ###", value)
			el2 := page.MustElement(prevPath)
			el2.MustClick()
			el2.MustSelect(prevValue)
			time.Sleep(5 * time.Second)
		}
	}
}
