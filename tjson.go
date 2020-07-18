package tjson


//		主要函数Value  可变参数 1个到2个 3个以上当2个处理
//		参数1个为获取值 string 或 int
//		参数两个为设置值 第一个参数还是 string 或 int  第二个参数可以是string int bool object array
//		当调用Value时 设置值或者获取值的时候 最好是判断类型，如果不判断直接调用，遇到值不匹配的情况下 ，被强制清0 并且强制转换成你正在调用的类型

//		判断是否为空的函数 IsNull
//		这个函数有两种形式 不带参数 和 带一个key或者index参数
//		比如jsonObject.IsNull() 是判断当前对象是不是Null
//		jsonObject.IsNull(string("xxx")) 或者 sonObject.IsNull(int(xxx)) 判断当前对象的xxx数据是否为null
//		但是也可以这样判断 那string举例 jsonObject.Value("xxx").IsNull()  和上面的差不多 。唯一区别是不管这个xxx是否为null 都会创建一个null对象 上面那种则不会

import (
	"encoding/json"
	"fmt"
)

type Type int

const (
	//空值
	Null Type = iota
	// 布尔
	Bool
	// 数字 统一64位
	Number
	// 字符串
	String
	//数组
	Array
	// json对象
	Object

	//数组初始化和扩大参考值
	arrayInitSize = 5
)

type Value struct {
	t Type

	//josn对象是一个map[string] *Value  如果是数组 暂时用切片 切片对数据访问速度是很高效的
	//如果数据是插入删除居多。。就要换成list 修改几个数组操作函数就可以了
	value interface{}
	arrayItemCount int
}

func New(v ...interface{}) *Value {
	if len(v) == 0 {
		return &Value{t:Null}
	}
	switch keyValue := v[0].(type) {
	case string:
		return parse(keyValue)
		break
	}
	return &Value{t:Null}
}

func (this* Value)Copy() *Value {

	//暂时用效率低的方法 懒得写
	if this.t == Null {
		return &Value{t:Null}
	}else if this.t == Number || this.t == Bool || this.t == String {
		return this.createValue(this.value)
	}

	return New(this.ToString())
}

func (this* Value)createValue(v interface{}) *Value {
	var result *Value
	switch keyValue := v.(type) {
	case int:
		result = &Value{t:Number, value: int64(keyValue)}
		break
	case int64:
		result = &Value{t:Number, value: keyValue}
		break
	case string:
		result = &Value{t:String, value: keyValue}
		break
	case bool:
		result = &Value{t:Bool, value: keyValue}
		break
	case Value:
		result = &keyValue
		break
	case *Value:
		result = keyValue
		break
	case int32:
		result = &Value{t:Number, value: int64(keyValue)}
		break
	case nil:
		result = &Value{}
		break
	default:
		break
	}
	return result
}

//数组以外设置值
func (this* Value)setValue(key string, v *Value)  {
	var mapValue map[string]*Value
	if this.t == Null {
		this.t = Object
		mapValue = map[string]*Value{}
		this.value = mapValue
	}else if this.t == Object{
		mapValue = this.value.(map[string]*Value)
	}else{
		return
	}

	mapValue[key] = v
}

func (this* Value)getValue(key string) *Value {
	var mapValue map[string]*Value
	//如果自己本身是个空值  把自己转化为Object
	if this.t != Object {
		this.t = Object
		mapValue = map[string]*Value{}
		this.value = mapValue
	}else{
		mapValue = this.value.(map[string]*Value)
	}

	if _, ok := mapValue[key]; ok {
		//如果子key存在直接返回
		return mapValue[key]
	} else {
		//否则创建一个Null对象返回
		result := &Value{t:Null}
		mapValue[key] = result
		return result
	}
	return nil
}

func (this* Value)Clean() {
	this.arrayItemCount = 0
	this.t = Null
	this.value = nil
}

func (this* Value)Insert(index int, v interface{}) bool {
	return this.arrayInsert(index, this.createValue(v))
}

func (this* Value)Remove(v interface{}) bool {
	if this.t == Object {
		//移除对象数据
		switch keyValue := v.(type) {
		case string:
			mapTmp := this.value.(map[string]*Value)
			delete(mapTmp, keyValue)
			return true
		}
	}else{
		//移除数组数据
		switch keyValue := v.(type) {
		case int:
			return this.arrayRemove(keyValue)
		}
	}
	return false
}

func (this* Value)arrayRemove(index int) bool {
	if this.t != Array {
		return false
	}
	if index < 0 || index >= this.arrayItemCount {
		return false
	}
	array := this.value.([]*Value)

	array = append(array[:index], array[index+1:]...)

	this.value = array

	this.arrayItemCount--

	return true
}

