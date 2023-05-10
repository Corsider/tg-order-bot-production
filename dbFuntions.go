package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func First[T, U any](val T, _ U) T {
	return val
}

func CreateUser(tgUserID string) {
	// Create user if only it doesn't exist
	counter := 0
	q := fmt.Sprintf("SELECT count(*) FROM users WHERE tguser_id='%s'", tgUserID)
	//log.Println(counter)
	_ = DB.QueryRow(q).Scan(&counter)

	if counter == 0 {
		query := fmt.Sprintf("INSERT INTO users (cart, order_history_id, tguser_id, address, mark) VALUES ('%s', '%s', '%s', '%s', '%s')", "{}", "{}", tgUserID, "", "0")
		_, err := DB.Exec(query)
		if err != nil {
			log.Println(err)
		}
	}
}

func AddToUserCart(id int64, newItem int) {
	query := fmt.Sprintf("SELECT cart FROM users WHERE tguser_id='%s'", strconv.FormatInt(id, 10))
	row := DB.QueryRow(query)
	oldCart := ""
	row.Scan(&oldCart)
	cartTrim, _ := strings.CutSuffix(oldCart, "}")
	if oldCart == "{}" {
		cartTrim += strconv.Itoa(newItem) + "}"
	} else {
		cartTrim += "," + strconv.Itoa(newItem) + "}"
	}
	query2 := fmt.Sprintf("UPDATE users SET cart='%s' WHERE tguser_id='%s'", cartTrim, strconv.FormatInt(id, 10))
	_, err := DB.Exec(query2)
	if err != nil {
		//log.Println("error in AddToUserCart")
		log.Println(err)
	}
}

func CartContains(id int64, itemID int) bool {
	cart := GetUserCart(id)
	for _, el := range cart {
		if el == itemID {
			return true
		}
	}
	return false
}

func RemoveFromUserCart(id int64, itemID int) {
	query := fmt.Sprintf("SELECT cart FROM users WHERE tguser_id='%s'", strconv.FormatInt(id, 10))
	row := DB.QueryRow(query)
	oldCart := ""
	row.Scan(&oldCart)

	cartTrim := strings.Replace(oldCart, strconv.Itoa(itemID), "", 1) //{1,2,3} ->(2)-> {1,,3}; {1,2,3} ->(3)-> {1,2,} OR {,2,3}
	cartTrim1 := strings.Replace(cartTrim, "{,", "{", 1)
	cartTrim2 := strings.Replace(cartTrim1, ",}", "}", 1)
	cartTrim3 := strings.Replace(cartTrim2, ",,", ",", 1)
	q := fmt.Sprintf("UPDATE users SET cart='%s' WHERE tguser_id='%s'", cartTrim3, strconv.FormatInt(id, 10))
	_, err := DB.Exec(q)
	if err != nil {
		log.Println(err)
	}
}

func GetUserCart(id int64) []int {
	query := fmt.Sprintf("SELECT cart FROM users WHERE tguser_id='%s'", strconv.FormatInt(id, 10))
	row := DB.QueryRow(query)
	itemsString := ""

	err := row.Scan(&itemsString)
	if err != nil {
		log.Println(err)
	}

	if len(itemsString) == 2 {
		return nil
	}
	ns, _ := strings.CutPrefix(itemsString, "{")
	nss, _ := strings.CutSuffix(ns, "}")
	items := strings.Split(nss, ",")
	out := []int{}
	for _, v := range items {
		out = append(out, First(strconv.Atoi(v)))
	}
	return out
}

// AddToUserOrderHistory When placing order, this function adds current order to order_history field and removes current cart.
func AddToUserOrderHistory(id int64) {
	query := fmt.Sprintf("SELECT cart, order_history_id FROM users WHERE tguser_id='%s'", strconv.FormatInt(id, 10))
	row := DB.QueryRow(query)
	currentCart := ""
	currentOrderHistory := ""
	row.Scan(&currentCart, &currentOrderHistory)

	OrderHistoryIDToUser := ""
	query1 := fmt.Sprintf("INSERT INTO orders (order_list) VALUES ('%s') RETURNING order_id", currentCart)
	err := DB.QueryRow(query1).Scan(&OrderHistoryIDToUser)
	if err != nil {
		log.Println(err)
	}
	nc := ""
	if currentOrderHistory == "{}" {
		nc = "{" + OrderHistoryIDToUser + "}"
	} else {
		nc, _ = strings.CutSuffix(currentOrderHistory, "}")
		nc += "," + OrderHistoryIDToUser + "}"
	}
	_, _ = DB.Exec("UPDATE users SET cart='{}', order_history_id=$1 WHERE tguser_id=$2", nc, strconv.FormatInt(id, 10))
}

