package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/golang/go/src/pkg/strconv"
		"time"
)

var client *redis.Client
var templates *template.Template
var db *sql.DB
var connStr = "root:YOURSQLUSERNAME@tcp(YOURSQLPORT)/shoes"

var shoeArray = getShoesIndex()

type Shoe struct {
	Id string
	ShoeName string
	Designer string
	Price string
	SizeScore string
}
type SizeScoring struct {
	shoeID string
	score string
	date string
}
type AverageSizeScoring struct {
	shoeID string
	averageScore string
}


func main(){
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexHanlder).Methods("GET")
	r.HandleFunc("/addSizeReview", addSizeReview).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func indexHanlder(w http.ResponseWriter, r *http.Request){
	var shoeCatalog = getShoeCatalog()
	templates.ExecuteTemplate(w, "index.html", shoeCatalog)
}

func addSizeReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	trueToSizeScore := r.FormValue("trueToSizeScore")
	shoeID := r.FormValue("shoeID")

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var date = time.Now().Local().String()
	var queryString = "insert into true_to_size_scoring (shoeID, score, date) values ('" + shoeID + "', '" +
		trueToSizeScore + "', '" + date + "')"
	db.Query(queryString)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getShoeCatalog() []Shoe {
	var sizingScore = getSizeScoring()

	for index := range shoeArray{
		var id = shoeArray[index].Id
		for e := range sizingScore {
			if sizingScore[e].shoeID == id {
				shoeArray[index].SizeScore = sizingScore[e].averageScore
			}
		}
	}

	return shoeArray
}

func getShoeIds() []string {
	var shoeIdsArray []string

	for e := range shoeArray {
		shoeIdsArray = append(shoeIdsArray, shoeArray[e].Id)
	}

	return shoeIdsArray
}

func getShoesIndex() []Shoe {
	var shoeArray []Shoe
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM shoe_catalog")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var shoe Shoe
		err = results.Scan(&shoe.Id, &shoe.ShoeName, &shoe.Designer, &shoe.Price)
		if err != nil {
			panic(err.Error())
		}

		shoeArray = append(shoeArray, shoe)
	}

	return shoeArray
}

func getSizeScoring() []AverageSizeScoring{
	var sizeScoring []SizeScoring
	var averageSizeScoring []AverageSizeScoring
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM true_to_size_scoring")
	if err != nil {
		panic(err.Error())
	}
	//if results.Next() == false {
	//	return averageSizeScoring
	//}

	for results.Next(){
		var sizeScore SizeScoring
		err = results.Scan(&sizeScore.shoeID, &sizeScore.score, &sizeScore.date)
		if err != nil {
			panic(err.Error())
		}
		sizeScoring = append(sizeScoring, sizeScore)
	}

	var idArray = getShoeIds()
	for e := range idArray {
		var rawTotal = 0
		var counter = 0
		for f := range sizeScoring {
			if (idArray[e] == sizeScoring[f].shoeID){
				counter++
				i, err := strconv.Atoi(sizeScoring[f].score)
				if err != nil {
					panic(err.Error())
				}
				rawTotal += i
			}
		}
		var average = rawTotal / counter
		var averageScore AverageSizeScoring
		averageScore.shoeID = idArray[e]
		averageScore.averageScore = strconv.Itoa(average)

		averageSizeScoring = append(averageSizeScoring, averageScore)
	}

	return averageSizeScoring
}





