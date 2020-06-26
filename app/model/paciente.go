package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Paciente struct {
	ID             primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Nombre         string                 `json:"nombre,omitempty" bson:"nombre,omitempty"`
	Apellido       string                 `json:"apellido,omitempty" bson:"apellido,omitempty"`
	Identificacion string                 `json:"identificacion,omitempty" bson:"identificacion,omitempty"`
	Prestador      string                 `json:"prestador,omitempty" bson:"prestador,omitempty"`
	Afiliado       string                 `json:"afiliado,omitempty" bson:"afiliado,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // data opcional cualquier key:value

}

func NewPaciente(nombre, apellido, prestador, afiliado string, data map[string]interface{}) *Paciente {
	return &Paciente{
		Nombre:    nombre,
		Apellido:  apellido,
		Prestador: prestador,
		Afiliado:  afiliado,
		Data:      data,
	}
}
