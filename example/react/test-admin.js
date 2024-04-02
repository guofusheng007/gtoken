import React from 'react';
import {Link } from 'react-router-dom'

export default function Admin() {
    return (
        <div>
            <Link to="/read">查看记录</Link><br />
            <Link to="/write">增改删</Link><br />
        </div>
    )
}