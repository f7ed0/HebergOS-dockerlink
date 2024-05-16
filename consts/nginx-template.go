package consts

const NGINX_TEMPLATE string = `
server {
	server_name %s.insash.org;
	location / {
		proxy_pass http://localhost:%v;
	}

    listen 80;
}
server {
	server_name ssh.%s.insash.org;
        location / {
                proxy_pass http://localhost:%v;
        }

    listen 80;
}
`
