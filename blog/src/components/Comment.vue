<template>
  <div>
    <div class="comment-title"><i class="iconfont iconpinglunzu"/>评论</div>
    <!-- 评论框 -->
    <div class="comment-input-wrapper">
      <div style="display:flex">
        <v-avatar size="40">
          <img
              v-if="this.$store.state.avatar"
              :src="this.$store.state.avatar"
          />
          <img
              v-else
              :src="this.$store.state.blogInfo.websiteConfig.touristAvatar"
          />
        </v-avatar>
        <div style="width:100%" class="ml-3">
          <div class="comment-input">
            <textarea
                class="comment-textarea"
                v-model="commentContent"
                placeholder="留下点什么吧..."
                auto-grow
                dense
            />
          </div>
          <!-- 操作按钮 -->
          <div class="emoji-container">
            <span
                :class="chooseEmoji ? 'emoji-btn-active' : 'emoji-btn'"
                @click="chooseEmoji = !chooseEmoji"
            >
              <i class="iconfont iconbiaoqing"/>
            </span>
            <button
                @click="insertComment"
                class="upload-btn v-comment-btn"
                style="margin-left:auto"
            >
              提交
            </button>
          </div>
          <!-- 表情框 -->
          <emoji @addEmoji="addEmoji" :chooseEmoji="chooseEmoji"/>
        </div>
      </div>
    </div>
    <!-- 评论详情 -->
    <div v-if="count > 0 && reFresh">
      <!-- 评论数量 -->
      <div class="count">{{ count }} 评论</div>
      <!-- 评论列表 -->
      <div
          style="display:flex"
          class="pt-5"
          v-for="(item, index) of commentList"
          :key="item.id"
      >
        <!-- 头像 -->
        <v-avatar size="40" class="comment-avatar">
          <img :src="item.avatar"/>
        </v-avatar>
        <div class="comment-meta">
          <!-- 用户名 -->
          <div class="comment-user">
            <span v-if="!item.webSite">{{ item.nickname }}</span>
            <a v-else :href="item.webSite" target="_blank">
              {{ item.nickname }}
            </a>
            <span class="blogger-tag" v-if="item.userId == 1">博主</span>
          </div>
          <!-- 信息 -->
          <div class="comment-info">
            <!-- 楼层 -->
            <span style="margin-right:10px">{{ count - index }}楼</span>
            <!-- 发表时间 -->
            <span style="margin-right:10px">{{ item.createTime | date }}</span>
            <!-- 点赞 -->
            <span
                :class="isLike(item.id) + ' iconfont icondianzan'"
                @click="like(item)"
            />
            <span v-show="item.likeCount > 0"> {{ item.likeCount }}</span>
            <!-- 回复 -->
            <span class="reply-btn" @click="replyComment(index, item)">
              回复
            </span>
          </div>
          <!-- 评论内容 -->
          <p v-html="item.commentContent" class="comment-content"></p>
          <!-- 回复人 -->
          <div
              style="display:flex"
              v-for="reply of item.replyDTOList"
              :key="reply.id"
          >
            <!-- 头像 -->
            <v-avatar size="36" class="comment-avatar">
              <img :src="reply.avatar"/>
            </v-avatar>
            <div class="reply-meta">
              <!-- 用户名 -->
              <div class="comment-user">
                <span v-if="!reply.webSite">{{ reply.nickname }}</span>
                <a v-else :href="reply.webSite" target="_blank">
                  {{ reply.nickname }}
                </a>
                <span class="blogger-tag" v-if="reply.userId == 1">博主</span>
              </div>
              <!-- 信息 -->
              <div class="comment-info">
                <!-- 发表时间 -->
                <span style="margin-right:10px">
                  {{ reply.createTime | date }}
                </span>
                <!-- 点赞 -->
                <span
                    :class="isLike(reply.id) + ' iconfont icondianzan'"
                    @click="like(reply)"
                />
                <span v-show="reply.likeCount > 0"> {{ reply.likeCount }}</span>
                <!-- 回复 -->
                <span class="reply-btn" @click="replyComment(index, reply)">
                  回复
                </span>
              </div>
              <!-- 回复内容 -->
              <p class="comment-content">
                <!-- 回复用户名 -->
                <template v-if="reply.replyUserId != item.userId">
                  <span v-if="!reply.replyWebSite" class="ml-1">
                    @{{ reply.replyNickname }}
                  </span>
                  <a
                      v-else
                      :href="reply.replyWebSite"
                      target="_blank"
                      class="comment-nickname ml-1"
                  >
                    @{{ reply.replyNickname }}
                  </a>
                  ，
                </template>
                <span v-html="reply.commentContent"/>
              </p>
            </div>
          </div>
          <!-- 回复数量 -->
          <div
              class="mb-3"
              style="font-size:0.75rem;color:#6d757a"
              v-show="item.replyCount > 3"
              ref="check"
          >
            共
            <b>{{ item.replyCount }}</b>
            条回复，
            <span
                style="color:#00a1d6;cursor:pointer"
                @click="checkReplies(index, item)"
            >
              点击查看
            </span>
          </div>
          <!-- 回复分页 -->
          <div
              class="mb-3"
              style="font-size:0.75rem;color:#222;display:none"
              ref="paging"
          >
            <span style="padding-right:10px">
              共{{ Math.ceil(item.replyCount / 5) }}页
            </span>
            <paging
                ref="page"
                :totalPage="Math.ceil(item.replyCount / 5)"
                :index="index"
                :commentId="item.id"
                @changeReplyCurrent="changeReplyCurrent"
            />
          </div>
          <!-- 回复框 -->
          <Reply :type="type" ref="reply" @reloadReply="reloadReply"/>
        </div>
      </div>
      <!-- 加载按钮 -->
      <div class="load-wrapper">
        <v-btn outlined v-if="count > commentList.length" @click="listComments">
          加载更多...
        </v-btn>
      </div>
    </div>
    <!-- 没有评论提示 -->
    <div v-else style="padding:1.25rem;text-align:center">
      来发评论吧~
    </div>
  </div>
