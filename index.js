const fs = require('fs');
const { getAllMaterias } = require("./src/main.js");
const GRUPOS = require("./data/grupos.json");

async function WorkerCarrera(codigo) {
	const data = await getAllMaterias(codigo);
	return data;
}

async function main() {
	const index = process.argv[2];

	const grupoAsignado = GRUPOS[
		parseInt(index) - 1
	];
	console.log("Grupo asignado:", grupoAsignado);

	const DATA = {};
	const promises = grupoAsignado.map(async (carrera) => {
		console.log("Extrayendo:", carrera.nombre);
		const data = await WorkerCarrera(carrera.codigo);
		DATA[carrera.nombre] = data;
	});

	await Promise.all(promises);
	fs.writeFileSync(`${index}.json`, JSON.stringify(DATA));

	console.log("Done!");
}

main();