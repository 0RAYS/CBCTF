import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  loading: {
    global: false,
    effects: {}, // 用于存储具体操作的loading状态
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

// 选择器
export const selectEffectLoading = (effect) => (state) => state.app.loading.effects[effect];
export const selectError = (state) => state.app.error;

export default appSlice.reducer;