func (this* Value)arrayInsert(index int, v *Value) bool {
	var array []*Value
	if this.t != Array {
		//创建一个数组 并把json类型设置为Array
		array = make([]*Value, arrayInitSize)
		this.t = Array
		this.value = array
		this.arrayItemCount = 0
	}
	if this.t == Array {

		array = this.value.([]*Value)
		var arrayTmp []*Value
		if len(array) == this.arrayItemCount {
			//数组尺寸不够 扩展
			arrayTmp = make([]*Value, len(array) + this.arrayItemCount)
			//覆盖之前的切片
			copy(arrayTmp, array)
			this.value = arrayTmp
		}else{
			arrayTmp = array
		}

		//将数据统一往后面移动
		if index < 0 || index > this.arrayItemCount {
			index = this.arrayItemCount
		}
		for i := this.arrayItemCount; i > index; i-- {
			arrayTmp[i] = arrayTmp[i - 1]
		}
		arrayTmp[index] = v
		this.arrayItemCount++


	}else{
		return false
	}


	return true
}

func (this* Value)getArrayValue(index int) *Value {

	var array []*Value
	if this.t != Array {
		//创建一个数组 并把json类型设置为Array
		array = make([]*Value, arrayInitSize)
		this.t = Array
		this.value = array
		this.arrayItemCount = 0
	}

	if index >= 0 && index < this.arrayItemCount {
		return this.value.([]*Value)[index]
	}

	result := &Value{t:Null}
	this.arrayInsert(-1, result)

	return result
}


func (this* Value)Value(v ...interface{}) *Value {
	paramCount := len(v)
	if paramCount == 0 {
		return nil
	}else if paramCount == 1 {
		//参数只有1个的情况下是取值  统一返回指针 方便修改
		key := v[0]
		switch keyValue := key.(type) {
		case int:
			return this.getArrayValue(keyValue)
		case string:
			return this.getValue(keyValue)
		default:

			break
		}

		return nil
	}else{
		//统一当成两个参数解释
		key := v[0]
		switch keyValue := key.(type) {
		case int:
			//是一个数组
			if this.t == Array {
				array := this.value.([]*Value)
				if keyValue >= 0 && keyValue < this.arrayItemCount {
					valueTmp := this.createValue(v[1])
					if valueTmp != nil {
						array[keyValue] = valueTmp
					}
				}else{
					valueTmp := this.createValue(v[1])
					if valueTmp != nil {
						this.arrayInsert(keyValue, valueTmp)
					}
				}
			}else{
				valueTmp := this.createValue(v[1])
				if valueTmp != nil {
					this.arrayInsert(keyValue, valueTmp)
				}
			}
			return this
		case string:

			valueTmp := this.createValue(v[1])
			if valueTmp != nil {
				this.setValue(keyValue, valueTmp)
			}

			return this
		default:
			_=keyValue
			break
		}
	}

	//获取key对应的对象

	return nil
}

func (this* Value)IsNull(v ...interface{}) bool {
	paramCount := len(v)
	if paramCount == 0 {
		//判断当前是否是空
		if this.Type() == Null {
			return true
		}else{
			return false
		}
	}else {
		key := v[0]
		switch keyValue := key.(type) {
		case int:
			//数组判断
			if this.t == Array {
				if keyValue >= 0 && keyValue < this.arrayItemCount {
					arrayTmp := this.value.([]*Value)
					if arrayTmp[keyValue].t != Null {
						return false
					}
				}
			}
			break
		case string:
			//对象判断
			if this.t == Object {
				mapTmp := this.value.(map[string]*Value)
				if _, ok := mapTmp[keyValue]; ok {
					//如果子key存在直接返回
					if mapTmp[keyValue].t != Null {
						return false
					}
				}
			}
			break
		default:
			return true
		}
	}
	return true
}

func (this* Value)Type() Type {
	return this.t
}

func (this* Value)ArraySize() int {
	if this.t == Array {
		return this.arrayItemCount
	}
	return 0
}

func escapeString(str string) string{
	result := ""
	last := 0
	i := 0
	for ; i < len(str); i++  {
		if str[i] <= 0x22 {
			if str[i] == 0x22 {
				result += str[last:i] + `\"`
				last = i + 1
			}else if str[i] == '\b' {
				result += str[last:i] + `\b`
				last = i + 1
			}else if str[i] == '\f' {
				result += str[last:i] + `\f`
				last = i + 1
			}else if str[i] == '\n' {
				result += str[last:i] + `\n`
				last = i + 1
			}else if str[i] == '\r' {
				result += str[last:i] + `\r`
				last = i + 1
			}else if str[i] == '\t' {
				result += str[last:i] + `\t`
				last = i + 1
			}
		}else if str[i] == '\\' {
			result += str[last:i] + `\\`
			last = i + 1
		}

	}
	result += str[last:i]

	return result
}

