import React from 'react';
import { ApolloProvider } from 'react-apollo';
import { Route } from 'react-router-dom';
import 'semantic-ui-css/semantic.min.css';
import './App.css';
import { XLayout, XUnimplemented } from './components/layout/';
import XClient from './config/apollo';
import Routes from './config/routes';
import { XMultiJobView, XMultiTargetView, XTargetView, XTaskView } from './views';



function App() {
  return (
    <ApolloProvider client={XClient}>
      <XLayout routeMap={Routes} className='App'>
        <Route exact path='/' component={XUnimplemented} />
        <Route exact path='/targets' component={XMultiTargetView} />
        <Route exact path='/jobs' component={XMultiJobView} />

        <Route path='/targets/:id' component={XTargetView} />
        <Route path='/tasks/:id' component={XTaskView} />

      </XLayout>
    </ApolloProvider>
  );

}

export default App;

