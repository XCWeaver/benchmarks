CHANGE MASTER TO
    MASTER_HOST='mysql1',
    MASTER_USER='repluser',
    MASTER_PASSWORD='replpassword',
    MASTER_LOG_FILE='{{MASTER_LOG_FILE}}',
    MASTER_LOG_POS={{MASTER_LOG_POS}};
START SLAVE;