package consts

const NGINX_TEMPLATE string = `
server {
	server_name %s.insash.fr;
	location / {
		proxy_pass http://localhost:%v;
	}

    listen 80;
}
server {
	server_name ssh.%s.insash.fr;
        location / {
                proxy_pass http://localhost:%v;
        }

    listen 80;
}
`