import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React from 'react';
import { Container, Label, Loader, Menu } from 'semantic-ui-react';
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
            <Label.Group tag size='big'>
                {data.tags.map((tag, index) =>
                    <Label key={index}>{tag.name}</Label>
                )}
            </Label.Group>
        );
    };

    return (
        <Container style={{ padding: '10px' }}>
            <Menu secondary>
                <Menu.Item position='right'><XTagCreateModal /></Menu.Item>
            </Menu>
            <Container fluid style={{ padding: '20px' }}>
                <Loader disabled={!called || !loading} />
                {showList()}
                <XErrorMessage title='Error Loading Tags' err={error} />
            </Container>
        </Container>
    );
}

export default XMultiTagView;