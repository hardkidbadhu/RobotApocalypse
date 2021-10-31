CREATE TABLE IF NOT EXISTS tbl_survivor_inventory (
    ID INT(11) NOT NULL UNIQUE AUTO_INCREMENT,
    Water INT(11) NOT NULL default 0,
    Food INT(11) NOT NULL default 0,
    Medication INT(11) NOT NULL default 0,
    Ammunition INT(11) NOT NULL default 0,
    UserID INT(11) NOT NULL,
    FOREIGN KEY (UserID) REFERENCES tbl_survivor(id) ,
    PRIMARY KEY (ID)
    )ENGINE=INNODB DEFAULT CHARSET=utf8;