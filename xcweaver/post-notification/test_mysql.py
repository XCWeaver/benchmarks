import mysql.connector

# Establish a connection to the MySQL database
connection = mysql.connector.connect(
    host="34.65.60.146",
    user="root",
    password="root_password",
    database="post_notification",
    port=3307
)

cursor = connection.cursor()

# Define the SQL query to delete all rows from the table
delete_query = "DELETE FROM posts"

# Execute the query
cursor.execute(delete_query)

# Commit the changes
connection.commit()

print(f"{cursor.rowcount} rows were deleted.")

# Close the cursor and connection
cursor.close()
connection.close()