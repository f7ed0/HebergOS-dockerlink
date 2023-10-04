# HebergOS docker link

## .env

fields required :

* `docker_sock` : path to your docker socket

## API endpoints

### `GET` v1/container/

Retrieve information about containers

#### params

* (optional multiple) `id` : list of the container ids to get
* (optional) `running` : set to true to get only running container (overridden by id)

#### Response

type : `application/json`

```
{
    "[container_id]" : {
        "name" : "[name]",
        "state": "[running|exited|paused|restarting|removing|dead]
		"exit_code" : [value if exited]
		"started_at" : [unix timestamp if running]
    },
	...
}
```

---

### `GET` v1/container/stats

Retrieve the timestamped stat of the container with id `id` since `since`

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

---

### `POST` v1/container/start

Start the container of id `id`

#### params

* `id` : the container id that you want to start

#### Response

`plain\text` or `no content`

---

### `POST` v1/container/stop

Stop the container of id `id`

#### params

* `id` : the container id that you want to stop

#### Response

`plain\text` or `no content`