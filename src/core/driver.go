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
			Sel(page, Paths.Nivel, codigo.Nivel)
			println("Nivel seleccionado...")
			Sel(page, Paths.Sede, codigo.Sede)
			println("Sede seleccionada...")
			Sel(page, Paths.Facultad, codigo.Facultad)
			println("Facultad seleccionada...")
			Sel(page, Paths.Carrera, codigo.Carrera)
			println("Carrera seleccionada...")
			Sel(page, Paths.Tipologia, codigo.Tipologia)
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

func Sel(page *rod.Page, path string, value string) {
	// Wait for done

	page.Timeout(15 * time.Second).MustWaitStable().CancelTimeout()

	// Get element
	el := page.MustElement(path)

	// Click and select

	el.MustClick().MustSelect(value)
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
