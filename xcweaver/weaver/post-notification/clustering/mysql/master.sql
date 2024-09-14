CREATE USER IF NOT EXISTS 'repluser'@'%' IDENTIFIED WITH 'mysql_native_password' BY 'replpassword';
GRANT REPLICATION SLAVE ON *.* TO 'repluser'@'%';
FLUSH PRIVILEGES;