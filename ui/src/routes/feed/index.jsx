import { h, Component } from "preact";
import { Link } from "preact-router";
import { bind } from "decko";

import Modal from "@/components/modal";
import FolderList from "@/components/folder_list";
import PostList from "@/components/post_list";
import CreateFolderForm from "@/components/create_folder_form";
import CreateFeedForm from "@/components/create_feed_form";

import { listFolders, listFeeds } from "@/http";

import style from "./style.css";
import textboxStyle from "@/styles/textbox.css";

const initialState = {
  newFeedModal: false,
  newFolderModal: false,

  folders: []
};

export default class Feed extends Component {
  constructor(props) {
    super(props);

    this.setState(initialState);
  }

  async componentDidMount(props) {
    try {
      let folders = await listFolders({ apiKey: this.props.apiKey });

      // TODO(fortytw2): push this down to the API layer
      folders = await Promise.all(
        folders.map(async f => {
          const feeds = await listFeeds({
            apiKey: this.props.apiKey,
            folderId: f.id
          });
          return {
            ...f,
            feeds
          };
        })
      );

      this.setState({ folders: folders });
    } catch (e) {
      console.log(e);
    }
  }

  @bind
  openNewFeedModal(e) {
    e.preventDefault();
    this.setState({ newFeedModal: true });
  }

  @bind
  closeNewFeedModal(e) {
    e.preventDefault();
    this.setState({ newFeedModal: false });
  }

  @bind
  submitFeed({ name, id }) {
    this.setState({ newFeedModal: false });
    // TODO: add feed to lcoa list
  }

  @bind
  openNewFolderModal(e) {
    e.preventDefault();
    this.setState({ newFolderModal: true });
  }

  @bind
  closeNewFolderModal(e) {
    e.preventDefault();
    this.setState({ newFolderModal: false });
  }

  @bind
  submitFolder({ title, id }) {
    const folders = [{ id, title, feeds: [] }, ...this.state.folders];

    this.setState({ newFolderModal: false, folders: folders });
  }

  render(
    { apiKey, folderId, feedId, postId },
    { newFeedModal, newFolderModal, folders }
  ) {
    return (
      <div class={style.feedContainer}>
        <Modal close={this.closeNewFeedModal} open={newFeedModal}>
          <CreateFeedForm
            onSubmit={this.submitFeed}
            apiKey={apiKey}
            folderId={folderId}
          />
        </Modal>

        <Modal close={this.closeNewFolderModal} open={newFolderModal}>
          <CreateFolderForm onSubmit={this.submitFolder} apiKey={apiKey} />
        </Modal>

        <div class={style.folderList}>
          <div class={style.editBox}>
            <button class={style.editButton} onClick={this.openNewFeedModal}>
              +Feed
            </button>
            <button class={style.editButton} onClick={this.openNewFolderModal}>
              +Folder
            </button>
          </div>
          <FolderList folders={folders} folderId={folderId} feedId={feedId} />
        </div>

        <PostList
          apiKey={apiKey}
          folderId={folderId}
          feedId={feedId}
          postId={postId}
        />
      </div>
    );
  }
}
