package domain

type MealCard string

const (
	MealCardTicket   MealCard = "ticket"
	MealCardMultinet MealCard = "multinet"
	MealCardSodexo   MealCard = "sodexo"
	MealCardSetcard  MealCard = "setcard"
	MealCardMetropol MealCard = "metropol"
	MealCardPayeKart MealCard = "payekart"
)

func (m MealCard) IsValid() bool {
	switch m {
	case MealCardTicket, MealCardMultinet, MealCardSodexo, MealCardSetcard, MealCardMetropol, MealCardPayeKart:
		return true
	default:
		return false
	}
}
