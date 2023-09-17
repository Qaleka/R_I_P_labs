package models

type FIO struct {
	First_name  string `json:"first_name"`
	Second_name string `json:"second_name"`
	Third_name  string `json:"third_name"`
}

type Recipients struct {
	Id       int    `json:"id"`
	Name     FIO    `json:"name"`
	ImageSrc string `json:"image_src"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Adress   string `json:"adress"`
}
