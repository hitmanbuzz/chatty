import './index.css';
import { render } from 'solid-js/web';
import 'solid-devtools';

import { Route, Router } from '@solidjs/router';
import Auth from './pages/AuthPage';
import About from './pages/About';

const root = document.getElementById('root');

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    'Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got misspelled?',
  );
}

render(() => (
  <Router>
    <Route path="/" component={Auth} />
    <Route path="/about" component={About} />
  </Router>
), root!
);