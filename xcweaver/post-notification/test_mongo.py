from pymongo import MongoClient

client = MongoClient('mongodb://34.65.41.16:27017/?directConnection=true')
db = client['post-storage'] 
collection = db['posts'] 

result = collection.delete_many({})

# Print how many documents were deleted
print(f"{result.deleted_count} documents deleted.")

client.close()