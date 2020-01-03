import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React from 'react';
import { Card, Container, Loader } from 'semantic-ui-react';
import XTargetCard from '../components/target/XTargetCard';

const TARGETS_QUERY = gql`
{
targets {
        id
        name
        primaryIP
        lastSeen

        tags {
            id
            name
        }

        tasks {
            id

            queueTime
            claimTime
            execStartTime
            execStopTime

            job {
                id
                name
            }
        }
    }
}`

const XMultiTargetView = () => {
    const now = Math.floor(Date.now() / 1000);

    const { loading, error, data } = useQuery(TARGETS_QUERY);
    console.log(loading)
    console.log(error)
    console.log(data)

    if (loading) return (<Loader active />);
    if (error) return (`${error}`);
    if (!data || !data.targets || data.targets.length < 1) {
        return (<h1>No targets found!</h1>);
    }

    return (
        <Container fluid style={{ padding: '20px' }}>
            <Card.Group centered itemsPerRow={4}>
                {data.targets.map(target => (<XTargetCard key={target.id} {...target} />))}
            </Card.Group>
        </Container>
    );
}

// XTargetCardGroup.propTypes = {
//     targets: PropTypes.arrayOf(PropTypes.shape({
//         id: PropTypes.number.isRequired,
//         name: PropTypes.string.isRequired,
//         tags: PropTypes.arrayOf(PropTypes.string),
//     })).isRequired,
// };

export default XMultiTargetView