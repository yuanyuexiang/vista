// App.jsx - React Expo 应用示例

import React from 'react';
import { View, StyleSheet } from 'react-native';
import WechatLogin from './components/WechatLogin';
import useWechatAuth from './hooks/useWechatAuth';

const App = () => {
  // 使用微信授权Hook
  const {
    userInfo,
    isLoading,
    error,
    isAuthenticated,
    login,
    logout,
    refresh
  } = useWechatAuth({
    autoCheck: true,
    onAuthSuccess: (userInfo) => {
      console.log('授权成功:', userInfo);
      // 可以在这里处理授权成功后的逻辑
      // 比如跳转到主页面、发送数据到其他API等
    },
    onAuthError: (error) => {
      console.error('授权失败:', error);
      // 处理授权失败
    }
  });

  return (
    <View style={styles.container}>
      <WechatLogin />
      
      {/* 你的其他组件 */}
      {isAuthenticated && (
        <View style={styles.mainContent}>
          {/* 这里是登录后的主要内容 */}
          <h2>欢迎, {userInfo?.nickname}!</h2>
          <p>这里是你的应用主要内容...</p>
        </View>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  mainContent: {
    padding: 20,
    backgroundColor: 'white',
    margin: 20,
    borderRadius: 10,
  }
});

export default App;
