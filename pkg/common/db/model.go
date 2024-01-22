package db

import (
	"Open_IM/pkg/common/config"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"strings"

	//"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	go_redis "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2"
	"time"
)

var DB DataBases

type DataBases struct {
	MysqlDB    mysqlDB
	mgoSession *mgo.Session
	//redisPool   *redis.Pool
	mongoClient *mongo.Client
	RDB         go_redis.UniversalClient
	Rc          *rockscache.Client
	WeakRc      *rockscache.Client
	Pool        *redsync.Redsync
}

type RedisClient struct {
	client  *go_redis.Client
	cluster *go_redis.ClusterClient
	go_redis.UniversalClient
	enableCluster bool
}

func key(dbAddress, dbName string) string {
	return dbAddress + "_" + dbName
}

func init() {
	env := os.Getenv("OPEN_DEV")
	if env == "dev" {
		return
	}
	//log.NewPrivateLog(constant.LogFileName)
	var mongoClient *mongo.Client
	var err1 error
	//mysql init
	initMysqlDB()
	// mongo init
	// "mongodb://sysop:moon@localhost/records"
	uri := "mongodb://sample.host:27017/?maxPoolSize=20&w=majority"
	if config.Config.Mongo.DBUri != "" {
		// example: mongodb://$user:$password@mongo1.mongo:27017,mongo2.mongo:27017,mongo3.mongo:27017/$DBDatabase/?replicaSet=rs0&readPreference=secondary&authSource=admin&maxPoolSize=$DBMaxPoolSize
		uri = config.Config.Mongo.DBUri
	} else {
		//mongodb://mongodb1.example.com:27317,mongodb2.example.com:27017/?replicaSet=mySet&authSource=authDB
		mongodbHosts := ""
		for i, v := range config.Config.Mongo.DBAddress {
			if i == len(config.Config.Mongo.DBAddress)-1 {
				mongodbHosts += v
			} else {
				mongodbHosts += v + ","
			}
		}

		if config.Config.Mongo.DBPassword != "" && config.Config.Mongo.DBUserName != "" {
			// clientOpts := options.Client().ApplyURI("mongodb://localhost:27017,localhost:27018/?replicaSet=replset")
			//mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
			//uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d&authSource=admin&replicaSet=replset",
			uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d&authSource=admin",
				config.Config.Mongo.DBUserName, config.Config.Mongo.DBPassword, mongodbHosts,
				config.Config.Mongo.DBDatabase, config.Config.Mongo.DBMaxPoolSize)
		} else {
			uri = fmt.Sprintf("mongodb://%s/%s/?maxPoolSize=%d&authSource=admin",
				mongodbHosts, config.Config.Mongo.DBDatabase,
				config.Config.Mongo.DBMaxPoolSize)
		}
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println(" mongo.Connect  failed, try ", utils.GetSelfFuncName(), err.Error(), uri)
		time.Sleep(time.Duration(30) * time.Second)
		mongoClient, err1 = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err1 != nil {
			fmt.Println(" mongo.Connect retry failed, panic", err.Error(), uri)
			panic(err1.Error())
		}
	}
	fmt.Println("mongo driver client init success: ", uri)
	// mongodb create index
	if err := createMongoIndex(mongoClient, cSendLog, false, "send_id", "-send_time"); err != nil {
		fmt.Println("send_id", "-send_time", "index create failed", err.Error())
		panic(err.Error())
	}
	if err := createMongoIndex(mongoClient, cChat, false, "uid"); err != nil {
		fmt.Println("uid", " index create failed", err.Error())
		//panic(err.Error())
	}
	if err := createMongoIndex(mongoClient, cWorkMoment, true, "-create_time", "work_moment_id"); err != nil {
		fmt.Println("-create_time", "work_moment_id", "index create failed", err.Error())
		panic(err.Error())
	}
	if err := createMongoIndex(mongoClient, cWorkMoment, true, "work_moment_id"); err != nil {
		fmt.Println("work_moment_id", "index create failed", err.Error())
		panic(err.Error())
	}
	if err := createMongoIndex(mongoClient, cWorkMoment, false, "user_id", "-create_time"); err != nil {
		fmt.Println("user_id", "-create_time", "index create failed", err.Error())
		panic(err.Error())
	}
	if err := createMongoIndex(mongoClient, cTag, false, "user_id", "-create_time"); err != nil {
		fmt.Println("user_id", "-create_time", "index create failed", err.Error())
		panic(err.Error())
	}
	if err := createMongoIndex(mongoClient, cTag, true, "tag_id"); err != nil {
		fmt.Println("tag_id", "index create failed", err.Error())
		panic(err.Error())
	}
	fmt.Println("createMongoIndex success")
	DB.mongoClient = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if config.Config.Redis.EnableCluster {
		DB.RDB = go_redis.NewClusterClient(&go_redis.ClusterOptions{
			Addrs:    config.Config.Redis.DBAddress,
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			PoolSize: 50,
		})
		_, err = DB.RDB.Ping(ctx).Result()
		if err != nil {
			panic(err.Error())
		}
		clientPool := goredis.NewPool(DB.RDB)
		DB.Pool = redsync.New(clientPool)
	} else {
		DB.RDB = go_redis.NewClient(&go_redis.Options{
			Addr:     config.Config.Redis.DBAddress[0],
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			DB:       0,                              // use default DB
			PoolSize: 100,                            // 连接池大小
		})
		clientPool := goredis.NewPool(DB.RDB)
		DB.Pool = redsync.New(clientPool)
		_, err = DB.RDB.Ping(ctx).Result()
		if err != nil {
			panic(err.Error())
		}
	}
	// 强一致性缓存，当一个key被标记删除，其他请求线程会被锁住轮询直到新的key生成，适合各种同步的拉取, 如果弱一致可能导致拉取还是老数据，毫无意义
	DB.Rc = rockscache.NewClient(DB.RDB, rockscache.NewDefaultOptions())
	DB.Rc.Options.StrongConsistency = true

	// 弱一致性缓存，当一个key被标记删除，其他请求线程直接返回该key的value，适合高频并且生成很缓存很慢的情况 如大群发消息缓存的缓存
	DB.WeakRc = rockscache.NewClient(DB.RDB, rockscache.NewDefaultOptions())
	DB.WeakRc.Options.StrongConsistency = false
	fmt.Println("init kv db ok")
}

func createMongoIndex(client *mongo.Client, collection string, isUnique bool, keys ...string) error {
	db := client.Database(config.Config.Mongo.DBDatabase).Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	indexView := db.Indexes()
	keysDoc := bson.D{}

	// 复合索引
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = append(keysDoc, bson.E{Key: strings.TrimLeft(key, "-"), Value: -1})
		} else {
			keysDoc = append(keysDoc, bson.E{Key: key, Value: 1})
		}
	}

	// 创建索引
	index := mongo.IndexModel{
		Keys: keysDoc,
	}
	if isUnique == true {
		index.Options = options.Index().SetUnique(true)
	}
	result, err := indexView.CreateOne(
		context.Background(),
		index,
		opts,
	)
	if err != nil {
		return utils.Wrap(err, result)
	}
	return nil
}
