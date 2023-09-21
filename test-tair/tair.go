package main

import (
	"context"
	"fmt"
	"github.com/alibaba/tair-go/tair"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var ctx = context.Background()

var tairClient *tair.TairClient

func init() {
	tairClient = tair.NewTairClient(&redis.Options{
		Addr:     "xxx.redis.rds.aliyuncs.com:6379",
		Password: "xxx",
		DB:       0,
	})
}

func main() {
	err := tairClient.TrSetBit(ctx, "clock20230901", 522809, 1).Err()
	if err != nil {
		fmt.Printf("TrSetBit err: %v\n", err)
	}
	isClock := tairClient.TrGetBit(ctx, "clock20230901", 522809).Val()
	fmt.Printf("TrGetBit isClock: %v\n", isClock)
	// 输出: TrGetBit isClock: 1
}

func testRoaring() {
	startTs := time.Now().UnixMicro() - 10000
	startTsStr := strconv.FormatInt(startTs, 10)
	points := []*tair.ExTsDataPoint{
		(&tair.ExTsDataPoint{}).SetSKey("cpu").SetTs(startTsStr).SetValue(0.3),
		(&tair.ExTsDataPoint{}).SetSKey("memory").SetTs(startTsStr).SetValue(0.2),
	}
	err := tairClient.ExTsMAdd(ctx, "pk", points).Err()
	if err != nil {
		fmt.Printf("ExTsMAdd err: %v\n", err)
	}
	cpu := tairClient.ExTsGet(ctx, "pk", "cpu").Val()
	memory := tairClient.ExTsGet(ctx, "pk", "memory").Val()
	fmt.Printf("ExTsGet cpu: %v memory: %v\n", cpu, memory)
	// 输出: ExTsGet cpu: [1695292325437112 0.3] memory: [1695292325437112 0.2]
}

func testSearch() {
	tairClient.TftCreateIndex(ctx, "index1", "{\"mappings\":{\"dynamic\":\"false\",\"properties\":{\"f0\":{\"type\":\"text\"},\"f1\":{\"type\":\"text\"}}}}")
	tairClient.TftCreateIndex(ctx, "index2", "{\"mappings\":{\"dynamic\":\"false\",\"properties\":{\"f0\":{\"type\":\"text\"},\"f1\":{\"type\":\"text\"}}}}")
	tairClient.TftCreateIndex(ctx, "index3", "{\"mappings\":{\"dynamic\":\"false\",\"properties\":{\"f0\":{\"type\":\"text\"},\"f1\":{\"type\":\"text\"}}}}")
	tairClient.TftAddDocWithId(ctx, "index1", "{\"f0\":\"v0\",\"f1\":\"3\"}", "1")
	tairClient.TftAddDocWithId(ctx, "index2", "{\"f0\":\"v1\",\"f1\":\"3\"}", "2")
	tairClient.TftAddDocWithId(ctx, "index3", "{\"f0\":\"v3\",\"f1\":\"3\"}", "3")
	tairClient.TftAddDocWithId(ctx, "index1", "{\"f0\":\"v3\",\"f1\":\"4\"}", "4")
	tairClient.TftAddDocWithId(ctx, "index2", "{\"f0\":\"v3\",\"f1\":\"5\"}", "5")

	res, err := tairClient.TftMSearch(ctx, 3, "{\"query\":{\"match\":{\"f1\":\"3\"}}}", "index1", "index2", "index3").Result()
	if err != nil {
		fmt.Printf("TftMSearch err: %v\n", err)
	}
	fmt.Printf("TftMSearch res: %v\n", res)
	// 输出: TftMSearch res: {"hits":{"hits":[{"_id":"1","_index":"index1","_score":1.0,"_source":{"f0":"v0","f1":"3"}},{"_id":"2","_index":"index2","_sco
	//re":1.0,"_source":{"f0":"v1","f1":"3"}},{"_id":"3","_index":"index3","_score":0.094159,"_source":{"f0":"v3","f1":"3"}}],"max_score":1.0,"tota
	//l":{"relation":"eq","value":3}},"aux_info":{"index_crc64":52600736426816810}}
}

func testJson() {
	json := "[{\"name\":\"mike\",\"age\":25},{\"name\":\"john\",\"age\":18}]"
	err := tairClient.JsonSet(ctx, "testJson", ".", json).Err()
	if err != nil {
		fmt.Printf("JsonSet err: %v\n", err)
	}
	res, err := tairClient.JsonGet(ctx, "testJson").Result()
	if err != nil {
		fmt.Printf("JsonGet err: %v\n", err)
	}
	fmt.Printf("JsonGet res: %v\n", res)
	// 输出: JsonGet res: [{"name":"mike","age":25},{"name":"john","age":18}]
	err = tairClient.JsonArrAppend(ctx, "testJson", ".", "{\"name\":\"tim\",\"age\":35}").Err()
	if err != nil {
		fmt.Printf("JsonArrAppend err: %v\n", err)
	}
	res, err = tairClient.JsonGet(ctx, "testJson").Result()
	if err != nil {
		fmt.Printf("JsonGet err: %v\n", err)
	}
	fmt.Printf("JsonGet res: %v\n", res)
	// 输出: JsonGet res: [{"name":"mike","age":25},{"name":"john","age":18},{"name":"tim","age":35}]
}

func testBloom() {
	var capacity int64
	capacity = 100    //容量
	errorRate := 0.01 //误判率
	err := tairClient.BfReserve(ctx, "testBloom", capacity, errorRate).Err()
	if err != nil {
		fmt.Printf("BfReserve err: %v\n", err)
	}
	for i := 1; i <= 100; i++ {
		err = tairClient.BfAdd(ctx, "testBloom", string(i)).Err()
		if err != nil {
			fmt.Printf("BfAdd err: %v\n", err)
		}
	}
	exists, err := tairClient.BfExists(ctx, "testBloom", "1").Result()
	if err != nil {
		fmt.Printf("BfExists err: %v\n", err)
	}
	fmt.Printf("BfExists 1 exists res %v\n", exists)
	// BfExists 1 exists res true

	exists, err = tairClient.BfExists(ctx, "testBloom", "101").Result()
	if err != nil {
		fmt.Printf("BfExists err: %v\n", err)
	}
	fmt.Printf("BfExists 101 exists res %v\n", exists)
	// BfExists 1 exists res true
}

func GisAdd() {
	// 设置多边形经纬度 (例如图中配送范围形成的四边形各角的经纬度)
	err := tairClient.GisAdd(ctx, "sellArea", "pingShuiJie", "POLYGON ((120.11205 30.30644, 120.11900 30.30918, 120.11976 30.29999, 120.11310 30.29987))").Err()
	if err != nil {
		fmt.Printf("GisAdd err: %v\n", err)
	}
	// 判断指定的点、线或面是否包含在目标area的多边形中，若包含，则返回目标area中命中的多边形数量与多边形信息。
	res, err := tairClient.GisContains(ctx, "sellArea", "POINT (120.11533 30.30085)").Result()
	if err != nil {
		fmt.Printf("GisContains err: %v\n", err)
	}
	fmt.Printf("GisContains res: %v\n", res)
	// 输出: GisContains res: map[pingShuiJie:POLYGON((120.11205 30.30644,120.119 30.30918,120.11976 30.29999,120.1131 30.29987))]
	// 返回area信息,则查询经纬度在sellArea范围内
}

func ExZAddMember() {
	m1 := tair.ExZAddMember{
		Score:  "1#2#1",
		Member: "m1",
	}
	m2 := tair.ExZAddMember{
		Score:  "2#1#3",
		Member: "m2",
	}
	m3 := tair.ExZAddMember{
		Score:  "1#2#2",
		Member: "m3",
	}
	//说明
	//如需实现多维度的排序，各维度的分数之间使用井号（#）分隔，例如111#222#121，且要求该Key中所有成员的分数格式必须相同。
	err := tairClient.ExZAddManyMember(ctx, "testExZAdd", m1, m2, m3).Err()
	if err != nil {
		fmt.Printf("ExZAdd err: %v\n", err)
	}
	res, err := tairClient.ExZRangeByScore(ctx, "testExZAdd", "0#0#0", "999#999#999").Result()
	if err != nil {
		fmt.Printf("ExZRangeByScore err: %v\n", err)
	}
	fmt.Printf("ExZRangeByScore res: %v\n", res)
	// 输出: ExZRangeByScore res: [m1 m3 m2]
	// 根据分数优先级从左至右优先级由高到底来看, m1,m1第一维度的分数相等并小于m2, m1和m3第一和第二维度分数相等, 根据第三维度分数比较m1小于m3
	// 得出升序排序下的结果为 m1 m3 m2
}

func ExZset() {
	exArgs := tair.ExHSetArgs{}.New()
	exArgs.Ex(time.Minute * 1)
	err := tairClient.ExHSetArgs(ctx, "testHSet", "field1", "val1", exArgs).Err()
	if err != nil {
		fmt.Printf("ExHSetArgs err: %v\n", err)
	}

	exArgs.Ex(time.Minute * 10)
	err = tairClient.ExHSetArgs(ctx, "testHSet", "field2", "val2", exArgs).Err()
	if err != nil {
		fmt.Printf("ExHSetArgs err: %v\n", err)
	}
	exGetRes, err := tairClient.ExHGet(ctx, "testHSet", "field1").Result()
	ttl1, err1 := tairClient.ExHTTL(ctx, "testHSet", "field1").Result()
	fmt.Printf("ExGet res: %v, ttl: %v, err: %v, err1 :%v \n", exGetRes, ttl1, err, err1)
	// 输出: ExGet res: val1, ttl: 59, err: <nil>, err1 :<nil>
	exGetRes, err = tairClient.ExHGet(ctx, "testHSet", "field2").Result()
	ttl2, err2 := tairClient.ExHTTL(ctx, "testHSet", "field2").Result()
	fmt.Printf("ExGet res: %v, ttl: %v, err: %v, err2 :%v \n", exGetRes, ttl2, err, err2)
	// 输出: ExGet res: val2, ttl: 599, err: <nil>, err2 :<nil>

}

func ExIncrBy() {
	exArgs := tair.ExIncrByArgs{}.New()
	exArgs.Min(1)
	// 设置TairString Value的最小值
	exArgs.Max(2)
	// 设置TairString Value的最大值
	exIncrByRes, err := tairClient.ExIncrByArgs(ctx, "testIncrBy", 1, exArgs).Result()
	fmt.Printf("ExIncrByArgs res: %v, err: %v\n", exIncrByRes, err)
	// 输出: ExIncrByArgs res: 1, err: <nil>
	// ExIncrBy 命令 和redis incrby命令一样, 会将自增后value值返回
	exIncrByRes, err = tairClient.ExIncrByArgs(ctx, "testIncrBy", 2, exArgs).Result()
	fmt.Printf("ExIncrByArgs res: %v, err: %v\n", exIncrByRes, err)
	// 输出: ExIncrByArgs res: 0, err: ERR increment or decrement would overflow
	// 当设置value值超过最大值时tair会返回错误信息
	exIncrByRes, err = tairClient.ExIncrByArgs(ctx, "testIncrBy", -2, exArgs).Result()
	fmt.Printf("ExIncrByArgs res: %v, err: %v\n", exIncrByRes, err)
	// 输出: ExIncrByArgs res: 0, err: ERR increment or decrement would overflow
	// 最小值限制同理
}

func ExString() {
	exArgs := tair.ExSetArgs{}.New()
	//exArgs.Nx()
	// Nx 只在key存在时写入
	//exArgs.Xx()
	// Xx 只在key不存在时写入
	exArgs.Ex(time.Minute)
	exArgs.Ver(1)
	//如果key存在，和当前版本号做比较：
	//	如果相等，写入，且版本号加1。
	//  如果不相等，返回异常。
	//如果key不存在或者key当前版本为0，忽略传入的版本号直接设置value，成功后版本号变为1。
	//exArgs.Abs(1)
	//	Abs：绝对版本号。设置后，无论key当前的版本号是多少，完成写入并将key的版本号覆盖为该选项中设置的值。
	err := tairClient.ExSetArgs(ctx, "exStrTest", "hello world !", exArgs).Err()
	fmt.Printf("ExSetArgs err: %v\n", err)
	// 输出: ExSetArgs err: redis: nil
	exGetRes, err := tairClient.ExGet(ctx, "exStrTest").Result()
	fmt.Printf("ExGetArgs res: %v, err: %v\n", exGetRes, err)
	// 输出: ExGetArgs res: [hello world ! 1], err: <nil>
	err = tairClient.ExSetArgs(ctx, "exStrTest", "hello world !2", exArgs).Err()
	fmt.Printf("ExSetArgs err: %v\n", err)
	exGetRes, err = tairClient.ExGet(ctx, "exStrTest").Result()
	fmt.Printf("ExGetArgs res: %v, err: %v\n", exGetRes, err)
	// 输出: ExGetArgs res: [hello world !2 2], err: <nil>
	// ExGet 会返回两个值, 第一个是value 第二个则是版本号
	// 两次写入的版本号均为1, 此时value更新了, 并且版本号自增1, 版本号变为2并且返回
	err = tairClient.ExSetArgs(ctx, "exStrTest", "hello world !3", exArgs).Err()
	fmt.Printf("ExSetArgs err: %v\n", err)
	// 输出: ExSetArgs err: ERR update version is stale
	// 版本号已经变为2, 此时ExSet传入版本号与当前版本号不相等, 就会返回错误
	exArgs.Ver(2)
	err = tairClient.ExSetArgs(ctx, "exStrTest", "hello world !4", exArgs).Err()
	fmt.Printf("ExSetArgs err: %v\n", err)
	// 输出: ExSetArgs err: <nil>
	// 更新成功
	exGetRes, err = tairClient.ExGet(ctx, "exStrTest").Result()
	fmt.Printf("ExGetArgs res: %v, err: %v\n", exGetRes, err)
	// 输出: ExGetArgs res: [hello world !4 3], err: <nil>
	// value更新成功, 版本号此时等于3
}
