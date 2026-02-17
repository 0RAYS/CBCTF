import { createSlice } from '@reduxjs/toolkit';
import { getUserInfo, getAdminInfo } from '../api/user';
import i18n from '../i18n';

const initialState = {
  user: null,
  token: localStorage.getItem('token') || null,
  isAuthenticated: false,
  loading: true,
  error: null,
  isAdmin: false,
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
      localStorage.removeItem('token');
      localStorage.removeItem('userType');
    },
    clearError: (state) => {
      state.error = null;
    },
  },
});

export const { setUser, setAdmin, setLoading, setError, logout, clearError } = userSlice.actions;

/**
 * 处理API响应的通用函数
 */
const handleApiResponse = (response, dispatch, isAdmin) => {
  if (response.code === 200) {
    if (isAdmin) {
      dispatch(setAdmin(response.data));
    } else {
      dispatch(setUser(response.data));
    }
    return true;
  } else {
    const errorMsg = response.msg || i18n.t(isAdmin ? 'toast.user.fetchAdminFailed' : 'toast.user.fetchFailed');
    dispatch(setError(errorMsg));
    dispatch(logout());
    return false;
  }
};

/**
 * 获取用户信息
 * @param {boolean} isAdmin - 是否为管理员
 */
export const fetchUserInfo =
  (isAdmin = false) =>
  async (dispatch) => {
    dispatch(setLoading(true));

    try {
      // 根据用户类型请求不同API
      const apiMethod = isAdmin ? getAdminInfo : getUserInfo;
      const response = await apiMethod();

      // 处理API响应
      handleApiResponse(response, dispatch, isAdmin);
    } catch (error) {
      dispatch(setError(error.message));
      dispatch(logout());
    } finally {
      dispatch(setLoading(false));
    }
  };

/**
 * 登出用户
 */
export const logoutUser = () => (dispatch) => {
  dispatch(logout());
};

export default userSlice.reducer;
