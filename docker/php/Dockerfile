FROM php:8.0-apache-buster

RUN apt update && apt install -y git 
RUN docker-php-ext-configure mysqli && \
    docker-php-ext-install mysqli && \
    docker-php-ext-enable mysqli

EXPOSE 80

ENTRYPOINT mkdir tmp && git clone ${GIT_URL} tmp && \
    cp -r tmp/${SUBPATH}/* /var/www/html/ && rm -rf tmp && \
    docker-php-entrypoint apache2-foreground