package driver

import (
	"sia-extractor/src/core"
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

	page.Timeout(timeoutPage).MustNavigate(core.SIA_URL).MustWaitStable().CancelTimeout()

	return page
}

func (driver *Driver) loadPageWithRetry(codigo core.Codigo) *rod.Page {
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

func (driver *Driver) LoadPageCarrera(codigo core.Codigo) {
	driver.loadPageWithRetry(codigo)
	driver.selectAllCheckboxes()
}

func (driver *Driver) SelectElectivas(codigo core.Codigo, codigoElectiva core.PathElectiva) {
	driver.selectWithSleep(core.PathsElectiva.Por, codigoElectiva.Por, core.Paths.Tipologia, codigo.Tipologia)
	driver.selectWithSleep(core.PathsElectiva.SedePor, codigoElectiva.SedePor, core.PathsElectiva.Por, core.ValuesElectiva.Por)
	driver.selectWithSleep(core.PathsElectiva.FacultadPor, codigoElectiva.FacultadPor, core.PathsElectiva.SedePor, core.ValuesElectiva.SedePor)
	driver.selectWithSleep(core.PathsElectiva.CarreraPor, codigoElectiva.CarreraPor, core.PathsElectiva.FacultadPor, core.ValuesElectiva.FacultadPor)
}
