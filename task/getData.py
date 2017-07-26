import  pymongo

client = pymongo.MongoClient('mongodb://mongodb.t2.daoapp.io:61257/')
db = client['task']
posts = db['result']


servers = posts.aggregate([{"$group" : {"_id" : "$status", "num_tutorial" : {"$sum":1 }}}])
for server in servers:
    print server


res = posts.find({"url":{'$regex':'.'}})
for i in res:
    print i["url"]

print posts.count()
