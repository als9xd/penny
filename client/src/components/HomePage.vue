<template>
    <v-container>
        <v-layout row wrap>
            <v-flex xs12>
              <v-list two-line>
                <template v-for="(post, index) in posts">
                  <v-subheader
                    v-if="post.header"
                    :key="post.header"
                  >
                    {{ post.header }}
                  </v-subheader>

                  <v-divider
                    v-else-if="post.divider"
                    :inset="post.inset"
                    :key="index"
                  ></v-divider>

                  <v-list-tile
                    v-else
                    :key="post.title"
                    avatar
                    @click="changeRoute(`/p/${post.Id}`)"
                  >
                    <v-list-tile-avatar>
                      <img :src="post.avatar">
                    </v-list-tile-avatar>

                    <v-list-tile-content>
                      <v-list-tile-title v-html="post.title"></v-list-tile-title>
                      <v-list-tile-sub-title v-html="post.subtitle"></v-list-tile-sub-title>
                    </v-list-tile-content>
                  </v-list-tile>
                </template>
              </v-list>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script>
export default {
  name: 'Home',

  methods: {
    changeRoute: function(path){
      this.$router.push({ path }) 
    },
    
  },
  data () {
      return {
        posts: []
      };
  },
  mounted: function(){
    // this.axios.get("http://localhost:3000/api/v1/p").then((response) => {
    //   console.log(response)
    // });
  }
}
</script>

<style scoped>

</style>
