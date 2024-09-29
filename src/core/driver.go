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
			page.Timeout(15 * time.Second).MustWaitStable().CancelTimeout().MustElement(Paths.Nivel).MustClick().MustSelect(codigo.Nivel)
			println("Nivel seleccionado...", codigo.Carrera)
			page.Timeout(15 * time.Second).MustWaitStable().CancelTimeout().MustElement(Paths.Sede).MustClick().MustSelect(codigo.Sede)
			println("Sede seleccionada...", codigo.Carrera)
			page.Timeout(15 * time.Second).MustWaitStable().CancelTimeout().MustElement(Paths.Facultad).MustClick().MustSelect(codigo.Facultad)
			println("Facultad seleccionada...", codigo.Carrera)

			page.Timeout(15 * time.Second).MustWaitStable().CancelTimeout()
			i := 0
			for {
				if i > 5 {
					panic("Error al cargar la pagina, timeout")
				}

				selectCarrera := page.MustElement(Paths.Carrera)
				options := selectCarrera.MustElements("option")

				if len(options) != 0 {
					selectCarrera.MustClick().MustSelect(codigo.Carrera)
					println("Carrera seleccionada...", codigo.Carrera)
					break
				}

				if len(options) == 0 {
					i++
					println("### Pooling again ###")
					page.MustElement(Paths.Facultad).MustClick().MustSelect(codigo.Facultad)
					page.MustWaitStable()
					println("Facultad seleccionada...", codigo.Carrera)
				}
			}

			page.Timeout(15 * time.Second).MustWaitStable().CancelTimeout().MustElement(Paths.Tipologia).MustClick().MustSelect(codigo.Tipologia)
			println("Tipologia seleccionada...", codigo.Carrera)
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
