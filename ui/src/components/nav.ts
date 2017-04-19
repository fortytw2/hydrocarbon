import m from 'mithril'
import * as Mithril from 'mithril'

export default {
    view (vnode: Mithril.Vnode<{}, {}>) {
		return m('div',
			m('a', {href: '/', oncreate: m.route.link}, "Home"),
			m('span', " | "),
			m('a', {href: '/about', oncreate: m.route.link}, "About")
		)
	}
} as Mithril.Component<{},{}>
