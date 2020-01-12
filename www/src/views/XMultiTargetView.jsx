import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React from 'react';
import { Card, Container, Loader, Menu } from 'semantic-ui-react';
import { XErrorMessage } from '../components/messages';
import { XTargetCreateModal } from '../components/target';
import XTargetCard from '../components/target/XTargetCard';

export const MULTI_TARGET_QUERY = gql`
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
    const { called, loading, error, data } = useQuery(MULTI_TARGET_QUERY);

    const showCards = () => {
        if (!data || !data.targets || data.targets.length < 1) {
            return (
                <h1>No targets found!</h1>
            );
        }
        return (<Card.Group centered itemsPerRow={4}>
            {data.targets.map(target => (<XTargetCard key={target.id} {...target} />))}
        </Card.Group>);
    };


    return (
        <Container style={{ padding: '10px' }}>
            <Menu secondary>
                <Menu.Item position='right'><XTargetCreateModal /></Menu.Item>
            </Menu>
            <Container fluid style={{ padding: '20px' }}>
                <Loader disabled={!called || !loading} />
                <XErrorMessage title='Error Loading Targets' err={error} />

                {showCards()}
            </Container>
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