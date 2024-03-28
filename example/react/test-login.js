import React from 'react';
import {Link} from 'react-router-dom'
import { useNavigate } from 'react-router-dom'
import cookie from 'react-cookies'

export default function Login() {
  const Token = "react_submit_token_8483"
  const navigate = useNavigate();
  function handleSubmit(event) {
    event.preventDefault();   // 阻止表单提交
    const jsonData = 
       {
         username: document.getElementById('username').value,
         password: document.getElementById('password').value,
       }
    //console.log(JSON.stringify(jsonData, null, 2)); //输出到浏览器的console接口
  
    //识别提交的按扭
    var url = ''
    if (event.target.id === 'login') {
        url = 'http://192.168.3.110:8080/login'
        console.log("login被单击")
      }
    if (event.target.id === 'register') {
        url = 'http://192.168.3.110:8080/login'
        console.log("register被单击")
    }
  
    //写入数据
    fetch(url,{
        method: 'POST',
        headers: {
         'Content-Type': 'application/json',
         'Token': Token,
        },
         body: JSON.stringify(jsonData)
    })
    .then((response) => {
       if (!response.ok) {
           console.log('Failed to submit data');
            throw new Error('Network response was not ok.');
        }
        //提取后台提供的token,并写入cookie
        const token = response.headers.get('token')
        console.log("token值:",token)
        var ttl = new Date(new Date().getTime() + 24*60*60*1000)   //1000表示1秒。
        cookie.save('TOKEN', token,{expires: ttl });
        return response.json()  //将url的返回body信息转换为json格式
    })
    .then((data) => {       //此部分测试正常
        console.log(data)   //将返回信息回显到浏览器的console接口
        document.getElementById('opmsg').textContent = data.info   //将post时返回信息显示到指定标签上。
        if (data.info === '验证成功') {
          document.getElementById('opmsg').textContent = "验证成功,请跳转到管理admin页面"
          navigate("/admin",{ replace: false });   //跳转到管理admin页面
        }
    })
    .catch(error => {
      console.error('Error:', error);
      document.getElementById('opmsg').textContent = '后台服务异常若其它'
    });
    // console.log("test2") 此行的执行次序早于上面的 console.log("test1")
  }

  return (
    <div>
      <Link to="/">Home</Link><br />
      <form action="">
        用户名称:<input type="text" id="username" /><br/>
        用户密码:<input type="password" id="password" /><br/>
        <button id="login" type="submit" onClick={handleSubmit}>登录</button>
        <button id="register" type="submit" onClick={handleSubmit}>注册</button>
      </form>
      <p id="opmsg"></p> {/* 显示提示信息 */}
    </div>
  )
 }
