FROM php:8.0-apache-buster

# docker run -e SUBPATH=docker/php/src -e GIT_URL=https://github.com/haisum/k99s.git \
# -e DB_HOST=host.docker.internal -e DB_USER=root -e DB_PASSWORD='' -e DB_NAME=cdb_api_dev \
# -e APP_URL=php.k99s.com -d -p 8080:80  localhost:5000/php:1.0

RUN apt update && apt install -y git 
RUN docker-php-ext-configure mysqli && \
    docker-php-ext-install mysqli && \
    docker-php-ext-enable mysqli

EXPOSE 80

ENTRYPOINT mkdir tmp && git clone ${GIT_URL} tmp && \
    cp -r tmp/${SUBPATH}/* /var/www/html/ && rm -rf tmp && \
    docker-php-entrypoint apache2-foreground