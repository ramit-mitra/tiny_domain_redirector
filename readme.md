# Tiny Domain Redirector
## tiny_domain_redirector

A simple Golang web application that can be used to redirect HTTP requests for a domain (or subdomain) to another domain (or subdomain or URL). The application uses the `http` package in Go to create an HTTP server that receives HTTP requests and redirects the client to the specified destination URL. The application can be used to redirect domains for a variety of purposes, such as migrating a website to a new domain, redirecting a domain to a new subdomain, or redirecting a domain to a different URL.

Built as a personal project. Please review the code before using. Suggestions and contributions are welcome.

Thank you

## Run as a service

- Create a `Systemd` service. Example below:

```
[Unit]
Description      = redirector

[Service]
Type             = simple
User             = root
Group            = root
ExecStart        = /path/to/this/binary
WorkingDirectory = /path/to/this/golang/project
Restart          = always
RestartSec       = 5s

[Install]
WantedBy         = multi-user.target
```

- Create the file in `/lib/systemd/system/` (usually)

- Save the service file as `redirector.service`

- Enable the service: `sudo systemctl enable redirector.service`

- Start the service: `sudo systemctl start redirector.service`

- Check the status of the service: `sudo systemctl status redirector.service`

## Setup as nginx reverse proxy (example)

Filename: `/etc/nginx/sites-enabled/contact.ramit.io.conf`

```
server {
    listen 443 ssl http2;
    server_name contact.ramit.io;
    client_max_body_size 128k;

    # logging
    access_log /var/log/nginx/redirector.access.log;
    error_log /var/log/nginx/redirector.error.log warn;

    # reverse proxy
    location / {
        proxy_pass http://127.0.0.1:9990;
    }
    
    # ...
}
```

## Test in local 

- Run the app
- Use curl to verify if the redirect works. Example: `curl -skIXGET --connect-to ::localhost:9990 contact.ramit.io` 

## Motivation

Read more [here](https://ramit-mitra.medium.com/about-tiny-domain-redirector-bb943c72fd7a)