import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React from 'react';
import { useParams } from 'react-router-dom';
import { Loader } from 'semantic-ui-react';

const TARGET_QUERY = gql`
    query Target($id: ID!) {
    target(id: $id) {
        id
        name
    }
  }
`;


const XTargetView = () => {
    let { id } = useParams();

    const { loading, error, data } = useQuery(TARGET_QUERY, {
        variables: { id },
    });
    if (loading) return (<Loader active />);
    if (error) return (`${error}`);

    return (
        <div>
            <h1>TARGET: {id}</h1>
            <h2>{data}</h2>
        </div>
    );
}

export default XTargetView;
