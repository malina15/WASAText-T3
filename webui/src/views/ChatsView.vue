<script>
export default {
  data(){
    return { errormsg: null, peers: [] }
  },
  methods:{
    async load(){
      try{
        const id = localStorage.getItem('token')
        const res = await this.$axios.get(`/users/${id}/chats`)
        this.peers = res.data || []
      }catch(e){ this.errormsg = e.toString() }
    },
    open(peer){
      const id = peer && (peer.user_id || peer.id_user || peer.IdUser || peer.id || `${peer}`)
      this.$router.push(`/chats/${encodeURIComponent(id)}`)
    }
  },
  async mounted(){ await this.load() }
}
</script>

<template>
  <div class="container mt-4">
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
    <h3 class="mb-3">Chats</h3>
    <div v-if="peers.length===0" class="text-muted">You are not following anyone yet.</div>
    <ul class="list-group">
      <li v-for="(u,i) in peers" :key="i" class="list-group-item d-flex justify-content-between align-items-center" @click="open(u)">
        <span>{{ u.nickname || u.IdUser || u.id_user }}</span>
        <button class="btn btn-sm btn-primary">Open</button>
      </li>
    </ul>
  </div>
  
</template>

<style scoped>
.list-group-item{ cursor:pointer }
</style>


