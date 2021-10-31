CREATE TABLE IF NOT EXISTS tbl_survivor(
    id INT(11) NOT NULL UNIQUE AUTO_INCREMENT,
    name VARCHAR (100) DEFAULT '',
    age  VARCHAR (100) DEFAULT '',
    gender  VARCHAR (50) DEFAULT '',
    lastLocation POINT NOT NULL,
    createdOn TIMESTAMP NOT NULL default current_timestamp,
    SPATIAL INDEX `SPATIAL` (`lastLocation`)
    PRIMARY KEY (id)
    )ENGINE=INNODB DEFAULT CHARSET=utf8;