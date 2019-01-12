<template>
  <div v-if="thread">
    <v-tabs-items v-model="model">
      <v-layout row wrap id="thread-header">
        <v-container>
          <v-layout row wrap>
            <v-flex md2 style="text-align:center">
              <v-layout row wrap>
                <v-flex>
                  <img
                    style="height:64px;border-radius:64px;width:64px;"
                    v-if="thread.avatar"
                    :src="`http://localhost:3000/uploads/${thread.avatar}`"
                  >
                  <v-icon v-else size="64">people</v-icon>
                </v-flex>
                <v-flex>
                  <v-btn
                    @click="submitThreadSubscribe"
                    v-if="subscriptions && !subscriptions.threads.filter(t => t.id === threadId).length"
                    color="primary"
                  >Subscribe</v-btn>
                  <v-btn
                    v-else
                    @click="submitThreadSubscribe"
                    color="primary"
                    v-on:mouseover="subscribedHoverText='Unsubscribe'"
                    v-on:mouseleave="subscribedHoverText='Subscribed'"
                  >{{subscribedHoverText}}</v-btn>
                </v-flex>
              </v-layout>
            </v-flex>
            <v-flex md8>
              <div id="thread-name">{{thread.name}}</div>
              {{thread.description}}
            </v-flex>
          </v-layout>
        </v-container>
        <v-flex xs12>
          <v-tabs slot="extension" v-model="model" color="secondary" slider-color="primary">
            <v-tab href="#tab-posts">Posts</v-tab>
            <v-tab href="#tab-rules">Rules</v-tab>
          </v-tabs>
        </v-flex>
      </v-layout>
      <v-tab-item value="tab-posts">
        <v-container>
          <v-layout row wrap>
            <v-flex>
              <ThreadPosts :posts="posts"/>
            </v-flex>
          </v-layout>
        </v-container>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script>
// import ThreadSettings from "../components/ThreadSettings";
import ThreadPosts from "../components/ThreadPosts";

import store from "../store";

export default {
  name: "Thread",
  components: {
    // ThreadSettings,
    ThreadPosts
  },
  computed: {
    profile: function() {
      return store.state.profile;
    },
    subscriptions: function() {
      return store.state.subscriptions;
    },
    jwt() {
      return store.state.jwt;
    }
  },
  props: ["threadId"],
  watch: {
    threadId: {
      handler(val) {
        this.getThreadInfo(val);
      }
    }
  },
  mounted: function() {
    this.getThreadInfo(this.threadId);
  },
  data() {
    return {
      subscribedHoverText: "Subscribed",
      model: "tab-posts",
      thread: null,
      posts: []
    };
  },
  methods: {
    getThreadInfo: function(threadId) {
      if (threadId) {
        this.axios
          .get(`http://localhost:3000/api/v1/t/${threadId}`)
          .then(response => {
            this.thread = response.data.data.thread;
            this.posts = response.data.data.posts;
          });
      } else {
        if (!this.thread) {
          this.$router.push({ path: "/" });
        }
      }
    },
    submitThreadSubscribe: function() {
      this.axios
        .post(
          `http://localhost:3000/api/v1/protected/subscribe/t/${this.threadId}`,
          {},
          {
            headers: {
              Authorization: `Bearer ${this.jwt}`,
              "Content-Type": "application/json"
            }
          }
        )
        .then(response => {
          const subscriptions = Object.assign({}, this.subscriptions);
          if (response.data.data.subscribed) {
            store.commit(
              "setSnackBarText",
              `Subscribed to ${this.thread.name}`
            );
            subscriptions.threads = [this.thread, ...subscriptions.threads];
          } else {
            store.commit(
              "setSnackBarText",
              `Unsubscribed from ${this.thread.name}`
            );
            subscriptions.threads = subscriptions.threads.filter(
              t => t.id !== response.data.data.subscription.thread_id
            );
          }
          store.commit("setSubscriptions", subscriptions);
        })
        .catch(err => {
          if (err.response && err.response.data.error.message) {
            store.commit("setSnackBarText", err.response.data.error.message);
            return;
          }
          throw err;
        });
    }
  }
};
</script>

<style scoped>
#thread-header {
  background-color: #eee;
}
#thread-name {
  font-family: "Roboto", sans-serif;
  font-weight: 400;
  font-size: 2.4rem;
  line-height: 3rem;
}
</style>
