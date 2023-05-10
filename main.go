package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"strings"
)

/*

	When user starts the bot, it automatically adds user's data to database.

*/

// DB Connection...
func ConnectDB() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=5432 user=postgres password=%s dbname=tgbotdb sslmode=disable", os.Getenv("DBHOST"), os.Getenv("POSTGRES_PASS"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var DB *sql.DB

// KEYBOARDS...
var keyboardMain = CreateInlineKeyboardMarkup(
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Собрать заказ", "dat1"),
		CreateInlineKeyboardButton("Мои заказы", "dat2"),
	),
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Установить адрес", "dat3"),
		CreateInlineKeyboardButton("Обратная связь", "dat4"),
	),
)

var keyboardOrder = CreateInlineKeyboardMarkup(
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Пицца пеперони", "p1"),
		CreateInlineKeyboardButton("Пицца 4 сыра", "p2"),
	),
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Добрый Кола 0.5л", "p3"),
		CreateInlineKeyboardButton("Киндер-Сюрприз", "p4"),
	),
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Заказать", "p5"),
	),
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Назад", "toMain"),
		CreateInlineKeyboardButton("Корзина", "cart"),
	),
)

var keyboardRate = CreateInlineKeyboardMarkup(
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("👍", "good"),
		CreateInlineKeyboardButton("👎", "bad"),
	),
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Назад", "toMain"),
	),
)

var keyboardCart = CreateInlineKeyboardMarkup(
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Назад", "toMainOrder"),
		CreateInlineKeyboardButton("Удалить элемент из корзины", "deleteItem"),
	),
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Заказать", "p5"), //dat1
	),
)

var keyboardMyOrders = CreateInlineKeyboardMarkup(
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Назад", "toMain"),
		CreateInlineKeyboardButton("Новый заказ", "dat1"),
	),
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Очистить историю заказов", "clearHist"),
	),
)

var keyboardAddress = CreateInlineKeyboardMarkup(
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Назад", "toMain"),
		CreateInlineKeyboardButton("Установить адрес", "setAddress"),
	),
)

var keyboardDelete = CreateInlineKeyboardMarkup(
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Пицца пеперони", "dp1"),
		CreateInlineKeyboardButton("Пицца 4 сыра", "dp2"),
	),
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Добрый Кола 0.5л", "dp3"),
		CreateInlineKeyboardButton("Киндер-Сюрприз", "dp4"),
	),
	CreateInlineKeyboardRow(
		CreateInlineKeyboardButton("Назад", "toMainOrder"),
	),
)

