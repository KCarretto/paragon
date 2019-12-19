import gql from 'graphql-tag';
import React from 'react';
import { Query } from 'react-apollo';
import { Card, Container, Menu } from 'semantic-ui-react';
import { XJobQueueModal } from '../components/job';

const JOBS_QUERY = gql`
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

const XMultiJobView = () => {
    const jobs = [
        {
            id: 1,
            name: 'Deployment (initial)',
            tags: [
                {
                    id: 10,
                    name: 'linux'
                },
                {
                    id: 21,
                    name: 'team-1'
                }
            ],
        }
    ]
    return (
        <div style={{ padding: '10px' }}>
            <Menu secondary>
                <Menu.Item position='right'><XJobQueueModal /></Menu.Item>
            </Menu>
            <Container fluid style={{ padding: '20px' }}>
                <Card.Group centered itemsPerRow={4}>
                    <Query query={JOBS_QUERY}>
                        {() => jobs.map((job) => <XJobTemplateCard key={job.id} {...job} />)}
                    </Query>
                </Card.Group>
            </Container>
        </div>
    );
}

// XTargetCardGroup.propTypes = {
//     targets: PropTypes.arrayOf(PropTypes.shape({
//         id: PropTypes.number.isRequired,
//         name: PropTypes.string.isRequired,
//         tags: PropTypes.arrayOf(PropTypes.string),
//     })).isRequired,
// };

export default XMultiJobView