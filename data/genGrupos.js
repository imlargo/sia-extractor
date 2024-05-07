const fs = require('fs');
const CODIGOS = require("./codigos.json");
const { CONFIG } = require("../src/config.js");

function agrupar() {
	// Unir todas las carreras en un solo objeto
	const carreras = Object.values(CODIGOS).reduce((acc, curr) => ({ ...acc, ...curr }) );
	const entries = Object.entries(carreras)

	// Dividir las carreras en arrays de 5 carreras
	const grupos = [];
	for (let i = 0; i < entries.length; i += CONFIG.cantidadPorGrupo) {
		const grupo = entries.slice(i, i + CONFIG.cantidadPorGrupo).map(([nombre, codigo]) => ({nombre, codigo}));
		grupos.push(grupo);
	}

	console.log("Grupos:", grupos.length);

	fs.writeFile("grupos.json", JSON.stringify(grupos), (err) => {});
}

agrupar();