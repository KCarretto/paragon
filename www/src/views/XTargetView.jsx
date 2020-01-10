import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import moment from 'moment';
import React from 'react';
import { useParams } from 'react-router-dom';
import { Card, Container, Icon, Label, Loader } from 'semantic-ui-react';
import { XCredentialSummary } from '../components/credential';
import { XTaskSummary } from '../components/task';



const TARGET_QUERY = gql`
    query Target($id: ID!) {
    target(id: $id) {
        id
        name
        primaryIP
        machineUUID
        publicIP
        primaryMAC
        hostname
        lastSeen
        tasks {
            id
            queueTime
            claimTime
            execStartTime
            execStopTime
            error
            job {
                id
                name
            }
        }
        tags {
            id
            name
        }
        credentials {
            id
            principal
            secret
            fails
        }
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
        <Container fluid style={{ padding: '20px' }}>
            <Card fluid centered>
                <Card.Content>
                    <Card.Header>{data.target.name}</Card.Header>
                    {
                        moment(data.target.lastSeen).isBefore(moment().subtract(5, 'minutes')) ?
                            <Label corner='right' size='large' icon='times circle' color='red' />
                            : <Label corner='right' size='large' icon='check circle' color='green' />
                    }
                    <Card.Meta>
                        <a>
                            <i aria-hidden="true" className="clock icon"></i>
                            Last Seen: {data.target.lastSeen ? moment(data.target.lastSeen).fromNow() : 'Never'}<br />
                        </a>
                        {data.target.primaryIP ?
                            <a>
                                <i aria-hidden="true" className="user icon"></i>
                                Primary IP: {data.target.primaryIP}<br />
                            </a>
                            :
                            <div></div>
                        }
                        {data.target.hostname ?
                            <a>
                                <i aria-hidden="true" className="user icon"></i>
                                Hostname: {data.target.hostname}<br />
                            </a>
                            :
                            <div></div>
                        }
                        <Icon name='tags' /> {data.target.tags && data.target.tags.length != 0 ? data.target.tags.map(tag => tag.name).join(', ') : 'None'}
                    </Card.Meta>
                    <Card.Description>
                        <XTaskSummary tasks={data.target.tasks} />
                        <XCredentialSummary credentials={data.target.credentials} />
                    </Card.Description>
                </Card.Content>
                {data.target.primaryMAC || data.target.publicIP || data.target.machineUUID ?
                    <Card.Content extra>
                        {data.target.primaryMAC ?
                            <a>
                                <i aria-hidden="true" className="user icon"></i>
                                Primary MAC: {data.target.primaryMAC}<br />
                            </a>
                            :
                            <div></div>
                        }
                        {data.target.publicIP ?
                            <a>
                                <i aria-hidden="true" className="user icon"></i>
                                Public IP: {data.target.publicIP}<br />
                            </a>
                            :
                            <div></div>
                        }
                        {data.target.machineUUID ?
                            <a>
                                <i aria-hidden="true" className="user icon"></i>
                                MachineUUID: {data.target.machineUUID}<br />
                            </a>
                            :
                            <div></div>
                        }
                    </Card.Content>
                    :
                    <div></div>
                }
            </Card>
        </Container>
    );
}

export default XTargetView;
