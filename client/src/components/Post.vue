<template>
  <v-container>
    <v-layout row wrap v-if="post && post.created">
      <v-flex xs12>
        <h1>{{post.name}}</h1>
        {{ post.created | moment("from") }}
      </v-flex>
    </v-layout>
    <v-layout row wrap>
      <v-flex xs12>
        <div v-if="profile">
          <v-layout row wrap>
            <v-flex xs12>
              <v-textarea placeholder="Make a public comment..." v-model="newCommentValue" @focus="newCommentActive = true" hide-details auto-grow rows="1"/>
            </v-flex>
          </v-layout>
          <v-layout row wrap>
            <v-flex xs12 :style="{display:newCommentActive?'block':'none'}">
              <v-btn class="ml-0" @click="submitComment()" color="primary" style="float:right" small>Comment</v-btn>
              <v-btn color="secondary" style="float:right;color:black" small @click="newCommentActive = false;newCommentValue = '';">Cancel</v-btn>
            </v-flex>
          </v-layout>
        </div>
        <v-treeview :open="comments.map(c=>c.id)" :items="comments">
          <template slot="prepend" slot-scope="{ item }">
            <Comment :comment="item" @update:comment="updateComments"/>
          </template>
        </v-treeview>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import store from "../store";

import Comment from "../components/Comment";

var searchObject = function(
  object,
  matchCallback,
  currentPath,
  result,
  searched
) {
  currentPath = currentPath || "";
  result = result || [];
  searched = searched || [];
  if (searched.indexOf(object) !== -1 && object === Object(object)) {
    return;
  }
  searched.push(object);
  if (matchCallback(object)) {
    result.push({ path: currentPath, value: object });
  }
  try {
    if (object === Object(object)) {
      for (var property in object) {
        if (property.indexOf("$") !== 0) {
          //if (Object.prototype.hasOwnProperty.call(object, property)) {
          searchObject(
            object[property],
            matchCallback,
            currentPath + "." + property,
            result,
            searched
          );
          //}
        }
      }
    }
  } catch (e) {
    throw e;
  }
  return result;
};

export default {
  name: "Post",
  computed: {
    activePostId() {
      return store.state.activePostId;
    },
    profile() {
      return store.state.profile;
    }
  },
  components: {
    Comment
  },
  props: ["postId"],
  mounted: function() {
    this.axios
      .get(`http://localhost:3000/api/v1/p/${this.postId}`)
      .then(response => {
        this.post = response.data.data.post;

        this.comments = response.data.data.comments
          .map(c => {
            c.children = response.data.data.comments.filter(
              child => child.parent_comment_id === c.id
            );
            return c;
          })
          .filter(c => c.parent_comment_id === null);
      });
  },
  data() {
    return {
      newCommentActive:false,
      replyToComment: null,
      replyText: {},
      replyDialog: false,
      post: {},
      comments: [],
      newCommentValue: ""
    };
  },
  methods: {
    updateComments: function(comment) {
          if (typeof comment.children === "undefined") comment.children = [];
      if (comment && comment.parent_comment_id === null) {
        this.comments.push(comment)
        return;
      }
      searchObject(this.comments, value => {
        if (value !== null && value.id === comment.parent_comment_id) {
          value.children.unshift(comment);
        }
        return value;
      });
    },
    submitComment: function() {
      const fd = new FormData();
      fd.append("value", this.newCommentValue);
      fd.append("post_id", this.postId);

      this.axios
        .post("http://localhost:3000/api/v1/protected/c", fd, {
          headers: {
            Authorization: `Bearer ${store.state.jwt}`,
            "Content-Type": "multipart/form-data"
          }
        })
        .then(response => {
          console.log(response.data.data.comment);
          this.updateComments(response.data.data.comment);
        });
    }
  }
};
</script>

<style>
.v-treeview-node {
  margin-bottom: 12px;
  margin-top: 12px;
  border: 1px solid #bdbdbd;
  box-shadow: 0 3px 1px -2px rgba(0, 0, 0, 0.2), 0 2px 2px 0 rgba(0, 0, 0, 0.14),
    0 1px 5px 0 rgba(0, 0, 0, 0.12) !important;
  padding: 8px;
  background: white;
  overflow: hidden;
  border-radius: 5px;
  margin-bottom: 0px;
}

.v-treeview-node {
  margin-left: 0 !important;
}
.v-treeview-node .v-treeview-node {
  margin-left: 10px !important;
}

.v-treeview-node__content,
.v-treeview-node__label {
  flex-shrink: 1 !important;
}

.v-treeview-node__label {
}
.v-treeview-node__root {
  height: auto !important;
}

.v-treeview-node__content {
  align-items: start !important;
}

.post-time {
  color: #657786;
}

.reply-textarea textarea {
  margin-top: 0 !important;
}

.post-time::before {
  content: "\00b7";
  color: #657786;
  margin-right: 4px !important;
}
</style>
