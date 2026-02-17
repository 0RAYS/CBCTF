import { configureStore } from '@reduxjs/toolkit';
import userReducer from './user';
import appReducer from './app';

export const store = configureStore({
  reducer: {
    user: userReducer,
    app: appReducer,
  },
});
