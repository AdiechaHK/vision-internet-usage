package schedule

import (
	"fmt"
	"hello-heroku/data"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

func makeLogin(url string, fdt map[string]string) {
	c := colly.NewCollector()

	loggedIn := false

	c.OnHTML("#ContentPlaceHolder1_tbAdditonal_tpPlan_gdPlan", func(ele *colly.HTMLElement) {
		remain := ele.ChildText("#ContentPlaceHolder1_tbAdditonal_tpPlan_gdPlan_lblRemain_0")
		used := ele.ChildText("#ContentPlaceHolder1_tbAdditonal_tpPlan_gdPlan_lblUsed_0")
		re := regexp.MustCompile(`^\d+`)

		u, err := strconv.Atoi(re.FindString(used))
		if err != nil {
			panic(err)
		}
		r, err := strconv.Atoi(re.FindString(remain))
		if err != nil {
			panic(err)
		}
		fmt.Println("Used: " + used + " and Remain: " + remain)
		data.StoreData(u, r)
		if loggedIn {
			fmt.Println("Logged in")
		} else {
			fmt.Println("Not logged in")
		}
	})

	c.OnResponse(func(r *colly.Response) {

		if loggedIn {
			fmt.Println("-------------------------PLAN PAGE-----------------------------------")
			// fmt.Println("START ------------------------------------------------------------")
			// fmt.Println(string(r.Body))

			// op, err := os.Create("op.html")
			// if err != nil {
			// 	panic(err)
			// }
			// defer op.Close()

			// lineCount, err := op.Write(r.Body)
			// if err != nil {
			// 	panic(err)
			// }
			// fmt.Printf("wrote %d bytes\n", lineCount)
			// fmt.Println("END ------------------------------------------------------------")
		} else {
			loggedIn = true
			url = r.Request.AbsoluteURL("AccountViewPlan.aspx")
			fmt.Println("Going to " + url)
			c.Visit(url)
		}

	})
	c.Post(url, fdt)
}

func TestFun() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("form", func(e *colly.HTMLElement) {

		fdt := make(map[string]string)
		e.ForEach("input", func(_ int, elem *colly.HTMLElement) {
			name := elem.Attr("name")
			value := elem.Attr("value")
			fmt.Println(name, " >> ", value)
			if name == "txtUserName" {
				fdt[name] = "HARIKRUSHNA_KRUSHNANAGAR"
			} else if name == "txtPassword" {
				fdt[name] = "123456"
			} else {
				fdt[name] = value
			}
		})

		url := e.Request.AbsoluteURL(e.Attr("action"))

		fmt.Println(url, fdt)
		makeLogin(url, fdt)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://3.7.153.175/vtpl/customer")
}
