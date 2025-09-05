// components/WechatLogin.jsx - 微信登录组件

import React, { useState, useEffect } from 'react';
import { 
  isWechatBrowser, 
  checkAuthFromURL, 
  getUserInfo, 
  clearUserInfo, 
  startWechatAuth,
  formatGender,
  formatLocation 
} from '../utils/wechatAuth';

const WechatLogin = () => {
  const [userInfo, setUserInfo] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    initAuth();
  }, []);

  // 初始化授权检查
  const initAuth = async () => {
    setIsLoading(true);
    setError(null);

    try {
      // 1. 检查URL中是否有授权回调信息
      const authResult = checkAuthFromURL();
      if (authResult.success) {
        setUserInfo(authResult.userInfo);
        setIsLoading(false);
        return;
      }

      // 2. 检查本地存储中是否有用户信息
      const localUserInfo = getUserInfo();
      if (localUserInfo) {
        setUserInfo(localUserInfo);
        setIsLoading(false);
        return;
      }

      // 3. 没有找到用户信息
      setIsLoading(false);
    } catch (err) {
      console.error('授权初始化失败:', err);
      setError(err.message);
      setIsLoading(false);
    }
  };

  // 处理登录
  const handleLogin = () => {
    if (!isWechatBrowser()) {
      setError('请在微信中打开此页面');
      return;
    }

    setError(null);
    startWechatAuth();
  };

  // 处理登出
  const handleLogout = () => {
    clearUserInfo();
    setUserInfo(null);
    setError(null);
  };

  // 刷新用户信息
  const handleRefresh = () => {
    initAuth();
  };

  if (isLoading) {
    return (
      <div style={styles.container}>
        <div style={styles.loading}>
          <div style={styles.spinner}></div>
          <p>正在检查登录状态...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div style={styles.container}>
        <div style={styles.error}>
          <p>❌ {error}</p>
          <button style={styles.button} onClick={handleRefresh}>
            重试
          </button>
        </div>
      </div>
    );
  }

  if (!userInfo) {
    return (
      <div style={styles.container}>
        <div style={styles.loginCard}>
          <h2 style={styles.title}>🔐 微信登录</h2>
          <p style={styles.description}>
            {isWechatBrowser() 
              ? '点击下方按钮进行微信授权登录' 
              : '请在微信中打开此页面进行授权'}
          </p>
          
          {isWechatBrowser() ? (
            <button style={styles.loginButton} onClick={handleLogin}>
              微信授权登录
            </button>
          ) : (
            <div style={styles.wechatTip}>
              <p>📱 请在微信中扫码或分享链接打开</p>
              <code style={styles.url}>{window.location.href}</code>
            </div>
          )}
        </div>
      </div>
    );
  }

  return (
    <div style={styles.container}>
      <div style={styles.userCard}>
        <h2 style={styles.title}>✅ 登录成功</h2>
        
        <div style={styles.userProfile}>
          <img 
            src={userInfo.headimgurl || '/default-avatar.png'} 
            alt="头像"
            style={styles.avatar}
            onError={(e) => {
              e.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iODAiIGhlaWdodD0iODAiIHZpZXdCb3g9IjAgMCA4MCA4MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPGNpcmNsZSBjeD0iNDAiIGN5PSI0MCIgcj0iNDAiIGZpbGw9IiNmMGYwZjAiLz4KPHN2ZyB3aWR0aD0iODAiIGhlaWdodD0iODAiIHZpZXdCb3g9IjAgMCA4MCA4MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHRleHQgeD0iNDAiIHk9IjQ1IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LXNpemU9IjE0IiBmaWxsPSIjOTk5Ij7op6blg488L3RleHQ+Cjwvc3ZnPgo=';
            }}
          />
          
          <div style={styles.userDetails}>
            <h3 style={styles.nickname}>{userInfo.nickname || '未知用户'}</h3>
            <p style={styles.detail}>性别: {formatGender(userInfo.sex)}</p>
            <p style={styles.detail}>
              地区: {formatLocation(userInfo.country, userInfo.province, userInfo.city)}
            </p>
            <p style={styles.openid}>OpenID: {userInfo.openid}</p>
          </div>
        </div>

        <div style={styles.userInfo}>
          <h4>📋 详细信息</h4>
          <pre style={styles.jsonDisplay}>
            {JSON.stringify(userInfo, null, 2)}
          </pre>
        </div>

        <div style={styles.actions}>
          <button style={styles.button} onClick={handleRefresh}>
            🔄 刷新
          </button>
          <button style={styles.logoutButton} onClick={handleLogout}>
            🚪 退出登录
          </button>
        </div>
      </div>
    </div>
  );
};

