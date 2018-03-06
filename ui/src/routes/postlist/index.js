import { h, Component } from "preact";
import style from "./style";

export default class PostList extends Component {
  constructor(props) {
    super(props);

    this.setState({
      loading: true,
      currentPostIdx: 0,
      posts: []
    });
  }

  componentWillReceiveProps(nextProps) {
    if (!this.state.loading && this.props.feedID === nextProps.feedID) {
      return;
    }

    this.setState({ loading: true });

    let key = window.localStorage.getItem("hydrocarbon-key");

    fetch(window.baseURL + "/v1/post/list", {
      method: "POST",
      headers: {
        "x-hydrocarbon-key": key
      },
      body: JSON.stringify({
        feed_id: this.props.feedID
      })
    })
      .then(res => {
        if (res.ok) {
          return res.json();
        }
      })
      .then(json => {
        this.setState({ loading: false, posts: json.posts });
      });
  }

  render({ feedID }, { loading, currentPostIdx, posts }) {
    if (loading) {
      return <div class={style.content}>loading..</div>;
    }

    if (posts.length === 0) {
      return <div class={style.content}> no posts </div>;
    }

    return (
      <div class={style.content}>
        <h1>{posts[currentPostIdx].title}</h1>
        <p>{posts[currentPostIdx].body}</p>
      </div>
    );
  }
}