func main() {
	bot, err := CreateBot(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("SERVER STARTED")

	DB = ConnectDB()
	log.Println("DB CONNECTION ESTABLISHED")

	defer DB.Close()

	waitingInput := false
	FoodArray := []string{"Пицца Пеперони", "Пицца 4 сыра", "Добрый Кола 0.5л", "Киндер-Сюрприз"}

	updateParameters := UpdateParams{
		Limit:   0,
		Offset:  0,
		Timeout: 30,
	}

	updates := bot.UpdatesChannel(updateParameters)
	for update := range updates {
		if update.MyChatMember != nil { // simple defense from banning bot
			continue
		}
		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "dat1":
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Итак, выбери, что ты закажешь:", &keyboardOrder)
			case "dat2":
				orderStr := GetUserOrderHistory(update.CallbackQuery.Message.Chat.Id)
				orders := strings.Split(orderStr, ";")
				if orders[0] == "" {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Заказов пока нет!\n\n", &keyboardMyOrders)
					continue
				}
				out := ""
				for i, order := range orders {
					subout := strconv.Itoa(i+1) + "-й заказ:\n"
					food := strings.Split(order, ",")
					foodmap := make(map[int]int)
					foodmap[0] = 0
					foodmap[1] = 0
					foodmap[2] = 0
					foodmap[3] = 0
					for k := 0; k <= 4; k++ {
						for _, v := range food {
							if First(strconv.Atoi(v)) == k {
								foodmap[k] += 1
							}
						}
					}
					count := 1
					for k := 0; k <= 4; k++ {
						if foodmap[k] != 0 {
							if foodmap[k] > 1 {
								subout += "\t" + strconv.Itoa(count) + ") " + FoodArray[k] + " x" + strconv.Itoa(foodmap[k]) + "\n"
							} else {
								subout += "\t" + strconv.Itoa(count) + ") " + FoodArray[k] + "\n"
							}
							count += 1
						}
					}
					out += subout + "\n"
				}
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Твои заказы:\n\n"+out, &keyboardMyOrders)
			case "cart":
				outstr := "Вот что сейчас в твоей корзине:\n\n"
				items := GetUserCart(update.CallbackQuery.From.ID)
				if items == nil {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Корзина пуста! Выбери что-нибудь:", &keyboardOrder)
				} else {
					str := ""
					foodmap := make(map[int]int)
					foodmap[0] = 0
					foodmap[1] = 0
					foodmap[2] = 0
					foodmap[3] = 0
					for i := 0; i <= 4; i++ {
						for _, v := range items {
							if v == i {
								foodmap[i] += 1
							}
						}
					}
					count := 0
					for i := 0; i <= 4; i++ {
						if foodmap[i] != 0 {
							str += strconv.Itoa(count+1) + ") " + FoodArray[i] + " - x" + strconv.Itoa(foodmap[i]) + "\n"
							count += 1
						}
					}
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, outstr+str, &keyboardCart)
				}
			case "deleteItem":
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Выбери, что ты хочешь удалить из корзины?", &keyboardDelete)
			case "dp1":
				cart := GetUserCart(update.CallbackQuery.From.ID)
				if cart == nil {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Корзина пуста!", &keyboardOrder)
				} else {
					if CartContains(update.CallbackQuery.From.ID, 0) {
						RemoveFromUserCart(update.CallbackQuery.From.ID, 0)
						_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Удалено.", &keyboardOrder)
						continue
					}
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Этого продукта не было в корзине!", &keyboardOrder)
				}
			case "dp2":
				cart := GetUserCart(update.CallbackQuery.From.ID)
				if cart == nil {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Корзина пуста!", &keyboardOrder)
				} else {
					if CartContains(update.CallbackQuery.From.ID, 1) {
						RemoveFromUserCart(update.CallbackQuery.From.ID, 1)
						_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Удалено.", &keyboardOrder)
						continue
					}
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Этого продукта не было в корзине!", &keyboardOrder)
				}
			case "dp3":
				cart := GetUserCart(update.CallbackQuery.From.ID)
				if cart == nil {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Корзина пуста!", &keyboardOrder)
				} else {
					if CartContains(update.CallbackQuery.From.ID, 2) {
						RemoveFromUserCart(update.CallbackQuery.From.ID, 2)
						_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Удалено.", &keyboardOrder)
						continue
					}
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Этого продукта не было в корзине!", &keyboardOrder)
				}
			case "dp4":
				cart := GetUserCart(update.CallbackQuery.From.ID)
				if cart == nil {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Корзина пуста!", &keyboardOrder)
				} else {
					if CartContains(update.CallbackQuery.From.ID, 3) {
						RemoveFromUserCart(update.CallbackQuery.From.ID, 3)
						_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Удалено.", &keyboardOrder)
						continue
					}
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Этого продукта не было в корзине!", &keyboardOrder)
				}
			case "clearHist":
				if GetUserOrderHistory(update.CallbackQuery.From.ID) == "" {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Упс! История заказов и так пуста!", &keyboardMain)
				} else {
					ClearOrderHistory(update.CallbackQuery.From.ID)
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Удалено.", &keyboardMain)
				}
			case "dat3":
				address := GetUserAddress(update.CallbackQuery.From.ID)
				if address == "" {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Ты еще не устанавливал свой адрес. Пожалуйста, выбери одно из действий ниже:", &keyboardAddress)
				} else {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Твой адресс: "+address, &keyboardAddress)
				}
			case "setAddress":
				_, _ = bot.sendMessage(update.CallbackQuery.Message.Chat.Id, "Пришли мне свой адрес:", nil)

				waitingInput = true
			case "dat4":
				if GetUserMark(update.CallbackQuery.From.ID) == 0 {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Оцените работу бота:", &keyboardRate)
				} else {
					mark := GetUserMark(update.CallbackQuery.From.ID)
					markstr := ""
					if mark == 1 {
						markstr = "👍"
					} else {
						markstr = "👎"
					}
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Ты уже оценивал наш сервис как "+markstr+". Хочешь изменить оценку?", &keyboardRate)
				}

			case "good":
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Спасибо за оценку!", &keyboardMain)
				UpdateUserMark(update.CallbackQuery.From.ID, 1)
			case "bad":
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Мы будем стараться лучше, спасибо!", &keyboardMain)
				UpdateUserMark(update.CallbackQuery.From.ID, -1)
			case "toMain":
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Выбери, что ты хочешь сделать:", &keyboardMain)
			case "toMainOrder":
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Выбери, что ты закажешь:", &keyboardOrder)
			case "p1":
				//adding food
				AddToUserCart(update.CallbackQuery.From.ID, PizzaPepperony)
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Добавлено! Что-то еще?", &keyboardOrder)
			case "p2":
				AddToUserCart(update.CallbackQuery.From.ID, PizzaFourCheese)
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Добавлено! Что-то еще?", &keyboardOrder)
			case "p3":
				AddToUserCart(update.CallbackQuery.From.ID, Cola)
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Добавлено! Что-то еще?", &keyboardOrder)
			case "p4":
				AddToUserCart(update.CallbackQuery.From.ID, KinderSurprise)
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Добавлено! Что-то еще?", &keyboardOrder)
			case "p5":
				if GetUserAddress(update.CallbackQuery.From.ID) == "" {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Упс! Кажется ты забыл указать нам свой адрес! Ты можешь сделать это в меню ниже:", &keyboardMain)
					continue
				}
				items := GetUserCart(update.CallbackQuery.From.ID)
				if items == nil {
					_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Упс! Твоя корзина пуста. Выбери что-нибудь:", &keyboardOrder)
					continue
				}
				str := ""
				foodmap := make(map[int]int)
				foodmap[0] = 0
				foodmap[1] = 0
				foodmap[2] = 0
				foodmap[3] = 0
				for i := 0; i <= 4; i++ {
					for _, v := range items {
						if v == i {
							foodmap[i] += 1
						}
					}
				}
				count := 0
				for i := 0; i <= 4; i++ {
					if foodmap[i] != 0 {
						str += strconv.Itoa(count+1) + ") " + FoodArray[i] + " - x" + strconv.Itoa(foodmap[i]) + "\n"
						count += 1
					}
				}
				AddToUserOrderHistory(update.CallbackQuery.From.ID)
				_, _ = bot.sendEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageID, "Заказ принят! Мы уже начали готовить его:\n\n"+str, &keyboardMain)

			default:
			}
			_, _ = bot.sendAnswerCallbackQuery(update.CallbackQuery.ID)
			continue
		}
		if update.Message.CheckIfCommand() {
			waitingInput = false
			switch update.Message.Text {
			case "/start":
				_, err = bot.sendMessage(update.Message.Chat.Id, "Привет! Здесь ты можешь сделать заказ для нашего ресторана, а мы доставим тебе его в течение часа!\nДля этого используй команду /my", nil)
				if err != nil {
					log.Println(err)
				}
				// Add user to DB when /start engaged
				CreateUser(strconv.FormatInt(update.Message.From.ID, 10))
			case "/my":
				_, _ = bot.sendDeleteMessage(update.Message.Chat.Id, GetLastMessageID(update.Message.From.ID))
				msg, errr := bot.sendMessage(update.Message.Chat.Id, "Выбери, что ты хочешь сделать:", &keyboardMain)
				SetLastMessageID(update.Message.From.ID, msg.MessageID)
				if errr != nil {
					log.Println(err)
				}
			case "/help":
				_, err = bot.sendMessage(update.Message.Chat.Id, "Используй кнопку Меню слева чтобы выбрать доступную команду.\nЧтобы задать вопрос, можешь писать в телеграмм @Corsider", nil)
				if err != nil {
					log.Println(err)
				}
			default:
				_, err = bot.sendMessage(update.Message.Chat.Id, "Нет такой команды!", nil)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			if !waitingInput {
				_, err = bot.sendMessage(update.Message.Chat.Id, "Нет такой команды!", nil)
				if err != nil {
					log.Println(err)
				}
			} else {
				address := update.Message.Text
				t := "Ок, вот твой адрес: " + address
				_, err = bot.sendMessage(update.Message.Chat.Id, t, &keyboardMain)
				SetUserAddress(update.Message.From.ID, address)
				if err != nil {
					log.Println(err)
				}
				waitingInput = false
			}
		}
	}
}
