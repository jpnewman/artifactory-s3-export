
# Artifactory S3 Export Testing

A MySQL dump of the Artifactory database can be imported into a docker container for testing.

## Create MySQL docker container

- <https://hub.docker.com/_/mysql>

~~~
mkdir -p mysql

docker run -d --rm --name=art-mysql -v ${PWD}/mysql:/var/lib/mysql -e MYSQL_ROOT_PASSWORD="T/*V@:L2edb5Yb+e" -e MYSQL_DATABASE=artdb -p 3306:3306 mysql:8.0.15 --default-authentication-plugin=mysql_native_password
~~~

## Docker ps

~~~
docker ps -a --no-trunc
~~~

## MySQL logs

~~~
docker logs art-mysql
~~~

## Connect to MySQL

~~~
mysql -h 127.0.0.1 -P 3306 -uroot -p"T/*V@:L2edb5Yb+e"
~~~

## MySQL Bash

~~~
docker exec -it art-mysql bash
~~~

### MySQL commands

~~~
mysql -uroot -p"T/*V@:L2edb5Yb+e"
~~~

## List databases

~~~
docker exec -i art-mysql /usr/bin/mysql -u root --password="T/*V@:L2edb5Yb+e" -e "SHOW DATABASES;"
~~~

## List tables

~~~
docker exec -i art-mysql /usr/bin/mysql -u root --password="T/*V@:L2edb5Yb+e" -e "USE artdb;SHOW TABLES;"
~~~

## List Artifactory Nodes

~~~
docker exec -i art-mysql /usr/bin/mysql -u root --password="T/*V@:L2edb5Yb+e" -e "SELECT node_id,repo,node_name FROM artdb.nodes WHERE node_name <> '.' AND node_path = '.' AND repo = '3rd-party-libs' LIMIT 10;"
~~~

## Import data

- <https://gist.github.com/spalladino/6d981f7b33f6e0afe6bb>

~~~
cat data/artifactory-01_artdb_2019-03-27_14-02.sql | docker exec -i art-mysql /usr/bin/mysql -u root --password="T/*V@:L2edb5Yb+e" artdb
~~~

> Importing data can take about 20 minutes.

## Docker stop and remove

~~~
docker stop art-mysql
~~~
