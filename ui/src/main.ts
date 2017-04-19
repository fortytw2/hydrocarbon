// App entry point

import m from 'mithril';
import home from './components/home';
import about from './components/about';

m.route.prefix("")

m.route(document.body, '/', {
	'/': home,
	'/about': about
})
