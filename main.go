package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"main.go/news"
)

//main data
var maintpl = template.Must(template.ParseFiles("templates/main.html"))

//calc data
var calcTmp = template.Must(template.ParseFiles("templates/calc.html"))
type CData struct {
	Result float64
}

//news data
var newsTmp = template.Must(template.ParseFiles("templates/news.html"))
var newsapi *news.Client
type NData struct {
	Var1	   int
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
}

var ndata = &NData{
	Var1: 0,
	Query:      "",
	NextPage:   1,
	TotalPages: 0,
	Results:    &news.Results{},
}




func homePage(w http.ResponseWriter,r *http.Request){
	
	buf := &bytes.Buffer{}
	err := maintpl.Execute(buf,nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}


func main(){
	fmt.Print("started!")

	r := mux.NewRouter()


	fs := http.FileServer(http.Dir("assets/"))
	// http.Handle("/assets/",http.StripPrefix("/assets/",fs))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", fs))
	
//news data
	apikey := "fa1ced5f40ad4c86b0a9339736aa4d67"
	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, apikey, 20)
	
	r.HandleFunc( "/", homePage)
	r.HandleFunc("/calc/{action}", calcHandler)
	r.HandleFunc("/news", newsInit)
	r.HandleFunc("/news/search", searchHandler(newsapi))

	http.ListenAndServe(":8080", r)
	
}


