import { h, Component } from "preact";
import { listPosts } from "../../http";
import Button from "preact-material-components/Button";
import "preact-material-components/Button/style.css";
import style from "./style";
import { bind } from "decko";

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

  @bind
  advancePage() {
    this.setState({
      ...this.state,
      currentPostIdx: this.state.currentPostIdx + 1
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
        <Button accept onClick={this.advancePage}>
          Next Page
        </Button>
        <h1>{posts[currentPostIdx].title}</h1>
        <p
          class={style.body}
          dangerouslySetInnerHTML={{ __html: posts[currentPostIdx].body }}
        />
      </div>
    );
  }
}
