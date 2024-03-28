import React from 'react';  
import ReactDOM from 'react-dom/client';
import {BrowserRouter,Routes,Route} from 'react-router-dom'

import Login from './test-login.js';
import UpdateToken from './test-updatetoken.js';


const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <BrowserRouter>
     <Routes>
       <Route path = '/' element = {<Home />} />
       <Route path = '/login' element = {<Login />} />
       <Route path = '/token' element = {<UpdateToken />} />
     </Routes>
    </BrowserRouter>
)

