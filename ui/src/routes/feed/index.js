import { h, Component } from "preact";
import Drawer from "preact-material-components/Drawer";
import List from "preact-material-components/List";
import "preact-material-components/Drawer/style.css";
import "preact-material-components/List/style.css";
import styles from "./style";
import { route } from "preact-router";

import PostList from "../postlist";

export default class Feed extends Component {
  componentWillMount() {
    console.log("fetch feed list here");

    this.setState({
      feeds: [
        {
          name: "ars technica",
          unread: 14,
          id: "arse-technical"
        },
        {
          name: "fucker jones",
          unread: 4,
          id: "fu-jonz"
        },
        {
          name: "test 42",
          unread: 234,
          id: "test-42"
        }
      ]
    });
  }

  linkTo = path => () => {
    route(path);
  };

  goToFeed(id) {
    return this.linkTo("/feed/" + id);
  }

  getContent(id) {
    if (id !== "0") {
      return <PostList id={id} />;
    }
    return <h1>Select a Post!</h1>;
  }

  render({ id }, { feeds }) {
    return (
      <div class={styles.content}>
        <Drawer.PermanentDrawer spacer={false}>
          <Drawer.PermanentDrawerContent>
            <List>
              {feeds.map(f => {
                return (
                  <List.LinkItem onClick={this.goToFeed(f.id)}>
                    {f.name}
                  </List.LinkItem>
                );
              })}
            </List>
          </Drawer.PermanentDrawerContent>
        </Drawer.PermanentDrawer>
        <div>{this.getContent(id)}</div>
      </div>
    );
  }
}
