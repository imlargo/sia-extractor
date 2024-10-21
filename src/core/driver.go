package core

import (
	"time"

	"github.com/go-rod/rod"
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
