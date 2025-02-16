## MSDS4343 Module 7: api-experiment-go

A demonstration of how a containerized RESTful API might be implemented in Go.

### Replication

```shell
# Clone and enter the repo...
git clone git@github.com:kgeidel/api-experiment-go.git && cd api-experiment-go

# Add required ENV vars (warning! will wipe an existing .env in this dir)
echo DB_NAME=goapi > .env
echo DB_USER=[YOUR DB USER NAME] >> .env
echo DB_PASSWORD=[YOUR DB PASSWORD] >> .env
echo DB_HOST=localhost >> .env
echo DB_PORT=5432 >> .env

# Initialize the DB (requires postgres or pg based RDS)
psql -f db_init.sql
```

```shell
# Launch the go container with...
docker compose up

# ---OR---

# run the go package manually

# Install Go dependencies to root folder
cd api-experiment-go
go get github.com/gorilla/mux
go get github.com/joho/godotenv
go get github.com/lib/pq

# Run the go package
go run main.go
```

### Troubleshooting

* Is Go properly installed?
* Is docker (and docker compose) properly installed?
* Is postgres properly installed and configured on the host listed in `.env`?
* Are your DB credentials correct? (`cat .env` to verify.)
* Does your firewall allow inbound traffic on port 8000?

### Testing the API

```python
# With api-experiment-go (either Go package or docker compose) running...
# Send GET and POST requests to the API
# We're using Python's requests module here but any method of HTTP request will do (browser, curl, postman, etc) 

import requests, os, json
from pprint import pprint

api_host = 'http://localhost:8000'

r = requests.get(f'{api_host}')

pprint(json.loads(r.content))
```

This should print the contents of the `product` table (which may be empty if this is the first time running this.)

```python

# Use post request to add records.

data = {
    'name': 'Pixel', 
    'description': 'Made by Google but works with Graphene OS',
    'price': 399.99,
    'available_flag': True
}

r = requests.post(f'{api_host}', json=data,)

pprint(r.content)
```

You should get a confirmation with the new record's pk. Rerunning the GET request above should show the new entry.