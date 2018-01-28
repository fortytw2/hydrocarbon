import { h, Component } from "preact";
import Drawer from "preact-material-components/Drawer";
import List from "preact-material-components/List";
import "preact-material-components/Drawer/style.css";
import "preact-material-components/List/style.css";
import styles from "./style";
import { route } from "preact-router";

import PostList from "../postlist";

export default class Feed extends Component {
  constructor(props) {
    super(props);

    this.setState({
      loading: true,
      feeds: []
    });
  }

  componentDidMount() {
    let key = window.localStorage.getItem("hydrocarbon-key");

    fetch("/v1/folder/list", {
      method: "POST",
      headers: {
        "x-hydrocarbon-key": key
      }
    })
      .then(res => {
        if (res.ok) {
          return res.json();
        }
      })
      .then(json => {
        this.setState({ loading: false, feeds: json[0].feeds });
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

  render({ id }, { loading, feeds }) {
    if (loading) {
      return <div class={styles.content}>Loading Feeds...</div>;
    }

    return (
      <div class={styles.content}>
        <Drawer.PermanentDrawer spacer={false}>
          <Drawer.PermanentDrawerContent>
            <List>
              {feeds.map(f => {
                return (
                  <List.LinkItem onClick={this.goToFeed(f.id)}>
                    {f.title} {f.unread}
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
