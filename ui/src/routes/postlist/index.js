import { h, Component } from "preact";
import { listPosts } from "../../http";
import style from "./style";

export default class PostList extends Component {
  constructor(props) {
    super(props);

    this.setState({
      loading: false,
      currentPostIdx: 0,
      posts: []
    });
  }

  componentDidMount() {
    this.updatePosts(this.props.feedID);
  }

  componentWillReceiveProps({ feedID }) {
    if (this.props.feedID === feedID) {
      return;
    }

    this.updatePosts(feedID);
  }

  updatePosts = feedID => {
    this.setState({
      loading: true,
      currentPostIdx: 0,
      posts: []
    });

    listPosts({ feedID })
      .then(res => {
        if (res.ok) {
          return res.json();
        }
      })
      .then(json => {
        this.setState({ loading: false, posts: json.data.posts });
      });
  };

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
        <p class={style.body}>{posts[currentPostIdx].body}</p>
      </div>
    );
  }
}
