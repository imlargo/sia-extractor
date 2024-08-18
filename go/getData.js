() => {

	/*
	// Ultimo elemento con la clase margin-t af_panelGroupLayout
	const allContainers = document.querySelectorAll(".margin-t.af_panelGroupLayout");
	const rawRequisitos = Array.from(allContainers[allContainers.length - 1].querySelectorAll("div"))
	const isRequisito = rawRequisitos[0].textContent.includes("Condición");

	const requisitos = isRequisito ? rawRequisitos.splice(1).map(div => {
		const dataSpans = Array.from(div.querySelector("span").querySelectorAll("span"));
		const [codigo, nombre] = dataSpans.map(span => span.textContent);
		return { codigo, nombre }
	}) : [];
	*/

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

	// Recorrer cada elemento de grupo
	for (const elementoGrupo of elementosGrupo) {
		// Verificar si el elemento contiene información de prerrequisitos o correquisitos
		if (elementoGrupo.getElementsByClassName("margin-t")[1] == undefined) {
			break; // Evitar información de prerrequisitos o correquisitos
		}

		// Obtener los datos del grupo
		const datosGrupo = elementoGrupo.getElementsByClassName("margin-t")[1].children;

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

	/*

	const prerrequisitos = elementosGrupo[elementosGrupo.length-1];
	const hasRequisitos = prerrequisitos.getElementsByClassName("margin-t")[1] == undefined;

	const requisitos = hasRequisitos ? prerrequisitos.parentElement.querySelectorAll(".borde.salto.af_panelGroupLayout").map(container => {
		const listado = Array.from(container.firstChild.childNodes);
		const datos = listado.shift();
		return listado.map(req => Array.from(req.firstChild.childNodes).map(node => node.textContent))
	}) : [];

	if (hasRequisitos) {
		console.log(requisitos);
	}
	*/
	
	// Crear objeto para la materia con los datos procesados
	const materia = {
		nombre: nombreMateria.toString(),
		codigo: codigo.toString(),
		tipologia: tipologia.toString(),
		creditos: creditos.toString(),
		facultad: facultad.toString(),
		fechaExtraccion: fechaExtraccion.toString(),
		cuposDisponibles: cuposDisponibles.toString(),
		grupos: grupos,
	};

	return materia;
}