</template>

<script>
import Reply from "./Reply";
import Paging from "./Paging";
import Emoji from "./Emoji";
import EmojiList from "../assets/js/emoji";

export default {
  components: {
    Reply,
    Emoji,
    Paging
  },
  props: {
    type: {
      type: Number
    }
  },
  created() {
    //  console.log('Component mounted, fetching comments...');
    this.listComments();
  },
  data: function () {
    return {
      reFresh: true,
      commentContent: "",
      chooseEmoji: false,
      current: 1,
      commentList: [],
      count: 0
    };
  },
  methods: {
    replyComment(index, item) {
      this.$refs.reply.forEach(item => {
        item.$el.style.display = "none";
      });
      //传值给回复框
      this.$refs.reply[index].commentContent = "";
      this.$refs.reply[index].nickname = item.nickname;
      this.$refs.reply[index].replyUserId = item.userId;
      this.$refs.reply[index].parentId = this.commentList[index].id;
      this.$refs.reply[index].chooseEmoji = false;
      this.$refs.reply[index].index = index;
      this.$refs.reply[index].$el.style.display = "block";
    },
    addEmoji(key) {
      this.commentContent += key;
    },
    checkReplies(index, item) {
      this.axios
          .get("/api/comments/" + item.id + "/replies", {
            params: {current: 1, size: 5}
          })
          .then(({data}) => {
            this.$refs.check[index].style.display = "none";
            item.replyDTOList = data.data;
            //超过1页才显示分页
            if (Math.ceil(item.replyCount / 5) > 1) {
              this.$refs.paging[index].style.display = "flex";
            }
          });
    },
    changeReplyCurrent(current, index, commentId) {
      //查看下一页回复
      this.axios
          .get("/api/comments/" + commentId + "/replies", {
            params: {current: current, size: 5}
          })
          .then(({data}) => {
            this.commentList[index].replyDTOList = data.data;
          });
    },

    listComments() {
      const path = this.$route.path;
      const arr = path.split("/");
      const param = {
        current: this.current,
        type: this.type,
        size: 10,
        _timestamp: new Date().getTime(),
        topicId: this.type === 1 || this.type === 3 ? arr[2] : undefined
      };

      //console.log("请求参数：", param);

      const token = this.$store.state.token || localStorage.getItem('token');
      this.axios.get("/api/comments", {
        headers: {
          'Authorization': `Bearer ${token}`
        },
        params: param
      }).then(({data}) => {
//        console.log("获取的评论数据：", data.data.recordList);

        if (data.data && data.data.recordList.length > 0) {
          const transformedComments = data.data.recordList.map(comment => ({
            id: comment.id,
            userId: comment.userId,
            nickname: comment.nickname,
            avatar: comment.avatar,
            commentContent: comment.commentContent,
            createTime: comment.createTime,
            likeCount: comment.likeCount,
            replyCount: comment.replyCount,
            replyDTOList: comment.replyDTOList || []
          }));

          if (this.current === 1) {
            this.commentList = transformedComments;
          } else {
            this.commentList.push(...transformedComments);
          }
          this.current++;
          this.count = data.data.count;
          this.$emit("getCommentCount", this.count);

          this.$nextTick(() => {
          //  console.log("评论列表更新完成：", this.commentList);
          });
        } else {
          console.warn("没有获取到评论数据，recordList 为空。");
        }
      }).catch(error => {
        console.error("请求评论数据时出错：", error);
      });
    }

    ,
    insertComment() {
      const token = this.$store.state.token || localStorage.getItem('token');
      if (!this.$store.state.userId) {
        console.log('User ID is not set, setting loginFlag to true');
        this.$store.state.loginFlag = true;
        this.$toast({type: "error", message: "请先登录"});
        return false;
      }

      if (this.commentContent.trim() === "") {
        this.$toast({type: "error", message: "评论不能为空"});
        return false;
      }

      const reg = /\[.+?\]/g;
      const formattedContent = this.commentContent.replace(reg, (str) => {
        return (
            "<img src='" +
            EmojiList[str] +
            "' width='24' height='24' style='margin: 0 1px; vertical-align: text-bottom'/>"
        );
      });

      const path = this.$route.path;
      const arr = path.split("/");
      const comment = {
        commentContent: formattedContent,
        type: this.type
      };
      if (this.type === 1 || this.type === 3) {
        comment.topicId = arr[2];
      }

      //console.log("提交的评论数据：", comment);

      this.axios.post("/api/comments", comment, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      }).then(({data}) => {

        if (data.data && data.data.recordList && data.data.recordList.length) {
          const newComment = data.data.recordList.map(comment => ({
            id: comment.id,
            userId: comment.userId,
            nickname: comment.nickname,
            avatar: comment.avatar,
            commentContent: comment.commentContent,
            createTime: comment.createTime,
            likeCount: comment.likeCount,
            replyCount: comment.replyCount,
            replyDTOList: comment.replyDTOList || []
          }))[0]; // 获取新评论

          this.commentList.unshift(newComment); // 将新评论添加到列表的顶部

          this.commentContent = "";

          this.count++;
          this.$emit("getCommentCount", this.count);

          const isReview = this.$store.state.blogInfo.websiteConfig.isCommentReview;
          if (isReview) {
            this.$toast({type: "warning", message: "评论成功，正在审核中"});
          } else {
            this.$toast({type: "success", message: "评论成功"});
          }
        } else {
          this.$toast({type: "error", message: "评论提交失败"});
        }
      }).catch(error => {
        console.error("提交评论时出错：", error);
        this.$toast({type: "error", message: "评论提交失败，请稍后再试"});
      });
    }

    ,
    like(comment) {
      // 判断是否已登录
      if (!this.$store.state.userId) {
        this.$toast({type: "error", message: "请先登录"});
        this.$store.state.loginFlag = true;
        return false;
      }
      const token = this.$store.state.token || localStorage.getItem('token');

      // 发送请求
      this.axios
          .post(`/api/comments/${comment.id}/like`, {}, {
            headers: {
              'Authorization': `Bearer ${token}`
            },
          })
          .then(({data}) => {
            if (data.flag) {
              const newLikeCount = data.data; // 后端返回的新点赞数量
              this.$set(comment, "likeCount", newLikeCount); // 更新点赞数量

              // 更新 Vuex 中的点赞状态
              if (this.$store.state.commentLikeSet.includes(comment.id.toString())) {
                this.$store.commit("commentLike", comment.id); // 取消点赞
              } else {
                this.$store.commit("commentLike", comment.id); // 点赞
              }
            }
          })
          .catch(error => {
            console.error("Error processing like:", error);
          });
    }

    ,
    reloadReply(index) {
      const token = this.$store.state.token || localStorage.getItem('token');
      this.axios
          .get("/api/comments/" + this.commentList[index].id + "/replies", {
            headers: {
              'Authorization': `Bearer ${token}`
            },
            params: {
              current: 1,
              size: 5
            }
          })
          .then(({data}) => {
            console.log("Replies data:", data); // Debug log
            if (data.flag && data.data) {
              // 更新 replyDTOList
              this.$set(this.commentList[index], 'replyDTOList', data.data);

              this.commentList[index].replyCount = data.data.length;

              // 更新分页按钮
              this.$nextTick(() => {
                if (this.commentList[index].replyCount > 5) {
                  this.$refs.paging[index].style.display = "flex";
                }
                this.$refs.check[index].style.display = "none";
                this.$refs.reply[index].$el.style.display = "none";
              });
            } else {
              console.error("Failed to load replies:", data.message || "Unknown error");
            }
          }).catch(error => {
        console.error("Error loading replies:", error); // Add this line
      });
    }
  },
  computed: {
    isLike() {
      return (commentId) => {
        var commentLikeSet = this.$store.state.commentLikeSet;
        return commentLikeSet.indexOf(commentId.toString()) !== -1 ? "like-active" : "like";
      };
    }
  },
  watch: {
    commentList() {
     // console.log('CommentList updated in watcher:', this.commentList);
      this.reFresh = false;
      this.$nextTick(() => {
        this.reFresh = true;
      });
    }
  }
};
</script>

