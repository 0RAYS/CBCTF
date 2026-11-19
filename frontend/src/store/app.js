import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  loading: {
    global: false,
    effects: {},
  },
  error: {
    message: null,
    type: null, // 'error' | 'warning' | 'info'
    details: null,
  },
};

const appSlice = createSlice({
  name: 'app',
  initialState,
  reducers: {
    setGlobalLoading: (state, action) => {
      state.loading.global = action.payload;
    },
    setEffectLoading: (state, action) => {
      const { effect, loading } = action.payload;
      state.loading.effects[effect] = loading;
    },
    setError: (state, action) => {
      const { message, type = 'error', details = null } = action.payload;
      state.error = { message, type, details };
    },
    clearError: (state) => {
      state.error = initialState.error;
    },
  },
});

export const { setGlobalLoading, setEffectLoading, setError, clearError } = appSlice.actions;

export default appSlice.reducer;
