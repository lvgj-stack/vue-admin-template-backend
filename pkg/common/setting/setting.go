package setting

import (
	"container/list"
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"sync/atomic"
	"unsafe"

	"github.com/Mr-LvGJ/jota/log"
	"github.com/Mr-LvGJ/jota/models"

	"github.com/spf13/viper"
)

var gConfig = unsafe.Pointer(&Config{})

type Config struct {
	RunMode   string                 `yaml:"RunMode"   default:"debug"`
	Addr      string                 `yaml:"Addr"      default:":8080"`
	Log       *log.Config            `yaml:"Log"`
	Database  *models.DatabaseConfig `yaml:"Database"`
	Jwt       JwtConfig              `yaml:"Jwt"`
	AccessLog *log.Config            `yaml:"AccessLog"`
}

type JwtConfig struct {
	Key         string `yaml:"Key"`
	IdentityKey string `yaml:"IdentityKey"`
}

func C() *Config {
	return (*Config)(atomic.LoadPointer(&gConfig))
}

func InitConfig(configFile string) {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // 读取匹配的环境变量
	SetDefault(C())
	log.Info(context.Background(), "after init config", "config", C())
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	if err := viper.Unmarshal(C()); err != nil {
		fmt.Fprintln(os.Stderr, "unmarshal error")
	}

}

func SetDefault(v any) error {
	setDefaultByTagWithPrefix(context.Background(), v)
	typeOf := reflect.TypeOf(v).Elem()
	valueOf := reflect.ValueOf(v).Elem()
	return subSetDefault(typeOf, valueOf)
}

func typeElem(t reflect.Type) (ret reflect.Type) {
	ret = t
	for ret.Kind() == reflect.Ptr || ret.Kind() == reflect.Uintptr {
		ret = ret.Elem()
	}
	return
}

func subSetDefault(typeOf reflect.Type, valueOf reflect.Value) error {

	for i := 0; i < typeOf.NumField(); i++ {
		tField := typeOf.Field(i)
		tField.Type = typeElem(tField.Type)
		if tField.Type.Kind() != reflect.Struct {
			continue
		}
		vField := valueOf.Field(i)
		switch tField.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			defaultVal, ok := tField.Tag.Lookup("default")
			if !ok {
				continue
			}
			defaultValInt, err := strconv.ParseInt(defaultVal, 10, 64)
			if err != nil {
				return err
			}
			vField.SetInt(defaultValInt)
		case reflect.String:
			defaultVal, ok := tField.Tag.Lookup("default")
			if !ok {
				continue
			}
			vField.SetString(defaultVal)
		case reflect.Struct, reflect.Pointer:
			err := subSetDefault(tField.Type, vField)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type ptype struct {
	typ reflect.Type
	vt  reflect.Value
}

func setDefaultByTagWithPrefix(ctx context.Context, v any) {
	l := list.New()
	typ := reflect.TypeOf(v).Elem()
	vt := reflect.ValueOf(v).Elem()

	l.PushBack(&ptype{typ: typ, vt: vt})
	for ele := l.Front(); ele != nil; ele = ele.Next() {
		e := ele.Value.(*ptype)
		e.typ = typeElem(e.typ)
		if e.typ.Kind() != reflect.Struct {
			log.Info(ctx, "skip unsupported type", "type", e.typ)
			continue
		}
		for i := 0; i < e.typ.NumField(); i++ {
			var vField reflect.Value
			f := e.typ.Field(i)
			ft := typeElem(f.Type)
			if e.vt.Kind() == reflect.Ptr || e.vt.Kind() == reflect.Uintptr {
				vField = e.vt.Elem().Field(i)
			} else {
				vField = e.vt.Field(i)
			}
			if vField.Kind() == reflect.Ptr && vField.IsNil() {
				vField = reflect.New(vField.Type().Elem())
				e.vt.Field(i).Set(vField)
			}
			switch ft.Kind() {
			case reflect.Bool,
				reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64,
				reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32,
				reflect.Uint64,
				reflect.Float32,
				reflect.Float64,
				reflect.String:
				//setFieldDefaultByTag(ctx, f, tf)
				setFieldDefaultByTag(ctx, f, vField)
			case reflect.Array, reflect.Slice:
				switch typeElem(ft.Elem()).Kind() {
				case reflect.Bool,
					reflect.Int,
					reflect.Int8,
					reflect.Int16,
					reflect.Int32,
					reflect.Int64,
					reflect.Uint,
					reflect.Uint8,
					reflect.Uint16,
					reflect.Uint32,
					reflect.Uint64,
					reflect.Float32,
					reflect.Float64,
					reflect.String:
					//setFieldDefaultByTag(ctx, f, tf)
					setFieldDefaultByTag(ctx, f, vField)
					continue
				}
				log.Info(ctx, "skip unsupported type of Array or Slice", "field", f.Name, "type", ft.Elem())
			case reflect.Struct:
				//l.PushBack(&pairPrefixType{key, ft})
				l.PushBack(&ptype{typ: f.Type, vt: vField})

			default:
				log.Info(ctx, "skip unsupported type", "field", f.Name, "type", f.Type)
			}
		}
	}
}

func setFieldDefaultByTag(_ context.Context, e reflect.StructField, vt reflect.Value) {
	switch e.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		defaultVal, _ := e.Tag.Lookup("default")
		defaultValInt, _ := strconv.ParseInt(defaultVal, 10, 64)
		vt.SetInt(defaultValInt)
	case reflect.String:
		defaultVal, _ := e.Tag.Lookup("default")
		vt.SetString(defaultVal)
	}
}
