package domain

import "brickrecon/lego"

type Stock = map[lego.PartId]map[lego.ColorId]int

func GetStock(stock Stock, part lego.PartId, color lego.ColorId) int {
	colorStock, found := stock[part]
	if !found {
		return 0
	}

	return colorStock[color]
}

func AddStock(stock Stock, part lego.PartId, color lego.ColorId, quantity int) {
	if quantity == 0 {
		return
	}

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
