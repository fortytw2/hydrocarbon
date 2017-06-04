import * as reducers from "./Reducers";

import { applyMiddleware, combineReducers, compose, createStore } from "redux";
import { autoRehydrate, persistStore } from "redux-persist";
import { routerReducer, syncHistoryWithStore } from "preact-router-redux";

import createBrowserHistory from "history/createBrowserHistory";

let initialState = {
  notifications: [],
  login: { apiKey: "", email: "" }
};

let Store = compose(autoRehydrate())(createStore)(
  combineReducers({ ...reducers, routing: routerReducer }),
  window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__()
);

// Create an enhanced history that syncs navigation events with the store
const History = syncHistoryWithStore(createBrowserHistory(), Store);

export { Store, History };
