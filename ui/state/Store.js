import * as reducers from "./Reducers";

import { combineReducers, createStore } from "redux";
import { routerReducer, syncHistoryWithStore } from "preact-router-redux";

import createBrowserHistory from "history/createBrowserHistory";

let h = createBrowserHistory();

let Store = createStore(
  combineReducers({
    reducers,
    routing: routerReducer
  }),
  window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__()
);

// Create an enhanced history that syncs navigation events with the store
const history = syncHistoryWithStore(h, Store);

export { Store, history };
