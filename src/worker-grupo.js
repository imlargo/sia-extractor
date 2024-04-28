const fs = require('fs');
const { getAllMaterias } = require("./main.js");

// import CODIGOS from "./src/codigos.json";

async function WorkerCarrera(nombreCarrera, codigoCarrera) {
	//const data = await getAllMaterias(codigoCarrera);
    fs.writeFile(`${nombreCarrera}.json`, JSON.stringify(data));
}

(() => {
    // Get args and call WorkerCarrera
    const args = process.argv.slice(2);
    const datagrupo = JSON.parse(args[0]);

    for (const carrera of datagrupo) {
        WorkerCarrera(carrera.nombre, carrera.codigo);
    }
})();