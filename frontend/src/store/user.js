import { createSlice } from '@reduxjs/toolkit';
import { getUserInfo, getAccessibleRoutes } from '../api/user';
import i18n from '../i18n';

const initialState = {
  user: null,
  token: localStorage.getItem('token') || null,
  isAuthenticated: false,
  loading: true,
  error: null,
  isAdmin: false,
  routes: [],
};

const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setUser: (state, action) => {
      state.user = action.payload;
      state.isAuthenticated = true;
      state.isAdmin = false;
      localStorage.setItem('userType', 'user');
    },
    setAdmin: (state, action) => {
      state.user = action.payload;
      state.isAdmin = true;
      state.isAuthenticated = true;
      localStorage.setItem('userType', 'admin');
    },
    setRoutes: (state, action) => {
      state.routes = action.payload;
    },
    setLoading: (state, action) => {
      state.loading = action.payload;
    },
    setError: (state, action) => {
      state.error = action.payload;
    },
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
      state.isAdmin = false;
      state.routes = [];
      localStorage.removeItem('token');
      localStorage.removeItem('userType');
    },
    clearError: (state) => {
      state.error = null;
    },
  },
});

export const { setUser, setAdmin, setRoutes, setLoading, setError, logout, clearError } = userSlice.actions;

/**
 * 获取用户信息
 */
export const fetchUserInfo = () => async (dispatch) => {
  dispatch(setLoading(true));

  try {
    const response = await getUserInfo();

    if (response.code === 200) {
      if (response.data.is_admin) {
        dispatch(setAdmin(response.data));
      } else {
        dispatch(setUser(response.data));
      }
    } else {
      const errorMsg = response.msg || i18n.t('toast.user.fetchFailed');
      dispatch(setError(errorMsg));
      dispatch(logout());
    }
  } catch (error) {
    dispatch(setError(error.message));
    dispatch(logout());
  } finally {
    dispatch(setLoading(false));
  }
};

/**
 * 获取当前用户可访问的 API 路由列表，存入 store
 */
export const fetchAccessibleRoutes = () => async (dispatch) => {
  try {
    const response = await getAccessibleRoutes();
    if (response.code === 200) {
      dispatch(setRoutes(response.data));
    }
  } catch {
    // 静默失败，routes 保持为空数组
  }
};

/**
 * 登出用户
 */
export const logoutUser = () => (dispatch) => {
  dispatch(logout());
};

export default userSlice.reducer;
