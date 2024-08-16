package models

type Permission struct {
	Admin bool `json:"admin" bson:"admin"`
	Write bool `json:"write" bson:"write"`
	Read  bool `json:"read" bson:"read"`
}

func (p *Permission) SetDefaultPermissions() {
	p.Admin = false
	p.Write = false
	p.Read = true
}