// format: "0,0,1;1,1;3,2" ,-order items separator, ;-orders separator
func GetUserOrderHistory(id int64) string {
	query := fmt.Sprintf("SELECT order_history_id FROM users WHERE tguser_id='%s'", strconv.FormatInt(id, 10))
	row := DB.QueryRow(query)
	orderHist := ""
	row.Scan(&orderHist)

	if orderHist == "{}" {
		return ""
	} else {
		ors, _ := strings.CutPrefix(orderHist, "{")
		ors2, _ := strings.CutSuffix(ors, "}")
		ordersIDs := strings.Split(ors2, ",")
		out := ""
		for _, el := range ordersIDs {
			q := fmt.Sprintf("SELECT order_list FROM orders WHERE order_id='%s'", el)
			roww := DB.QueryRow(q)
			str := ""
			roww.Scan(&str)
			str1, _ := strings.CutPrefix(str, "{")
			str2, _ := strings.CutSuffix(str1, "}")
			out += str2 + ";"
		}
		out1, _ := strings.CutSuffix(out, ";")
		return out1
	}
}

func GetUserAddress(id int64) string {
	query := fmt.Sprintf("SELECT address FROM users WHERE tguser_id='%s'", strconv.FormatInt(id, 10))
	row := DB.QueryRow(query)
	var address string
	err := row.Scan(&address)
	if err != nil {
		log.Println(err)
	}
	return address
}

func SetUserAddress(id int64, adr string) {
	query := fmt.Sprintf("UPDATE users SET address='%s' WHERE tguser_id='%s'", adr, strconv.FormatInt(id, 10))
	_, err := DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func GetUserMark(id int64) int {
	query := fmt.Sprintf("SELECT mark FROM users WHERE tguser_id='%s'", strconv.FormatInt(id, 10))
	row := DB.QueryRow(query)
	if row == nil {
		return 0
	}
	var m string
	row.Scan(&m)
	return First(strconv.Atoi(m))
}

func UpdateUserMark(id int64, mark int) {
	query := fmt.Sprintf("UPDATE users SET mark='%s' WHERE tguser_id='%s'", strconv.Itoa(mark), strconv.FormatInt(id, 10))
	_, err := DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func SetLastMessageID(id int64, messageID int) {
	query := fmt.Sprintf("UPDATE users SET last_inserted_id='%s' WHERE tguser_id='%s'", strconv.Itoa(messageID), strconv.FormatInt(id, 10))
	_, err := DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func GetLastMessageID(id int64) int {
	query := fmt.Sprintf("SELECT last_inserted_id FROM users WHERE tguser_id='%s'", strconv.FormatInt(id, 10))
	row := DB.QueryRow(query)
	var out int
	row.Scan(&out)
	return out
}

func ClearOrderHistory(id int64) {
	query := fmt.Sprintf("SELECT order_history_id FROM users WHERE tguser_id='%s'", strconv.FormatInt(id, 10))
	row := DB.QueryRow(query)
	var orders string
	row.Scan(&orders)
	or1, _ := strings.CutSuffix(orders, "}")
	or2, _ := strings.CutPrefix(or1, "{")
	ord := strings.Split(or2, ",")
	for _, el := range ord {
		q := fmt.Sprintf("DELETE FROM orders WHERE order_id='%s'", el)
		_, err := DB.Exec(q)
		if err != nil {
			log.Println(err)
		}
	}
	qu := fmt.Sprintf("UPDATE users SET order_history_id='{}' WHERE tguser_id='%s'", strconv.FormatInt(id, 10))
	_, err := DB.Exec(qu)
	if err != nil {
		log.Println(err)
	}
}
