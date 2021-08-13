package db

import (
	"context"
	"fw/src/core"
	"fw/src/core/log"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// ============================================================================

var reg *bsoncodec.Registry

// ============================================================================

type D = bson.D
type M = bson.M
type A = bson.A
type E = primitive.E

type Database struct {
	ctx context.Context // root ctx
	c   *mongo.Client   // client
	db  *mongo.Database // database
}

// ============================================================================

func init() {
	reg = bson.NewRegistryBuilder().
		RegisterTypeDecoder(
			reflect.TypeOf(time.Time{}),
			bsoncodec.NewTimeCodec(bsonoptions.TimeCodec().SetUseLocalTimeZone(true)),
		).
		Build()
}

func NewDatabase() *Database {
	return &Database{
		ctx: context.Background(),
	}
}

// ============================================================================

func (self *Database) Open(cnnstr string, maxPoolSize ...uint64) {
	cnnstr_obj, err := connstring.Parse(cnnstr)
	if err != nil {
		core.Panic("open database failed:", err)
	}

	// opts
	opts := options.Client().ApplyURI(cnnstr).SetRegistry(reg)
	if len(maxPoolSize) > 0 {
		opts.SetMaxPoolSize(maxPoolSize[0])
	}

	// connect
	c, err := mongo.Connect(self.ctx, opts)
	if err != nil {
		core.Panic("open database failed:", err)
	}

	// set
	self.c = c
	self.db = c.Database(cnnstr_obj.Database)
}

func (self *Database) Close() {
	self.c.Disconnect(self.ctx)

	self.c = nil
	self.db = nil
}

// ============================================================================

func (self *Database) HasDB() bool {
	names, err := self.c.ListDatabaseNames(self.ctx, bson.D{}) // 3.2 does NOT support filter
	if err != nil {
		log.Error("<HasDB>:", err)
		return true // return true for safety
	}

	for _, v := range names {
		if v == self.db.Name() {
			return true
		}
	}

	return false
}

func (self *Database) HasCollection(name string) bool {
	names, err := self.db.ListCollectionNames(self.ctx, bson.D{{"name", name}}) // 3.2 HAS filter
	if err != nil {
		log.Error("<HasCollection>:", err)
		return true // return true for safety
	}

	return len(names) > 0
}

func (self *Database) HasIndex(coll string, name string) bool {
	var names []*struct {
		Name string
	}

	cursor, err := self.db.Collection(coll).Indexes().List(self.ctx)
	if err == nil {
		err = cursor.All(self.ctx, &names)
	}
	if err != nil {
		return false
	}

	for _, v := range names {
		if v.Name == name {
			return true
		}
	}

	return false
}

func (self *Database) CreateCappedCollection(coll string, size int) {
	// check
	if self.HasCollection(coll) {
		return
	}

	// create
	err := self.Run(bson.D{
		{"create", coll},
		{"capped", true},
		{"size", size},
	}, nil)
	if err != nil {
		log.Error("<CreateCappedCollection>:", err)
	}
}

func (self *Database) CreateTTLIndex(coll string, name string, key string, sec int32) {
	// check
	if self.HasIndex(coll, name) {
		return
	}

	// prepare
	model := mongo.IndexModel{bson.D{{key, 1}}, options.Index().SetName(name).SetExpireAfterSeconds(sec)}

	// create
	_, err := self.db.Collection(coll).Indexes().CreateOne(self.ctx, model)
	if err != nil {
		log.Error("<CreateTTLIndex>:", err, coll, name)
	}
}

func (self *Database) CreateIndex(coll string, name string, keys []string, unique bool) {
	// check
	if self.HasIndex(coll, name) {
		return
	}

	// prepare
	var ascending_keys bson.D
	for _, v := range keys {
		ascending_keys = append(ascending_keys, E{v, 1})
	}
	model := mongo.IndexModel{ascending_keys, options.Index().SetName(name).SetUnique(unique)}

	// create
	_, err := self.db.Collection(coll).Indexes().CreateOne(self.ctx, model)
	if err != nil {
		log.Error("<CreateIndex>:", err, coll, name)
	}
}

// ============================================================================

func (self *Database) GetObject(coll string, id interface{}, obj interface{}) error {
	err := self.db.Collection(coll).FindOne(self.ctx, bson.D{{"_id", id}}).Decode(obj)
	if is_critical(err) {
		log.Error("db.GetObject():", err, coll, id)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) GetObjectByCond(coll string, cond interface{}, obj interface{}) error {
	err := self.db.Collection(coll).FindOne(self.ctx, cond).Decode(obj)
	if is_critical(err) {
		log.Error("db.GetObjectByCond():", err, coll, cond)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) GetProjection(coll string, id interface{}, proj interface{}, obj interface{}) error {
	opts := options.FindOne().SetProjection(proj)
	err := self.db.Collection(coll).FindOne(self.ctx, bson.D{{"_id", id}}, opts).Decode(obj)
	if is_critical(err) {
		log.Error("db.GetProjection():", err, coll, id, proj)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) GetProjectionByCond(coll string, cond interface{}, proj interface{}, obj interface{}) error {
	opts := options.FindOne().SetProjection(proj)
	err := self.db.Collection(coll).FindOne(self.ctx, cond, opts).Decode(obj)
	if is_critical(err) {
		log.Error("db.GetProjectionByCond():", err, coll, cond, proj)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) GetAllObjects(coll string, obj interface{}) error {
	cursor, err := self.db.Collection(coll).Find(self.ctx, bson.D{})
	if err == nil {
		err = cursor.All(self.ctx, obj)
	}
	if is_critical(err) {
		log.Error("db.GetAllObjects():", err, coll)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) GetAllObjectsByCond(coll string, cond interface{}, obj interface{}) error {
	cursor, err := self.db.Collection(coll).Find(self.ctx, cond)
	if err == nil {
		err = cursor.All(self.ctx, obj)
	}
	if is_critical(err) {
		log.Error("db.GetAllObjectsByCond():", err, coll, cond)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) GetAllProjectionsByCond(coll string, cond interface{}, proj interface{}, obj interface{}) error {
	opts := options.Find().SetProjection(proj)
	cursor, err := self.db.Collection(coll).Find(self.ctx, cond, opts)
	if err == nil {
		err = cursor.All(self.ctx, obj)
	}
	if is_critical(err) {
		log.Error("db.GetAllProjectionsByCond():", err, coll, cond, proj)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) Insert(coll string, doc interface{}) error {
	_, err := self.db.Collection(coll).InsertOne(self.ctx, doc)
	if is_critical(err) {
		log.Error("db.Insert():", err, coll, doc)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) InsertMany(coll string, docs ...interface{}) error {
	_, err := self.db.Collection(coll).InsertMany(self.ctx, docs)
	if is_critical(err) {
		log.Error("db.InsertMany():", err, coll, docs)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) Remove(coll string, id interface{}) error {
	_, err := self.db.Collection(coll).DeleteOne(self.ctx, bson.D{{"_id", id}})
	if is_critical(err) {
		log.Error("db.Remove():", err, coll, id)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) RemoveByCond(coll string, cond interface{}) error {
	_, err := self.db.Collection(coll).DeleteOne(self.ctx, cond)
	if is_critical(err) {
		log.Error("db.RemoveByCond():", err, coll, cond)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) RemoveAll(coll string, cond interface{}) error {
	_, err := self.db.Collection(coll).DeleteMany(self.ctx, cond)
	if is_critical(err) {
		log.Error("db.RemoveAll():", err, coll, cond)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) Update(coll string, id interface{}, doc interface{}) error {
	_, err := self.db.Collection(coll).UpdateOne(self.ctx, bson.D{{"_id", id}}, doc)
	if is_critical(err) {
		log.Error("db.Update():", err, coll, id)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) UpdateByCond(coll string, cond interface{}, doc interface{}) error {
	_, err := self.db.Collection(coll).UpdateOne(self.ctx, cond, doc)
	if is_critical(err) {
		log.Error("db.UpdateByCond():", err, coll, cond, doc)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) Upsert(coll string, id interface{}, doc interface{}) error {
	opts := options.Update().SetUpsert(true)
	_, err := self.db.Collection(coll).UpdateOne(self.ctx, bson.D{{"_id", id}}, doc, opts)
	if is_critical(err) {
		log.Error("db.Upsert():", err, coll, id)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) UpsertByCond(coll string, cond interface{}, doc interface{}) error {
	opts := options.Update().SetUpsert(true)
	_, err := self.db.Collection(coll).UpdateOne(self.ctx, cond, doc, opts)
	if is_critical(err) {
		log.Error("db.UpsertByCond():", err, coll, cond)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) UpdateAll(coll string, cond interface{}, doc interface{}) error {
	_, err := self.db.Collection(coll).UpdateMany(self.ctx, cond, doc)
	if is_critical(err) {
		log.Error("db.UpdateAll():", err, coll, cond, doc)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) Replace(coll string, id interface{}, doc interface{}) error {
	opts := options.Replace().SetUpsert(true)
	_, err := self.db.Collection(coll).ReplaceOne(self.ctx, bson.D{{"_id", id}}, doc, opts)
	if is_critical(err) {
		log.Error("db.Replace():", err, coll, id)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) ReplaceByCond(coll string, cond interface{}, doc interface{}) error {
	opts := options.Replace().SetUpsert(true)
	_, err := self.db.Collection(coll).ReplaceOne(self.ctx, cond, doc, opts)
	if is_critical(err) {
		log.Error("db.ReplaceByCond():", err, coll, cond)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) FindAndUpdate(coll string, cond interface{}, doc interface{}, proj interface{}, obj interface{}) error {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetProjection(proj)
	err := self.db.Collection(coll).FindOneAndUpdate(self.ctx, cond, doc, opts).Decode(obj)
	if is_critical(err) {
		log.Error("db.FindAndUpdate():", err, coll, cond, doc, proj)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) FindAndUpsert(coll string, cond interface{}, doc interface{}, proj interface{}, obj interface{}) error {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetProjection(proj).SetUpsert(true)
	err := self.db.Collection(coll).FindOneAndUpdate(self.ctx, cond, doc, opts).Decode(obj)
	if is_critical(err) {
		log.Error("db.FindAndUpsert():", err, coll, cond, doc, proj)
		log.Error(core.Callstack())
	}
	return err
}

func (self *Database) Count(coll string, cond interface{}) (int64, error) {
	n, err := self.db.Collection(coll).CountDocuments(self.ctx, cond)
	if is_critical(err) {
		log.Error("db.Count():", err, coll, cond)
		log.Error(core.Callstack())
	}
	return n, err
}

func (self *Database) Aggregate(coll string, pipeline interface{}, obj interface{}) error {
	cursor, err := self.db.Collection(coll).Aggregate(self.ctx, pipeline)
	if err == nil {
		err = cursor.All(self.ctx, obj)
	}
	if is_critical(err) {
		log.Error("db.Aggregate():", err, coll, pipeline)
		log.Error(core.Callstack())
	}
	return err
}

// ============================================================================

func (self *Database) Run(cmd bson.D, obj interface{}) error {
	r := self.db.RunCommand(self.ctx, cmd)
	if obj == nil {
		return r.Err()
	} else {
		return r.Decode(obj)
	}
}

func (self *Database) Execute(f func(*mongo.Database)) {
	f(self.db)
}

// ============================================================================

func IsNotFound(err error) bool {
	return err == mongo.ErrNoDocuments
}

func IsDup(err error) bool {
	return strings.Contains(err.Error(), "E11000")
}

// ============================================================================

func is_critical(err error) bool {
	return err != nil && !IsNotFound(err) && !IsDup(err)
}