// 样式
const styles = {
  container: {
    maxWidth: '600px',
    margin: '0 auto',
    padding: '20px',
    fontFamily: 'Arial, sans-serif'
  },
  loading: {
    textAlign: 'center',
    padding: '40px 20px'
  },
  spinner: {
    width: '40px',
    height: '40px',
    border: '4px solid #f3f3f3',
    borderTop: '4px solid #07c160',
    borderRadius: '50%',
    animation: 'spin 1s linear infinite',
    margin: '0 auto 20px'
  },
  error: {
    textAlign: 'center',
    padding: '40px 20px',
    backgroundColor: '#fff2f0',
    border: '1px solid #ffccc7',
    borderRadius: '8px',
    color: '#f5222d'
  },
  loginCard: {
    backgroundColor: 'white',
    padding: '40px 30px',
    borderRadius: '12px',
    boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
    textAlign: 'center'
  },
  userCard: {
    backgroundColor: 'white',
    padding: '30px',
    borderRadius: '12px',
    boxShadow: '0 4px 20px rgba(0,0,0,0.1)'
  },
  title: {
    color: '#333',
    marginBottom: '20px',
    fontSize: '24px'
  },
  description: {
    color: '#666',
    marginBottom: '30px',
    lineHeight: '1.6'
  },
  loginButton: {
    backgroundColor: '#07c160',
    color: 'white',
    border: 'none',
    padding: '15px 30px',
    borderRadius: '8px',
    fontSize: '16px',
    cursor: 'pointer',
    transition: 'background-color 0.3s'
  },
  wechatTip: {
    backgroundColor: '#f6f8fa',
    padding: '20px',
    borderRadius: '8px',
    border: '1px solid #e1e4e8'
  },
  url: {
    display: 'block',
    backgroundColor: '#f1f3f4',
    padding: '10px',
    borderRadius: '4px',
    fontSize: '12px',
    wordBreak: 'break-all',
    marginTop: '10px'
  },
  userProfile: {
    display: 'flex',
    alignItems: 'center',
    marginBottom: '20px',
    paddingBottom: '20px',
    borderBottom: '1px solid #f0f0f0'
  },
  avatar: {
    width: '80px',
    height: '80px',
    borderRadius: '50%',
    marginRight: '20px',
    objectFit: 'cover'
  },
  userDetails: {
    flex: 1
  },
  nickname: {
    margin: '0 0 10px 0',
    fontSize: '20px',
    color: '#333'
  },
  detail: {
    margin: '5px 0',
    color: '#666',
    fontSize: '14px'
  },
  openid: {
    margin: '10px 0 0 0',
    fontSize: '12px',
    color: '#999',
    fontFamily: 'monospace',
    wordBreak: 'break-all'
  },
  userInfo: {
    marginBottom: '20px'
  },
  jsonDisplay: {
    backgroundColor: '#f6f8fa',
    border: '1px solid #e1e4e8',
    borderRadius: '6px',
    padding: '15px',
    fontSize: '12px',
    overflow: 'auto',
    maxHeight: '200px'
  },
  actions: {
    display: 'flex',
    gap: '10px',
    justifyContent: 'center'
  },
  button: {
    backgroundColor: '#1890ff',
    color: 'white',
    border: 'none',
    padding: '10px 20px',
    borderRadius: '6px',
    cursor: 'pointer',
    fontSize: '14px'
  },
  logoutButton: {
    backgroundColor: '#f5222d',
    color: 'white',
    border: 'none',
    padding: '10px 20px',
    borderRadius: '6px',
    cursor: 'pointer',
    fontSize: '14px'
  }
};

export default WechatLogin;
