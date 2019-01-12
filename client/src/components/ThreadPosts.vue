<template>
  <div>
    <template v-for="(result,index) in posts">
      <v-card
        class="post-list-card"
        :key="result.username"
      >
        <v-layout row wrap>
          <v-flex xs12 style="display:flex;align-items:center;"  class="thread-profile" @click="changeRoute(`/thread/${result.thread_id}`)">
            <img
              class="thread-avatar mr-2 mb-1"
              v-if="result.thread_avatar"
              :src="`http://localhost:3000/uploads/${result.thread_avatar}`"
            >
            <div class="thread-title">{{result.thread_name}}</div>
          </v-flex>
        </v-layout>
        <v-layout row wrap style="overflow:hidden">
          <v-flex xs4 style="min-width:250px;text-align:center;" v-if="result.avatar">
            <img
              class="post-avatar"
              v-if="result.avatar"
              :src="`http://localhost:3000/uploads/${result.avatar}`"
            >
          </v-flex>
          <v-flex :class="result.avatar?'md8':'md12'">
            <v-layout row wrap>
              <v-flex xs12>
                <div class="post-title" @click="changeRoute(`/post/${result.id}`)">{{result.name}}</div>
              </v-flex>
            </v-layout>
            <v-layout row wrap class="post-subtext">
              <v-flex xs12 class="mb-2">
                <div class="profile-username" @click="changeRoute(`/profile/${result.profile_id}`)">{{result.profile_username}}</div>
                <div class="post-time">{{ result.created | moment("from") }}</div>
              </v-flex>
              <v-flex
                xs12
                v-if="result.description!=='null'"
                class="post-description"
              >{{result.description}}</v-flex>
              <v-flex xs12 class="post-actions mt-2">
                <div class="post-action px-1"><v-icon class="mr-1" small>comment</v-icon><div style="font-weight:bold;display:inline;">Comments</div></div>
                <div class="post-action px-1"><v-icon class="mr-1" small>share</v-icon><div style="font-weight:bold;display:inline;">Share</div></div>
                <div class="post-action px-1"><v-icon class="mr-1" small>flag</v-icon><div style="font-weight:bold;display:inline;">Report</div></div>
              </v-flex>
            </v-layout>
          </v-flex>
        </v-layout>
      </v-card>
      <v-divider :key="`divider-${result.id}`" v-if="index!==posts.length-1"/>
    </template>
  </div>
</template>


<script>

export default {
  name: "ThreadPosts",
  props: ['posts'],
  methods: {
    changeRoute: function(path) {
      this.$router.push({ path });
    }
  },
};
</script>

<style scoped>

.post-actions {
    display:flex;
    align-items:center;
    cursor:pointer;
}

.post-action:hover {
    background:#e0e0e0;
    border-radius: 5px;
}

.thread-profile {
  cursor: pointer;
}

.post-list-card {
  padding: 1rem;
}

.post-list-card:hover {
    background: #eee;
}

.post-title {
  font-size: 20px;
  font-weight: 400;
  max-height: 4.8rem;
  overflow: hidden;
  cursor: pointer;
}

.thread-title {
    font-weight: 500;
     -webkit-touch-callout: none; /* iOS Safari */
    -webkit-user-select: none; /* Safari */
     -khtml-user-select: none; /* Konqueror HTML */
       -moz-user-select: none; /* Firefox */
        -ms-user-select: none; /* Internet Explorer/Edge */
            user-select: none; /* Non-prefixed version, currently
                                  supported by Chrome and Opera */
}

.thread-profile:hover {
    text-decoration: underline;
}

.thread-avatar {
  height: 32px;
  width: 32px;
  border-radius: 32px;
}

.post-avatar {
  height: 138px;
}

.profile-username{
    cursor: pointer;
}

.profile-username:hover {
    color: black;
}

.post-subtext {
  font-size: 13px;
  font-weight: 400;
  color: #606060;
}

.post-subtext .profile-username{
    display: inline;
}

.post-subtext .post-time {
    display: inline;
}

.post-subtext .post-time::before {
    content: '•';
    margin:  0 0.25rem;
}

.post-description {
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical; 
  overflow: hidden;
}
</style>
