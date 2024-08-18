const fs = require('fs/promises');
const CODIGOS = require("../data/codigos.json");
const { CONFIG } = require("./config.js");

async function loadJson(path) {
    const data = await fs.readFile(path, 'utf-8');
    return JSON.parse(data);
}

async function main() {

    const indexs = Array.from({ length: CONFIG.totalGrupos }, (_, i) => i + 1);

    const data = await Promise.all(
        indexs.map(async i => await loadJson(`../../artifacts/${i}.json`))
    );

    const merged = data.reduce((acc, curr) => ({ ...acc, ...curr }));

    const DATA = {};
    for (const entries of Object.entries(CODIGOS)) {
        const [facultad, carreras] = entries;

        const dataFacultad = {};
        Object.keys(carreras).forEach(carrera => {
            dataFacultad[carrera] = merged[carrera];
        });

        DATA[facultad] = dataFacultad;
    }

    await fs.writeFile('data.json', JSON.stringify(DATA));
}

main();