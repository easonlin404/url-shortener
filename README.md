# URL Shortener

This is a URL shortener service implemented in Go using the Gin framework, MongoDB, and Redis.

## Prerequisites

- Docker
- Docker Compose

## Running Locally

1. Clone the repository:

   ```sh
   git clone https://github.com/easonlin404/url-shortener.git
   cd url-shortener
    ```
2. Build and run the Docker container:  
    ```sh
    docker build -t url-shortener .
    docker run -p 8080:8080 -p 27017:27017 -p 6379:6379 url-shortener
    ```

The service will be available at http://localhost:8080.   

## API Endpoints

### Upload URL
- **URL:** `/api/v1/urls`
- **Method:** `POST`
- **Request Body:**
  ```json
  {
    "url": "<original_url>",
    "expireAt": "2021-02-08T09:20:41Z"
  }
   ```
- **Response:**
  ```json
  {
   "id": "<url_id>",
   "shortUrl": "http://localhost:8080/<url_id>"
   }
  ```


## Note
- why choose MongoDB as a database?
    - Because it is a document-based database and a good fit for storing JSON-like documents, we don't need to handle relationships in the case. It is easy to scale horizontally and can hold much data and traffic.
- Why choose the snowflake algorithm to generate the unique ID for the short URL?
  - It's easy to generate the globally unique ID, and it is a distributed unique ID generator that ensures cardinality to some extent.
  - Although MongoDB already guarantees the unique for '_id' if we don't specify, we here choose the snowflake algorithm to generate the unique id for the short URL because we can have better control over the id generation.
- Redis is used to cache the short URL and its original URL to reduce query to DB and speed up query.
- - Haven't had enough time to write unit tests for endpoints; need to mock some dependencies like MongoDB and Redis.




### Redirect URL
- **URL:** `/:id`
- **Method:** `GET`
- **Response:** Redirects to the original URL or returns 404 if the URL is expired or not found.
