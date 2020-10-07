Start MySQL server on Docker.

```
$ docker volume create mysql_vol1

$ docker run --name mysql-test1 -v mysql_vol1:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=abcd1234 -p 3306:3306 -d mysql:5.7.31
```

Setup MySQL's databases and users.

```
$ docker exec -it mysql-test1 bash -p

# mysql -u root -p -h 127.0.0.1

> create database foo;
> create database bar;

> create user myuser identified by 'abcd1234';

> grant all on foo.* to 'myuser';
> grant all on bar.* to 'myuser';
```
