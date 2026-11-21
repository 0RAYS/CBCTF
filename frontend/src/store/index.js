import { configureStore } from '@reduxjs/toolkit';
import userReducer from '@/store/user';
import appReducer from '@/store/app';
import brandingReducer from '@/store/branding';
import publicConfigReducer from '@/store/settings.js';

export const store = configureStore({
  reducer: {
    user: userReducer,
    app: appReducer,
    branding: brandingReducer,
    publicConfig: publicConfigReducer,
  },
});
