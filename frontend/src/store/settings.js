import { createSlice } from '@reduxjs/toolkit';
import { getPublicConfig } from '../api/settings.js';

const initialState = {
  data: {
    registration_enabled: false,
  },
  loading: false,
  loaded: false,
  error: null,
};

const publicConfigSlice = createSlice({
  name: 'publicConfig',
  initialState,
  reducers: {
    setPublicConfig: (state, action) => {
      state.data = {
        ...state.data,
        ...action.payload,
      };
      state.loaded = true;
      state.error = null;
    },
    setPublicConfigLoading: (state, action) => {
      state.loading = action.payload;
    },
    setPublicConfigError: (state, action) => {
      state.error = action.payload;
      state.loaded = true;
    },
  },
});

export const { setPublicConfig, setPublicConfigLoading, setPublicConfigError } = publicConfigSlice.actions;

export const fetchPublicConfig = () => async (dispatch) => {
  dispatch(setPublicConfigLoading(true));
  try {
    const response = await getPublicConfig();
    if (response.code === 200) {
      dispatch(setPublicConfig(response.data ?? {}));
    } else {
      dispatch(setPublicConfigError(response.msg || 'Failed to fetch public config'));
    }
  } catch (error) {
    dispatch(setPublicConfigError(error.message));
  } finally {
    dispatch(setPublicConfigLoading(false));
  }
};

export default publicConfigSlice.reducer;
