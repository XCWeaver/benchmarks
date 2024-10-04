import redis


client = redis.StrictRedis(host='34.65.41.16', port=6381, db=0, password='')

filename = 'posIds.txt'

with open(filename, 'r') as f:
    keys = [line.strip() for line in f] 
if keys:
    client.delete(*keys)
    print(f"Keys read from {filename} and deleted from Redis.")
else:
    print("No keys found in the file.")