func do(a *float64, b *float64, op *string) float64{
	var c float64 = 0
	switch *op{
	case "+":
		c= *a + *b
	case "-":
		c= *a - *b
	case "*":
		c= *a * *b
	case "/":
		c= *a / *b
	case "":
		return 0

	}

	return c
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	data := &CData{
		Result: 0,
	}
	var val1 float64 = 0
	var val2 float64 = 0
	var check bool = true
	op := ""
	Cval1, err := r.Cookie("Val1")
    if err != nil {
        switch {
        case errors.Is(err, http.ErrNoCookie):
            // http.Error(w, "cookie not found", http.StatusBadRequest)
        default:
            log.Println(err)
            http.Error(w, "server error", http.StatusInternalServerError)
        }
        // return
    } else {
		val1,_ = strconv.ParseFloat(Cval1.Value,64)
	}
	Cval2, err := r.Cookie("Val2")
    if err != nil {
        switch {
        case errors.Is(err, http.ErrNoCookie):
            // http.Error(w, "cookie not found", http.StatusBadRequest)
        default:
            log.Println(err)
            http.Error(w, "server error", http.StatusInternalServerError)
        }
        // return
    } else {
		val2,_ = strconv.ParseFloat(Cval2.Value,64)
	}
	Cres, err := r.Cookie("Res")
    if err != nil {
        switch {
        case errors.Is(err, http.ErrNoCookie):
            // http.Error(w, "cookie not found", http.StatusBadRequest)
        default:
            log.Println(err)
            http.Error(w, "server error", http.StatusInternalServerError)
        }
        // return
    } else {
		data.Result,_ = strconv.ParseFloat(Cres.Value,64)
	}
	Cop, err := r.Cookie("Op")
    if err != nil {
        switch {
        case errors.Is(err, http.ErrNoCookie):
            // http.Error(w, "cookie not found", http.StatusBadRequest)
        default:
            log.Println(err)
            http.Error(w, "server error", http.StatusInternalServerError)
        }
        // return
    } else {
		op = Cop.Value
	}
	Ccheck, err := r.Cookie("Check")
    if err != nil {
        switch {
        case errors.Is(err, http.ErrNoCookie):
            // http.Error(w, "cookie not found", http.StatusBadRequest)
        default:
            log.Println(err)
            http.Error(w, "server error", http.StatusInternalServerError)
        }
        // return
    } else {
		check,_ = strconv.ParseBool(Ccheck.Value)
	}	




	vars := mux.Vars(r)

	a, err := strconv.Atoi(vars["action"])
	if err != nil {
		fmt.Printf("fuck me")
	}
	switch a {
	case 10:
		if val1 == 0 {
			data.Result = math.Sqrt(data.Result)
		} else { 
			val2 = data.Result
			data.Result = do(&val1,&val2,&op)
		}
		op = "%"
	case 11:
		data.Result = math.Sqrt(data.Result)
		op = "r"
	case 12:
		data.Result = math.Floor(data.Result/10)
		op = "ce"
	case 13:
		data.Result = 0
		val1 = 0
		val2 = 0
		op = ""
		check = false
	case 14:
		if val1 == 0 {
			val1 = data.Result
			data.Result = 0
		} else { 
			val2 = data.Result
			data.Result = do(&val1,&val2,&op)
			val1 = data.Result
			val2 = 0
			check = true
		}
		op = "-"
	case 15:
		if val1 == 0 {
			val1 = data.Result
			data.Result = 0
		} else { 
			val2 = data.Result
			data.Result = do(&val1,&val2,&op)
			val1 = data.Result
			val2 = 0
			check = true
		}
		op = "/"
	case 16:
		if val1 == 0 {
			val1 = data.Result
			data.Result = 0
		} else { 
			val2 = data.Result
			data.Result = do(&val1,&val2,&op)
			val1 = data.Result
			val2 = 0
			check = true
		}
		op = "*"
//todo-make the . right
	case 17:
		
		op = "d"
	case 18:
		if val1 == 0 {
			
		} else {
			val2 = data.Result
			data.Result = do(&val1,&val2,&op)
			val1 = data.Result
			val2 = 0
			check = true
		}
		// op = ""
	case 19:
		if val1 == 0 {
			val1 = data.Result
			data.Result = 0
		} else { 
			val2 = data.Result
			data.Result = do(&val1,&val2,&op)
			val1 = data.Result
			val2 = 0
			check = true
		}
		op = "+"
	default:
		// data.Result,err = strconv.ParseFloat(strconv.FormatFloat(data.Result*10, 'f', -1, 64)  + strconv.Itoa(a),64)
		// data.Result = data.Result * 10 + float64(a)
		// if val1 == 0 {
		// 	data.Result = data.Result * 10 + float64(a)
		// } else { 
		// 	if data.Result == val1 {
		// 		data.Result = 0
		// 	}
		// 	data.Result = data.Result * 10 + float64(a)
		// }
		if check {
			data.Result = 0
			check = false
		}
		data.Result = data.Result * 10 + float64(a)
		
	if err != nil {
		fmt.Printf("fuck me")
	}
	}


	Csval1 := http.Cookie{
		Name:     "Val1",
		Value:     strconv.FormatFloat(val1,'g', -1, 64),
	}
	Csval2 := http.Cookie{
		Name:    "Val2",
		Value:   strconv.FormatFloat(val2, 'g', -1, 64),
	}
	Csres := http.Cookie{
		Name:     "Res",
		Value:	  strconv.FormatFloat(data.Result,'g',-1,64),
	}
	Csop := http.Cookie{
		Name:     "Op",
		Value:	  op,
	}
	Cscheck := http.Cookie{
		Name:     "Check",
		Value:	  strconv.FormatBool(check),
	}
	http.SetCookie(w, &Csval1)
	http.SetCookie(w, &Csval2)
	http.SetCookie(w, &Csres)
	http.SetCookie(w, &Csop)
	http.SetCookie(w, &Cscheck)

	buf := &bytes.Buffer{}
	err = calcTmp.Execute(buf, data)
	if err != nil {
		// fmt.Print("hey")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}


func newsInit(w http.ResponseWriter , r *http.Request){

	
	buf := &bytes.Buffer{}
	err := newsTmp.Execute(buf,nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}
func (s *NData) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}
func (s *NData) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}
	return s.NextPage-1
}
func (s *NData) PreviousPage() int{
	return s.CurrentPage()-1
}
func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		// fmt.Println("Search Query is: ", searchQuery)

		// fmt.Println("page is: ", page)

		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// fmt.Printf("%+v",results)
		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ndata = &NData{
			Var1: 		ndata.Var1,
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results:    results,
		}

		if ok := !ndata.IsLastPage();ok {
			ndata.NextPage++
		}


		buf := &bytes.Buffer{}
		err = newsTmp.Execute(buf, ndata)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)

	}
}
