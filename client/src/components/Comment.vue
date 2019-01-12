<template>
  <div style="width:100%">
    <v-layout row wrap style="overflow:hidden;width:100%;" class="pl-1">
      <v-flex style="max-width:64px;">
        <v-avatar size="48" class="my-2" style="cursor:pointer;" @click="$router.push({path:`/profile/${comment.profile_id}`})">
          <img
            v-if="comment.profile_avatar"
            :src="`http://localhost:3000/uploads/${comment.profile_avatar}`"
          >
          <v-icon v-else dark style="background:#bdbdbd" large>people</v-icon>
        </v-avatar>
      </v-flex>
      <v-flex class="pl-1" style="width: calc(100% - 64px)">
        <v-layout row wrap>
          <v-flex xs12>
            <div
              class="mr-1"
              style="display:inline"
              
            >
              <strong @click="$router.push({path:`/profile/${comment.profile_id}`})" style="cursor:pointer;">{{comment.profile_username}}</strong>
            </div>
            <div class="post-time" style="display:inline">{{ comment.created | moment("from") }}</div>
          </v-flex>
          <v-flex xs12 class="mt-1">
            <div class="post-title" @click="changeRoute(`/post/${comment.id}`)">{{comment.value}}</div>
          </v-flex>
        </v-layout>
        <v-layout row wrap class="post-subtext">
          <div style="cursor:pointer;" class="my-1" @click="dialog = true;replyText = replyText || '';">
            <v-icon class="comment-action" small>reply</v-icon>
            <small class="mr-1">({{findValuesHelper(comment,'children',[]).length-1}})</small>
          </div>
          <v-icon style="cursor:pointer;" class="mx-1 my-1" small>save</v-icon>
          <v-icon style="cursor:pointer;" class="mx-1 my-1" small>flag</v-icon>
        </v-layout>
      </v-flex>
    </v-layout>
    <v-dialog v-model="dialog" width="500" class="vert-dialogue">
      <v-card v-if="comment">
        <v-card-title class="justify-center py-2" primary-title>
          <h3>Reply to {{comment.profile_username}}</h3>
        </v-card-title>
        <v-divider/>

        <v-layout row wrap style="overflow:hidden" class="px-2">
          <v-flex style="max-width:64px;">
            <v-avatar size="48" class="mx-2 my-2">
              <img
                class="post-avatar"
                v-if="comment.profile_avatar"
                :src="`http://localhost:3000/uploads/${comment.profile_avatar}`"
              >
              <v-icon v-else dark style="background:#bdbdbd" large>people</v-icon>
            </v-avatar>
          </v-flex>
          <v-flex class="my-2" style="width: calc(100% - 64px)">
            <v-layout row wrap>
              <v-flex xs12>
                <div
                  class="mr-1"
                  style="display:inline"
                  @click="changeRoute(`/profile/${comment.profile_id}`)"
                >
                  <strong>{{comment.profile_username}}</strong>
                </div>
                <div class="post-time" style="display:inline">{{ comment.created | moment("from") }}</div>
              </v-flex>
              <v-flex xs12>
                <div
                  class="post-title"
                  @click="changeRoute(`/post/${comment.id}`)"
                >{{comment.value}}</div>
              </v-flex>
            </v-layout>
            <v-layout row wrap class="post-subtext">
              <v-flex
                xs12
                v-if="comment.description!=='null'"
                class="post-description"
              >{{comment.description}}</v-flex>
            </v-layout>
          </v-flex>
        </v-layout>
        <v-textarea v-model="replyText" auto-grow outline hide-details class="reply-textarea mx-2"/>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="secondary"
            dark
            @click="()=>{replyText = '';dialog = false;}"
            style="color:black;"
            small
          >Cancel</v-btn>
          <v-btn color="primary" dark @click="submitReply()" small>Reply</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>


<script>
import store from "../store";

export default {
  name: "Comment",
  props: ["comment"],
  data() {
    return {
      dialog: false,
      replyText: ""
    };
  },
  methods: {
    findValuesHelper: function(obj, key, list) {
      if (!obj) return list;
      if (obj instanceof Array) {
        for (var i in obj) {
          list = list.concat(this.findValuesHelper(obj[i], key, []));
        }
        return list;
      }
      if (obj[key]) list.push(obj[key]);

      if (typeof obj == "object" && obj !== null) {
        var children = Object.keys(obj);
        if (children.length > 0) {
          for (i = 0; i < children.length; i++) {
            list = list.concat(
              this.findValuesHelper(obj[children[i]], key, [])
            );
          }
        }
      }
      return list;
    },
    submitReply: function() {
      const fd = new FormData();
      fd.append("value", this.replyText);
      fd.append("post_id", this.comment.post_id);
      fd.append("parent_comment_id", this.comment.id);

      this.axios
        .post("http://localhost:3000/api/v1/protected/c", fd, {
          headers: {
            Authorization: `Bearer ${store.state.jwt}`,
            "Content-Type": "multipart/form-data"
          }
        })
        .then(response => {
          this.dialog = false;
          this.replyText = "";
          this.$emit("update:comment", response.data.data.comment);
        });
    }
  }
};
</script>
