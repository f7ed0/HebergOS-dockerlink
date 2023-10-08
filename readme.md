# HebergOS docker link

## .env

fields required :

* `docker_sock` : path to your docker socket
* `port_area_size` : size of the port area for each docker (recommended 5000)
* `base_image` : id of the image you use by default to run a new docker container

## API endpoints

### 

---
---

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
		"host_port_root": [root of port on host],
        "name" : "[name]",
		"ports": {
			"[port]/[tcp|udp]": {},
			...
		},
        "state": "[running|exited|paused|restarting|removing|dead|created]
		"exit_code" : [value if exited]
		"started_at" : [unix timestamp if running]
    },
	...
}
```

---

### `PUT` v1/container

creates a new container for HerbergOS usage (creates a volume on www/var)

(port 22,80 and 443 are forwarded by default)

#### body

`json`

```
{
	"name" : [name of the container],
	"host_port_root" : [number basis for port forwarding],
	"memory" : [limit of memory in Go (can be decimal)],
	"cpulimit" : [limit in number of cpu (can be decimal)],
	"ports"(optional) : [
		"[port]/[tcp|udp]",
		...
	]
	"image"(optional) : [image id],
	"commands"(optional) : [
		[command to launch at start]
	]
}
```

#### Response

`plain\text` on error

`application\json` on success

```
{
	"Id" : [id of the new container]
	"warning" : [
		[warnings coming from docker],
		...
	]
}
```

---
### `DELETE` v1/container

Delete a container by id

#### params

* id

#### Response

`plain\text`

---

### `PATCH` v1/container

Update memory and cpulimit

#### params

* id

#### body

```
{
	"cpulimit"(optional) : [limit in number of cpu (can be decimal)],
	"memory"(optional): [limit in Go(can be decimal)],
}
```

---

### `GET` v1/container/stats

Retrieve the timestamped stat of the container with id `id` since `since`

#### params

* `id` : the container id that you want to retrieve the stat
* `since` : Unix seconds timestamp (0 for all stats)

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
	  "limit" : [cpu limit usage in percent (100% = 1 core) (can be NaN)]
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

----
----

### `GET` v1/git

execute `git pull` on the git local repository

#### params

* `id` : id of the container
* (optional) `path` : indicates the path of the git repo

#### response

`plain\text` : git output

----

### `PUT` v1/git

execute `git init` on the git local repository

#### params

* `id` : id of the container
* (optional) `path` : indicates the path of the git repo

#### response

on error : `plain\text`

on success : no content

----

### `GET` v1/git/head

returns information about the git head

#### params

* `id` : id of the container
* (optional) `path` : indicates the path of the git repo

#### response

`plain\text` : git output

----

### `GET` v1/git/branches

fetch all the branch of a given directory if the directory is a git repository
(distant branches are also fetched)

#### params

* `id` : id of the container
* (optional) `path` : indicates the path of the git repo

#### response

`plain\text` on error, 
`application\json` on success

```
[
	"[branch 1]",
	"[branch 2]",
	...,
]
```

---

### `GET` v1/git/branch

Get the name of the current branch

#### params

* `id` : id of the container
* (optional) `path` : indicates the path of the git repo
* `branch` : name of the branch

#### response

`plain\text` : git output

---

### `POST` v1/git/branch

Checkout to a new branch (if possible)

#### params

* `id` : id of the container
* (optional) `path` : indicates the path of the git repo
* `branch` : name of the branch

#### response

`plain\text` : git output
