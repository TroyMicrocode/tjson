# tjson
<<<<<<< HEAD
//暂时没写解析代码 
=======
>>>>>>> 50585fac277c348667490d3363dc3afcc632806d

主要函数Value  可变参数 1个到2个 3个以上当2个处理
参数1个为获取值 string 或 int
参数两个为设置值 第一个参数还是 string 或 int  第二个参数可以是string int bool object array
当调用Value时 设置值或者获取值的时候 最好是判断类型，如果不判断直接调用，遇到值不匹配的情况下 ，被强制清0 并且强制转换成你正在调用的类型

判断是否为空的函数 IsNull
这个函数有两种形式 不带参数 和 带一个key或者index参数
比如jsonObject.IsNull() 是判断当前对象是不是Null
jsonObject.IsNull(string("xxx")) 或者 sonObject.IsNull(int(xxx)) 判断当前对象的xxx数据是否为null
但是也可以这样判断 那string举例 jsonObject.Value("xxx").IsNull()  和上面的差不多 。唯一区别是不管这个xxx是否为null 都会创建一个null对象 上面那种则不会


例子

  var jsonTest tjson.Value = tjson.Value{}
	var jsonTest2 tjson.Value = tjson.Value{}

	//如果是数组 -1就是插入到末尾
	jsonTest2.Value("aaa").Value(-1, "111")
	jsonTest2.Value("aaa").Value(-1, "222")
	jsonTest2.Value("aaa").Value(2, "333")
	jsonTest2.Value("aaa").Insert(1, "ins11") //插入到中间
	jsonTest2.Value("aaa").Remove(0) //把第一个删除了

	jsonTest.Value("111").Value("s111", "v1")
	jsonTest.Value("111").Value("s222", "v2")
	jsonTest.Value("111").Value("s333", "v3")
	jsonTest.Value("111").Remove("s222") //把s222删除掉
	jsonTest.Value("222", "v2")
	jsonTest.Value("333", true)
	jsonTest.Value("obj", jsonTest2)

	jsonStr := jsonTest.ToString()
	_=jsonStr
  
<<<<<<< HEAD
  //jsonStr得到下面这个字符串 
=======
  //jsonStr得到下面这个字符串 转出来是紧凑排列的 这里为了方便观察 在json网页上转了下
>>>>>>> 50585fac277c348667490d3363dc3afcc632806d

{
    "111":{
        "s111":"v1",
        "s333":"v3"
    },
    "222":"v2",
    "333":true,
    "obj":{
        "aaa":[
            "ins11",
            "222",
            "333"
        ]
    }
}
