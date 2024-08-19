() => {

	// Obtener el nombre de la materia
	const rawName = document.getElementsByTagName("h2")[0].textContent;
	const codigo = rawName.match(/\(([^)]+)\)/)[1].trim();
	const nombreMateria = rawName.replace(`(${codigo})`, "").trim();

	// Obtener tipologia, creditos y facultad
	const tipologia = document.querySelector(".detass-tipologia").querySelector("span").textContent.trim();
	const creditos = document.querySelector(".detass-creditos").querySelector("span").textContent.trim();
	const facultad = document.querySelector(".detass-centro").textContent.replace("Facultad: ", "").trim();

	// Obtener la marca de tiempo y fecha de extracción
	const marcaDeTiempo = new Date();
	const opciones = {
		timeZone: 'America/Bogota',
		hour12: true,
		hour: "2-digit",
		minute: "2-digit",
	};
	const fechaExtraccion = `${marcaDeTiempo.toLocaleDateString('es-CO', { timeZone: 'America/Bogota' })} - ${marcaDeTiempo.toLocaleTimeString('es-CO', opciones)}`;

	// Inicializar variables para los grupos y los cupos disponibles
	const grupos = [];
	let cuposDisponibles = 0;

	// Obtener los elementos de grupo en la página
	const elementosGrupo = document.querySelectorAll(".borde.salto:not(.ficha-docente)");


	const prerequisitos = []

	// Recorrer cada elemento de grupo
	for (const elementoGrupo of elementosGrupo) {

		const childItems = elementoGrupo.getElementsByClassName("margin-t")

		// Verificar si el elemento contiene información de prerrequisitos o correquisitos
		if (childItems[1] == undefined) {
			/*
				Tipo de prerrequisito implica. 
				M - no se puede matricular la asignatura sin superar el prerrequisito. 
				O - podrá matricular, pero no ser calificado sin la superación del prerrequisito. 
				E - o matricula el prerrequisito simultáneamente, o lo ha matriculado alguna vez. 
				A - anulación por incompatibilidad. Si se matricula de las dos asignaturas afectadas por el prerrequisito y no supera la asignatura llave, las asignaturas afectadas por el prerrequisito aparecerán como anuladas.
			*/

			const subElementos = childItems[0].childNodes

			const asignaturas = []
			// Recorrer asignaturas
			for (let i = 1; i < subElementos.length; i++) {
				const asignatura = subElementos[i].firstChild.childNodes
				const codigo = asignatura[0].textContent.trim()
				const nombre = asignatura[1].textContent.trim()

				asignaturas.push({ codigo, nombre })
			}

			const encabezado = subElementos[0].firstChild.childNodes

			const tipo = encabezado[3].textContent.trim()
			const isTodas = encabezado[5].textContent === "[S]"
			const cantidad = isTodas ? asignaturas.length : parseInt(encabezado[7].textContent.replace(/\[|\]/g, ""))

			prerequisitos.push({
				tipo: tipo,
				isTodas: isTodas,
				cantidad: cantidad,
				asignaturas: asignaturas
			})
			
			continue;
		}

		// Obtener los datos del grupo
		const datosGrupo = childItems[1].children;

		// Extraer información del grupo
		const nombreGrupo = elementoGrupo.getElementsByClassName("af_showDetailHeader_title-text0 ")[0].textContent.replace(/\(.*?\)/g, '').trim();
		const profesor = datosGrupo[0].textContent.split(": ")[1].trim();
		const duracion = datosGrupo[3].textContent.split(": ")[1].trim();
		const jornada = datosGrupo[4].textContent.split(": ")[1].trim();

		// Inicializar arreglo para los horarios del grupo
		const horarios = [];

		// Obtener los datos de horario del grupo
		let datosHorario = datosGrupo[2].getElementsByClassName("af_panelGroupLayout")[0];
		if (datosHorario.childElementCount > 0) {
			datosHorario = [].slice.call(datosHorario.children[0].children).slice(2);
			for (const elementoHorario of datosHorario) {
				// Extraer información de cada horario
				const informacionHorario = elementoHorario.children[0].textContent.split(" ");
				const itemHorario = {
					dia: informacionHorario[0],
					inicio: informacionHorario[2],
					fin: informacionHorario[4].replace(".", ""),
				};
				horarios.push(itemHorario);
			}
		}

		// Obtener el número de cupos disponibles
		let cupos = "NaN";
		if (datosGrupo[5] !== undefined) {
			cupos = parseInt(datosGrupo[5].textContent.split(": ")[1]);
			cuposDisponibles += parseInt(datosGrupo[5].textContent.split(": ")[1]);
		}

		// Crear objeto para el grupo y agregarlo al arreglo de grupos
		const grupo = {
			grupo: nombreGrupo,
			cupos: cupos,
			profesor: profesor,
			duracion: duracion,
			jornada: jornada,
			horarios: horarios,
		};

		grupos.push(grupo);
	}
	
	// Crear objeto para la materia con los datos procesados
	const materia = {
		nombre: nombreMateria,
		codigo: codigo,
		tipologia: tipologia,
		creditos: parseInt(creditos),
		facultad: facultad,
		fechaExtraccion: fechaExtraccion,
		cuposDisponibles: cuposDisponibles,
		prerequisitos: prerequisitos,
		grupos: grupos,
	};

	return materia;
}