import { h, Component } from "preact";
import { bind } from "decko";
import { markRead, listPosts } from "@/http";
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
    await this.fetchData();
  }

  async componentDidUpdate(prevProps) {
    if (this.props.feedId !== prevProps.feedId) {
      this.setState(initialState);
      await this.fetchData();
    } else if (this.props.postId && this.props.postId !== prevProps.postId) {
      const currentPost = this.state.posts.filter(
        p => p.id === this.props.postId
      )[0];
      if (currentPost.read) {
        return;
      }

      try {
        markRead({ apiKey: this.props.apiKey, postId: this.props.postId }).then(
          () => {
            const newPosts = this.state.posts.map(p => {
              if (p.id === this.props.postId) {
                p.read = true;
              }
              return p;
            });
            this.setState({ posts: newPosts });
          }
        );
      } catch (e) {
        console.warn("could not mark as read", e);
        return;
      }
    }
  }

  @bind
  async fetchData(feedId) {
    if (!this.props.feedId) {
      this.setState({ loading: false });
      return;
    }

    const feed = await listPosts({
      apiKey: this.props.apiKey,
      feedId: this.props.feedId
    });
    this.setState({ loading: false, posts: feed.posts });
  }

  @bind
  renderPost(folderId, feedId, postId, post) {
    const friendlyTime = DateTime.fromISO(post.posted_at);
    let displayTime = "";
    if (friendlyTime.year > 1000) {
      displayTime = friendlyTime.toLocaleString(DateTime.DATETIME_SHORT);
    }

    let postStyle = style.post;
    if (!post.read) {
      postStyle = style.unreadPost;
    }

    if (post.id === postId) {
      postStyle = [postStyle, style.activePost].join(" ");
    }

    return (
      <Link
        tabIndex="0"
        activeClassName={style.activeLink}
        href={`/feed/${folderId}/${feedId}/${post.id}`}
      >
        <li class={postStyle}>
          <span class={style.postTitle}>
            {this.truncatePostTitle(post.title)}
          </span>
          <span class={style.postTime}>{displayTime}</span>
        </li>
      </Link>
    );
  }

  truncatePostTitle(title) {
    if (title.length > 37) {
      return title.substring(0, 37) + "...";
    }
    return title;
  }

  @bind
  getActivePost(postId, posts) {
    let post;
    if (!postId) {
      if (posts.length > 0) {
        post = posts[0];
      } else {
        return <h1> no post selected </h1>;
      }
    } else {
      post = posts.filter(p => p.id === postId)[0];
    }

    const friendlyTime = DateTime.fromISO(post.posted_at);
    return (
      <div class={style.postInnerContent}>
        <div class={style.postSubHeader}>
          <h4>{post.title}</h4>
          <h4>{friendlyTime.toLocaleString(DateTime.DATETIME_MED)}</h4>
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
