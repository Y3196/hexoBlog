// axios.js
import axios from 'axios';
import store from './store'; // 根据你的项目结构导入 Vuex store

// 创建一个 axios 实例
const instance = axios.create();

// 请求拦截器
instance.interceptors.request.use(
    config => {
        // 从 Vuex store 或 localStorage 中获取 token
        const token = localStorage.getItem('token') || store.state.token;

        // 如果 token 存在，则将其添加到请求头
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }

        return config;
    },
    error => {
        return Promise.reject(error);
    }
);

// 使用 Vue 原型链挂载 axios 实例
export default {
    install(Vue) {
        Vue.prototype.$axios = instance;
    },
}