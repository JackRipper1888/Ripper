package main

import (
	"fmt"
)

//go语言的interface类型 本质是一个指针类型
type AnimalIF interface {
	Sleep()        //让动物睡觉
	Color() string //获取动物颜色
	Type() string  //获取动物类型
}

// ===========================
type cat struct {
	color string
}

func (this *cat) Sleep() {
	fmt.Println("Cat is Sleep")
}

func (this *cat) Color() string {
	return this.color
}

func (this *cat) Type() string {
	return "Cat"
}

// =================
type dog struct {
	color string
}

func (this *dog) Sleep() {
	fmt.Println("Dog is Sleep")
}

func (this *dog) Color() string {
	return this.color
}

func (this *dog) Type() string {
	return "Dog"
}

//创建一个工厂
func Factory(kind string) AnimalIF {
	switch kind {
	case "dog":
		//生产一个狗
		return &dog{"yellow"}
	case "cat":
		//生产一个猫
		cat := new(cat)
		cat.color = "green"
		return cat
	default:
		fmt.Println("kind is wrong")
		return nil
	}
}

func showAnimal(animal AnimalIF) {
	animal.Sleep()              //多态现象
	fmt.Println(animal.Color()) //多态
	fmt.Println(animal.Type())  //多态
}

func main() {
	//var animal AnimalIF
	//animal = Factory("cat")
	//count := 0
	//for {
	//	count++
	//	Get(
	//		"http://183.60.143.82:3030/peer/select/tracker","0001","0008",map[string]string{
	//		"streamid":"us",
	//	})
	//	if count == 10000{
	//		return
	//	}
	//}
	//animal.Sleep()
	//showAnimal(animal)
	//t2,_ := time.Parse("2016-01-02 15:05:05", "2018-04-23 00:00:06")
	//t1 := time.Date(2018, 1, 2, 15, 5, 0, 0, time.Local)
	//t2 := time.Date(2018, 1, 2, 15, 0, 0, 0, time.Local)
	//
	//logs.Error(t1.Unix()%(24*3600)%300,t2.Unix()%300)
}
