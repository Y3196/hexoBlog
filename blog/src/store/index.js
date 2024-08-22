import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    searchFlag: false,
    loginFlag: false,
    registerFlag: false,
    forgetFlag: false,
    emailFlag: false,
    drawer: false,
    loginUrl: "",
    userId: null,
    avatar: null,
    nickname: null,
    intro: null,
    webSite: null,
    loginType: null,
    email: null,
    articleLikeSet: [],
    commentLikeSet: [],
    talkLikeSet: [],
    blogInfo: {},
    token: localStorage.getItem('token') || '', // 从 localStorage 中获取 token

  },
  mutations: {
    setToken(state, token) {
      state.token = token;
      localStorage.setItem('token', token); // 将 token 存储到 localStorage
    },
    updateArticleLikeSet(state, articleId) {
      state.articleLikeSet.push(articleId);
      console.log("Mutation - articleLikeSet updated:", state.articleLikeSet);
    },
    login(state, user) {
      console.log("Logging in user, received data:", user);
      state.userId = user.userInfoId;
      state.avatar = user.avatar;
      state.nickname = user.nickname;
      state.intro = user.intro;
      state.webSite = user.webSite;
      state.articleLikeSet = Array.isArray(user.articleLikeSet) ? user.articleLikeSet : [];
      state.commentLikeSet = Array.isArray(user.commentLikeSet) ? user.commentLikeSet : [];
      state.talkLikeSet = Array.isArray(user.talkLikeSet) ? user.talkLikeSet : [];
      state.email = user.email;
      state.loginType = user.loginType;
      console.log("User logged in, talkLikeSet updated:", state.talkLikeSet);

    },
    logout(state) {
      console.log('Before logout:', JSON.stringify(state)); // 添加日志输出
      state.userId = null;
      state.avatar = null;
      state.nickname = null;
      state.intro = null;
      state.webSite = null;
      state.articleLikeSet = [];
      state.commentLikeSet = [];
      state.talkLikeSet = [];
      state.email = null;
      state.loginType = null;
      state.loginFlag = false; // 添加此行，确保登出时更新 loginFlag
      localStorage.removeItem('token'); // 清除 token
    },
    saveLoginUrl(state, url) {
      state.loginUrl = url;
    },
    saveEmail(state, email) {
      console.log("Updating email in Vuex:", email); // 调试输出
      state.email = email;
    },
    updateUserInfo(state, user) {
      state.nickname = user.nickname;
      state.intro = user.intro;
        state.webSite = user.webSite;
    },
    savePageInfo(state, pageList) {
      state.pageList = pageList;
    },
    updateAvatar(state, avatar) {
      state.avatar = avatar;
    },
    checkBlogInfo(state, blogInfo) {
      state.blogInfo = blogInfo;
    },
    closeModel(state) {
      state.registerFlag = false;
      state.loginFlag = false;
      state.searchFlag = false;
      state.emailFlag = false;
    },
    articleLike(state, articleId) {
      const articleIdStr = articleId.toString();
      const index = state.articleLikeSet.indexOf(articleIdStr);

      if (index !== -1) {
        // 如果已经点赞，则取消点赞
        state.articleLikeSet.splice(index, 1);
      } else {
        // 如果没有点赞，则添加点赞
        state.articleLikeSet.push(articleIdStr);
      }
    }
,
    commentLike(state, commentId) {
      const commentLikeSet = state.commentLikeSet;
      const index = commentLikeSet.indexOf(commentId.toString());

      if (index !== -1) {
        // 如果已经点赞，移除点赞状态
        commentLikeSet.splice(index, 1);
      } else {
        // 如果没有点赞，添加点赞状态
        commentLikeSet.push(commentId.toString());
      }

      console.log("Updated commentLikeSet:", state.commentLikeSet);  // 调试输出
    },
    talkLike(state, talkId) {
      const talkIdStr = talkId.toString();
      const index = state.talkLikeSet.indexOf(talkIdStr);
      if (index === -1) {
        state.talkLikeSet.push(talkIdStr);
        console.log(`Added talkId: ${talkIdStr} to talkLikeSet.`);
      } else {
        state.talkLikeSet.splice(index, 1);
        console.log(`Removed talkId: ${talkIdStr} from talkLikeSet.`);
      }
    }
  },
  actions: {},
  modules: {},
  plugins: [
    createPersistedState({
      storage: window.localStorage,
    })
  ]
});