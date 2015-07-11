package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
)

const (
	PORT = 4444
)

var (
	baseUrl = fmt.Sprintf("http://localhost:%v/admin", PORT)
	driver  *agouti.WebDriver
	page    *agouti.Page
)

func TestMain(m *testing.M) {
	var t *testing.T
	var err error

	driver = agouti.ChromeDriver() // choose browser driver
	driver.Start()

	go Start(PORT) // start our program

	page, err = driver.NewPage() // get page object from driver, this is what we will use to perform browser testing
	if err != nil {
		t.Error("Failed to open page.")
	}

	RegisterTestingT(t)
	test := m.Run() // start test

	driver.Stop() // close driver after test
	os.Exit(test)
}

func StopDriverOnPanic() {
	var t *testing.T
	if r := recover(); r != nil {
		debug.PrintStack()
		fmt.Println("Recovered in f", r)
		driver.Stop()
		t.Fail()
	}
}

func TestEnv(t *testing.T) {
	Expect(page.Navigate(fmt.Sprintf("%v/user", baseUrl))).To(Succeed())
}

func TestCreateUser(t *testing.T) {
	var user User
	userName := "user name"

	Expect(page.Navigate(fmt.Sprintf("%v/user", baseUrl))).To(Succeed()) // visit user page
	Expect(page.Find("#plus").Click()).To(Succeed())                     // click add user button

	page.Find("#QorResourceName").Fill(userName) // fill in user name

	page.FindByButton("Save").Click() // submit form

	DB.Last(&user) // query the user we just created

	if user.Name != userName { // assert it created as we expected
		t.Error("user name not set")
	}
}
