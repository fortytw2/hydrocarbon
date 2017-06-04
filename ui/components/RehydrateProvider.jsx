import { Component, h } from "preact";

import Redux from "preact-redux";
import { persistStore } from "redux-persist";

class RehydrateProvider extends Component {
  constructor() {
    super();
    this.state = { rehydrated: false };
  }

  componentWillMount() {
    persistStore(this.props.store, { blacklist: ["routing"] }, () => {
      this.setState({ rehydrated: true });
    });
  }

  render(props, state) {
    if (!this.state.rehydrated) {
      return <div>Loading...</div>;
    }

    return (
      <Redux.Provider store={this.props.store}>
        {this.props.children}
      </Redux.Provider>
    );
  }
}

export default RehydrateProvider;
