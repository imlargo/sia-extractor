const puppeteer = require('puppeteer');
const fs = require('fs');
const { JSDOM } = require("jsdom");
const fetch = require("node-fetch");

const sleep = ms => new Promise(res => setTimeout(res, ms));

/**
 * Obtiene las carreras desde una API externa.
 * @returns {Promise<Object>} Un objeto que contiene los datos de las carreras.
 */
async function getCarreras() {
	const res = await fetch("https://www.bettercampus.com.co/api/careers?fileName=searchCodes");
	/*
	const data = await res.json();
	return data;
	*/
	const data = await res.text();
	return JSON.parse(data);
}

/**
 * Carga una página web en un navegador Puppeteer y devuelve el navegador y la página cargada.
 * @param {string} url - La URL de la página web a cargar.
 * @param {object} options - Opciones adicionales para la carga de la página.
 * @param {boolean} options.consola - Indica si se deben mostrar los mensajes de la consola del navegador.
 * @param {boolean} options.req - Indica si se deben mostrar las solicitudes enviadas desde la página.
 * @param {boolean} options.res - Indica si se deben mostrar las respuestas recibidas desde la página.
 * @param {string} options.selector - Selector CSS para esperar hasta que un elemento esté presente en la página.
 * @returns {Promise<Array>} - Una promesa que se resuelve con un arreglo que contiene el navegador y la página cargada.
 */
async function load(url, options = { consola: false, req: false, res: false, selector: false }) {
	const browser = await puppeteer.launch({
		headless: "new",
		args: ["--no-sandbox", "--disable-setuid-sandbox", '--incognito']
	});
	const page = await browser.newPage();

	try {
		console.log("Cargando página...");

		await page.goto(url, {
			waitUntil: "networkidle0",
			timeout: 120000,
		});

		if (options.selector) {
			await page.waitForSelector(options.selector);
		}

		// Opciones
		if (options.consola) {
			page.on('console', message => {
				console.log(`Mensaje del navegador: ${message.text()}`);
			});
		}

		if (options.req) {
			page.on('request', request => {
				console.log(`Solicitud enviada: ${request.url()}`);
				console.log(request.headers());
			});
		}

		if (options.res) {
			page.on('response', async response => {
				console.log(`Respuesta recibida: ${response.url()}`);
			});
		}

		console.log("	> Página cargada");

		return [browser, page];
	} catch (error) {
		console.error("Error al cargar la página:", error);
		if (browser) {
			await browser.close();
		}
		throw error;
	}
}

/**
 * Obtiene los valores de búsqueda a partir de un código de búsqueda y una indicación de si es electiva.
 * @param {string} searchCode - El código de búsqueda en el formato "nivel-sede-facultad-carrera".
 * @param {boolean} isElectiva - Indica si la búsqueda es para una electiva.
 * @returns {object} - Un objeto con los valores de búsqueda separados por nivel, sede, facultad, carrera y electiva.
 */
function getSearchValues(searchCode, isElectiva) {
	const [nivel, sede, facultad, carrera] = searchCode.split("-");
	const electiva = isElectiva ? "7" : "0";
	return { nivel, sede, facultad, carrera, electiva };
}

/**
 * Procesa una tabla de materias y devuelve un arreglo de objetos con los datos de cada materia.
 * @param {Array<Array<string>>} materias - La tabla de materias a procesar.
 * @returns {Array<Object>} - Un arreglo de objetos con los datos de cada materia.
 */
function processTable(materias) {
	return materias.map(rawMateria => {
		const [
			codigo, nombre, creditos, tipo, descripcion
		] = rawMateria.map(raw => raw.trim());
		return {
			codigo: codigo,
			nombre: nombre,
			creditos: creditos,
			tipo: tipo,
			descripcion: descripcion,
		};
	});
}

/**
 * Selecciona una opción en un elemento de selección desplegable en una página web.
 * @param {Page} page - La página de Puppeteer en la que se encuentra el elemento.
 * @param {string} selectId - El ID del elemento de selección desplegable.
 * @param {string} optionValue - El valor de la opción que se desea seleccionar.
 * @returns {Promise<void>} - Una promesa que se resuelve cuando se ha seleccionado la opción.
 */
async function selectOption(page, selectId, optionValue) {
	const selectElement = await page.$(`#${selectId}`);
	await selectElement.click();
	await selectElement.select(optionValue);
}

/**
 * Procesa la información de una materia y devuelve un objeto con los datos procesados.
 * @returns {Object} Objeto con los datos procesados de la materia.
 */