<style scoped>
.blogger-tag {
  background: #ffa51e;
  font-size: 12px;
  display: inline-block;
  border-radius: 2px;
  color: #fff;
  padding: 0 5px;
  margin-left: 6px;
}

.comment-title {
  display: flex;
  align-items: center;
  font-size: 1.25rem;
  font-weight: bold;
  line-height: 40px;
  margin-bottom: 10px;
}

.comment-title i {
  font-size: 1.5rem;
  margin-right: 5px;
}

.comment-input-wrapper {
  border: 1px solid #90939950;
  border-radius: 4px;
  padding: 10px;
  margin: 0 0 10px;
}

.count {
  padding: 5px;
  line-height: 1.75;
  font-size: 1.25rem;
  font-weight: bold;
}

.comment-meta {
  margin-left: 0.8rem;
  width: 100%;
  border-bottom: 1px dashed #f5f5f5;
}

.reply-meta {
  margin-left: 0.8rem;
  width: 100%;
}

.comment-user {
  font-size: 14px;
  line-height: 1.75;
}

.comment-user a {
  color: #1abc9c !important;
  font-weight: 500;
  transition: 0.3s all;
}

.comment-nickname {
  text-decoration: none;
  color: #1abc9c !important;
  font-weight: 500;
}

.comment-info {
  line-height: 1.75;
  font-size: 0.75rem;
  color: #b3b3b3;
}

.reply-btn {
  cursor: pointer;
  float: right;
  color: #ef2f11;
}

.comment-content {
  font-size: 0.875rem;
  line-height: 1.75;
  padding-top: 0.625rem;
  white-space: pre-line;
  word-wrap: break-word;
  word-break: break-all;
}

.comment-avatar {
  transition: all 0.5s;
}

.comment-avatar:hover {
  transform: rotate(360deg);
}

.like {
  cursor: pointer;
  font-size: 0.875rem;
}

.like:hover {
  color: #eb5055;
}

.like-active {
  cursor: pointer;
  font-size: 0.875rem;
  color: #eb5055;
}

.load-wrapper {
  margin-top: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.load-wrapper button {
  background-color: #49b1f5;
  color: #fff;
}
</style>
