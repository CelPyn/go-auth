package domain

type User struct {
	Name string `bson:"name"`
	Password string `bson:"password"`
	Grants []string `bson:"grants"`
}
