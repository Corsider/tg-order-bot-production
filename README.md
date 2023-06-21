# tg-order-bot

## Общее
Проект был создан для отбора на стажировку VK. \
Бот служит для отправки заказов в ресторан. Пользователь может добавлять товары в корзину и удалять их оттуда, а так же устанавливать свой адрес для доставки заказа. Есть возможность оставить оценку работы бота.

## Технические детали
Бот написан польностью на Go 1.20.3. Сторонние библиотеки использованы не были, кроме lib/pq для связи с базой данных. База данных проекта - Postgresql 15.
Когда пользователь начинает общение с ботом (команда /start, которая запускается автоматически) бот проверяет, есть ли пользователь в базе данных (по UserID). Если пользователя нет, то бот добавляет его в БД.

## Структура БД
## Таблица users:
| Столбец  | Тип данных |
|---------:|-----------|
|user_id|integer (PK)|
|cart|intger[]|
|order_history_id| integer[]|
|tguser_id|character varying|
|address|character varying|
|last_inserted_id|integer|
|mark|integer|

## Таблица orders:
|Столбец|Тип данных|
|-----:|-------|
|order_id|integer (PK)|
|order_list|integer[]|

## 
Таблица users хранит пользователей. cart - массив текущих элементов корзины. Каждый продукт (всего их 4) имеет свой ID.\
Пицца Пеперони - 0, Пицца 4 сыра - 1 и т.д.\
order_history_id - массив ID строк таблицы orders. Здесь хранится история заказов пользователя (в order_list таблицы orders хранится конкретный заказ пользователя).\
tguser_id - ID пользователя телеграмма.\
address - адрес. last_inserted_id - ID последнего сообщения бота, содержавшего inline кнопки (по команде /my), нужно для удаления сообщения от бота при повторной отправке
команды /my.\
mark - оценка пользователя. 0 - нет оценки, 1 - положительная; -1 - отрицательная оценки.

## Env
Переменные окружения, которые использует бот:\
TG_BOT_TOKEN - токен бота\
DBHOST - хост БД\
POSTGRES_PASS - пароль от postgres пользователя

## Дальнейшие улучшения
Хранить корзину пользователя используя redis.
