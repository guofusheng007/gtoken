package main

import (
	"fmt"
	"time"
	"github.com/guofusheng007/gtoken"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

type TestController struct {
	beego.Controller
}

//Token配置
var key = []byte("012345678901234511111111")   // token共享key,其长度为16或24或32个字符切片

//-----------------用户端提供的信息-----------------------
type UserCredentials struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
  }
  //存储页面录入的信息
var user UserCredentials = UserCredentials{}


//基本配置
func init() {
	//webadmin
	beego.BConfig.Listen.EnableAdmin = true 
	beego.BConfig.Listen.AdminAddr = "localhost"
    beego.BConfig.Listen.AdminPort = 8088

	//BindXXX网页数据解析
	beego.BConfig.CopyRequestBody = true
	//关闭自动渲染
	beego.BConfig.WebConfig.AutoRender = false


	//CORS安全配置
	var corsconf = cors.Options {
		AllowAllOrigins: true,
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		AllowMethods:[]string{"GET","POST","PUT","DELETE","OPTIONS"},
		AllowHeaders:[]string{
			"Origin",
			"Authorization",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"token",
			"tokenid",
		},
		ExposeHeaders:[]string{
			"Content-Length",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"token",
			"tokenid",
		},
	}
	//所有url需要命中该中间函数，即打开url前先运行该配置。
	beego.InsertFilter("/*", beego.BeforeRouter, cors.Allow(&corsconf))
}

func main() {
	//认证处理/token更新
	var auth = new(TestController)
	beego.Router("/login", auth, "post:Login_v2")
	beego.Router("/updatetoken", auth, "get:UpdateToken")   //查看cookies和headers

	//运行beego服务
	beego.Run()
}
//------------------------------------------------

//API 用户login
//认证后，把凭证写入cookie或session
func (c *TestController) Login_v2() {  //react提交post用户帐号和密码以加密方式提交给接口
	//react提交post用户帐号和密码
	c.BindJSON(&user)
	fmt.Printf("%+v\n",user) //查看用户端提供的body数据
	//c.ShowHeadCookie()     //查看header、cookie.回显到服务端console口

	//通过验证后产生新token返回组用户
	if ((user.Username == "guofs") && (user.Password == "123321")) {
		//token信息，如下值不能作为公共变量，否则时间不会变化。
		var token = gtoken.Gtoken{
			UserID: 10,
			UserName: "guofs",
			UserEmail: "guofs@139.com",
			UserMobile: "13700000000",
			ExpiresAt: time.Now().Add(time.Minute * 6),   //第一次签发token时Token的TTL为：5分钟
			IssuedAt: time.Now(),
			Issuer: "guofs",
			Subject: "webtest",
			TokenID: gtoken.RandomString(12),
		}
		//fmt.Println(token.TokenID)
		token_enc := gtoken.CreateToken(&token,key)
		c.Ctx.Output.Header("Token",token_enc)  
		c.Ctx.WriteString(`{"info":"验证成功"}`)
		//fmt.Printf("token:%#v\n",token)
		////fmt.Printf("token:%v\n",token)
		
	} else {
		c.Ctx.Output.Header("Token","")
		c.Ctx.WriteString(`{"info":"验证失败"}`)
	}
}

//供react客户端刷新token(即续签过程)
func (c *TestController) UpdateToken() {
	//从headers中获取客户端提取token认息
	token_enc := c.Ctx.Request.Header.Get("Token")
	token_id := c.Ctx.Request.Header.Get("Tokenid")
	//fmt.Printf("旧Token:%v\n",token_enc)
	//fmt.Printf("旧TokenID:%v\n",token_id)

	//判断token的有效性
	if ((token_enc == "undefined") || (token_enc == "")) {
		fmt.Println("Token已过期,不能续签,请重新认证后产生新Token")
		c.Ctx.Output.Header("Token","")  
		c.Ctx.WriteString(`{"info":"Token已过期,不能续签,请重新认证后产生新Token"}`) 
		return
	} 

	//查验客户端提交的token是否合法，若不能解析，则直接跳出。
	oldToken, ExpiresTime,err := gtoken.CheckToken(token_enc,key) 
	if err != nil {
		c.Ctx.Output.Header("Token","")  
		c.Ctx.WriteString(`{"info":"Token非法,不能解析,请重新认证后产生新Token"}`) 
		fmt.Println("Token非法,不能解析,请重新认证后产生新Token")
		return
	} 
	//fmt.Printf("旧Token将在 %d 分钟后过期\n",ExpiresTime)

	//token失效处理
	if (ExpiresTime <= 0) {
		fmt.Println("Token已过期,不能续签,请重新认证后产生新Token")
		c.Ctx.Output.Header("Token","")  
		c.Ctx.WriteString(`{"info":"Token已过期,不能续签,请重新认证后产生新Token"}`) 
		return
	} 

	//当用户提交 tokenID与token密文中的tokenid不相同时，禁止更新token.
	if (token_id != oldToken.TokenID ) {
		//该check是为了token泄露后而造在的安全风险。
		c.Ctx.Output.Header("Token",token_enc) 
		c.Ctx.WriteString(`{"info":"用户提交的 TokenID 有误"}`)
		fmt.Println("用户提交的 TokenID 有误")
		return
	}

	//当ttl > 2 时暂时不需要update token
	if (ExpiresTime > 2 ) {
		//该check是为了减少不必要的计算，减省cpu计算压力
		c.Ctx.Output.Header("Token",token_enc) 
		c.Ctx.WriteString(`{"info":"Token TTL大于2分钟,暂时不需要续签"}`)
		fmt.Println("Token TTL大于2分钟,暂时不需要续签")
		return
	} 


	//token更新。当token的TTL小于或等2分钟时，需及时续签。
	if ((ExpiresTime <= 2 ) && (token_id == oldToken.TokenID )) {
		//token信息
		var token = gtoken.Gtoken{
			UserID: 10,
			UserName: "guofs",
			UserEmail: "guofs@139.com",
			UserMobile: "13700000000",
			ExpiresAt: time.Now().Add(time.Minute * 5),  //续签时Token的TTL为：5分钟
			IssuedAt: time.Now(),
			Issuer: "guofs",
			Subject: "webtest",
			TokenID: gtoken.RandomString(12),
		}
		fmt.Println(token.TokenID)
		//新token密文
		token_enc_New := gtoken.CreateToken(&token,key)
		fmt.Println("新token密文:",token_enc_New)
		//返回新token
	    c.Ctx.Output.Header("Token",token_enc_New)  
	    c.Ctx.WriteString(`{"info":"Token更新成功"}`) 
	} 
}


//遍历所有Header和cookie
func (c *TestController) ShowHeadCookie() {
	//header打印
	//单一值
	//fmt.Println("Content-Type:", c.Ctx.Request.Header.Get("Content-Type"))
    //遍历所有Header
	fmt.Println("-------------Header-------------------")
	for k, v := range c.Ctx.Request.Header {
		fmt.Printf("k:%v,v:%+v\n", k, v)
	}
    //cookie查看
	fmt.Println("-------------Cookie-------------------")
	//单一值
	//ck, _ := c.Ctx.Request.Cookie("username")
	//遍历所有cookie
	for _, v := range c.Ctx.Request.Cookies() {
		// fmt.Printf("%v=%v,%s\n", v.Name,v.Value,v.Expires.Format("2006-01-02 15:04:05"))
		fmt.Printf("%v=%v\n", v.Name,v.Value)
	}
}

