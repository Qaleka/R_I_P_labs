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

func GetCardsInfo() []Recipients {
	return []Recipients{
		{
			Id: 0,
			Name: FIO{
				First_name:  "Олег",
				Second_name: "Орлов",
				Third_name:  "Никитович",
			},
			ImageSrc: "http://localhost:8080/image/men1.jpg",
			Email:    "OlegO@mail.ru",
			Age:      27,
			Adress:   "Москва, ул. Измайловская, д.13, кв.54",
		},
		{
			Id: 1,
			Name: FIO{
				First_name:  "Василий",
				Second_name: "Гречко",
				Third_name:  "Валентинович",
			},
			ImageSrc: "http://localhost:8080/image/men2.jpg",
			Email:    "Grechko_101@mail.ru",
			Age:      31,
			Adress:   "Москва, ул. Тверская, д.25, кв.145",
		},
		{
			Id: 2,
			Name: FIO{
				First_name:  "Александр",
				Second_name: "Лейко",
				Third_name:  "Кириллович",
			},
			ImageSrc: "http://localhost:8080/image/men3.jpg",
			Email:    "Alek221@mail.ru",
			Age:      37,
			Adress:   "Москва, ул. Изюмская, д.15, кв.89",
		},
	}
}
