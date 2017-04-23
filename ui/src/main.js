// App entry point

import m from 'mithril';
import home from './components/home';
import about from './components/about';
import login from './components/login';

m.route.prefix("")

m.route(document.body, '/', {
	'/': home,
	'/about': about,
	'/login': login
})
