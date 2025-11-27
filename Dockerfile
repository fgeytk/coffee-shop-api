FROM mysql:8.0

# Copier le script d'initialisation
COPY init.sql /docker-entrypoint-initdb.d/

# Variables d'environnement par défaut
ENV MYSQL_ROOT_PASSWORD=root123
ENV MYSQL_DATABASE=coffee_shop
ENV MYSQL_USER=coffee_user
ENV MYSQL_PASSWORD=coffee123

EXPOSE 3306
