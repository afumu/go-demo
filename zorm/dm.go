package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"gitee.com/chunanyong/dm"
	"gitee.com/chunanyong/zorm"
	"io"
	"reflect"
	"strconv"
)

func initDm(url string) {
	var err error
	dbConfig := &zorm.DataSourceConfig{
		DSN: url,
		//sql.Open(DriverName,DSN) DriverName就是驱动的sql.Open第一个字符串参数,根据驱动实际情况获取
		DriverName:            "dm",
		Dialect:               "dm",
		SlowSQLMillis:         0, //慢sql的时间阈值,单位毫秒.小于0是禁用SQL语句输出;等于0是只输出SQL语句,不计算执行时间;大于0是计算SQL执行时间,并且>=SlowSQLMillis值
		MaxOpenConns:          0, //数据库最大连接数,默认50
		MaxIdleConns:          0, //数据库最大空闲连接数,默认50
		ConnMaxLifetimeSecond: 0, //连接存活秒时间. 默认600(10分钟)后连接被销毁重建.
		//避免数据库主动断开连接,造成死连接.MySQL默认wait_timeout 28800秒(8小时)
		DefaultTxOptions: nil, //事务隔离级别的默认配置,默认为nil
	}
	db, err = zorm.NewDBDao(dbConfig)
	if err != nil {
		fmt.Errorf("数据库连接异常 %v", err)
		panic(err)
	}

	// 注册达梦TEXT类型转string插件,dialectColumnType 值是 Dialect.字段类型 ,例如 dm.TEXT
	zorm.RegisterCustomDriverValueConver("dm.TEXT", CustomDMText{})
}

// CustomDMText 实现ICustomDriverValueConver接口,扩展自定义类型,例如 达梦数据库TEXT类型,映射出来的是dm.DmClob类型,无法使用string类型直接接收
type CustomDMText struct{}

// GetDriverValue 根据数据库列类型,返回driver.Value的实例,struct属性类型
// map接收或者字段不存在,无法获取到structFieldType,会传入nil
func (dmtext CustomDMText) GetDriverValue(ctx context.Context, columnType *sql.ColumnType, structFieldType *reflect.Type) (driver.Value, error) {
	// 如果需要使用structFieldType,需要先判断是否为nil
	// if structFieldType != nil {
	// }
	return &dm.DmClob{}, nil
}

// ConverDriverValue 数据库列类型,GetDriverValue返回的driver.Value的临时接收值,struct属性类型
// map接收或者字段不存在,无法获取到structFieldType,会传入nil
// 返回符合接收类型值的指针,指针,指针!!!!
func (dmtext CustomDMText) ConverDriverValue(ctx context.Context, columnType *sql.ColumnType, tempDriverValue driver.Value, structFieldType *reflect.Type) (interface{}, error) {
	// 如果需要使用structFieldType,需要先判断是否为nil
	// if structFieldType != nil {
	// }

	// 类型转换
	dmClob, isok := tempDriverValue.(*dm.DmClob)
	if !isok {
		return tempDriverValue, errors.New("->ConverDriverValue-->转换至*dm.DmClob类型失败")
	}
	if dmClob == nil || !dmClob.Valid {
		return new(string), nil
	}
	// 获取长度
	dmlen, errLength := dmClob.GetLength()
	if errLength != nil {
		return dmClob, errLength
	}

	// int64转成int类型
	strInt64 := strconv.FormatInt(dmlen, 10)
	dmlenInt, errAtoi := strconv.Atoi(strInt64)
	if errAtoi != nil {
		return dmClob, errAtoi
	}

	// 读取字符串
	str, errReadString := dmClob.ReadString(1, dmlenInt)

	// 处理空字符串或NULL造成的EOF错误
	if errReadString == io.EOF {
		return new(string), nil
	}

	return &str, errReadString
}
