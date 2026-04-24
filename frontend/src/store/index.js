import { configureStore } from '@reduxjs/toolkit';
import userReducer from './user';
import appReducer from './app';
import brandingReducer from './branding';
import publicConfigReducer from './settings.js';

export const store = configureStore({
  reducer: {
    user: userReducer,
    app: appReducer,
    branding: brandingReducer,
    publicConfig: publicConfigReducer,
  },
});
