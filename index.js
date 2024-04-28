const { getCarreras } = require("./src/util.js");
const { getListMaterias, getMateria, getAllMaterias } = require("./src/main.js");

// import CODIGOS from "./src/codigos.json";

async function main() {
	const data = await getAllMaterias("0-6-5-13");
	console.log(data);
}

main();