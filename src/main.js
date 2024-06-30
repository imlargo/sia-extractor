const puppeteer = require('puppeteer');
const fs = require('fs');
const { JSDOM } = require("jsdom");
const fetch = require("node-fetch");
const {
	load, sleep,
	getSearchValues, processTable, selectOption,
	procesarMateria
} = require('./util.js');

const GLOBALS = {
	url: 'https://sia.unal.edu.co/Catalogo/facespublico/public/servicioPublico.jsf?taskflowId=task-flow-AC_CatalogoAsignaturas',
	sede: {},
};

/**
 * Carga una carrera en el sistema.
 * @param {string} codigoCarrera - El código de la carrera a cargar.
 * @returns {Promise<Array>} - Una promesa que se resuelve con un arreglo que contiene el navegador, la página y la tabla cargada.
 */
async function loadCarrera(codigoCarrera, carreraName) {
	try {
		// Inicializar el navegador
		const [browser, page] = await load(
			GLOBALS.url,
			{
				consola: true,
				req: false,
				res: false
			}
		);

		// Obtener los valores de búsqueda
		const searchValues = getSearchValues(codigoCarrera, false)

		// IDs de los selectores
		const selectIds = {
			nivel: "pt1\\:r1\\:0\\:soc1\\:\\:content",
			sede: "pt1\\:r1\\:0\\:soc9\\:\\:content",
			facultad: "pt1\\:r1\\:0\\:soc2\\:\\:content",
			carrera: "pt1\\:r1\\:0\\:soc3\\:\\:content",
			electiva: "pt1\\:r1\\:0\\:soc4\\:\\:content",
		};

		console.log("	> Seleccionando elementos...");

		// IDs sin escape
		const rawIds = {
			nivel: "pt1:r1:0:soc1::content",
			sede: "pt1:r1:0:soc9::content",
			facultad: "pt1:r1:0:soc2::content",
			carrera: "pt1:r1:0:soc3::content",
			electiva: "pt1:r1:0:soc4::content",
		};

		// Seleccionar nivel
		await selectOption(page, selectIds.nivel, searchValues.nivel);
		console.log("Nivel seleccionado");

		try {
			// Esperar a la respuesta y seleccionar sede
			await page.waitForResponse(res => res.url().includes("sia"));
			await page.waitForFunction(() => {
				const element = document.getElementById("pt1:r1:0:soc9::content");
				const disabled = element.disabled;
				return disabled == false;
			}, { timeout: 15000 });
			await selectOption(page, selectIds.sede, searchValues.sede);
		} catch (error) {
			console.error("Error al seleccionar sede:", error);
			throw error;
		}

		try {
			// Esperar a la respuesta y seleccionar facultad
			await page.waitForResponse(res => res.url().includes("sia"));
			await sleep(1000);
			await page.waitForFunction(() => {
				const element = document.getElementById("pt1:r1:0:soc2::content");
				const disabled = element.disabled;
				return disabled == false;
			}, { timeout: 15000 });
			await selectOption(page, selectIds.facultad, searchValues.facultad);
		} catch (error) {
			console.error("Error al seleccionar facultad:", error);
			throw error;
		}

		try {
			// Esperar a la respuesta y seleccionar carrera
			await page.waitForResponse(res => res.url().includes("sia"));
			await page.waitForFunction(() => {
				const element = document.getElementById("pt1:r1:0:soc3::content");
				const isDisabled = element.disabled;
				return !isDisabled;
			}, { timeout: 15000 });
			await selectOption(page, selectIds.carrera, searchValues.carrera);
		} catch (error) {
			console.error("Error al seleccionar carrera: ", carreraName,  error);
			throw error;
		}



		try {
			// Esperar a la respuesta y seleccionar electiva
			await page.waitForResponse(res => res.url().includes("sia"));
			await page.waitForFunction(() => {
				const element = document.getElementById("pt1:r1:0:soc4::content");
				const disabled = element.disabled;
				return disabled == false;
			}, { timeout: 15000 });
			await selectOption(page, selectIds.electiva, searchValues.electiva);
		} catch (error) {
			console.log("Error: ", carreraName);
			console.error("Error al seleccionar electiva:", error);
			throw error;
		}


		await sleep(1000);
		// Seleccionar días
		await page.evaluate(() => {
			const checkboxes = document.querySelectorAll(".af_selectBooleanCheckbox_native-input");
			checkboxes.forEach((checkbox) => checkbox.checked = true);
		});

		// Imprimir mensaje de selección exitosa
		console.log("	> Todos los elementos seleccionados");

		// Hacer clic en el botón para ejecutar la búsqueda
		const button = await page.$(".af_button_link");
		button.click();

		// Esperar a la respuesta y cargar la tabla
		await page.waitForResponse(res => res.url() == "https://sia.unal.edu.co/Catalogo/afr/wk-column-select.cur");
		await sleep(1000);

		// Imprimir mensaje de carga exitosa
		console.log("	> Datos cargados con éxito");

		// Obtener la tabla
		const table = await page.$(".af_table_data-table-VH-lines");

		return [browser, page, table];
	} catch (error) {
		console.error("Error loading carrera:", error);
		throw error;
	}
}

