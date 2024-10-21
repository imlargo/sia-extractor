package core

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
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

func (driver *Driver) selectOptions(codigo Codigo) {
	driver.selectWithSleep(Paths.Nivel, codigo.Nivel, "", "")
	driver.selectWithSleep(Paths.Sede, codigo.Sede, Paths.Nivel, codigo.Nivel)
	driver.selectWithSleep(Paths.Facultad, codigo.Facultad, Paths.Sede, codigo.Sede)
	driver.selectWithSleep(Paths.Carrera, codigo.Carrera, Paths.Facultad, codigo.Facultad)
	driver.selectWithSleep(Paths.Tipologia, codigo.Tipologia, Paths.Carrera, codigo.Carrera)
}

func (driver *Driver) GetTable() rod.Elements {

	var rows rod.Elements

	for {
		println("Buscando tabla...")

		table := driver.Page.MustElement(".af_table_data-table-VH-lines")
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
