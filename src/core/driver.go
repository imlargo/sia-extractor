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
			Sel(page, Paths.Nivel, codigo.Nivel, codigo.Carrera)
			println("Nivel seleccionado...", codigo.Carrera)
			Sel(page, Paths.Sede, codigo.Sede, codigo.Carrera)
			println("Sede seleccionada...", codigo.Carrera)
			Sel(page, Paths.Facultad, codigo.Facultad, codigo.Carrera)
			println("Facultad seleccionada...", codigo.Carrera)
			Sel(page, Paths.Carrera, codigo.Carrera, codigo.Carrera)
			println("Carrera seleccionada...", codigo.Carrera)
			Sel(page, Paths.Tipologia, codigo.Tipologia, codigo.Carrera)
			println("Campos seleccionados...", codigo.Carrera)
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

func Sel(page *rod.Page, path string, value string, carrera string) {
	// Wait for done

	println("Waiting for stable...", carrera, value)
	page.Timeout(15 * time.Second).MustWaitStable().CancelTimeout()
	println("Stable...", carrera, value)

	// Get element
	println("Elemento encontrado...", carrera, value)

	// Verificar que el elemento tenga options
	intentos := 0
	for {
		options := page.MustElement(path).MustElements("option")
		if len(options) != 0 {
			break
		}
		intentos += 1
		println("Waiting for options...")
		if intentos > 10 {
			time.Sleep(3 * time.Second)
			// panic("Error al cargar la pagina, timeout")
		}
	}

	// Click and select
	el := page.MustElement(path)
	el.MustClick()
	println("Clicked...", carrera, value)

	el.MustSelect(value)
	println("Seleccionado...", carrera, value)
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
