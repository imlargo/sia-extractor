const fs = require('fs');
const CODIGOS = require("./codigos.json");

const CONFIG = {
	cantidad: 5,
};

function agrupar() {
	// Unir todas las carreras en un solo objeto
	const carreras = Object.values(CODIGOS).reduce((acc, curr) => ({ ...acc, ...curr }) );
	const entries = Object.entries(carreras)

	// Dividir las carreras en arrays de 5 carreras
	const grupos = [];
	let grupo = [];
	for (let i = 0; i < entries.length; i++) {
		const [nombre, codigo] = entries[i];
		grupo.push({nombre, codigo});
		if (grupo.length === CONFIG.cantidad) {
			grupos.push(grupo);
			grupo = [];
		}
	}

	fs.writeFile("grupos.json", JSON.stringify(grupos), (err) => {});
}

agrupar();