package domain

import "brickrecon/lego"

type Stock = map[lego.PartId]map[lego.ColorId]int

func AddStock(stock Stock, part lego.PartId, color lego.ColorId, quantity int) {
	if _, found := stock[part]; !found {
		stock[part] = map[lego.ColorId]int{}
	}

	stock[part][color] = stock[part][color] + quantity
}

func RemoveStock(stock Stock, part lego.PartId, color lego.ColorId, quantity int) {
	if _, found := stock[part]; !found {
		return
	}

	stock[part][color] = stock[part][color] - quantity

	if stock[part][color] <= 0 {
		delete(stock[part], color)
	}

	if len(stock[part]) == 0 {
		delete(stock, part)
	}
}
