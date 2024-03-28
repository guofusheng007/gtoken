import cookie from 'react-cookies'
import axios from 'axios'
import { Decrypt } from './crypto.js';

//更新token
function UpdateToken() {
    //react解密使用，由后端提供。前后端务必保持一致。
    const keyDec = '012345678901234511111111'    

    //从cookie中读取token
    const token = cookie.load('TOKEN')
    //判断token是否存在，若不存在，就退出
    if ((token === undefined) || (token === '')) {
        console.log('cookie中没有token记录。请跳转到login页面重新生成token');
        return
    } 

    //若token存在，则解密token
    const txt_token = JSON.parse(Decrypt(token,keyDec))

    //const ExpiresAt = "2024-03-29T11:12:37.4536506+08:00";
    const tokenExpires = new Date(txt_token.expire);
    const timeDiff = Math.round((tokenExpires.valueOf() - (new Date()).valueOf()) / (60*1000))
    console.log("解密当前token:",txt_token)
    //console.log("username:",txt_token.username)

    //当ttl > 2 时，暂时不需要更新，减少服务端计算压力。
    if (timeDiff > 2) {
        console.log(`token TTL ${timeDiff} 分钟，暂时不用更新.当ttl小于或等2时才更新`)
        return
    } 

    //当ttl <= 0 时，token已失效，直接跳转到login页面。
    if (timeDiff <= 0) {
        console.log(`token已失效,请跳转到login页面重新生成token。`)
        return
    } 

    //更新token的模块,如下axios是一个整体。
    axios({
        url:'http://192.168.3.110:8080/updatetoken',
        method: 'get',    //get或post等html方法。不区分大小写
        headers: {
           'Content-Type':'application/json',
           'Token': token,
           'Tokenid': txt_token.tokenid,
        },
    }).then(response => {
        if(response && response.status === 200){
            const token =  response.headers['token'];   
            var ttl = new Date(new Date().getTime() + 24*60*60*1000)
            cookie.save('TOKEN', token,{expires: ttl });      
            console.log("token:",token) 
            //console.log(response.data) 
        }else{
           // 响应失败
           console.log('Failed to submit data');
        }
    })

 }

 export default UpdateToken