function procesarMateria() {

	/*
	// Ultimo elemento con la clase margin-t af_panelGroupLayout
	const allContainers = document.querySelectorAll(".margin-t.af_panelGroupLayout");
	const rawRequisitos = Array.from(allContainers[allContainers.length - 1].querySelectorAll("div"))
	const isRequisito = rawRequisitos[0].textContent.includes("Condición");

	const requisitos = isRequisito ? rawRequisitos.splice(1).map(div => {
		const dataSpans = Array.from(div.querySelector("span").querySelectorAll("span"));
		const [codigo, nombre] = dataSpans.map(span => span.textContent);
		return { codigo, nombre }
	}) : [];
	*/

	// Obtener el nombre de la materia
	const rawName = document.getElementsByTagName("h2")[0].textContent;
	const codigo = rawName.match(/\(([^)]+)\)/)[1].trim();
	const nombreMateria = rawName.replace(`(${codigo})`, "").trim();

	// Obtener tipologia, creditos y facultad
	const tipologia = document.querySelector(".detass-tipologia").querySelector("span").textContent.trim();
	const creditos = document.querySelector(".detass-creditos").querySelector("span").textContent.trim();
	const facultad = document.querySelector(".detass-centro").textContent.replace("Facultad: ", "").trim();

	// Obtener la marca de tiempo y fecha de extracción
	const marcaDeTiempo = new Date();
	const opciones = {
		timeZone: 'America/Bogota',
		hour12: true,
		hour: "2-digit",
		minute: "2-digit",
	};
	const fechaExtraccion = `${marcaDeTiempo.toLocaleDateString('es-CO', { timeZone: 'America/Bogota' })} - ${marcaDeTiempo.toLocaleTimeString('es-CO', opciones)}`;

	// Inicializar variables para los grupos y los cupos disponibles
	const grupos = [];
	let cuposDisponibles = 0;

	// Obtener los elementos de grupo en la página
	const elementosGrupo = document.querySelectorAll(".borde.salto:not(.ficha-docente)");

	// Recorrer cada elemento de grupo
	for (const elementoGrupo of elementosGrupo) {
		// Verificar si el elemento contiene información de prerrequisitos o correquisitos
		if (elementoGrupo.getElementsByClassName("margin-t")[1] == undefined) {
			break; // Evitar información de prerrequisitos o correquisitos
		}

		// Obtener los datos del grupo
		const datosGrupo = elementoGrupo.getElementsByClassName("margin-t")[1].children;

		// Extraer información del grupo
		const nombreGrupo = elementoGrupo.getElementsByClassName("af_showDetailHeader_title-text0 ")[0].textContent.replace(/\(.*?\)/g, '').trim();
		const profesor = datosGrupo[0].textContent.split(": ")[1].trim();
		const duracion = datosGrupo[3].textContent.split(": ")[1].trim();
		const jornada = datosGrupo[4].textContent.split(": ")[1].trim();

		// Inicializar arreglo para los horarios del grupo
		const horarios = [];

		// Obtener los datos de horario del grupo
		let datosHorario = datosGrupo[2].getElementsByClassName("af_panelGroupLayout")[0];
		if (datosHorario.childElementCount > 0) {
			datosHorario = [].slice.call(datosHorario.children[0].children).slice(2);
			for (const elementoHorario of datosHorario) {
				// Extraer información de cada horario
				const informacionHorario = elementoHorario.children[0].textContent.split(" ");
				const itemHorario = {
					dia: informacionHorario[0],
					inicio: informacionHorario[2],
					fin: informacionHorario[4].replace(".", ""),
				};
				horarios.push(itemHorario);
			}
		}

		// Obtener el número de cupos disponibles
		let cupos = "NaN";
		if (datosGrupo[5] !== undefined) {
			cupos = parseInt(datosGrupo[5].textContent.split(": ")[1]);
			cuposDisponibles += parseInt(datosGrupo[5].textContent.split(": ")[1]);
		}

		// Crear objeto para el grupo y agregarlo al arreglo de grupos
		const grupo = {
			grupo: nombreGrupo,
			cupos: cupos,
			profesor: profesor,
			duracion: duracion,
			jornada: jornada,
			horarios: horarios,
		};

		grupos.push(grupo);
	}

	// Crear objeto para la materia con los datos procesados
	const materia = {
		nombre: nombreMateria,
		codigo: codigo,
		tipologia: tipologia,
		creditos: creditos,
		facultad: facultad,
		fechaExtraccion: fechaExtraccion,
		cuposDisponibles: cuposDisponibles,
		grupos: grupos,
	};

	return materia;
}

module.exports = {
	load,
	getCarreras,
	sleep,
	getSearchValues,
	processTable,
	selectOption,
	procesarMateria
};