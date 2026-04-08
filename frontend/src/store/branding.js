import { createSlice } from '@reduxjs/toolkit';
import { getBranding } from '../api/branding';
import { DEFAULT_BRANDING, mergeBranding } from '../config/branding';

const brandingSlice = createSlice({
  name: 'branding',
  initialState: {
    data: DEFAULT_BRANDING,
    loading: false,
    loaded: false,
    error: null,
  },
  reducers: {
    setBranding: (state, action) => {
      state.data = mergeBranding(action.payload);
      state.loaded = true;
      state.error = null;
    },
    setBrandingLoading: (state, action) => {
      state.loading = action.payload;
    },
    setBrandingError: (state, action) => {
      state.error = action.payload;
      state.loaded = true;
    },
  },
});

export const { setBranding, setBrandingLoading, setBrandingError } = brandingSlice.actions;

export const fetchBranding = () => async (dispatch) => {
  dispatch(setBrandingLoading(true));
  try {
    const response = await getBranding();
    if (response.code === 200) {
      dispatch(setBranding(response.data));
    } else {
      dispatch(setBrandingError(response.msg || 'Failed to fetch branding'));
    }
  } catch (error) {
    dispatch(setBrandingError(error.message));
  } finally {
    dispatch(setBrandingLoading(false));
  }
};

export default brandingSlice.reducer;
