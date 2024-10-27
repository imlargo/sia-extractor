package driver

import (
	"fmt"
	"sia-extractor/src/core"
	"time"

	"github.com/go-rod/rod"
)

const (
	timeoutDuration = 10 * time.Second
)

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
		selectEl := driver.Page.Timeout(timeoutDuration).MustElement(path).CancelTimeout()
		options := selectEl.MustElements("option")

		if len(options) != 0 {
			selectEl.MustClick()
			Sel(selectEl, value)
			return
		}

		if i == maxRetries {
			panic("Error al cargar la pagina, timeout")
		}

		el2 := driver.Page.Timeout(timeoutDuration).MustElement(prevPath).CancelTimeout()
		el2.MustClick()
		Sel(el2, prevValue)
		time.Sleep(sleepSelect)
	}
}

func (driver *Driver) selectOptions(codigo core.Codigo) {
	driver.selectWithSleep(core.Paths.Nivel, codigo.Nivel, "", "")
	driver.selectWithSleep(core.Paths.Sede, codigo.Sede, core.Paths.Nivel, codigo.Nivel)
	driver.selectWithSleep(core.Paths.Facultad, codigo.Facultad, core.Paths.Sede, codigo.Sede)
	driver.selectWithSleep(core.Paths.Carrera, codigo.Carrera, core.Paths.Facultad, codigo.Facultad)
	driver.selectWithSleep(core.Paths.Tipologia, codigo.Tipologia, core.Paths.Carrera, codigo.Carrera)
}

func (driver *Driver) GetTable() rod.Elements {

	var rows rod.Elements

	for {
		println("Buscando tabla...")

		table := driver.Page.Timeout(timeoutDuration).MustElement(".af_table_data-table-VH-lines").CancelTimeout()
		time.Sleep(3 * time.Second)

		if table == nil {
			continue
		}

		if !table.MustHas("tbody") {
			break
		}

		tbody := table.MustElement("tbody")
		if tbody == nil {
			continue
		}

		rows = tbody.MustElements("tr")
		if rows == nil || len(rows) > 200 {
			continue
		}

		break
	}

	return rows
}

func (driver *Driver) getPage() *rod.Page {

	page := driver.Browser.MustIncognito().MustPage("")

	router := InterceptRequests(page)
	go router.Run()

	page.Timeout(timeoutPage).MustNavigate(core.SIA_URL).MustWaitStable().CancelTimeout()

	return page
}
