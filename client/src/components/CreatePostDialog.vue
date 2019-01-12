<template>
  <v-dialog v-model="createPostDialog" width="500" class="vert-dialogue">
    <v-list-tile slot="activator" avatar>
      <v-list-tile-avatar>
        <v-icon>chat</v-icon>
      </v-list-tile-avatar>
      <v-list-tile-content>
        <v-list-tile-title>Post</v-list-tile-title>
      </v-list-tile-content>
    </v-list-tile>
    <v-card>
      <v-card-title class="headline grey lighten-2" primary-title>
        <v-icon class="mr-2">chat</v-icon>Create Post
      </v-card-title>

      <v-card-text>
        <v-form v-model="createPostValid">
          <v-text-field label="Name" v-model="postName" required outline/>
          <v-textarea label="Description" v-model="postDescription" outline/>
          <v-text-field
            outline
            label="Image"
            @click="pickFile"
            v-model="imageName"
            prepend-icon="attach_file"
          ></v-text-field>
          <input
            type="file"
            style="display: none"
            ref="image"
            accept="image/*"
            @change="onFilePicked"
          >
          <img :src="imageUrl" height="150" v-if="imageUrl">
        </v-form>
      </v-card-text>

      <v-divider></v-divider>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" dark @click="submitPost()">Submit</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import store from "../store";

export default {
  name: "CreatePostDialog",
  computed: {
    jwt() {
      return store.state.jwt;
    },
  },
  data() {
    return {
      createPostDialog: false,
      createPostValid: true,

      postName: null,
      postDescription: null,

      imageName: "",
      imageFile: "",
      imageUrl: ""
    };
  },
  methods: {
    submitPost() {
      const { postName, postDescription, imageFile, jwt } = this;
      const fd = new FormData();
      fd.append("name", postName);
      fd.append("description", postDescription);
      fd.append("avatar", imageFile);
      fd.append("thread_id", this.$route.params.id);
      this.axios
        .post("http://localhost:3000/api/v1/protected/p", fd, {
          headers: {
            Authorization: `Bearer ${jwt}`,
            "Content-Type": "multipart/form-data"
          }
        })
        .then(response => {
          store.commit(
            "setSnackBarText",
            `Created post ${response.data.data.post.name}`
          );
          this.createPostDialog = false;
        });
    },
    pickFile() {
      this.$refs.image.click();
    },
    onFilePicked(e) {
      const files = e.target.files;
      if (files[0] !== undefined) {
        this.imageName = files[0].name;
        if (this.imageName.lastIndexOf(".") <= 0) {
          return;
        }
        const fr = new FileReader();
        fr.readAsDataURL(files[0]);
        fr.addEventListener("load", () => {
          this.imageUrl = fr.result;
          this.imageFile = files[0]; // this is an image file that can be sent to server...
        });
      } else {
        this.imageName = "";
        this.imageFile = "";
        this.imageUrl = "";
      }
    }
  }
};
</script>

<style scoped>
</style>
