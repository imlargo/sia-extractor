package core

type Prerequisito struct {
	Tipo        string              `json:"tipo"`
	IsTodas     bool                `json:"isTodas"`
	Cantidad    int                 `json:"cantidad"`
	Asignaturas []map[string]string `json:"asignaturas"`
}

type Horario struct {
	Inicio string `json:"inicio"`
	Fin    string `json:"fin"`
	Dia    string `json:"dia"`
}

type Grupo struct {
	Grupo    string    `json:"grupo"`
	Cupos    int       `json:"cupos"`
	Profesor string    `json:"profesor"`
	Duracion string    `json:"duracion"`
	Jornada  string    `json:"jornada"`
	Horarios []Horario `json:"horarios"`
}

type Asignatura struct {
	Nombre           string         `json:"nombre" bson:"nombre"`
	Codigo           string         `json:"codigo" bson:"codigo"`
	Tipologia        string         `json:"tipologia" bson:"tipologia"`
	Creditos         int            `json:"creditos" bson:"creditos"`
	Facultad         string         `json:"facultad" bson:"facultad"`
	Carrera          string         `json:"carrera" bson:"carrera"`
	FechaExtraccion  string         `json:"fechaExtraccion" bson:"fechaExtraccion"`
	CuposDisponibles int            `json:"cuposDisponibles" bson:"cuposDisponibles"`
	Prerequisitos    []Prerequisito `json:"prerequisitos" bson:"prerequisitos"`
	Grupos           []Grupo        `json:"grupos" bson:"grupos"`
}

type Codigo struct {
	Nivel     string
	Sede      string
	Facultad  string
	Carrera   string
	Tipologia string
}

type Path struct {
	Nivel     string
	Sede      string
	Facultad  string
	Carrera   string
	Tipologia string
}

type PathTipologia struct {
	Por         string
	SedePor     string
	FacultadPor string
	CarreraPor  string
}
