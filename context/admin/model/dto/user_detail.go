package admin_dto

import card_model "com.xalixcolabs.trading-card-album/context/card/model"

type UserDetail struct {
	User   User            `json:"user"`
	Albums []Album         `json:"albums"`
	Cards  []card_model.Card `json:"cards"`
}