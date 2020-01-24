import React, { useState } from 'react';
import { Button, Form, Grid, Header, Image, Input, Segment } from 'semantic-ui-react';

const XLogin = () => {
    const [name, setName] = useState('');

    const handleSignUp = (e) => {
        e.preventDefault();
        fetch(window.location.origin + '/oauth/signup', {
            method: "POST",
            body: JSON.stringify({ username: name })
        }).then(resp => {
            console.log("Fetch response", resp);
            return resp.json()
        }, err => {
            console.log("Fetch failed!", err);
            alert("Fetch failed!" + err);
        }).then(data => {
            console.log("JSON", data)
            window.location.href = data.url;
        }, err => {
            console.log("JSON failed!", err);
            alert("Failed!" + err);
        });
    }

    return (
        <Grid textAlign='center' style={{ height: '100vh' }} verticalAlign='middle'>
            <Grid.Column style={{ maxWidth: 450 }}>
                <Header as='h2' color='blue' textAlign='center'>
                    <Image src='/logo512.png' /> Log-in or Sign up
            </Header>
                <Form size='large'>
                    <Segment stacked>
                        <Input
                            fluid
                            icon='user'
                            iconPosition='left'
                            placeholder='Choose a Username'
                            value={name}
                            onChange={(e, { value }) => setName(value)}
                        />
                        <Button icon='google' color='green' fluid size='large' onClick={handleSignUp}>
                            Sign Up
                    </Button>
                    </Segment>
                </Form>
            </Grid.Column>
        </Grid>
    );
}

export default XLogin