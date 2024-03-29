package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
	"strconv"

	"github.com/gorilla/mux"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

type PlaceOrderRequest struct { 
	Origin []string `json:"origin"`
	Destination []string `json:"destination"`
}

type PlaceOrderSuccess struct{
	Id int `json:"id"`
	Distance string `json:"distance"`
	Status string `json:"status"`
}

type TakeOrder struct{
	Status string `json:"status"`
}

type Error struct{
	ErrorMessage string `json:"error"`
}

//structs for google map api
type Result struct{
	Destination_addresses []string `json:"destination_addresses"`
	Origin_addresses []string `json:"origin_addresses"`
	Rows []Rows `json:"rows"`
	Status string `json:"status"`
}

type Rows struct{
	Elements []Elements `json:"elements"`
}

type Elements struct{
	Distance Distance `json:"distance"`
	Duration Duration `json:"duration"`
	Status string  `json:"status"`
}
type Distance struct{
	Text string `json:"text"`
	Value int `json:"value"`
}
type Duration struct{
	Text string `json:"text"`
	Value int `json:"value"`
}

var APIKey = "<Your Google Map API Key>"

var errorMessage = map[string]string{
    "PageNumberExceeds": "Page number exceeds total number of pages.",
	"PageNotInt":"page is not a valid integer.",
	"LimitNotInt":"limit is not a valid integer.",
	"WrongRequestFormat":"The format of the request body is wrong.",
	"WrongStatus":"Status in request body is not in the correct value.",
	"ZeroResults":"Route not found. Order will not be placed.",
	"DbSelect":"Database SELECT failed.",
	"DbInsert":"Database INSERT failed.",
	"DbUpdate":"Database UPDATE failed.",
	"OrderTaken":"The order was taken by the others.",
	"OrderNotFound":"The order id is not found.",
	"APIKey":"Please configure the Google Map API Key in app.go var APIKey.",
}

func (app *App) SetupRouter() {
	app.Router.
		Methods("GET").
		Path("/endpoint/{id}").
		HandlerFunc(app.getFunction)
		
	app.Router.
		Methods("POST").
		Path("/endpoint").
		HandlerFunc(app.postFunction)
		
	app.Router.
		Methods("POST").
		Path("/orders").
		HandlerFunc(app.placeOrder)
	
	app.Router.
		Methods("PATCH").
		Path("/orders/{id}").
		HandlerFunc(app.takeOrder)
		
	app.Router.
		Methods("GET").
		Path("/orders").
		Queries("page","{page}").
		Queries("limit","{limit}").
		HandlerFunc(app.listOrder)
}

func (app *App) getFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["page"]
	if !ok {
		log.Fatal("No ID in the path")
	}

	dbdata := &DbData{}
	err := app.Database.QueryRow("SELECT ORDER_ID, TOTAL_DISTANCE, STATUS FROM `Orders` WHERE ORDER_ID = ?", id).Scan(&dbdata.ID, &dbdata.Distance, &dbdata.Status)
	if err != nil {
		log.Println(errorMessage["DbSelect"])
		error := Error{errorMessage["DbSelect"]}
		errorJson,err := json.Marshal(&error)
		if err!= nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(errorJson)
		return
		
	}

	log.Println("You fetched a thing! 1")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dbdata); err != nil {
		//panic(err)
		log.Println(err)
	}
}

func (app *App) postFunction(w http.ResponseWriter, r *http.Request) {
	_, err := app.Database.Exec("INSERT INTO `Orders` VALUES (NULL,22.277627,114.173463,22.2783034,114.1796477,1.1,'TAKEN',NULL)")
	if err != nil {
		log.Println("Error:")
		log.Println(err)
		log.Println(errorMessage["DbInsert"])
		
		error := Error{errorMessage["DbInsert"]}
		errorJson,err := json.Marshal(&error)
		if err!= nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(errorJson)
		return
	}

	log.Println("You called a thing!")
	w.WriteHeader(http.StatusOK)
}

