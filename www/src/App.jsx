import React from 'react';
import { Route } from 'react-router-dom';
import 'semantic-ui-css/semantic.min.css';
import './App.css';
import { XFileUploadModal } from './components/file';
import { XLayout, XUnimplemented } from './components/layout/';
import Routes from './config/routes';
import { XGraphProvider } from './graphql';
import { XJobView, XMultiJobView, XMultiTagView, XMultiTargetView, XTargetView, XTaskView } from './views';

const App = () => (
  <XGraphProvider>
    <XLayout routeMap={Routes} className='App'>
      <Route exact path='/' component={XUnimplemented} />
      <Route exact path='/news_feed' component={XUnimplemented} />
      <Route exact path='/profile' component={XUnimplemented} />

      <Route exact path='/targets' component={XMultiTargetView} />
      <Route exact path='/jobs' component={XMultiJobView} />
      <Route exact path='/tags' component={XMultiTagView} />
      <Route exact path='/payloads' component={XFileUploadModal} />

      <Route path='/targets/:id' component={XTargetView} />
      <Route path='/tasks/:id' component={XTaskView} />
      <Route path='/jobs/:id' component={XJobView} />

    </XLayout>
  </XGraphProvider>
);

export default App;

