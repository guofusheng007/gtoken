import React, { useEffect, useState} from 'react';
import {Link,useNavigate } from 'react-router-dom'
import { UpdateToken } from './test-update.js';



//--------动态态数据(从api中读取)-------------------------------
export default function Read() {
  //--------验证token---------------------
  const navigate = useNavigate();
  //调用updatetoken函数
  var auth = UpdateToken()
  console.log("authInfo",auth)
  //识别返回信息
  if ( !auth.Auth ) {
     console.log("认证过期,",auth.Msg)
     navigate("/login",{ replace: false });     //跳转到login页面,重新生成token
  } else {
    console.log("认证未过期,",auth.Msg)
  }
 
  //--------业务模块---------------------
  //定义状态变量
  const [data, setData] = useState([]);
  //更新状态，即单击按扭时执行一次fetch
  function change() {
    // 异步获取数据的逻辑
    fetch('http://192.168.3.110:8080/read')
      .then(response => response.json())
      .then(data => {
        setData(data);
        document.getElementById('opmsg').textContent = ''
      })
      .catch(error => {
        console.error('Error:', error);
        document.getElementById('opmsg').textContent = 'error,超时，可能网络异常'
    });
  }

  //当App被调用一次时，useEffect就执行一次fetch
  // 或
  useEffect(() => {
    change()
  }, []);
 


  return (
    <div>
        <Link to="/">Home</Link><br />
        <button onClick={change}> 刷新 </button><p id="opmsg"></p>
        {data.map( (user) => (
                       <div key={user.id}>
                        {user.name},{user.tel}
                       </div>
        ))}
    </div>
  );
};