import Layout from "@/layout/index.vue";
import router from "../../router";
import store from "../../store";
import axios from "axios";
import Vue from "vue";

export function generaMenu() {
  const token = store.state.token || localStorage.getItem('token');
  axios.get("/api/admin/user/menus", {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }).then(({ data }) => {
    //console.log('API response:', data);
    if (data.data) {
      var userMenuList = data.data;
     // console.log('Menu list:', userMenuList);

      userMenuList.forEach(item => {
        if (item.icon != null) {
          item.icon = "iconfont " + item.icon;
        }
        if (item.component === "Layout") {
          item.component = Layout;
        } else {
          item.component = loadView(item.component);
        }
        if (item.children && item.children.length > 0) {
          item.children.forEach(route => {
            route.icon = "iconfont " + route.icon;
            route.component = loadView(route.component);
          });
        }
      });

      store.commit("saveUserMenuList", userMenuList);

      // 确保每个路由都被添加
      userMenuList.forEach(menu => {
        //console.log('Adding route:', menu);
        router.addRoute(menu);
      });

   //   console.log('All routes:', router.getRoutes());

    } else {
      Vue.prototype.$message.error(data.message);
      router.push({ path: "/login" });
    }
  }).catch(error => {
    Vue.prototype.$message.error('Failed to load user menus');
    console.error('Error loading user menus:', error);
  });
}

export const loadView = view => {
 // console.log('Loading view:', view);
  // 暂时禁用懒加载
  return require(`@/views${view}`).default;
};

