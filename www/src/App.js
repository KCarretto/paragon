import React from 'react';
import { Route } from 'react-router-dom';
import 'semantic-ui-css/semantic.min.css';
import { CardGroup, Icon } from 'semantic-ui-react';
import './App.css';
import { XTaskCard } from './components/cards';
import { XLayout, XUnimplemented } from './components/layout/';

const routes = [
  {
    title: 'Dashboard',
    icon: <Icon name='dashboard' />,
    link: '/dashboard',
    routes: [
      <Route key={0} exact path='/' component={XUnimplemented} />,
      <Route key={1} path='/dashboard' component={XUnimplemented} />,
    ],
  },
  {
    title: 'Tasks',
    icon: <Icon name='tasks' />,
    link: '/tasks',
    routes: [
      <Route key={0} path='/tasks'>
        <CardGroup>
          <XTaskCard name='List Connections' tags={['Windows', 'AD']} />
          <XTaskCard />
        </CardGroup>
      </Route>,
    ],
  },
  {
    title: 'Groups',
    icon: <Icon name='group' />,
    link: '/groups',
    routes: [
      <Route key={0} path='/groups' component={XUnimplemented} />,
    ],
  },
  {
    title: 'Payloads',
    icon: <Icon name='gift' />,
    link: '/payloads',
    routes: [
      <Route key={0} path='/payloads' component={XUnimplemented} />,
    ],
  },
]
function App() {
  return (
    <XLayout routeMap={routes} className='App' />
  );

}

export default App;

