package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var seatingCapacity = 10
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 3 * time.Second

func main() {
	rand.Seed(time.Now().UnixNano())

	color.Yellow("Sleeping Barber")
	color.Yellow("---------------")

	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create the barbershop
	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	color.Green("The shop is open for the day!")

	// add barbers
	for i := 1; i <= 5; i++ {
		shop.addBarber(fmt.Sprintf("Barber#%v", i))
	}

	// start the barbershop as a goroutine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	// add clients as a goroutine
	go func(i int) {
		for {
			// get a random number with average arrival rate
			randomMs := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMs)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}(1)

	// block until the barbershop is closed
	<-closed
}
