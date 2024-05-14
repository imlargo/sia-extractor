const fs = require('fs');

const history = [
    "46a3c87f44ed2ff4e6583cc1ece21e772ff35da9",
    "30b2fedacf5f23c23372c4c7495a7696a0a96c51",
    "4dea63613b0f73732cb3caabb5e8a0c40429b2aa",
    "944ab150998f75154482e3c314a1108c0d642d81",
    "ca6906333c0eaae13de941f35924ccca3593a122",
    "871bfa6bfd314006b58f9a37dbb16d795297c6f7",
    "f5f0e4a617eafbb757b79e6b8c81920d315990d1",
    "dc322a70271e0e58db8962231fb7e07092e7cd18",
    "e06ae8332e4aba5c1b4b1fc614cf1d7173cc054c",
    "c03953231cb5f0215996ce54545b2d586f472508",
    "c5876620e396afb6525c51a24f6b63a234ea0eda",
    "a639e576428d640d5bd699b5e4df5f6062b1e6af",
    "e4e41443c0ebb1bff63ec5dfdafe1a8f20de01b8",
    "5398e2e6c363ffd9fe00ca9cc26c3c6c7af0da67",
    "2f40f9a94b4faa57cd4b33a49c25b87877c16a5f",
    "3d7241871de2054c6d0998ae5f46cb8794e02985",
    "22b9fcff707ac35cb3cf27fead9f5ae3856cf45b",
    "e5ce44be1622aa89fe1e3cf3122f3d397e875c08",
    "db395a6cea8e05e21c25ea34b0d451c93aba54a7",
    "8dbf167cf9bfdcf063fcec64b76f87a88d8043c4",
    "6a70f3462e35f8491d26cc75610ae72f7acada77",
    "d323106011695be10737a8a4b4c8213482a66afc",
    "1073dc408a9103bb91008ee740ed039ce528a1a1", // Original
];

(async () => {
    const promises = history.map(async sha => {
        const link = `https://raw.githubusercontent.com/imlargo/api/${sha}/data.json`;
        return await fetch(link).then(res => res.json());
    })
    const allData = await Promise.all(promises).then(allData => allData.reverse());

    const compiled = {};
    allData.forEach(datosDia => {

        for (const entryFacultad of Object.entries(datosDia)) {
            const [facultad, datosFacultad] = entryFacultad;

            if (!compiled.hasOwnProperty(facultad)) {
                compiled[facultad] = {};
            }

            for (const entryCarrera of Object.entries(datosFacultad)) {

                const [carrera, datosCarrera] = entryCarrera;

                if (!compiled[facultad].hasOwnProperty(carrera)) {
                    compiled[facultad][carrera] = {};
                }

                for (const materia of datosCarrera) {
                    extractData(
                        materia, compiled[facultad][carrera], facultad, carrera
                    );
                }
            }
        }
    });

    fs.writeFile("analisis.json", JSON.stringify(compiled), (err) => { });

})();

function processDate(rawString) {
    let regex = /(\d{1,2}:\d{2}:\d{2})/;
    let match = rawString.match(regex);

    const times = match[0].split(':');
    const hours = parseInt(times[0]);
    const minutes = times[1];

    return `${hours}:${minutes} ${hours > 6 && hours < 12 ? 'a.m.' : 'p.m.'}`;
}

function extractData(materia, carrera, facultadName, carreraName) {

    const fecha = materia.fechaExtraccion.includes('1/28/2024') ? "6:00 a. m." : processDate(materia.fechaExtraccion);

    if (!carrera.hasOwnProperty(materia.codigo)) {
        carrera[materia.codigo] = {
            nombre: materia.nombre,
            tipologia: materia.tipologia,
            codigo: materia.codigo,
            grupos: {},
            total: {},
            facultad: facultadName,
            carrera: carreraName,
        };
    }

    const refMateria = carrera[materia.codigo];
    refMateria.total[fecha] = materia.cuposDisponibles;
    materia.grupos.forEach(grupoObj => {

        if (!refMateria.grupos.hasOwnProperty(grupoObj.grupo)) {
            refMateria.grupos[grupoObj.grupo] = {
                profesor: grupoObj.profesor,
                cupos: {}
            };
        }

        const refGrupo = refMateria.grupos[grupoObj.grupo];

        refGrupo.cupos[fecha] = grupoObj.cupos;
    });

    const recomendaciones = getRecomendaciones(refMateria);
    refMateria.recomendaciones = recomendaciones;
}

function getRecomendaciones(asignatura) {
    const grupos = Object.values(asignatura.grupos);

    // Agrupar por docente
    const agrupado = {};

    for (const grupo of grupos) {
        const { profesor, cupos } = grupo;
        if (!agrupado.hasOwnProperty(profesor)) {
            agrupado[profesor] = [];
        }

        const cuposArray = Object.values(cupos);
        agrupado[profesor].push(cuposArray);
    }

    const finalData = [];
    for (const [docente, data] of Object.entries(agrupado)) {
        // Sumar todos los cupos por index
        const cuposTotales = data.reduce((acc, curr) => {
            return acc.map((val, i) => val + curr[i]);
        });

        let totalInscritos = 0;
        const cambios = [];
        if (cuposTotales.length === 1) {
            cambios.push(cuposTotales[0]);
            totalInscritos = cuposTotales[0];
        } else {
            for (let i = 0; i < (cuposTotales.length - 1); i++) {
                const current = cuposTotales[i];
                const next = cuposTotales[i + 1];
                const diferencia = current - next
                totalInscritos += diferencia;
                cambios.push(diferencia / data.length);
                if (next === 0) break;
            }
        }

        const inscritos = cambios.reduce((acc, curr) => acc + curr);

        finalData.push({
            docente,
            inscritos: totalInscritos,
            puntaje: inscritos / cambios.length,
        })
    }

    // Normalizar puntaje
    const max = Math.max(...finalData.map(({ puntaje }) => puntaje));
    finalData.forEach((docente) => {
        docente.puntaje = +(Math.round((docente.puntaje * 10) / max + "e+1") + "e-1");
    });

    const ordenado = finalData.sort((a, b) => b.puntaje - a.puntaje);

    return ordenado;
}