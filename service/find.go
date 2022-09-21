package service

import (
	"chatDemo/conf"
	"chatDemo/model/ws"
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     int64  `json:"read"`
	CreateAt int64  `json:"createAt"`
}

func InsertMsg(database string, id string, content string, read int64, expire int64) (err error) {
	collection := conf.MongoDBClient.Database(database).Collection(id)
	comment := ws.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}

	_, err = collection.InsertOne(context.TODO(), comment)
	return
}

func FindMany(database string, sendId string, id string, time int64, pageSize int) (results []ws.Result, err error) {
	// 定义接受数据数组
	var resultsMe []ws.Trainer
	var resultsYou []ws.Trainer
	// 申明游标指针
	var cursor *mongo.Cursor

	// 获取id集合
	coll := conf.MongoDBClient.Database(database).Collection(id)
	cursor, err = coll.Find(context.TODO(), bson.D{}, options.Find().SetSort(bson.D{primitive.E{Key: "startTime", Value: -1}}), options.Find().SetLimit(int64(pageSize)))
	if err != nil {
		log.Println("Collection(id) coll.Find", err)
	}

	if err = cursor.All(context.TODO(), &resultsMe); err != nil {
		log.Println("Collection(id) cursor.All", err)
	}
	log.Println("Collection(id)", id)
	log.Println("Collection(id) resultsMe", len(resultsMe))

	// 获取sendId集合
	coll = conf.MongoDBClient.Database(database).Collection(sendId)
	cursor, err = coll.Find(context.TODO(), bson.D{}, options.Find().SetSort(bson.D{primitive.E{Key: "startTime", Value: -1}}), options.Find().SetLimit(int64(pageSize)))
	if err != nil {
		log.Println("Collection(sendId) coll.Find", err)
	}
	if err = cursor.All(context.TODO(), &resultsYou); err != nil {
		log.Println("Collection(sendId) cursor.All", err)
	}

	log.Println("Collection(sendId)", sendId)
	log.Println("Collection(sendId) resultsYou", len(resultsYou))

	results, _ = AppendAndSort(resultsMe, resultsYou)
	return
}

func FirstFindMsg(database string, sendId string, id string) (results []ws.Result, err error) {
	var resultMe []ws.Trainer
	var resultYou []ws.Trainer

	sendIdCollection := conf.MongoDBClient.Database(database).Collection(sendId)
	idCollection := conf.MongoDBClient.Database(database).Collection(id)

	filter := bson.M{"read": bson.M{
		"&all": []uint{0},
	}}
	sendIdCursor, err := sendIdCollection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{primitive.E{
		Key: "startTime", Value: 1}}), options.Find().SetLimit(1))
	if sendIdCursor != nil {
		return
	}
	var unRead []ws.Trainer
	err = sendIdCursor.All(context.TODO(), &unRead)
	if err != nil {
		log.Println("sendIdCursor err=", err)
	}

	if len(unRead) > 0 {
		timefilter := bson.M{
			"startTime": bson.M{
				"$gte": unRead[0].StartTime,
			}}
		cursor, err := sendIdCollection.Find(context.TODO(), timefilter)
		if err != nil {
			log.Println("sendIdCollection.Find err=", err)
		}
		if err = cursor.All(context.TODO(), &resultYou); err != nil {
			log.Println("sendIdCollection cursor.All err=", err)
		}

		cursor, err = idCollection.Find(context.TODO(), timefilter)
		if err != nil {
			log.Println("idCollection.Find err=", err)
		}

		if err = cursor.All(context.TODO(), &resultMe); err != nil {
			log.Println("idCollection cursor.All", err)
		}
		results, err = AppendAndSort(resultMe, resultYou)
		if err != nil {
			log.Println("AppendAndSort err", err)
		}
	} else {
		if results, err = FindMany(database, sendId, id, 9999999999, 10); err != nil {
			log.Println("FindMany err", err)
		}

	}
	overTimeFilter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{primitive.E{Key: "endTime", Value: bson.M{"&lt": time.Now().Unix()}}},
			bson.D{primitive.E{Key: "read", Value: bson.M{"$eq": 1}}},
		}},
	}
	_, _ = sendIdCollection.DeleteMany(context.TODO(), overTimeFilter)
	_, _ = idCollection.DeleteMany(context.TODO(), overTimeFilter)
	// 将所有的维度设置为已读
	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"$set": bson.M{"read": 1},
	})
	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"&set": bson.M{"ebdTime": time.Now().Unix() + int64(3*month)},
	})
	return
}

func AppendAndSort(resultMe, resultYou []ws.Trainer) (results []ws.Result, err error) {
	for _, r := range resultMe {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}

		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "me",
		}

		results = append(results, result)
	}

	for _, r := range resultYou {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}

		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "you",
		}
		results = append(results, result)
	}

	//排序
	sort.Slice(results, func(i, j int) bool { return results[i].StartTime < results[j].StartTime })
	return results, nil
}
