# go-movie-api
### Introduction
This is a RESTful API built using Golang, GORM, Gin, Redis and Postgres that provides information about Star Wars movies. It includes endpoints to list movies with their opening crawls and comment count, add new comments to a movie, list comments for a movie, and get a list of characters for a movie.

</br>

### Setup
Clone the repository to your local machine.
```bash
git clone https://github.com/theifedayo/go-movie-api.git
```
Ensure that you have Golang and Postgres installed on your machine.
Navigate to the root directory of the project in a terminal.
```bash
cd go-movie-api
```
Run the following command to install the necessary dependencies
```bash
go get ./... 
```
Add a .env file following sample.env file example with the values of each variable
```.env
POSTGRES_HOST=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_PORT=


REDIS_ADDRESS=
REDIS_PASSWORD=
REDIS_DB=

PORT=8080
CLIENT_ORIGIN=
```

</br>

### Running Server
#### Locally
This project depends on redis to cache list of movies. In order to have everything run well, you have to run redis server
```bash
redis-server
```
Run the following command to start the server:
```bash
go run main.go
```
The server will run on http://localhost:8080 by default

</br>

#### As a docker service
Ensure, you have docker installed already.
After adding .env file with the correct environment variables. Build the app as an image and name it `go-movie-api`. Run the command below in your go-movie-api directory
```bash
docker build -t go-movie-api . 
```
Verify it has been built with
```bash
docker images
```
and go through the list of images to find `go-movie-api`.
Now, we spin up 3 containers (redis, postgres and go-movie-api) as a service with
```bash
docker-compose up
```
The server will run on http://0.0.0.0:8080 by default

</br>

## Available Endpoints
### List Movies
#### GET api/v1/movies
Returns a list of movies with their opening crawls and comment count. The list is sorted by release date from earliest to newest.\
Request Parameters\
None.\
Response
* 200 OK on success

```json
{
    "status": "success",
    "data": [
        {
            "name": "The Empire Strikes Back",
            "opening_crawl": "It is a dark time for the\r\nRebellion. Although the Death\r\nStar has been ..",
            "comment_count": 1
        },
        {
            "name": "Return of the Jedi",
            "opening_crawl": "Luke Skywalker has returned to\r\nhis home planet of Tatooine in\r\nan attend...",
            "comment_count": 1
        },
        {
            "name": "The Phantom Menace",
            "opening_crawl": "Turmoil has engulfed the\r\nGalactic Republic. The taxation\r\nof trade ....",
            "comment_count": 8
        }
    ]
}
```

### Add Comment
#### POST api/v1/movies/:movieId/comments
Adds a new comment to the specified movie.\
Request Parameters
* `movieId` (int, required) - The ID of the movie to add a comment to.
Request Body
* `comment`(string, required) - The text of the comment.

Example request body:
```json
{ "comment": "This is a great movie! I would rate it 9/10" }
```
Response
* 201 Created on success
```json
{
    "data": {
        "id": 18,
        "movie_id": "3",
        "comment": "This is a great movie! I would rate it 9/10",
        "ip": "fe80::aede:48ff:fe00:1234",
        "created_at": "2023-03-05T00:39:02Z",
        "updated_at": "2023-03-05T00:39:02Z"
    },
    "status": "success"
}
```
If the specified movie with the movieId is not found, the HTTP status code in the response header is `404 Not Found`.

### List Comments
#### GET api/v1/movies/:movieId/comments
Returns a list of comments for the specified movie.\
Request Parameters
* `movieId` (int, required) - The ID of the movie to list comments for.
Response
* 200 OK on success
```json
{
    "data": [
        {
            "id": 13,
            "movie_id": "4",
            "comment": "what an interesting movie. I would rate it 9/10",
            "ip": "fe80::aede:48ff:fe00:1122",
            "created_at": "2023-03-04T06:15:25+01:00",
            "updated_at": "2023-03-04T06:15:25+01:00"
        },
        {
            "id": 12,
            "movie_id": "4",
            "comment": "what an interesting movie. I would rate it 9/10",
            "ip": "fe80::aede:48ff:fe00:1122",
            "created_at": "2023-03-04T06:12:44.126034+01:00",
            "updated_at": "2023-03-04T06:12:44.126034+01:00"
        },
        {
            "id": 11,
            "movie_id": "4",
            "comment": "what an interesting movie. I would rate it 6/10",
            "ip": "fe80::aede:48ff:fe00:1234",
            "created_at": "2023-03-04T06:04:14.673369+01:00",
            "updated_at": "2023-03-04T06:04:14.673369+01:00"
        },
    ],
    "results": 3,
    "status": "success"
}
```
If the specified movie with the movieId is not found, the HTTP status code in the response header is `404 Not Found`.

### Get Characters
#### GET api/v1/movies/:movieId/characters
Returns a list of characters for the specified movie.\
Request Parameters
* `movieId` (int, required) - The Id of the movie to get characters for.
Response
* `sort`(optional): The field to sort the characters by one of name, gender, or height.
* `order` (optional) - Use  `asc` or `desc` to sort in ascending or descending order, respectively. For example, ?sort=height&order=desc will sort by height in descending order, while ?sort=height&order=asc will sort by height in ascending order.
* `gender`(optional) - The filter criteria to apply to the characters to filter by male or female. For example, ?gender=male will filter by male characters and return only male characters and \?sort=height&order=desc&gender=female will filter by female characters, listing only female characters with their height in descending order
* 200 OK on success\
`api/v1/movies/:movieId/characters`
```json
{
    "data": [
        {
            "name": "Luke Skywalker",
            "height": "172",
            "gender": "male"
        },
        {
            "name": "C-3PO",
            "height": "167",
            "gender": "n/a"
        },
        {
            "name": "R2-D2",
            "height": "96",
            "gender": "n/a"
        },
        {
            "name": "Darth Vader",
            "height": "202",
            "gender": "male"
        },
        {
            "name": "Leia Organa",
            "height": "150",
            "gender": "female"
        },
        ....
        ....
    ],
    "metadata": {
        "total_count": 18,
        "total_height": {
            "cm": 3066,
            "feet": 100.59,
            "inch": 1207.08
        }
    },
    "status": "success"
}
```
`api/v1/movies/:movieId/characters?sort=height&order=desc&gender=female`
```json
{
    "data": [
        {
            "name": "Beru Whitesun lars",
            "height": "165",
            "gender": "female"
        },
        {
            "name": "Leia Organa",
            "height": "150",
            "gender": "female"
        }
    ],
    "metadata": {
        "total_count": 2,
        "total_height": {
            "cm": 315,
            "feet": 4.06,
            "inch": 124.01
        }
    },
    "status": "success"
}
```

If the specified movieId is not found, the HTTP status code in the response header is `404 Not Found`.
</br>

### Errors
The response for request failures or any other error are rather simple.
```json
{
	"status": "error",
	"message": "The error message"
}
```

</br>

### Conclusion
You can find additional documentation for this API, including request and response signatures, by visiting http://0.0.0.0:8080/api/v1/docs/index.html if running on your local server, or http://gomovie-api.herokuapp.com/api/v1/docs/index.html in your web browser.