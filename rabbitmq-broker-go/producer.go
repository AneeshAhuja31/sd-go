package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RestaurantOrder struct {
	Id int `json:"id"`
	Items []string  `json:"items"`
}


func publish_message(ch *amqp.Channel,q *amqp.Queue,json_body []byte,ctx context.Context){
	err := ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: json_body,
		},
	)
	errHelper(err)
}

func initQueue(ch *amqp.Channel)amqp.Queue{
	q,err := ch.QueueDeclare(
		"Orders",
		false,
		false,
		false,
		false,
		nil,
	)
	errHelper(err)
	return q
}

func producer(ch *amqp.Channel){
	q := initQueue(ch)
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	slice_of_order_structures := []RestaurantOrder{
		{1, []string{"Chicken", "Roti"}},
		{2, []string{"Pasta", "Lasagna"}},
		{3, []string{"Rice", "Curry"}},
		{4, []string{"Burger", "Fries"}},
		{5, []string{"Pizza", "Garlic Bread"}},
		{6, []string{"Biryani", "Raita"}},
		{7, []string{"Noodles", "Manchurian"}},
		{8, []string{"Sandwich", "Chips"}},
		{9, []string{"Tacos", "Salsa"}},
		{10, []string{"Paneer", "Naan"}},
		{11, []string{"Fish Fry", "Rice"}},
		{12, []string{"Soup", "Bread"}},
		{13, []string{"Salad", "Juice"}},
		{14, []string{"Idli", "Sambar"}},
		{15, []string{"Dosa", "Chutney"}},
		{16, []string{"Upma", "Tea"}},
		{17, []string{"Poha", "Coffee"}},
		{18, []string{"Paratha", "Curd"}},
		{19, []string{"Rajma", "Rice"}},
		{20, []string{"Chole", "Bhature"}},
		{21, []string{"Kebab", "Mint Chutney"}},
		{22, []string{"Pulao", "Salad"}},
		{23, []string{"Fried Rice", "Gobi"}},
		{24, []string{"Spring Roll", "Dip"}},
		{25, []string{"Momos", "Sauce"}},
		{26, []string{"Steak", "Mashed Potatoes"}},
		{27, []string{"Grilled Chicken", "Veggies"}},
		{28, []string{"Hot Dog", "Mustard"}},
		{29, []string{"Wrap", "Coleslaw"}},
		{30, []string{"Shawarma", "Hummus"}},
		{31, []string{"Falafel", "Pita"}},
		{32, []string{"Sushi", "Wasabi"}},
		{33, []string{"Ramen", "Egg"}},
		{34, []string{"Udon", "Tempura"}},
		{35, []string{"Pad Thai", "Peanuts"}},
		{36, []string{"Curry", "Bread"}},
		{37, []string{"Dal", "Rice"}},
		{38, []string{"Khichdi", "Ghee"}},
		{39, []string{"Pancakes", "Syrup"}},
		{40, []string{"Waffles", "Honey"}},
		{41, []string{"Omelette", "Toast"}},
		{42, []string{"Scrambled Eggs", "Bacon"}},
		{43, []string{"French Toast", "Butter"}},
		{44, []string{"Croissant", "Coffee"}},
		{45, []string{"Bagel", "Cream Cheese"}},
		{46, []string{"Donut", "Milk"}},
		{47, []string{"Cake", "Ice Cream"}},
		{48, []string{"Brownie", "Fudge"}},
		{49, []string{"Cupcake", "Icing"}},
		{50, []string{"Pie", "Whipped Cream"}},
		{51, []string{"Apple Tart", "Custard"}},
		{52, []string{"Cheesecake", "Berry Sauce"}},
		{53, []string{"Pudding", "Caramel"}},
		{54, []string{"Ice Cream", "Sprinkles"}},
		{55, []string{"Milkshake", "Cookies"}},
		{56, []string{"Smoothie", "Granola"}},
		{57, []string{"Fruit Bowl", "Yogurt"}},
		{58, []string{"Nachos", "Cheese Dip"}},
		{59, []string{"Quesadilla", "Sour Cream"}},
		{60, []string{"Enchiladas", "Beans"}},
		{61, []string{"Burrito", "Guacamole"}},
		{62, []string{"Chips", "Salsa"}},
		{63, []string{"Popcorn", "Butter"}},
		{64, []string{"Pretzel", "Cheese"}},
		{65, []string{"Hot Chocolate", "Marshmallow"}},
		{66, []string{"Tea", "Biscuits"}},
		{67, []string{"Coffee", "Muffin"}},
		{68, []string{"Latte", "Cookie"}},
		{69, []string{"Espresso", "Chocolate"}},
		{70, []string{"Cappuccino", "Croissant"}},
		{71, []string{"Mocha", "Brownie"}},
		{72, []string{"Cold Brew", "Donut"}},
		{73, []string{"Iced Tea", "Lemon"}},
		{74, []string{"Lemonade", "Mint"}},
		{75, []string{"Juice", "Toast"}},
		{76, []string{"Energy Drink", "Snack Bar"}},
		{77, []string{"Protein Shake", "Banana"}},
		{78, []string{"Oats", "Milk"}},
		{79, []string{"Cereal", "Fruits"}},
		{80, []string{"Granola", "Yogurt"}},
		{81, []string{"Avocado Toast", "Egg"}},
		{82, []string{"Quinoa Bowl", "Veggies"}},
		{83, []string{"Tofu Stir Fry", "Rice"}},
		{84, []string{"Veg Burger", "Sweet Potato Fries"}},
		{85, []string{"Stuffed Peppers", "Salad"}},
		{86, []string{"Lasagna", "Breadsticks"}},
		{87, []string{"Spaghetti", "Meatballs"}},
		{88, []string{"Mac and Cheese", "Garlic Bread"}},
		{89, []string{"Risotto", "Parmesan"}},
		{90, []string{"Gnocchi", "Pesto"}},
		{91, []string{"Focaccia", "Olive Oil"}},
		{92, []string{"Bruschetta", "Tomato"}},
		{93, []string{"Calzone", "Dip"}},
		{94, []string{"Stuffed Pasta", "Sauce"}},
		{95, []string{"BBQ Ribs", "Corn"}},
		{96, []string{"Pulled Pork", "Coleslaw"}},
		{97, []string{"Roast Chicken", "Potatoes"}},
		{98, []string{"Lamb Chops", "Mint Sauce"}},
		{99, []string{"Seafood Platter", "Lemon"}},
		{100, []string{"Chef Special", "Dessert"}},
	}

	for i,orders := range slice_of_order_structures{
		time.Sleep(2 * time.Second)
		order_json,err:=json.Marshal(orders)
		errHelper(err)
		publish_message(ch,&q,order_json,ctx)
		fmt.Printf("Published %d, %v\n",i,orders)
	}
}

