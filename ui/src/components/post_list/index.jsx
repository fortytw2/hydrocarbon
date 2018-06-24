import { h, Component } from "preact";
import { bind } from "decko";
import { listPosts } from "@/http";
import { Link } from "preact-router";
import { DateTime } from "luxon";

import style from "./style.css";

const initialState = {
  loading: true,
  error: null,
  posts: []
};

export default class PostList extends Component {
  constructor(props) {
    super(props);

    this.setState(initialState);
  }

  async componentDidMount() {
    const feed = await listPosts({
      apiKey: this.props.apiKey,
      feedId: this.props.feedId
    });
    this.setState({ loading: false, posts: feed.posts });
  }

  @bind
  renderPost(folderId, feedId, postId, post) {
    const friendlyTime = DateTime.fromISO(post.created_at);

    let postStyle = style.post;
    if (!post.read) {
      postStyle = style.unreadPost;
    }

    return (
      <Link
        tabIndex="0"
        activeClassName={style.activeLink}
        href={`/feed/${folderId}/${feedId}/${post.id}`}
      >
        <li class={postStyle}>
          <span class={style.postTitle}>{post.title}</span>
          <span class={style.postTime}>
            {friendlyTime.toFormat("ccc HH:mm")}
          </span>
        </li>
      </Link>
    );
  }

  @bind
  getActivePost(postId, posts) {
    const post = posts.filter(p => p.id === postId)[0];

    const friendlyTime = DateTime.fromISO(post.created_at);
    return (
      <div class={style.postInnerContent}>
        <div class={style.postHeader}>
          <h2 class={style.postTitleHeader}>{post.title}</h2>
        </div>
        <div class={style.postSubHeader}>
          <h4>{post.title}</h4>
          <h4>{friendlyTime.toFormat("ccc MMM YYYY, HH:mm")}</h4>
        </div>
        <div
          class={style.postBody}
          dangerouslySetInnerHTML={{ __html: post.body }}
        />
      </div>
    );
  }

  render({ folderId, feedId, postId }, { posts, loading, error }) {
    if (loading) {
      return (
        <div class={style.postList}>
          <ol class={style.postListInside}>
            <li>Loading posts...</li>
          </ol>
        </div>
      );
    }

    if (error) {
      return (
        <div class={style.postList}>
          <ol class={style.postListInside}>
            <li>{error}</li>
          </ol>
        </div>
      );
    }

    return (
      <div class={style.postView}>
        <div class={style.postList}>
          <ol class={style.postListInside}>
            {posts.map(p => this.renderPost(folderId, feedId, postId, p))}
          </ol>
        </div>
        <div class={style.postContent}>{this.getActivePost(postId, posts)}</div>
      </div>
    );
  }
}