/**
 * Obtiene una lista de materias para una carrera específica.
 * @param {string} codigoCarrera - El código de la carrera.
 * @returns {Promise<Array>} - Una promesa que se resuelve con una lista de materias.
 */
async function getListMaterias(codigoCarrera) {
	try {
		// Cargar carrera
		const [browser, page, table] = await loadCarrera(codigoCarrera);

		// Obtener datos de la tabla
		const data = await page.evaluate((table) => {
			const data = [];
			const rows = table.querySelectorAll("tr");
			for (const row of rows) {
				const cells = row.querySelectorAll("td");
				const rowData = Array.from(cells, (cell) => cell.textContent);
				data.push(rowData);
			}
			return data;
		}, table);

		// Cerrar navegador
		await browser.close();

		// Procesar y retornar la lista de materias
		return processTable(data);
	} catch (error) {
		// Manejar errores
		console.error("Error al obtener la lista de materias:", error);
		throw error;
	}
}

/**
 * Obtiene los detalles de una materia específica para una carrera dada.
 * @param {string} codigoCarrera - El código de la carrera.
 * @param {string} codigoMateria - El código de la materia.
 * @returns {Promise<Object>} - Una promesa que se resuelve con los detalles de la materia.
 */
async function getMateria(codigoCarrera, codigoMateria) {
	try {
		// Cargar carrera
		const [browser, page, table] = await loadCarrera(codigoCarrera);

		// Buscar y hacer clic en la materia específica
		await page.evaluate((table, codigoMateria) => {
			const courses = Array.from(table.querySelectorAll(".af_commandLink"));
			const course = courses.find((course) => course.textContent === codigoMateria);
			if (course) {
				course.click();
			} else {
				throw new Error(`No se encontró la materia con el código ${codigoMateria}`);
			}
		}, table, codigoMateria);

		// Esperar a que se carguen los detalles de la materia
		await page.waitForSelector(".af_showDetailHeader_content0", {
			timeout: 3000,
		});

		// Obtener los detalles de la materia
		const materia = await page.evaluate(procesarMateria);

		// Cerrar el navegador
		await browser.close();

		return materia;
	} catch (error) {
		console.error("Error al obtener los detalles de la materia:", error);
		throw error;
	}
}

/**
 * Obtiene todas las materias para una carrera específica.
 * @param {string} codigoCarrera - El código de la carrera.
 * @returns {Promise<Array>} - Una promesa que se resuelve con una lista de materias.
 */
async function getAllMaterias(codigoCarrera, facultadName, carreraName) {
	try {
		// Cargar carrera
		const [browser, page, table] = await loadCarrera(codigoCarrera, carreraName);

		const materias = [];

		// Obtener la lista de cursos
		const courses = await page.$$(".af_commandLink");
		const size = courses.length;

		console.log(`	> ${size} materias encontradas para el plan de estudios!`);

		for (let i = 0; i < size; i++) {

			console.log(`${i} - ${carreraName}`);

			// Obtener la lista de cursos nuevamente para evitar errores de referencia
			const courses = await page.$$(".af_commandLink");

			// Navegar a la página de la materia
			const element = courses[i];
			try {
				await element.click();
			} catch {
				console.error("No se pudo hacer clic en el enlace del curso");
			}

			// Esperar a que se cargue la información de la materia
			try {
				await page.waitForSelector(".af_showDetailHeader_content0", {
					timeout: 7000,
				});

				// Obtener los detalles de la materia
				const materia = await page.evaluate(procesarMateria);
				materia.facultad = facultadName;
				materia.carrera = carreraName;
				materias.push(materia);
				//console.log(materia.nombre);
			} catch (e) {
				console.log(`${i}: ¡NO SE ENCONTRÓ INFORMACIÓN!`);
			}

			// Regresar a la lista de cursos
			const backButton = await page.$(".af_button");
			await backButton.click();

			try {
				// Esperar a que se cargue la lista de cursos nuevamente
				await page.waitForSelector(".af_selectBooleanCheckbox_native-input", {
					timeout: 15000,
				});
			} catch {
				// Si no se carga la lista de cursos, hacer clic en el botón de regresar nuevamente
				await backButton.click();
			}

		}

		console.log("¡Todas las materias han sido procesadas con exito!");

		// Cerrar el navegador
		await browser.close();

		return materias;
	} catch (error) {
		console.error("Error al obtener todas las materias:", error);
		throw error;
	}
}

module.exports = {
	getListMaterias,
	getMateria,
	getAllMaterias,
}