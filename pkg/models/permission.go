package models

type Permission struct {
	Admin bool `bson:"admin"`
	Write bool `bson:"write"`
	Read  bool `bson:"read"`
}

func (p *Permission) SetDefaultPermissions() {
	p.Admin = false
	p.Write = false
	p.Read = true
}