func (app *App) placeOrder(w http.ResponseWriter, r *http.Request) {
	
	order := PlaceOrderRequest{} //initialize

	if r.Method == "POST" {
		/*bytesBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		json.Unmarshal(bytesBody,&order)*/
		
		err := json.NewDecoder(r.Body).Decode(&order)
		
		if err != nil{
			error := Error{errorMessage["WrongRequestFormat"]}
			errorJson,err := json.Marshal(&error)
			if err!= nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type","application/json")
			w.Write(errorJson)
			return
		}
		
		if len(order.Origin)!=2 || len(order.Destination)!=2{
			error := Error{errorMessage["WrongRequestFormat"]}
			errorJson,err := json.Marshal(&error)
			if err!= nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type","application/json")
			w.Write(errorJson)
			return			
			
		}
		
		//log.Println("Request info: ")
		//log.Println(order)
		//log.Println(order.Origin)
		//log.Println(len(order.Origin))
		//log.Println(order.Destination)
		//log.Println(len(order.Destination))
					
		
		//call google map API to get distaince
		if len(order.Origin)==2 && len(order.Destination)==2 {
		
			var originLat = order.Origin[0]
			var originLong = order.Origin[1]
			var destLat = order.Destination[0]
			var destLong = order.Destination[1]
			
			//latitude must be between -90.0 to 90.0, longitude -180 to 180
			originLatFloat, err := strconv.ParseFloat(originLat, 64)
			originLongFloat, err := strconv.ParseFloat(originLong, 64)
			destLatFloat, err := strconv.ParseFloat(destLat, 64)
			destLongFloat, err := strconv.ParseFloat(destLong, 64)
			
			//fmt.Println("originLatFloat:")
			//fmt.Println(originLatFloat)
			
			if originLatFloat < -90.0 || originLatFloat>90.0{
				error := Error{"origin latitude is not in the valid range."}
				errorJson,err := json.Marshal(&error)
				if err!= nil{
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type","application/json")
				w.Write(errorJson)
				return
			}
			if originLongFloat < -180.0 || originLongFloat>180.0{
				error := Error{"origin longitude is not in the valid range."}
				errorJson,err := json.Marshal(&error)
				if err!= nil{
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type","application/json")
				w.Write(errorJson)
				return
			}
			if destLatFloat < -90.0 || destLatFloat>90.0{
				error := Error{"destination latitude is not in the valid range."}
				errorJson,err := json.Marshal(&error)
				if err!= nil{
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type","application/json")
				w.Write(errorJson)
				return
			}
			if destLongFloat < -180.0 || destLongFloat>180.0{
				error := Error{"destination longitude is not in the valid range."}
				errorJson,err := json.Marshal(&error)
				if err!= nil{
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type","application/json")
				w.Write(errorJson)
				return
			}

					
			fmt.Println("originLat:"+ originLat)
			fmt.Println("originLong:"+ originLong)
			fmt.Println("destLat:"+ destLat)
			fmt.Println("destLong:"+ destLong)
			
			//check if API key is configured
			if APIKey=="<Your Google Map API Key>"{
				error := Error{errorMessage["APIKey"]}
				errorJson,err := json.Marshal(&error)
				if err!= nil{
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type","application/json")
				w.Write(errorJson)
				return				
				
			}
						
			response, err := http.Get("https://maps.googleapis.com/maps/api/distancematrix/json?origins="+originLat+","+originLong+"&destinations="+ destLat+ ","+ destLong+"&key="+ APIKey)
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
				error := Error{"The HTTP request failed with error."}
				errorJson,err := json.Marshal(&error)
				if err!= nil{
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type","application/json")
				w.Write(errorJson)
				return
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println("Called google map API!")
				//fmt.Println(string(data))
				
				//get distance
				var rs Result
				json.Unmarshal([]byte(data),&rs)
				fmt.Println("rs")
				fmt.Println(rs)		
				fmt.Println(rs.Rows[0].Elements[0].Distance.Value)
							
				d := strconv.Itoa(rs.Rows[0].Elements[0].Distance.Value)
				status := rs.Rows[0].Elements[0].Status
				fmt.Println("Status: "+status)
				
				if status == "ZERO_RESULTS"{
					error := Error{errorMessage["ZeroResults"]}
					errorJson,err := json.Marshal(&error)
					if err!= nil{
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					w.Header().Set("Content-Type","application/json")
					w.Write(errorJson)
					return
					
				}
				
				if status == "OK"{
					//insert into db
					_, err := app.Database.Exec("INSERT INTO `Orders` VALUES (NULL,"+originLat+","+originLong+","+destLat+","+destLong+","+d+",'UNASSIGNED',NULL)")
					
					if err != nil {
						log.Println("Error:")
						log.Println(err)
						log.Println(errorMessage["DbInsert"])
						
						error := Error{errorMessage["DbInsert"]}
						errorJson,err := json.Marshal(&error)
						if err!= nil{
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
						w.Header().Set("Content-Type","application/json")
						w.Write(errorJson)
						return
					}
					
					//get largest id
					
					var maxId int
					err2 := app.Database.QueryRow("SELECT MAX(ORDER_ID) FROM `Orders`").Scan(&maxId)
					if err2 != nil {
						error := Error{"Maximum ID is not found in the database."}
						errorJson,err := json.Marshal(&error)
						if err!= nil{
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
						w.Header().Set("Content-Type","application/json")
						w.Write(errorJson)
						return
					}
					
			
					//write success response body
					success := PlaceOrderSuccess{Id:maxId,Distance:d,Status:"UNASSIGNED"}
					successJson,err := json.Marshal(&success)
					if err!= nil{
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					w.Header().Set("Content-Type","application/json")
					w.Write(successJson)	
				}
					

				
			}

		log.Println("POST DONE")
		
		}
		
		//fmt.Fprint(w, d)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}


	log.Println("You called postOrder!")
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
}


func (app *App) takeOrder(w http.ResponseWriter, r *http.Request){
	//get order id
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Fatal("No ID in the path")
	}
	
	takeorder := TakeOrder{} //initialize
	error := Error{errorMessage["WrongRequestFormat"]}
	errorWrongStatus := Error{errorMessage["WrongStatus"]}

	if r.Method == "PATCH" {
		
		err := json.NewDecoder(r.Body).Decode(&takeorder)
		if err != nil{
			//panic(err)
			errorJson,err := json.Marshal(&error)
			if err!= nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type","application/json")
			w.Write(errorJson)	
			return
		}
		
		//check if id exist in database
		rows, err :=app.Database.Query("SELECT ORDER_ID, TOTAL_DISTANCE, STATUS FROM `Orders` WHERE ORDER_ID="+id)
		var count int
		for rows.Next(){
			count++
		}
		log.Println("Rows selected select ID:")
		log.Println(count)
		if count == 0 {
			log.Println(errorMessage["OrderNotFound"])
			error := Error{errorMessage["OrderNotFound"]}
			errorJson,err := json.Marshal(&error)
			if err!= nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type","application/json")
			w.Write(errorJson)
			return
		}
	
		var takeorder = takeorder.Status
		
		if takeorder== "TAKEN"{
			//update status to "TAKEN"
			
			res, err := app.Database.Exec("UPDATE `Orders` SET STATUS='TAKEN' WHERE ORDER_ID="+id+ " AND `STATUS` ='UNASSIGNED'")
			
			if err != nil {
				log.Println("Error:")
				log.Println(err)
				//log.Fatal("Database UPDATE failed")
				log.Println(errorMessage["DbUpdate"])
				error := Error{errorMessage["DbUpdate"]}
				errorJson,err := json.Marshal(&error)
				if err!= nil{
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type","application/json")
				w.Write(errorJson)
				return
				
				
				//TODO id not exist
			}
			rowCnt,err := res.RowsAffected()
			log.Println("Rows affected:")
			log.Println(rowCnt)
			if rowCnt ==0{
				log.Println(errorMessage["OrderTaken"])
				error := Error{errorMessage["OrderTaken"]}
				errorJson,err := json.Marshal(&error)
				if err!= nil{
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type","application/json")
				w.Write(errorJson)
				return
				
			}
			//update status to "SUCCESS"
			_,err2 := app.Database.Exec("UPDATE `Orders` SET STATUS='SUCCESS' WHERE ORDER_ID="+id+ " AND `STATUS` ='TAKEN'")
			
			if err2 != nil {
				log.Println("Error:")
				log.Println(err)
				log.Println(errorMessage["DbUpdate"])
				error := Error{errorMessage["DbUpdate"]}
				errorJson,err := json.Marshal(&error)
				if err!= nil{
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type","application/json")
				w.Write(errorJson)
				return
			}
			
			//write success response body
			success := TakeOrder{Status:"SUCCESS"} 
			successJson,err := json.Marshal(&success)
			if err!= nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type","application/json")
			w.Write(successJson)		
			
		}else{
			errorJson,err := json.Marshal(&errorWrongStatus)
			if err!= nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type","application/json")
			w.Write(errorJson)		
			
		}
		log.Println("PATCH DONE")
				
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}


	log.Println("You called patch order!")
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)		
}

func (app *App) listOrder(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	pageStr, ok := vars["page"]
	if !ok {
		log.Fatal("No page in the path")
	}
	limitStr, ok := vars["limit"]
	if !ok {
		log.Fatal("No limit in the path")
	}
	//parse page, limit to int
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		error := Error{errorMessage["PageNotInt"]}
		errorJson,err := json.Marshal(&error)
		if err!= nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(errorJson)
			
		return
	}
	
	limit, err2 := strconv.Atoi(limitStr)
	if err2 != nil {
		error := Error{errorMessage["LimitNotInt"]}
		errorJson,err := json.Marshal(&error)
		if err!= nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(errorJson)
			
		return
	}
	
	
	
	//dbresult := &PlaceOrderSuccess{}
	orders := make([]*PlaceOrderSuccess,0)
	rows, err :=app.Database.Query("SELECT ORDER_ID, TOTAL_DISTANCE, STATUS FROM `Orders`")
		
	//err := app.Database.QueryRow("SELECT ORDER_ID, TOTAL_DISTANCE, STATUS FROM `Orders`").Scan(&dbresult.Id, &dbresult.Distance, &dbresult.Status)
	
	if err != nil {
		log.Println(errorMessage["DbSelect"])
		error := Error{errorMessage["DbSelect"]}
		errorJson,err := json.Marshal(&error)
		if err!= nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(errorJson)
		return
	}
	defer rows.Close()
	
	for rows.Next(){
		order := new(PlaceOrderSuccess)
		err:= rows.Scan(&order.Id, &order.Distance, &order.Status)
		if err != nil {
			log.Println(errorMessage["DbSelect"])
			error := Error{errorMessage["DbSelect"]}
			errorJson,err := json.Marshal(&error)
			if err!= nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type","application/json")
			w.Write(errorJson)
			return
		}
		
		orders = append(orders,order)
		
	}
	if err := rows.Err(); err != nil {
		//panic(err)
		log.Println(err)
	}
	//assign orders to diffeerent pages
	//use an array to stores different pages or records
	
	var ordersAll [][]*PlaceOrderSuccess
	
	for i := 0; i< len(orders); i+= limit{
		end := i + limit
		
		if end> len(orders){
			end = len(orders)
		}
		
		ordersAll = append(ordersAll,orders[i:end])
		
	}

	log.Println("You called list orders!")
	log.Println("Page: "+pageStr)
	log.Println("Page int: "+strconv.Itoa(page))
	
	log.Println("Limit: "+limitStr)
	log.Println("Total no. of pages: "+strconv.Itoa(len(ordersAll)))
	
	if page> len(ordersAll){
		error := Error{errorMessage["PageNumberExceeds"]}
		errorJson,err := json.Marshal(&error)
		if err!= nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(errorJson)
			
		return
	}
	
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ordersAll[(page-1)]); err != nil {
		log.Println(err)
	}
	
}
