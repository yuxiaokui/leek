import  pymongo
import  threading
import  time

client = pymongo.MongoClient('mongodb://mongodb.t2.daoapp.io:61257/')
db = client['task']
posts = db['tasks']

def insert(task):
    posts.insert_one(task)

for tgt in open("urls.txt"):
    task = {"url": tgt,"time":time.time()}
    print task
    t = threading.Thread(target=insert, args=(task,))
    t.start()

t.join()
