const fs = require('fs/promises');
const CARRERAS = require("../../data/carreras.json");
const { CONFIG } = require("./config.js");

async function loadJson(path) {
    const data = await fs.readFile(path, 'utf-8');
    return JSON.parse(data);
}

function GroupBy(array, func) {
	return array.reduce((acc, obj) => {
		const key = func(obj);
		if (!acc[key]) {
			acc[key] = [];
		}
		acc[key].push(obj);
		return acc;
	}, {});
}

async function main() {

    const indexs = Array.from({ length: CONFIG.totalGrupos }, (_, i) => i + 1);

    const data = await Promise.all(
        indexs.map(async i => await loadJson(`../../artifacts/${i}.json`))
    );

    // Hacer un objeto con todas las carreras
    const merged = data.reduce((acc, curr) => ({ ...acc, ...curr }));

    // Agrupar carreas por facultad
    const agrupado = GroupBy(CARRERAS, (carrera) => carrera.facultad); 

    const DATA = {};
    for (const entries of Object.entries(agrupado)) {
        const [facultad, carreras] = entries;

        const dataFacultad = {};
        carreras.forEach(carrera => {
            dataFacultad[carrera.carrera] = merged[carrera.carrera];
        });

        DATA[facultad] = dataFacultad;
    }

    await fs.writeFile('data.json', JSON.stringify(DATA));
}

main();