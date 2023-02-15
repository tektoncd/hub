import React from 'react';
import { observer } from 'mobx-react';
import '@patternfly/react-core/dist/styles/base.css';
import { Page } from '@patternfly/react-core';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import HomePage from '../../components/HomePage';
import ResourceDetailsPage from '../../components/ResourceDetailsPage';
import Header from '../../containers/Header';
import Footer from '../../components/Footer';
import ParseUrl from '../ParseUrl';
import { createProvider } from '../../store/root';
import './App.css';

const Provider = createProvider();

const App: React.FC = observer(() => {
  return (
    <Provider>
      <Router>
        <ParseUrl />
        <Page header={<Header />} className="hub-page">
          <Routes>
            <Route path="/login" element={<Navigate to="/" replace />} />
            <Route path="/" element={<HomePage />}></Route>
            <Route path="/:catalog/:kind/:name/:version?" element={<ResourceDetailsPage />}></Route>
          </Routes>
          <Footer />
        </Page>
      </Router>
    </Provider>
  );
});

export default App;
