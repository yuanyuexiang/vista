// utils/wechatAuth.js - 微信授权工具函数

/**
 * 微信授权相关工具函数
 */

// 配置
const WECHAT_CONFIG = {
  // 后端服务地址
  API_BASE: 'https://carture.matrix-net.tech',
  // 本地存储键名
  STORAGE_KEY: 'wechat_user_info',
  // 授权完成标识
  AUTH_SUCCESS_KEY: 'wechat_auth_success'
};

/**
 * 检查是否在微信浏览器中
 */
export const isWechatBrowser = () => {
  const ua = navigator.userAgent.toLowerCase();
  return /micromessenger/.test(ua);
};

/**
 * 检查URL参数中是否包含授权信息
 */
export const checkAuthFromURL = () => {
  const urlParams = new URLSearchParams(window.location.search);
  const authStatus = urlParams.get('wechat_auth');
  const userInfoB64 = urlParams.get('user_info');
  
  if (authStatus === 'success' && userInfoB64) {
    try {
      // 解码用户信息
      const userInfoJson = atob(userInfoB64);
      const userInfo = JSON.parse(userInfoJson);
      
      // 保存到本地存储
      saveUserInfo(userInfo);
      
      // 清理URL参数
      cleanURLParams();
      
      return { success: true, userInfo };
    } catch (error) {
      console.error('解析授权信息失败:', error);
      return { success: false, error };
    }
  }
  
  return { success: false };
};

/**
 * 保存用户信息到本地存储
 */
export const saveUserInfo = (userInfo) => {
  try {
    localStorage.setItem(WECHAT_CONFIG.STORAGE_KEY, JSON.stringify(userInfo));
    localStorage.setItem(WECHAT_CONFIG.AUTH_SUCCESS_KEY, 'true');
    return true;
  } catch (error) {
    console.error('保存用户信息失败:', error);
    return false;
  }
};

/**
 * 从本地存储获取用户信息
 */
export const getUserInfo = () => {
  try {
    const userInfoStr = localStorage.getItem(WECHAT_CONFIG.STORAGE_KEY);
    const isAuthenticated = localStorage.getItem(WECHAT_CONFIG.AUTH_SUCCESS_KEY) === 'true';
    
    if (userInfoStr && isAuthenticated) {
      const userInfo = JSON.parse(userInfoStr);
      // 检查信息是否过期（24小时）
      const loginTime = userInfo.login_time || 0;
      const now = Math.floor(Date.now() / 1000);
      
      if (now - loginTime < 24 * 60 * 60) {
        return userInfo;
      } else {
        // 信息已过期，清除
        clearUserInfo();
      }
    }
    
    return null;
  } catch (error) {
    console.error('获取用户信息失败:', error);
    return null;
  }
};

/**
 * 清除用户信息
 */
export const clearUserInfo = () => {
  localStorage.removeItem(WECHAT_CONFIG.STORAGE_KEY);
  localStorage.removeItem(WECHAT_CONFIG.AUTH_SUCCESS_KEY);
};

/**
 * 发起微信授权
 */
export const startWechatAuth = () => {
  const currentURL = window.location.href;
  const authURL = `${WECHAT_CONFIG.API_BASE}/wechat/auth?redirect_url=${encodeURIComponent(currentURL)}`;
  
  // 在微信浏览器中直接跳转
  if (isWechatBrowser()) {
    window.location.href = authURL;
  } else {
    // 非微信浏览器，显示提示信息
    alert('请在微信中打开此页面进行授权');
  }
};

/**
 * 清理URL参数
 */
const cleanURLParams = () => {
  const url = new URL(window.location);
  url.searchParams.delete('wechat_auth');
  url.searchParams.delete('user_info');
  
  // 使用replace避免在浏览器历史中留下记录
  window.history.replaceState({}, document.title, url.toString());
};

/**
 * 从API获取用户信息
 */
export const fetchUserInfoFromAPI = async (openid) => {
  try {
    const response = await fetch(`${WECHAT_CONFIG.API_BASE}/api/user/${openid}`);
    const data = await response.json();
    
    if (data.code === 200) {
      return data.data.user_info;
    } else {
      throw new Error(data.message || '获取用户信息失败');
    }
  } catch (error) {
    console.error('API获取用户信息失败:', error);
    throw error;
  }
};

/**
 * 获取所有用户列表
 */
export const fetchAllUsers = async () => {
  try {
    const response = await fetch(`${WECHAT_CONFIG.API_BASE}/api/users`);
    const data = await response.json();
    
    if (data.code === 200) {
      return data.data.users;
    } else {
      throw new Error(data.message || '获取用户列表失败');
    }
  } catch (error) {
    console.error('API获取用户列表失败:', error);
    throw error;
  }
};

/**
 * 格式化性别
 */
export const formatGender = (sex) => {
  switch (sex) {
    case 1: return '男';
    case 2: return '女';
    default: return '未知';
  }
};

/**
 * 格式化地址
 */
export const formatLocation = (country, province, city) => {
  return [country, province, city].filter(Boolean).join(' ') || '未知';
};

export default {
  isWechatBrowser,
  checkAuthFromURL,
  saveUserInfo,
  getUserInfo,
  clearUserInfo,
  startWechatAuth,
  fetchUserInfoFromAPI,
  fetchAllUsers,
  formatGender,
  formatLocation,
  WECHAT_CONFIG
};
