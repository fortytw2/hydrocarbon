import { applyMiddleware, combineReducers, compose, createStore } from "redux";
import { autoRehydrate, persistStore } from "redux-persist";
import { routerReducer, syncHistoryWithStore } from "preact-router-redux";

import createBrowserHistory from "history/createBrowserHistory";
import { login } from "./login/reducers";
import { notifications } from "./notifications/reducers";

let Store = compose(autoRehydrate())(createStore)(
  combineReducers({
    login: login,
    notifications: notifications,
    routing: routerReducer
  }),
  window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__()
);

// Create an enhanced history that syncs navigation events with the store
const History = syncHistoryWithStore(createBrowserHistory(), Store);

export { Store, History };
