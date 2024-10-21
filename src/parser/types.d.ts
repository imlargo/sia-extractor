export type Prerequisito = {
    tipo: string
    isTodas: boolean
    cantidad: number
    asignaturas: Record<string, string>[]
}

export type PrerequisitoAsignatura = {
    codigo: string;
    nombre: string;
}

export type Horario = {
    inicio: string;
    fin: string;
    dia: string;
}

export type Grupo = {
    grupo: string;
    cupos: number;
    profesor: string;
    duracion: string;
    jornada: string;
    horarios: Horario[];
}

export type Asignatura = {
    nombre: string;
    codigo: string;
    tipologia: string;
    creditos: number;
    facultad: string;
    carrera: string;
    fechaExtraccion: string;
    cuposDisponibles: number;
    prerequisitos: Prerequisito[];
    grupos: Grupo[];
}