import gql from 'graphql-tag';
import React from 'react';
import { Query } from 'react-apollo';
import { Card, Container } from 'semantic-ui-react';
import XTargetCard from './XTargetCard';

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

const XTargetCardGroup = () => {
    const now = Math.floor(Date.now() / 1000);

    const targets = [
        {
            id: 1,
            name: 'Team 1 - Web',
            primaryIP: '10.1.1.10',
            lastSeen: now - 20,
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
            tasks: [
                {
                    id: 100,
                    queueTime: now - 300,
                    claimTime: now - 240,
                    execStartTime: now - 120,
                    execStopTime: now - 10,
                    job: {
                        id: 500,
                        name: 'Deployment (initial)',
                    }
                },
                {
                    id: 101,
                    queueTime: now - 1020,
                    claimTime: now - 740,
                    execStartTime: now - 620,

                    job: {
                        id: 501,
                        name: 'User Snapshot',
                    }
                },
                {
                    id: 102,
                    queueTime: now - 754,
                    claimTime: now - 600,
                    execStartTime: now - 520,
                    execStopTime: now - 412,
                    error: "error: No such file or directory '/root/.bash_history'",
                    job: {
                        id: 502,
                        name: 'Capture Bash History',
                    }
                },
                {
                    id: 102,
                    queueTime: now - 754,
                    job: {
                        id: 502,
                        name: 'Disable MySQL',
                    }
                },
            ]
        },
        {
            id: 2,
            name: 'Team 2 - AD',
            primaryIP: '10.2.1.60',
            lastSeen: now - 120,
            tags: [
                {
                    id: 11,
                    name: 'windows'
                },
                {
                    id: 22,
                    name: 'team-2'
                }
            ],
            tasks: [
                {
                    id: 111,
                    queueTime: now - 300,
                    claimTime: now - 240,
                    execStartTime: now - 120,
                    job: {
                        id: 500,
                        name: 'AD Flush DNS',
                    }
                },
                {
                    id: 112,
                    queueTime: now - 200,
                    job: {
                        id: 500,
                        name: 'Capture Screenshot',
                    }
                },
                {
                    id: 113,
                    queueTime: now - 300,
                    claimTime: now - 30,
                    job: {
                        id: 500,
                        name: 'Change Desktop Background',
                    }
                },
            ]
        },
        {
            id: 3,
            name: 'Team 2 - Web',
            primaryIP: '10.2.1.10',
            lastSeen: now - 624,
            tags: [
                {
                    id: 10,
                    name: 'linux'
                },
                {
                    id: 22,
                    name: 'team-2'
                }
            ],
        }
    ]
    return (
        <Container fluid style={{ padding: '20px' }}>
            <Card.Group centered itemsPerRow={4}>
                <Query query={TARGETS_QUERY}>
                    {() => targets.map((target) => <XTargetCard key={target.id} {...target} />)}
                </Query>
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

export default XTargetCardGroup