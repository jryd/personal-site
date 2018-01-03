# My Personal Website

This is the source for my personal website, currently hosted at https://bannister.me.

## Installation

This assumes you have your go environment setup correctly, if you don't then you'll need to do so before following the steps below.

1. `go get github.com/jryd/personal-site` this will copy everything you need into your `GOPATH/src/github.com/jryd/personal-site`
2. Switch into the newly created directory and run `cp .env.example .env` and replace the values contained within. The application uses the Mailgun API for sending contact messages and the PRODUCTION boolen to determine whether you are developing locally or not
3. Run `go build` to build the executable to run the site

To ensure the `.env` file can always be sourced correctly, please ensure you are running the executable from the same folder as that of the `personal-site` executable.

## Deployment

There are two options for deployment, as an application on a custom port or as a standalone web server.

### Application on Custom Port

To provision this, just run the installation steps above, and be sure to set the `SERVER_PORT` to a port that isn't in use.

If you are using Apache, you may wish to configure a proxy to pass requests on port 80 and 443 to our application. An example configuration is below:

```
<VirtualHost *:80>
    ServerName bannister.me
    ServerAlias www.bannister.me

    ProxyPreserveHost On
    ProxyRequests Off

    ProxyPass / http://localhost:8080/
    ProxyPassReverse / http://localhost:8080/

</VirtualHost>
```

### Standalone Web Server

If you aren't going to use the server to run any other applications, then you can spin this up as a standalone web server. By this I mean it will bind to port 80 and 443.

When running this in production, be sure to set the `PRODUCTION` boolean to true. This will ensure the application starts on port 443 too (as opposed to just port 80). It's important to set `DOMAIN_NAME` correctly and to ensure your DNS is setup to point your domain to the IP address of the server. This is because when the application runs in production, it will automatically provision SSL certificates through Let's Encrypt, and to do so it needs to be able to resolve the host back to the server it is being hosted on.