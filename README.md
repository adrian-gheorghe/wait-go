## wait-go

[![CircleCI](https://circleci.com/gh/adrian-gheorghe/wait-go.svg?style=svg)](https://circleci.com/gh/adrian-gheorghe/wait-go)

`wait-go` is the golang rewrite of https://github.com/adrian-gheorghe/wait
The script waits for a host or multiple hosts to respond on a TCP port but can also wait for a command to output a value. For example you can wait for a file to exist or contain something.

The script is mainly useful to link containers that dependend on one another to start. For example you can have a container that runs install scripts that will have to wait for the database to be accessible.

## Download

Download latest from the releases page: https://github.com/adrian-gheorghe/wait-go/releases

## Usage

```
wait-go --help
 -command value
    	Command that should be run when all waits are accessible. Multiple commands can be added.
  -interval int
    	Interval between calls (default 15)
  -timeout int
    	Timeout untill script is killed. (default 600)
  -version
    Prints current version
  -wait value
    	You can specify the HOST and TCP PORT using the format HOST:PORT, or you can specify a command that should return an output. Multiple wait flags can be added.
```

## Examples shell

```
$ wait-go --wait "database_host:3306" --wait "ls -al /var/www/html | grep docker-compose.yml" --command "Database is up and files exist"
$ wait-go --wait "database_host:3306" --wait "database_host2:3306" --command "echo \"Databases are up\""
```

You can set your own timeout with the `-t` or `--timeout=` option.  Setting the timeout value to 0 will disable the delay between requests:

```
$ wait-go --wait "database_host:3306" --wait "database_host2:3306" --command "echo \"Databases are up\"" --timeout 15
```
## Examples docker-compose

```
version: '3.3'
services:
  db:
    image: mysql:5.7
    deploy:
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: database
  wait:
    build:
      context: .
    command: "wait-go --wait \"db:3306\" --command \"ls -al\""
    
```

## Example docker multiple FROM

In the following example the Dockerfile adds the wait-go file from the adighe/wait-go container. 
The setup allows the running of database migrations only after the database is accessible and the volume is mounted

### Dockerfile
```

FROM adighe/wait-go as wait-go
FROM php:7.1.3-fpm

RUN curl -sS https://getcomposer.org/installer | php && \
    mv composer.phar /usr/local/bin/composer

COPY --from=wait-go /app/wait-go /app/wait-go

ENTRYPOINT ["docker-php-entrypoint"]
CMD ["php-fpm"]
    
```
### docker-compose.yml
```
version: '3.3'
services:
  db:
    image: mysql:5.7
    deploy:
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: database
  install:
      build:
        context: .
      command: "/app/wait-go --wait 'db:3306' --wait 'ls -al /var/www/html/ | grep composer.json' --command 'cd /var/www/html' --command 'ls -al' --command 'composer install' --command 'php /var/www/html/bin/console doctrine:migrations:migrate -n -vvv'"

```