import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React from 'react';
import { Container, Icon, List, Loader, Menu } from 'semantic-ui-react';
import { XErrorMessage } from '../components/messages';
import { XTagCreateModal } from '../components/tag';

export const MULTI_TAG_QUERY = gql`
    {
        tags {
            id
            name
        }
    }
`;

const XMultiTagView = () => {
    const { called, loading, error, data } = useQuery(MULTI_TAG_QUERY);

    const showList = () => {
        if (!data || !data.tags || data.tags.length < 1) {
            return (
                // TODO: Better styling
                <h1>No tags found!</h1>
            );
        }
        return (
            <List>{data.tags.map(({ name }) =>
                <List.Item>
                    <Icon name='tag' />
                    <List.Content>{name}</List.Content>
                </List.Item>
            )}</List>
        );
    };

    return (
        <Container style={{ padding: '10px' }}>
            <Menu secondary>
                <Menu.Item position='right'><XTagCreateModal /></Menu.Item>
            </Menu>
            <Container fluid style={{ padding: '20px' }}>
                <Loader disabled={!called || !loading} />
                <XErrorMessage title='Error Loading Tags' err={error} />
                {showList()}
            </Container>
        </Container>
    );
}

export default XMultiTagView;