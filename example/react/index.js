import React from 'react';  
import ReactDOM from 'react-dom/client';
import {BrowserRouter,Routes,Route} from 'react-router-dom'

import Admin from './test-admin.js';
import Login from './test-login.js';
import Read from './test-read.js';


const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <BrowserRouter>
     <Routes>
       <Route path = '/admin' element = {<Admin />} />
       <Route path = '/login' element = {<Login />} />
       <Route path = '/read' element = {<Read />} /> 
     </Routes>
    </BrowserRouter>
)


