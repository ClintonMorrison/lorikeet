# Lorikeet

Website: [https://lorikeet.clintonmorrison.com/](https://lorikeet.clintonmorrison.com/)

Lorikeet is a password manager web app, created with Golang, React and Materialize CSS. It uses client-side and server-side AES encryption.

## TODO:
- update frontend libs and refactor UI
- show UI validation error when registering with invalid chars
- test everything

## Running Lorikeet

### Running the Go server

The server runs on port 8080.
```
cd ./lorikeet/
go run cmd/server/main.go
```

### Running webpack
For frontend development:
```
cd ./lorikeet/ui
npm i
npm run start
```

The webpack dev server runs on port 3000.

### Building for production
```
# Server
cd ./lorikeet/
sh build.sh

# Frontend
cd ./lorikeet/ui
npm run build

# Building docker container
sh build.sh
```

### Creating a backup
A backup copy of the encrypted password data can be created with the backup script.

It copies the data in `./lorikeet/data` into a tarball in `./lorikeet/backup/`.

```
cd ./lorikeet/
go build -o doBackup cmd/backup/main.go &&
./doBackup
```

### Restoring from a backup
In the event of data loss or corruption, a previous snapshot of the data can be restored with the restore script.

The script takes the path to a tarball (created with the backup script) and copies it to the data folder. 
Any pre-existing data in the data folder  is moved to `./data-old-TIMESTAMP`.

```
cd ./lorikeet/
go build -o doRestore cmd/restore/main.go
./doRestore -file ./lorikeet/backup/backup-TIMESTAMP.tar.gz
```

### Serving with NGINX
The built frontend app and API can be served with NGINX.

```
http {
  server {
    server_name  lorikeet.clintonmorrison.com;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    location ~ ^/api/ {
      proxy_pass http://0.0.0.0:8080;
    }

    location ~ / {
      root /pathToRepo/lorikeet/ui/build/;
      rewrite ^$ /index.html;
      try_files $uri /index.html =400;
    }
  }


  # Or for local development:
  server {
    server_name  lorikeet.dev.ca;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    location ~ ^/api/ {
            proxy_pass http://0.0.0.0:8080;
    }

    location ~ / {
            proxy_pass http://0.0.0.0:3000;
    }
  }

}
```


## Architecture
The frontend is a single page React app. It makes AJAX requests to the Go server to handle authentication and CRUD operations.

When users login their hashed credentials are stored in browser session storage, and passed with all AJAX requests.

Before sending passwords to the server the frontend encrypts them with AES using a secret token derived from the user's credentials.
The server never sees this token so cannot decrypt the passwords.

When the server receives the passwords it encrypts them again using another derived token, and a user-specific salt.
The passwords are saved in a text file. The name of the text file is also derived from the token.

When retrieving passwords, the server determines which file contains the passwords, and decrypts them using the credentials in the AJAX request.
The passwords are still partially encrypted (from the frontend encryption). The frontend decrypts them fully using the secret user credentials which it has not shared with the server.

The result is that the passwords are stored securely, and it's not possible to decrypt them server-side without the user's original credentials.
