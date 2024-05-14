const fs = require('fs');
const CODIGOS = require("./codigos.json");
const { CONFIG } = require("../src/config.js");

function agrupar() {

	const allCarreras = [];
	for (const entryFacultad of Object.entries(CODIGOS)) {
		const [facultad, datosFacultad] = entryFacultad;

		for (const entryCarrera of Object.entries(datosFacultad)) {
			const [carrera, codigoCarrera] = entryCarrera;

			allCarreras.push({
				facultad: facultad,
				carrera: carrera,
				codigo: codigoCarrera,
			})
		}
	}
	// Dividir las carreras en arrays de 5 carreras
	const grupos = [];
	for (let i = 0; i < allCarreras.length; i += CONFIG.cantidadPorGrupo) {
		const grupo = allCarreras.slice(i, i + CONFIG.cantidadPorGrupo)
		grupos.push(grupo);
	}

	console.log("Grupos:", grupos.length);

	fs.writeFile("grupos.json", JSON.stringify(grupos), (err) => {});
}

agrupar();