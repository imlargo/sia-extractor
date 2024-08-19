# SIA Extractor

Este repositorio contiene un conjunto de herramientas para extraer, procesar y desplegar datos del Sistema de Información Académica (SIA) de la Universidad Nacional de Colombia para alimentar los datos de mi aplicacion web de horarios **"Pegaso"**.

Ejemplo de datos extraidos de una asignatura

```json
{
            "nombre": "Teoría de lenguajes de programación",
            "codigo": "3010426",
            "tipologia": "DISCIPLINAR OBLIGATORIA",
            "creditos": "3",
            "facultad": "3068 FACULTAD DE MINAS",
            "carrera": "3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA",
            "fechaExtraccion": "19/8/2024 - 10:48 a. m.",
            "cuposDisponibles": "17",
            "prerequisitos": [
                {
                    "tipo": "M",
                    "isTodas": true,
                    "cantidad": 2,
                    "asignaturas": [
                        {
                            "codigo": "3010435",
                            "nombre": "Fundamentos de programación"
                        },
                        {
                            "codigo": "3006906",
                            "nombre": "MATEMÁTICAS DISCRETAS"
                        }
                    ]
                }
            ],
            "grupos": [
                {
                    "grupo": "Grupo 1",
                    "cupos": 10,
                    "profesor": "Demetrio ...",
                    "duracion": "Semestral",
                    "jornada": "DIURNO",
                    "horarios": [
                        {
                            "inicio": "10:00",
                            "fin": "12:00",
                            "dia": "MIÉRCOLES"
                        },
                        {
                            "inicio": "08:00",
                            "fin": "10:00",
                            "dia": "VIERNES"
                        }
                    ]
                },
                {
                    "grupo": "Grupo 2",
                    "cupos": 6,
                    "profesor": "Demetrio ...",
                    "duracion": "Semestral",
                    "jornada": "DIURNO",
                    "horarios": [
                        {
                            "inicio": "10:00",
                            "fin": "12:00",
                            "dia": "MIÉRCOLES"
                        },
                        {
                            "inicio": "10:00",
                            "fin": "12:00",
                            "dia": "VIERNES"
                        }
                    ]
                },
                {
                    "grupo": "Grupo 3 REMOTA",
                    "cupos": 1,
                    "profesor": "Oscar ...",
                    "duracion": "Semestral",
                    "jornada": "DIURNO",
                    "horarios": [
                        {
                            "inicio": "08:00",
                            "fin": "10:00",
                            "dia": "MARTES"
                        },
                        {
                            "inicio": "08:00",
                            "fin": "10:00",
                            "dia": "JUEVES"
                        }
                    ]
                }
            ]
        }
```
