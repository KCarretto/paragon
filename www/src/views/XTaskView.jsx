import React from 'react';
import { useParams } from 'react-router-dom';

const XTaskView = () => {
    let { id } = useParams();

    return (
        <h1>TASK: {id}</h1>
    );
}

export default XTaskView;
