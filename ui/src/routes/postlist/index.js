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

  componentDidMount() {
    fetch(window.baseURL + "/v1/feed/list", {
      method: "POST",
      body: JSON.stringify({
        id: this.props.id
      }),
      headers: {
        "X-Hydrocarbon-Key": window.localStorage.getItem("hydrocarbon-key")
      }
    })
      .then(res => {
        if (res.ok) {
          return res.json();
        }
      })
      .then(json => {
        this.setState({
          loading: false,
          currentPostIdx: 0,
          posts: json.posts
        });
      });
  }

  render({ id }, { loading, currentPostIdx, posts }) {
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
