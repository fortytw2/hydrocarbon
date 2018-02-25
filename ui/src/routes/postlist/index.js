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

  componentDidMount() {}

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
