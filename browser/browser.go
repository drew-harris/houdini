package browser

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
)

func ClockIn(b *rod.Browser) error {
	page, err := signIn(b)

	// Value: 2 = out, 1 = in
	err = rod.Try(func() {
		err = page.MustElement("select.ps-dropdown").Select([]string{`[value="1"]`}, true, rod.SelectorTypeCSSSector)
		if err != nil {
			panic(err)
		}
		page.MustWaitStable()

		page.MustElementR(`a.ps-button[role="button"]`, "Submit").MustClick()
		time.Sleep(time.Second * 5)
		page.MustClose()
	})

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Clocked in successfully")
	return nil
}

func ClockOut(b *rod.Browser) error {
	page, err := signIn(b)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Value: 2 = out, 1 = in
	err = rod.Try(func() {
		err = page.MustElement("select.ps-dropdown").Select([]string{`[value="2"]`}, true, rod.SelectorTypeCSSSector)
		if err != nil {
			panic(err)
		}
		page.MustWaitStable()

		page.MustElementR(`a.ps-button[role="button"]`, "Submit").MustClick()
		time.Sleep(time.Second * 5)
		page.MustClose()
	})

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Clocked out successfully")
	return nil
}

func signIn(b *rod.Browser) (*rod.Page, error) {
	page := b.MustPage("https://my.smu.edu/psc/ps/EMPLOYEE/SA/c/NUI_FRAMEWORK.PT_LANDINGPAGE.GBL")

	page.MustWaitStable()

	if page.MustInfo().Title != "Homepage" {
		err := mySmuLogin(page)
		if err != nil {
			return nil, err
		}
	} else {
		// We are logged in already !!!
		fmt.Println("Logged in already")
	}

	fmt.Println("Going to time sheet")
	err := gotoTimeSheet(page)
	if err != nil {
		return page, err // Could not access time sheet
	}

	return page, nil
}

func mySmuLogin(page *rod.Page) error {
	fmt.Println("LOGGING IN WITH MYSMU")

	err := rod.Try(func() {
		page.MustElement("#username").Input(os.Getenv("username"))
		page.MustElement("#password").Input(os.Getenv("password"))

		page.MustElement("button[type=\"submit\"]").MustClick()

		ctx, cancel := context.WithCancel(context.Background())

		pageWithCancel := page.Context(ctx)

		go func() {
			time.Sleep(time.Minute * 5)
			cancel()
		}()

		pageWithCancel.MustElementR("button", "Yes, trust browser").MustClick()

		page.MustWaitStable()
		fmt.Println("Duo push approved")

	})

	if err != nil {
		return err
	}

	return nil
}

func gotoTimeSheet(page *rod.Page) error {
	err := rod.Try(func() {
		_, err := page.Timeout(time.Second * 25).Element(`a.ps-button[title="Homepage Selector"]`)
		if err != nil {
			fmt.Println("Could not find pt_notify")
			panic(err)
		}
		page.Timeout(time.Second * 5).MustElementX("/html/body/form/div[2]/div[4]/div[2]/div/div/div/div/div[4]/section/div/div[3]/div[1]/div[1]/div/span/a").MustClick()

		page.MustWaitStable()

		fmt.Println("Looking for self service button")
		page.MustElementR("a", "Employee Self Service").MustClick()

		page.MustWaitStable()

		fmt.Println("Looking for time reporting")
		page.MustElementR("div.ps_box-group[role=\"link\"]", "Time Reporting").MustClick()

		page.MustWaitStable()
	})

	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Could not access time sheet")
	}
	return nil
}
