import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React from 'react';
import { useParams } from 'react-router-dom';
import { Loader } from 'semantic-ui-react';

const TASK_QUERY = gql`
    query Task($id: ID!) {
    task(id: $id) {
        id
        output
    }
  }
`;

const XTaskView = () => {
    let { id } = useParams();

    const { loading, error, data } = useQuery(TASK_QUERY, {
        variables: { id },
    });
    if (loading) return (<Loader active />);
    if (error) return (`${error}`);

    return (
        <div>
            <h1>TASK: {id}</h1>
            <h2>{data}</h2>
        </div>
    );
}

export default XTaskView;
