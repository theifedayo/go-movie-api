# go-movie-api
This is a REST API built using Golang, GORM, Gin, and Postgres that provides information about Star Wars movies. It includes endpoints to list movies with their opening crawls and comment count, add new comments to a movie, list comments for a movie, and get a list of characters for a movie.
### Installation
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
Add a .env file following .env.sample file
```.env
POSTGRES_HOST=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_PORT=


REDIS_ADDRESS=
REDIS_PASSWORD=
REDIS_DB=

PORT=8000
CLIENT_ORIGIN=
```
Run the following command to start the server:
```bash
go run main.go
```
The server will run on http://localhost:8000 by default

### As a docker service
Ensure, you have docker installed already.
After adding .env file with the correct environment variables. Build the app as an image and name it `go-movie-api`. Run the command below in your go-movie-api directory
```bash
docker build -t go-movie-api . 
```
Verify it has been built with
```bash
docker images
```
and go through the list of images to find `go-movie-api`
Now, we spin up redis and postgres containers with our go-movie-api as a service with
```bash
docker-compose up
```
The server will run on http://0.0.0.0:8000 by default

## Available Endpoints
List Movies
bash


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
Add Comment
bash


POST /movies/:id/comments
Adds a new comment to the specified movie.
Request Parameters
* id (int, required) - The ID of the movie to add a comment to.
Request Body
* comment (string, required) - The text of the comment.
Example request body:
json


{ "comment": "This is a great movie!" }
Response
* 201 Created on success
List Comments
bash


GET /movies/:id/comments
Returns a list of comments for the specified movie.
Request Parameters
* id (int, required) - The ID of the movie to list comments for.
Response
* 200 OK on success
json


[ { "id": 1, "movie_id": 1, "comment": "This movie is awesome!" }, { "id": 2, "movie_id": 1, "comment": "I love this movie!" }, ... ]
Get Characters
bash


GET /movies/:id/characters
Returns a list of characters for the specified movie.
Request Parameters
* id (int, required) - The ID of the movie to get characters for.
Response
* 200 OK on success
json


[ "Luke Skywalker", "Han Solo", "Princess Leia", ... ]
Documentation
You can find additional documentation for this API, including request and response signatures, by visiting http://localhost:8000/swagger/index.html in your web browser.