package models

type Menu struct {
	Date string
	DailyDish Dish
	SteakOfTheWeek Dish
	BurgerOfTheWeek Dish
}
type Dish struct {
	Content string
	Price string
}
