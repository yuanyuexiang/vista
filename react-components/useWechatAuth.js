// hooks/useWechatAuth.js - 微信授权Hook

import { useState, useEffect, useCallback } from 'react';
import {
  checkAuthFromURL,
  getUserInfo,
  clearUserInfo,
  startWechatAuth,
  isWechatBrowser,
  fetchUserInfoFromAPI
} from '../utils/wechatAuth';

/**
 * 微信授权自定义Hook
 * @param {Object} options 配置选项
 * @param {boolean} options.autoCheck 是否自动检查授权状态
 * @param {Function} options.onAuthSuccess 授权成功回调
 * @param {Function} options.onAuthError 授权失败回调
 */
const useWechatAuth = (options = {}) => {
  const {
    autoCheck = true,
    onAuthSuccess,
    onAuthError
  } = options;

  const [userInfo, setUserInfo] = useState(null);
  const [isLoading, setIsLoading] = useState(autoCheck);
  const [error, setError] = useState(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // 检查授权状态
  const checkAuth = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      // 1. 检查URL参数中的授权信息
      const authResult = checkAuthFromURL();
      if (authResult.success) {
        setUserInfo(authResult.userInfo);
        setIsAuthenticated(true);
        onAuthSuccess?.(authResult.userInfo);
        setIsLoading(false);
        return authResult.userInfo;
      }

      // 2. 检查本地存储
      const localUserInfo = getUserInfo();
      if (localUserInfo) {
        setUserInfo(localUserInfo);
        setIsAuthenticated(true);
        onAuthSuccess?.(localUserInfo);
        setIsLoading(false);
        return localUserInfo;
      }

      // 3. 未找到用户信息
      setIsAuthenticated(false);
      setIsLoading(false);
      return null;
    } catch (err) {
      console.error('检查授权状态失败:', err);
      setError(err.message);
      setIsAuthenticated(false);
      setIsLoading(false);
      onAuthError?.(err);
      return null;
    }
  }, [onAuthSuccess, onAuthError]);

  // 发起登录
  const login = useCallback(() => {
    if (!isWechatBrowser()) {
      const error = new Error('请在微信中打开此页面');
      setError(error.message);
      onAuthError?.(error);
      return false;
    }

    setError(null);
    startWechatAuth();
    return true;
  }, [onAuthError]);

  // 登出
  const logout = useCallback(() => {
    clearUserInfo();
    setUserInfo(null);
    setIsAuthenticated(false);
    setError(null);
  }, []);

  // 刷新用户信息
  const refresh = useCallback(async () => {
    if (userInfo?.openid) {
      try {
        setIsLoading(true);
        const freshUserInfo = await fetchUserInfoFromAPI(userInfo.openid);
        setUserInfo(freshUserInfo);
        return freshUserInfo;
      } catch (err) {
        console.error('刷新用户信息失败:', err);
        setError(err.message);
        onAuthError?.(err);
        return null;
      } finally {
        setIsLoading(false);
      }
    } else {
      return checkAuth();
    }
  }, [userInfo, checkAuth, onAuthError]);

  // 获取特定用户信息
  const fetchUser = useCallback(async (openid) => {
    try {
      setIsLoading(true);
      const user = await fetchUserInfoFromAPI(openid);
      return user;
    } catch (err) {
      console.error('获取用户信息失败:', err);
      setError(err.message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // 清除错误
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // 自动检查授权状态
  useEffect(() => {
    if (autoCheck) {
      checkAuth();
    }
  }, [autoCheck, checkAuth]);

  return {
    // 状态
    userInfo,
    isLoading,
    error,
    isAuthenticated,
    isWechatBrowser: isWechatBrowser(),
    
    // 方法
    login,
    logout,
    refresh,
    checkAuth,
    fetchUser,
    clearError
  };
};

export default useWechatAuth;
