const fs = require('fs');
const { getAllMaterias } = require("./src/main.js");
const GRUPOS = require("./data/grupos.json");

async function WorkerCarrera(codigo, facultadName, carreraName) {
	const data = await getAllMaterias(codigo, facultadName, carreraName);
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
		console.log("Extrayendo:", carrera.carrera);
		const data = await WorkerCarrera(carrera.codigo, carrera.facultad, carrera.carrera);
		DATA[carrera.carrera] = data;
	});

	await Promise.all(promises);
	fs.writeFileSync(`${index}.json`, JSON.stringify(DATA));

	console.log("Done!");
}

main();