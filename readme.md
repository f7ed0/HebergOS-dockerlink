# HebergOS docker link

## .env

fields required :

* `docker_sock` : path to your docker socket

## API endpoints

### v1/container/list

Retrieve the list of all containers

#### Method

`GET`

#### params

no params

#### Response

type : `application/json`

```
{
    "[container_id]" : {
        "names" : [
            "name",
            "alias_1",
            ...
        ],
        "state": "running"|"exited"|...
    }
}
```

### v1/container/stats

Retrieve the timestamped stat of the container with id `id` since `since`

#### Method

`GET`

#### params

* `id` : the container id that you want to retrieve the stat
* `since` : Unix seconds timestamp

#### Response

type : `application\json`

```
{
  "[unix timestamp of the pull]" : {
	"memory" : {
	  "used" : [amount of memory used in Go],
	  "limit" : [limit of memory of the container in Go]
	},
	"cpu" : {
	  "usage_percent" : [cpu usage in percent (100% = 1 core)]
	},
	"net" : {
	  "up" : [upload since launch in ko],
	  "down" : [download since launch in ko],
	  "delta_up": [upload since last pull (10 seconds) in ko],
	  "delta_down": [download since last pull (10 seconds) in ko],
	}
  },
  ...
}
```