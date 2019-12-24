import React from 'react';
import { useParams } from 'react-router-dom';

const XTargetView = () => {
    let { id } = useParams();

    return (
        <h1>TARGET: {id}</h1>
    );
}

export default XTargetView;