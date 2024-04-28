const fs = require('fs');
const { getAllMaterias } = require("./src/main.js");
const GRUPOS = require("./data/grupos.json");

async function WorkerCarrera(codigo) {
	const data = await getAllMaterias(codigo);
	return data;
}

async function main() {
	const input = process.argv[2];
	const index = parseInt(input) - 1;

	const grupoAsignado = GRUPOS[index];
	console.log("Grupo asignado:", grupoAsignado);

	/*
	const DATA = {};
	const promises = grupoAsignado.map(async (carrera) => {
		const data = await WorkerCarrera(carrera.codigo);
		DATA[carrera.nombre] = data;
	});

	await Promise.all(promises);
	fs.writeFileSync(`${index}.json`, JSON.stringify(allData));
	*/

	console.log("Done!");
}

main();