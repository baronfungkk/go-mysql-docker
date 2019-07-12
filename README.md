# go-mysql-docker
This is my first time to use Golang to build an API (used node.JS to build before), using docker and MySql.
This project is used for placing order of delivery, taking order as well as listing all orders.
Google Map API is used to calculate the distance between the origin and destination.
Hope this can give you some insights.

## API ENDPOINTS


### Place order
- Method: `POST`
- Path : `/posts`
- Request body:

    ```
    {
        "origin": ["START_LATITUDE", "START_LONGTITUDE"],
        "destination": ["END_LATITUDE", "END_LONGTITUDE"]
    }
    ```
- Response:

    Header: `HTTP 200`
    Body:
     ```
      {
          "id": <order_id>,
          "distance": <total_distance>,
          "status": "UNASSIGNED"
      }
     ```
    or

    Header: `HTTP <HTTP_CODE>`
    Body:

     ```
      {
          "error": "ERROR_DESCRIPTION"
      }
     ```

### Take order
- Method: `PATCH`
- Path : `/orders/:id`
- Request body:
	```
    {
        "status": "TAKEN"
    }
    ```
- Response:
    Header: `HTTP 200`
    Body:
    ```
      {
          "status": "SUCCESS"
      }
    ```
    or

    Header: `HTTP <HTTP_CODE>`
    Body:
    ```
      {
          "error": "ERROR_DESCRIPTION"
      }
    ```

### List order
- Method: `GET`
- Path : `/orders?page=:page&limit=:limit`
- Response:
    Header: `HTTP 200`
    Body:
    ```
      [
          {
              "id": <order_id>,
              "distance": <total_distance>,
              "status": <ORDER_STATUS>
          },
          ...
      ]
    ```

    or

    Header: `HTTP <HTTP_CODE>` Body:

    ```
    {
        "error": "ERROR_DESCRIPTION"
    }
    ```


## Required Packages
- Database
    * [MySql](https://github.com/go-sql-driver/mysql)
- Routing
    * [mux](https://github.com/gorilla/mux)

## Quick Run Project
First clone the repo then go to go-mysql-docker/github.com/baronfung folder. After that build your image and run by docker. Make sure you have docker in your machine. 

```
git clone https://github.com/baronfungkk/go-mysql-docker

cd go-mysql-crud/github.com/baronfung/order

chmod +x start.sh
./start.sh

```

## Google Map API key
You will need a key for Google Map API in order to use their services.

To use your own Google Map API Key, go to /app/app.go change the variable APIKey to your key.

```
var APIKey = <API_Key>>

```