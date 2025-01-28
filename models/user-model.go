package models

import (
	"time"
	"todo-manager/utils"
)

type UserModel struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

func (this *UserModel) Equals(other UserModel) bool {
	createdAtMatch := utils.TimesMatch(this.CreatedAt, other.CreatedAt)
	updatedAtMatch := utils.TimesMatch(this.UpdatedAt, other.UpdatedAt)

	return this.Id == other.Id &&
		this.Name == other.Name &&
		this.Email == other.Email &&
		createdAtMatch &&
		updatedAtMatch
}
