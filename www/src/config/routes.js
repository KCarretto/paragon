import React from 'react';
import { Icon } from 'semantic-ui-react';

const routes = [
    {
        title: 'News Feed',
        link: '/news_feed',
        icon: <Icon name='newspaper' />,
    },
    {
        title: 'Targets',
        link: '/targets',
        icon: <Icon name='desktop' />,
    },
    {
        title: 'Jobs',
        link: '/jobs',
        icon: <Icon name='cubes' />
    },
    {
        title: 'Tags',
        link: '/tags',
        icon: <Icon name='tags' />,
    },
    {
        title: 'Payloads',
        link: '/payloads',
        icon: <Icon name='gift' />,
    },
    {
        title: 'Profile',
        link: '/profile',
        icon: <Icon name='user secret' />,
    }
]

export default routes;