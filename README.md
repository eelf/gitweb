

yarn add webpack webpack-cli\
    react react-dom @types/react-dom\
    css-loader style-loader\
    ts-loader typescript\
    grpc-web-client google-protobuf @improbable-eng/grpc-web ts-protoc-gen @types/google-protobuf


create user gitweb@localhost identified by 'gitweb';

grant all on gitweb.* to gitweb@localhost;

create database gitweb;

CREATE TABLE `user` (
`id` int unsigned NOT NULL AUTO_INCREMENT,
`name` varchar(255) NOT NULL,
`key` text NOT NULL,
PRIMARY KEY (`id`)
)


CREATE TABLE `access` (
`user_id` int unsigned NOT NULL,
`mode` enum('Read','Write') NOT NULL DEFAULT 'Read',
`repo` varchar(255) NOT NULL
)

insert into gitweb.user set name = 'name', `key` = 'ssh-rsa AAAA';

insert into gitweb.access (user_id, mode, repo) values (1, 'Write', '*');