func (this* Value)valueToString() string {

	var result string
	switch this.t {
	case Object:
		result = `{`
		i := 0
		//fmt.Println(reflect.TypeOf(this.value))
		mapSize := len(this.value.(map[string]*Value))
		for k, v := range this.value.(map[string]*Value) {
			switch v.t {
			case Null:
				result += `"` + escapeString(k) + `":` + "null"
				break
			case String:
				result += `"` + escapeString(k) + `":` + `"` + escapeString(v.value.(string)) + `"`
				break
			case Number:
				result += `"` + escapeString(k) + `":` + fmt.Sprintf("%v", v.value.(int64))
				break
			case Bool:
				result += `"` + escapeString(k) + `":` + fmt.Sprintf("%v", v.value.(bool))
				break
			case Object:
				result += `"` + escapeString(k) + `":` + v.valueToString()
				break
			case Array:
				result += `"` + escapeString(k) + `":` + v.valueToString()
				break
			}
			if i != mapSize - 1 {
				result += ","
			}
			i++
		}
		result += `}`
		break
	case Array:
		result = `[`
		arraySize := this.arrayItemCount
		for i, v := range this.value.([]*Value) {

			switch v.t {
			case Null:
				result += "null"
				break
			case String:
				result += `"` + escapeString(v.value.(string)) + `"`
				break
			case Number:
				result += fmt.Sprintf("%v", v.value.(int64))
				break
			case Bool:
				result += fmt.Sprintf("%v", v.value.(bool))
				break
			case Object:
				result += v.valueToString()
				break
			case Array:
				result += v.valueToString()
				break
			}
			if i != arraySize - 1 {
				result += ","
			}else{
				break
			}
		}
		result += `]`
		break
	}
	return result
}



func (this* Value)ToInt() int64 {

	switch this.t {
	case Number:
		return this.value.(int64)
	}

	return 0
}

func (this* Value)ToBool() bool {
	result := false

	switch this.t {
	case Number:
		if this.value.(int64) != 0 {
			result = true
		}
		break
	case Bool:
		result =this.value.(bool)
		break
	}

	return result
}

func (this* Value)ToString() string {

	var result string
	switch this.t {
	case Null:
		result = ""
		break
	case String:
		result = this.value.(string)
		break
	case Number:
		result = fmt.Sprintf("%v", this.value.(int64))
		break
	case Bool:
		result = fmt.Sprintf("%v", this.value.(bool))
		break
	case Object:
		result = this.valueToString()
		break
	case Array:
		result = this.valueToString()
		break
	}

	return result
}


/**********************************************************解析*****************************************************************/

func parseDoc(doc interface{}) *Value {

	result := New()
	//只会传进来数组或者对象
	switch docValue := doc.(type) {
	case map[string]interface{}:
		//如果是map
		m := map[string]*Value{}

		for key, value := range docValue {
			//fmt.Printf("%s  %v\n", key, reflect.TypeOf(value))
			switch keyValue := value.(type) {
			case map[string]interface{}:
				m[key] = parseDoc(keyValue)
				break
			case []interface{}:
				m[key] = parseDoc(keyValue)
				break
			case float64:
				m[key] = result.createValue(int(keyValue))
				break
			case string:
				m[key] = result.createValue(keyValue)
				break
			case bool:
				m[key] = result.createValue(keyValue)
			case nil:
				m[key] = result.createValue(keyValue)
				break
			}
		}
		if len(m) != 0 {
			result.value = m
			result.t = Object
		}
		break
	case []interface{}:
		//如果是数组
		a := make([]*Value, len(docValue))
		for i, value := range docValue  {
			switch keyValue := value.(type) {
			case map[string]interface{}:
				a[i] = parseDoc(keyValue)
				break
			case []interface{}:
				a[i] = parseDoc(keyValue)
				break
			case float64:
				a[i] = result.createValue(int64(keyValue))
				break
			case string:
				a[i] = result.createValue(keyValue)
				break
			case bool:
				a[i] = result.createValue(keyValue)
				break
			case nil:
				a[i] = result.createValue(keyValue)
				break
			}
		}
		if len(a) != 0 {
			result.value = a
			result.t = Array
			result.arrayItemCount = len(a)
		}
		break
	}
	return result
}
//暂时不自己做解析 用golang公开库
func parse(jsonString string) *Value {

	i := 0
	result := New()
	typeTmp := Null
	for ; i < len(jsonString); i++ {
		if jsonString[i] == '[' {
			typeTmp = Array
			break
		}else if jsonString[i] == '{' {
			typeTmp = Object
			break
		}
	}

	if typeTmp == Array {
		var a []interface{}
		err := json.Unmarshal([]byte(jsonString[i:]), &a)
		if err != nil {
			return result
		}

		return parseDoc(a)

	}else if typeTmp == Object {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(jsonString[i:]), &m)
		if err != nil {
			return result
		}
		return parseDoc(m)
	}

	return result
	


}