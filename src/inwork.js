
async function extracFacultad() {
	
}


(async () => {

	GLOBALS.sede = (await getCarreras())["Pregrado"]["1102 SEDE MEDELLÍN"];

	delete GLOBALS.sede["3066 FACULTAD DE CIENCIAS AGROPECUARIAS"];
	delete GLOBALS.sede["3 SEDE MEDELLÍN"];

	const args = process.argv.slice(2);

	const tipo = args[0];
	const codigoCarrera = args[1];

	try {
		if (tipo === "list") {
			const materias = await getListMaterias(codigoCarrera);
			fs.writeFileSync("data.json", JSON.stringify(materias));
			console.log("⋆｡°✩ Data saved! ⋆｡°✩");
			process.exit();
		} else if (tipo === "materia") {
			const codigoMateria = args[2];
			const materia = await getMateria(codigoCarrera, codigoMateria);
			fs.writeFileSync("data.json", JSON.stringify(materia));
			console.log("⋆｡°✩ Data saved! ⋆｡°✩");
			process.exit();
		} else if (tipo === "all") {
			const materias = await getAllMaterias(codigoCarrera);
			fs.writeFileSync("data.json", JSON.stringify(materias));
			console.log("⋆｡°✩ Data saved! ⋆｡°✩");
			process.exit();
		} else if (tipo === "facultad") {
			const sede = {
				"3068 FACULTAD DE MINAS": {
					"3515 INGENIERÍA ADMINISTRATIVA": "0-6-4-0",
					"3528 INGENIERÍA ADMINISTRATIVA": "0-6-4-1",
					"3527 INGENIERÍA AMBIENTAL": "0-6-4-2",
					"3529 INGENIERÍA AMBIENTAL": "0-6-4-3",
					"3516 INGENIERÍA CIVIL": "0-6-4-4",
					"3530 INGENIERÍA CIVIL": "0-6-4-5",
					"3517 INGENIERÍA DE CONTROL": "0-6-4-6",
					"3531 INGENIERÍA DE CONTROL": "0-6-4-7",
					"3518 INGENIERÍA DE MINAS Y METALURGIA": "0-6-4-8",
					"3532 INGENIERÍA DE MINAS Y METALURGIA": "0-6-4-9",
					"3519 INGENIERÍA DE PETRÓLEOS": "0-6-4-10",
					"3533 INGENIERÍA DE PETRÓLEOS": "0-6-4-11",
					"3520 INGENIERÍA DE SISTEMAS E INFORMÁTICA": "0-6-4-12",
					"3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA": "0-6-4-13",
					"3521 INGENIERÍA ELÉCTRICA": "0-6-4-14",
					"3535 INGENIERÍA ELÉCTRICA": "0-6-4-15",
					"3522 INGENIERÍA GEOLÓGICA": "0-6-4-16",
					"3536 INGENIERÍA GEOLÓGICA": "0-6-4-17",
					"3523 INGENIERÍA INDUSTRIAL": "0-6-4-18",
					"3537 INGENIERÍA INDUSTRIAL": "0-6-4-19",
					"3524 INGENIERÍA MECÁNICA": "0-6-4-20",
					"3538 INGENIERÍA MECÁNICA": "0-6-4-21",
					"3525 INGENIERÍA QUÍMICA": "0-6-4-22",
					"3539 INGENIERÍA QUÍMICA": "0-6-4-23"
				}
			};

			const facultades = Object.keys(sede);

			const allData = {};

			for (const nombreFacultad of facultades) {
				console.log(`\n\n > - - -	${nombreFacultad}	- - - <\n`);

				const facultad = sede[nombreFacultad];
				const carreras = Object.keys(facultad);
				const totalCarreras = carreras.length;

				const dataFacultad = {};
				for (const carrera of carreras) {
					console.log(`\n\n⋆｡°✩ ${carreras.indexOf(carrera) + 1}/${totalCarreras}. ${carrera} ⋆｡°✩\n`);

					try {
						dataFacultad[carrera] = await getAllMaterias(facultad[carrera]);
					} catch (error) {
						console.log(`!!!Error!!! -> ${carrera} ⋆｡°✩`);
					}
				}

				allData[nombreFacultad] = dataFacultad;
			}

			fs.writeFileSync("data.json", JSON.stringify(allData));
			console.log("⋆｡°✩ Data saved! ⋆｡°✩");
			process.exit();
		} else if (tipo === "sede") {
			const sede = GLOBALS.sede;
			const facultades = Object.keys(sede);
			const allData = {};
			for (const nombreFacultad of facultades) {
				console.log(`\n\n > - - -	${nombreFacultad}	- - - <\n`);

				const facultad = sede[nombreFacultad];
				const carreras = Object.keys(facultad);
				const totalCarreras = carreras.length;

				const dataFacultad = {};
				for (const carrera of carreras) {
					console.log(`\n\n⋆｡°✩ ${carreras.indexOf(carrera) + 1}/${totalCarreras}. ${carrera} ⋆｡°✩\n`);

					try {
						dataFacultad[carrera] = await getAllMaterias(facultad[carrera]);
					} catch (error) {
						console.log(`!!!Error!!! -> ${carrera} ⋆｡°✩`);
					}
				}

				allData[nombreFacultad] = dataFacultad;
			}

			fs.writeFileSync("data.json", JSON.stringify(allData));
			console.log("⋆｡°✩ Data saved! ⋆｡°✩");
			process.exit();
		} else if (tipo === "new-beta") {
			const sede = GLOBALS.sede;
			const facultades = Object.keys(sede);

			async function extractFacultad(keyFacultad) {
				console.log(`\n\n > - - -	${keyFacultad}	- - - <\n`);
				const objFacultad = sede[keyFacultad];
				const carreras = Object.keys(objFacultad);
				const totalCarreras = carreras.length;
				const dataFacultad = {};
				for (const carrera of carreras) {
					console.log(`\n\n⋆｡°✩ ${carreras.indexOf(carrera) + 1}/${totalCarreras}. ${carrera} ⋆｡°✩\n`);
					try {
						dataFacultad[carrera] = await getAllMaterias(objFacultad[carrera]);
					} catch (error) {
						console.log(`!!!Error!!! -> ${carrera} ⋆｡°✩`);
					}
				}
				return dataFacultad;
			}

			const allData = {};

			await Promise.all(facultades.map(async (keyFacultad) => {
				const data = await extractFacultad(keyFacultad);
				allData[keyFacultad] = data;
			}));

			fs.writeFileSync("data.json", JSON.stringify(allData));
			console.log("⋆｡°✩ Data saved! ⋆｡°✩");
			process.exit();
		} else if (tipo === "chunk-1") {
			process.setMaxListeners(0);
			const sede = GLOBALS.sede;
			delete sede["3068 FACULTAD DE MINAS"];
			delete sede["3065 FACULTAD DE CIENCIAS"];
			const facultades = Object.keys(sede);

			async function extractFacultad(keyFacultad) {
				console.log(`\n\n > - - -	${keyFacultad}	- - - <\n`);
				const objFacultad = sede[keyFacultad];
				const carreras = Object.keys(objFacultad);
				const totalCarreras = carreras.length;
				const dataFacultad = {};

				await Promise.all(carreras.map(async (carrera, index) => {
					console.log(`\n\n⋆｡°✩ ${index + 1}/${totalCarreras}. ${carrera} ⋆｡°✩\n`);
					try {
						const materias = await getAllMaterias(objFacultad[carrera]);
						dataFacultad[carrera] = materias;
					} catch (error) {
						console.log(`!!!Error!!! -> ${carrera} ⋆｡°✩`);
					}
				}));
				return dataFacultad;
			}

			const allData = {};

			await Promise.all(facultades.map(async (keyFacultad) => {
				const data = await extractFacultad(keyFacultad);
				allData[keyFacultad] = data;
			}));

			fs.writeFileSync("chunk-1.json", JSON.stringify(allData));
			console.log("⋆｡°✩ Data saved! ⋆｡°✩");
			process.exit();
		} else if (tipo === "chunk-4") {
			process.setMaxListeners(0);
			const sede = {
				"3065 FACULTAD DE CIENCIAS": GLOBALS.sede["3065 FACULTAD DE CIENCIAS"],
			};
			const facultades = Object.keys(sede);

			async function extractFacultad(keyFacultad) {
				console.log(`\n\n > - - -	${keyFacultad}	- - - <\n`);
				const objFacultad = sede[keyFacultad];
				const carreras = Object.keys(objFacultad);
				const totalCarreras = carreras.length;
				const dataFacultad = {};

				await Promise.all(carreras.map(async (carrera, index) => {
					console.log(`\n\n⋆｡°✩ ${index + 1}/${totalCarreras}. ${carrera} ⋆｡°✩\n`);
					try {
						const materias = await getAllMaterias(objFacultad[carrera]);
						dataFacultad[carrera] = materias;
					} catch (error) {
						console.log(`!!!Error!!! -> ${carrera} ⋆｡°✩`);
					}
				}));
				return dataFacultad;
			}

			const allData = {};

			await Promise.all(facultades.map(async (keyFacultad) => {
				const data = await extractFacultad(keyFacultad);
				allData[keyFacultad] = data;
			}));

			fs.writeFileSync("chunk-4.json", JSON.stringify(allData));
			console.log("⋆｡°✩ Data saved! ⋆｡°✩");
			process.exit();
		} else if (tipo === "chunk-2") {
			process.setMaxListeners(0);
			const entries = Object.entries(GLOBALS.sede["3068 FACULTAD DE MINAS"])
			const sliced = Object.fromEntries(
				entries.slice(0, 13)
			);

			const sede = {
				"3068 FACULTAD DE MINAS": sliced,
			}
			const facultades = Object.keys(sede);

			async function extractFacultad(keyFacultad) {
				console.log(`\n\n > - - -	${keyFacultad}	- - - <\n`);
				const objFacultad = sede[keyFacultad];
				const carreras = Object.keys(objFacultad);
				const totalCarreras = carreras.length;
				const dataFacultad = {};

				await Promise.all(carreras.map(async (carrera, index) => {
					console.log(`\n\n⋆｡°✩ ${index + 1}/${totalCarreras}. ${carrera} ⋆｡°✩\n`);
					try {
						const materias = await getAllMaterias(objFacultad[carrera]);
						dataFacultad[carrera] = materias;
					} catch (error) {
						console.log(`!!!Error!!! -> ${carrera} ⋆｡°✩`);
					}
				}));
				return dataFacultad;
			}

			const allData = {};

			await Promise.all(facultades.map(async (keyFacultad) => {
				const data = await extractFacultad(keyFacultad);
				allData[keyFacultad] = data;
			}));

			fs.writeFileSync("chunk-2.json", JSON.stringify(allData));
			console.log("⋆｡°✩ Data saved! ⋆｡°✩");
			process.exit();
		} else if (tipo === "chunk-3") {
			process.setMaxListeners(0);
			const entries = Object.entries(GLOBALS.sede["3068 FACULTAD DE MINAS"])
			const sliced = Object.fromEntries(
				entries.slice(13, 23)
			);
			const sede = {
				"3068 FACULTAD DE MINAS": sliced
			}
			const facultades = Object.keys(sede);

			async function extractFacultad(keyFacultad) {
				console.log(`\n\n > - - -	${keyFacultad}	- - - <\n`);
				const objFacultad = sede[keyFacultad];
				const carreras = Object.keys(objFacultad);
				const totalCarreras = carreras.length;
				const dataFacultad = {};

				await Promise.all(carreras.map(async (carrera, index) => {
					console.log(`\n\n⋆｡°✩ ${index + 1}/${totalCarreras}. ${carrera} ⋆｡°✩\n`);
					try {
						const materias = await getAllMaterias(objFacultad[carrera]);
						dataFacultad[carrera] = materias;
					} catch (error) {
						console.log(`!!!Error!!! -> ${carrera} ⋆｡°✩`);
					}
				}));
				return dataFacultad;
			}

			const allData = {};

			await Promise.all(facultades.map(async (keyFacultad) => {
				const data = await extractFacultad(keyFacultad);
				allData[keyFacultad] = data;
			}));

			fs.writeFileSync("chunk-3.json", JSON.stringify(allData));
			console.log("⋆｡°✩ Data saved! ⋆｡°✩");
			process.exit();
		} else {
			console.log("Invalid command");
		}
	} catch (error) {
		console.error("An error occurred:", error);
	}

})();