USE mysql;

CREATE USER user@'%' IDENTIFIED BY 'userpass';
GRANT ALL PRIVILEGES ON *.* TO user@'%';

CREATE USER user@'localhost' IDENTIFIED BY 'userpass';
GRANT ALL PRIVILEGES ON *.* TO user@'localhost';

CREATE USER viewer@'%' IDENTIFIED BY 'viewerpass';
GRANT SELECT ON *.* TO viewer@'%';

CREATE USER viewer@'localhost' IDENTIFIED BY 'viewerpass';
GRANT SELECT ON *.* TO viewer@'localhost';