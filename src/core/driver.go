package core

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

const (
	sleepSelect = 3 * time.Second
	timeoutPage = 15 * time.Second
	maxRetries  = 3
)

type Driver struct {
	Browser *rod.Browser
	Page    *rod.Page
}

func NewDriver() *Driver {
	browser := rod.New().MustConnect()

	return &Driver{
		Browser: browser,
	}
}

func (driver *Driver) getPage() *rod.Page {

	page := driver.Browser.MustIncognito().MustPage("")

	router := InterceptRequests(page)
	go router.Run()

	page.Timeout(timeoutPage).MustNavigate(SIA_URL).MustWaitStable().CancelTimeout()

	return page
}

func (driver *Driver) loadPageWithRetry(codigo Codigo) *rod.Page {
	var page *rod.Page
	for attempts := 0; attempts <= maxRetries; attempts++ {
		err := rod.Try(func() {
			page = driver.getPage()
			driver.Page = page
			driver.selectOptions(codigo)
		})
		if err == nil {
			return page
		}
		if attempts == maxRetries {
			driver.Page = page
			panic("Error al cargar la pagina, timeout")
		}
		println("Retrying...", codigo.Carrera)
		page.MustClose()
	}
	return nil
}

func (driver *Driver) LoadPageCarrera(codigo Codigo) *rod.Page {
	page := driver.loadPageWithRetry(codigo)
	driver.selectAllCheckboxes()
	return page
}

func (driver *Driver) SelectElectivas(codigo Codigo, codigoElectiva PathElectiva) {
	driver.selectWithSleep(PathsElectiva.Por, codigoElectiva.Por, Paths.Tipologia, codigo.Tipologia)
	driver.selectWithSleep(PathsElectiva.SedePor, codigoElectiva.SedePor, PathsElectiva.Por, ValuesElectiva.Por)
	driver.selectWithSleep(PathsElectiva.FacultadPor, codigoElectiva.FacultadPor, PathsElectiva.SedePor, ValuesElectiva.SedePor)
	driver.selectWithSleep(PathsElectiva.CarreraPor, codigoElectiva.CarreraPor, PathsElectiva.FacultadPor, ValuesElectiva.FacultadPor)
}

func (driver *Driver) selectOptions(codigo Codigo) {
	driver.selectWithSleep(Paths.Nivel, codigo.Nivel, "", "")
	driver.selectWithSleep(Paths.Sede, codigo.Sede, Paths.Nivel, codigo.Nivel)
	driver.selectWithSleep(Paths.Facultad, codigo.Facultad, Paths.Sede, codigo.Sede)
	driver.selectWithSleep(Paths.Carrera, codigo.Carrera, Paths.Facultad, codigo.Facultad)
	driver.selectWithSleep(Paths.Tipologia, codigo.Tipologia, Paths.Carrera, codigo.Carrera)
}

func (driver *Driver) selectWithSleep(path, value, prevPath, prevValue string) {
	time.Sleep(sleepSelect)
	driver.SelectWithRecover(path, value, prevPath, prevValue)
	println(fmt.Sprintf("%s seleccionado...", value))
}

func (driver *Driver) selectAllCheckboxes() {
	checkboxes := driver.Page.MustElements(".af_selectBooleanCheckbox_native-input")
	for _, checkbox := range checkboxes {
		checkbox.MustClick()
	}
}

func (driver *Driver) SelectWithRecover(path, value, prevPath, prevValue string) {
	for i := 0; i <= maxRetries; i++ {
		selectEl := driver.Page.MustElement(path)
		options := selectEl.MustElements("option")

		if len(options) != 0 {
			selectEl.MustClick()
			Sel(selectEl, value)
			return
		}

		if i == maxRetries {
			panic("Error al cargar la pagina, timeout")
		}

		el2 := driver.Page.MustElement(prevPath)
		el2.MustClick()
		Sel(el2, prevValue)
		time.Sleep(sleepSelect)
	}
}

func Sel(el *rod.Element, value string) error {
	regex := fmt.Sprintf("^%s$", value)
	return el.Select([]string{regex}, true, rod.SelectorTypeRegex)
}

func InterceptRequests(page *rod.Page) *rod.HijackRouter {
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

	return router
}
