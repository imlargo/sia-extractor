import type { Prerequisito, Asignatura, Horario, Grupo, PrerequisitoAsignatura } from "./types"

() => {

    function extraerMetadatos() {
        const tipologia = document.querySelector(".detass-tipologia")?.querySelector("span")?.textContent?.trim() ?? '';
        const creditos = document.querySelector(".detass-creditos")?.querySelector("span")?.textContent?.trim() ?? '';
        const facultad = document.querySelector(".detass-centro")?.textContent?.replace("Facultad: ", "").trim() ?? '';

        return {
            tipologia,
            creditos,
            facultad
        }
    }

    function extraerNombreYCodigo() {

        const rawName = document.getElementsByTagName("h2")[0]?.textContent ?? '';
        const codigoMatch = rawName.match(/\(([^)]+)\)/);
        const codigo = codigoMatch ? codigoMatch[1].trim() : '';
        const nombreMateria = rawName.replace(`(${codigo})`, "").trim();

        return {
            codigo,
            nombreMateria
        }
    }

    function getMarcaTemporal() {

        const currentDate = new Date();

        const fecha = currentDate.toLocaleDateString('es-CO', { timeZone: 'America/Bogota' });
        const hora = currentDate.toLocaleTimeString('es-CO', {
            timeZone: 'America/Bogota',
            hour12: true,
            hour: "2-digit",
            minute: "2-digit",
        })

        return `${fecha} - ${hora}`;
    }

    function extraerDatos() {
        // Inicializar variables para los grupos y los cupos disponibles
        let totalCupos = 0;
        const grupos: Grupo[] = [];
        const prerequisitos: Prerequisito[] = []

        // Recorrer cada elemento de grupo
        const elementosGrupo = document.querySelectorAll(".borde.salto:not(.ficha-docente)");
        for (const elementoGrupo of Array.from(elementosGrupo)) {

            const childItems = elementoGrupo.getElementsByClassName("margin-t")

            // Verificar si el elemento contiene información de prerrequisitos o correquisitos
            if (childItems[1] == undefined) {
                const prerequisito = getPrerequisito(childItems)
                prerequisitos.push(prerequisito)
                continue;
            }

            const datosGrupo = childItems[1].children;

            const { nombreGrupo, profesor, duracion, jornada } = getGroupData(datosGrupo, elementoGrupo)
            const horarios = getHorarios(datosGrupo)

            // Obtener el número de cupos disponibles
            const cuposText = datosGrupo[5]?.textContent?.split(": ")[1] ?? '';
            const cuposGrupo = cuposText === "" ? -1 : parseInt(cuposText);
            if (cuposGrupo !== -1) {
                totalCupos += cuposGrupo;
            }

            // Crear objeto para el grupo y agregarlo al arreglo de grupos
            const grupo = {
                grupo: nombreGrupo,
                cupos: cuposGrupo,
                profesor: profesor,
                duracion: duracion,
                jornada: jornada,
                horarios: horarios,
            };

            grupos.push(grupo);
        }

        return {
            grupos,
            prerequisitos,
            cuposDisponibles: totalCupos,
        }
    }

    function getGroupData(datosGrupo: HTMLCollection, elementoGrupo: Element) {
        // Obtener los datos del grupo

        // Extraer información del grupo
        const elEncabezado = elementoGrupo.getElementsByClassName("af_showDetailHeader_title-text0 ")[0];
        const nombreGrupo = elEncabezado && elEncabezado.textContent ? elEncabezado.textContent.replace(/\(.*?\)/g, '').trim() : '';

        const profesor = datosGrupo[0]?.textContent?.split(": ")[1]?.trim() ?? '';
        const duracion = datosGrupo[3]?.textContent?.split(": ")[1]?.trim() ?? '';
        const jornada = datosGrupo[4]?.textContent?.split(": ")[1]?.trim() ?? '';

        return {
            nombreGrupo,
            profesor,
            duracion,
            jornada,
        }

    }

    function getPrerequisito(childItems: HTMLCollectionOf<Element>): Prerequisito {
        /*
            Tipo de prerrequisito implica. 
            M - no se puede matricular la asignatura sin superar el prerrequisito. 
            O - podrá matricular, pero no ser calificado sin la superación del prerrequisito. 
            E - o matricula el prerrequisito simultáneamente, o lo ha matriculado alguna vez. 
            A - anulación por incompatibilidad. Si se matricula de las dos asignaturas afectadas por el prerrequisito y no supera la asignatura llave, las asignaturas afectadas por el prerrequisito aparecerán como anuladas.
        */

        const asignaturas: PrerequisitoAsignatura[] = []

        // Recorrer asignaturas
        const subElementos = childItems[0].childNodes
        for (let i = 1; i < subElementos.length; i++) {
            const firstChild = subElementos[i].firstChild;
            if (!firstChild) continue;

            const asignatura = firstChild.childNodes;
            const codigo = asignatura[0]?.textContent?.trim() ?? ''
            const nombre = asignatura[1]?.textContent?.trim() ?? ''

            asignaturas.push({
                codigo: codigo,
                nombre: nombre
            })
        }

        const childsEncabezado = subElementos[0]?.firstChild?.childNodes ?? []

        const tipo = childsEncabezado[3]?.textContent?.trim() ?? ''
        const isTodas = childsEncabezado[5].textContent === "[S]"
        const cantidad = isTodas ? asignaturas.length : parseInt(childsEncabezado[7]?.textContent?.replace(/\[|\]/g, "") ?? '-1')

        const prerequisito: Prerequisito = {
            tipo: tipo,
            isTodas: isTodas,
            cantidad: cantidad,
            asignaturas: asignaturas
        }

        return prerequisito;

    }

    function getHorarios(datosGrupo: HTMLCollection): Horario[] {
        // Inicializar arreglo para los horarios del grupo
        const horarios: Horario[] = [];

        // Obtener los datos de horario del grupo
        const datosHorario = datosGrupo[2].getElementsByClassName("af_panelGroupLayout")[0];

        if (datosHorario.childElementCount > 0) {
            const elementosHorario = Array.from(datosHorario.children[0].children).slice(2);

            for (const elementoHorario of elementosHorario) {

                // Extraer información de cada horario
                const informacionHorario = elementoHorario.children[0]?.textContent?.split(" ") ?? [];
                if (informacionHorario.length < 4) continue;

                const itemHorario = {
                    dia: informacionHorario[0],
                    inicio: informacionHorario[2],
                    fin: informacionHorario[4].replace(".", ""),
                };

                horarios.push(itemHorario);
            }
        }

        return horarios
    }

    // Obtener el nombre de la materia
    const { codigo, nombreMateria } = extraerNombreYCodigo()

    // Obtener tipologia, creditos y facultad
    const { tipologia, creditos, facultad } = extraerMetadatos()

    const { grupos, prerequisitos, cuposDisponibles } = extraerDatos()

    // Crear objeto para la materia con los datos procesados
    const asignatura: Asignatura = {
        nombre: nombreMateria,
        codigo: codigo,
        tipologia: tipologia,
        creditos: parseInt(creditos),
        carrera: "",
        facultad: facultad,
        fechaExtraccion: getMarcaTemporal(),
        cuposDisponibles: cuposDisponibles,
        prerequisitos: prerequisitos,
        grupos: grupos,
    };

    return asignatura;
}