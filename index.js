const CODIGOS = require("./src/codigos.json");
const child = require('child_process');

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

	return grupos
}


async function main() {
	const grupos = agrupar();

	for (const grupo of grupos) {
		// console.log(grupo);
		const proceso = child.spawn(
			'node', ['./src/worker-grupo.js', JSON.stringify(grupo)
		],
			{ detached: true }
		);

		proceso.stdout.on('data', (data) => {
			console.log(`stdout: ${data}`);
		});
	}
	
}

main();