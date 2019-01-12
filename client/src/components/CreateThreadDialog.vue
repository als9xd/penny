<template>
  <v-dialog v-model="createThreadDialog" width="500" class="vert-dialogue">
    <v-list-tile slot="activator" avatar>
      <v-list-tile-avatar>
        <v-icon>forum</v-icon>
      </v-list-tile-avatar>
      <v-list-tile-content>
        <v-list-tile-title>
          Thread
        </v-list-tile-title>
      </v-list-tile-content>
    </v-list-tile>
    <v-card>
      <v-card-title class="headline grey lighten-2" primary-title>
        <img class="penny-logo mr-2" alt="Penny logo" src="../assets/penny_new.png">Create Thread
      </v-card-title>

      <v-card-text>
        <v-form v-model="createThreadValid">
          <v-text-field label="Name" v-model="threadName" required outline/>
          <v-textarea label="Description" v-model="threadDescription" required outline/>
          <v-text-field
            outline
            label="Avatar"
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
        <v-btn color="primary" dark @click="submitThread()">Submit</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>

import store from "../store";

export default {
  name: "CreateThreadDialog",
  data() {
    return {
      createThreadDialog: false,
      createThreadValid: true,

      threadName: "",
      threadDescription: "",

      imageName: "",
      imageUrl: "",
      imageFile: ""
    };
  },
  computed: {
    jwt() {
      return store.state.jwt;
    },
  },
  methods: {
    submitThread() {
        const {threadName,threadDescription, imageFile, jwt} = this;
      const fd = new FormData();
      fd.append('name',threadName);
      fd.append('description',threadDescription);
      fd.append('avatar',imageFile);

      this.axios
        .post("http://localhost:3000/api/v1/protected/t",fd,{
            headers: {
            Authorization: `Bearer ${jwt}`,
            'Content-Type': 'multipart/form-data'
            },
        })
        .then(response => {
            store.commit(
                "setSnackBarText",
                `Created thread ${response.data.data.thread.name}`
              );
            this.createThreadDialog = false;
        })
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
